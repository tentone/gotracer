package main;

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"gotracer/graphics"
	"gotracer/hitable"
	"gotracer/vmath"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
);

// Max raytracing recursive depth
var MaxDepth int64 = 100;

//If true multiple rays are casted and blended for each pixel
var Antialiasing bool = false;

//If true the last n Frames are blended
var TemporalFilter bool = true;
var TemporalFilterSamples int = 64;
var Frames []*pixel.PictureData;

//If true splits the image generation into threads
var Multithreaded bool = true;
var MultithreadedTheads int = 4;
var SceneCopies []*hitable.HitableList;
var CameraCopies []*graphics.CameraDefocus;

func main() {
	//runtime.GOMAXPROCS(8);
	pixelgl.Run(run);
}

func run() {
	var width float64 = 320.0;
	var height float64 = 240.0;
	var upscale float64 = 2.0;

	var bounds = pixel.R(0, 0, width, height);
	var windowBounds = pixel.R(0, 0, width * upscale, height * upscale);

	// Prepare the scene
	var scene *hitable.HitableList = hitable.NewHitableList();
	scene.Add(hitable.NewSphere(500.0, vmath.NewVector3(0.0, -500.5, -1.0), hitable.NewLambertMaterial(vmath.NewVector3(0.4, 0.7, 0.0))));
	scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -3.0), hitable.NewNormalMaterial()));
	scene.Add(hitable.NewSphere(1.5, vmath.NewVector3(5.0, 1.0, -6.0), hitable.NewDieletricMaterial(1.3, vmath.NewVector3(0.90, 0.90, 0.90))));
	scene.Add(hitable.NewSphere(1.5, vmath.NewVector3(-1.0, 1.0, -3.0), hitable.NewMetalMaterial(vmath.NewVector3(0.6, 0.6, 0.6), 0.1)));

	// Place random objects
	for i := 0; i < 50; i++ {

		var min float64 = 15.0;
		var distance float64 = 30.0;

		var radius float64 = 0.4 + rand.Float64() * 0.2;
		var position *vmath.Vector3 = vmath.NewVector3(rand.Float64() * distance - min, radius - 0.5, rand.Float64() * distance - min);
		scene.Add(hitable.NewSphere(radius, position, hitable.NewLambertMaterial(vmath.NewRandomVector3(0.1, 1))));

		radius = 0.4 + rand.Float64() * 0.2;
		position = vmath.NewVector3(rand.Float64() * distance - min, radius - 0.5, rand.Float64() * distance - min);
		scene.Add(hitable.NewSphere(radius, position, hitable.NewMetalMaterial(vmath.NewRandomVector3(0.1, 1), rand.Float64())));

		radius = 0.4 + rand.Float64() * 0.2;
		position = vmath.NewVector3(rand.Float64() * distance - min, radius - 0.5, rand.Float64() * distance - min);
		scene.Add(hitable.NewSphere(radius, position, hitable.NewDieletricMaterial(2.0 * rand.Float64(), vmath.NewRandomVector3(0.95, 1.0))));
	}

	var camera *graphics.CameraDefocus = graphics.NewCameraDefocusBounds(bounds);

	if Multithreaded {
		for i := 0; i < MultithreadedTheads; i++ {
			SceneCopies = append(SceneCopies, scene.Clone());
			CameraCopies = append(CameraCopies, camera.Clone());
		}
	}

	var config = pixelgl.WindowConfig{
		Resizable: false,
		Undecorated: false,
		VSync: false,
		Title: "Gotracer",
		Bounds: windowBounds};

	var window, err = pixelgl.NewWindow(config);
	
	CheckError(err);

	var delta time.Duration;

	for !window.Closed() {
		
		var start time.Time = time.Now();

		window.Clear(colornames.Black);

		var picture *pixel.PictureData = Render(bounds, scene, camera);
		var sprite *pixel.Sprite;

		if TemporalFilter {

			// Add new frame to the list
			Frames = append(Frames, picture);
			if len(Frames) > TemporalFilterSamples {
				Frames = Frames[1:];
			}

			var final = pixel.MakePictureData(bounds);

			// Average the Frames in the list
			for i := 0; i < len(final.Pix); i++ {

				var r, g, b int;

				for j := 0; j < len(Frames); j++ {
					r += (int)(Frames[j].Pix[i].R);
					g += (int)(Frames[j].Pix[i].G);
					b += (int)(Frames[j].Pix[i].B);
				}

				final.Pix[i].R = (uint8)(r / len(Frames));
				final.Pix[i].G = (uint8)(g / len(Frames));
				final.Pix[i].B = (uint8)(b / len(Frames));
			}

			sprite = pixel.NewSprite(final, final.Bounds());
		} else {
			sprite = pixel.NewSprite(picture, picture.Bounds());
		}
		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()).Scaled(window.Bounds().Center(), upscale));

		//Keyboard input
		if window.Pressed(pixelgl.KeyRight) {
			camera.Position.X += 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyLeft) {
			camera.Position.X -= 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyUp) {
			camera.Position.Z -= 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyDown) {
			camera.Position.Z += 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyLeftControl) || window.Pressed(pixelgl.KeyRightControl) {
			camera.Position.Y -= 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeySpace) {
			camera.Position.Y += 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyW) {
			camera.Aperture += 0.1;
			UpdateCamera(camera);
		}
		if window.Pressed(pixelgl.KeyS) {
			camera.Aperture -= 0.1;
			UpdateCamera(camera);
		}

		window.Update();

		delta = time.Since(start);
		log.Printf("Frame time %s", delta);
	}
}

// Update the camera viewport
func UpdateCamera(camera *graphics.CameraDefocus){

	if TemporalFilter {
		Frames = nil;
	}

	if Multithreaded {
		for i := 0; i < MultithreadedTheads; i++ {
			CameraCopies[i].Copy(camera);
			CameraCopies[i].UpdateViewport();
		}
	} else {
		camera.UpdateViewport();
	}


}

//Render image the image
func Render(bounds pixel.Rect, scene *hitable.HitableList, camera *graphics.CameraDefocus) *pixel.PictureData {
	var size = bounds.Size();
	var picture *pixel.PictureData = pixel.MakePictureData(bounds);
	var nx int = int(size.X);
	var ny int = int(size.Y);
	var wg sync.WaitGroup;

	if Multithreaded {
		wg.Add(MultithreadedTheads);
		var wtx = nx / MultithreadedTheads;
		var itx = 0;
		for i := 0; i < MultithreadedTheads; i++ {
			go RaytraceThread(&wg, picture, SceneCopies[i], CameraCopies[i], MaxDepth, TemporalFilter, Antialiasing, size.X, size.Y, itx, 0, itx + wtx, ny);
			itx += wtx;
		}

		wg.Wait();
	} else {
		wg.Add(1);
		RaytraceThread(&wg, picture, scene, camera, MaxDepth, TemporalFilter, Antialiasing, size.X, size.Y, 0, 0, nx, ny);
	}

	return picture;
}

// Raytrace the picure in a thread and write it to the output object.
// The result is writen to the picture object passed as argument.
// This method is intended to be called multiple threads.
func RaytraceThread(wg *sync.WaitGroup, picture *pixel.PictureData, scene *hitable.HitableList, camera *graphics.CameraDefocus, depth int64, jitter bool, antialiasing bool, width float64, height float64, ix int, iy int, nx int, ny int) {
	for j := iy; j < ny; j++ {
		for i := ix; i < nx; i++ {
			var color *vmath.Vector3;

			//If using antialiasing jitter the UV and cast multiple rays
			if antialiasing {
				var samples int = 4;
				color = vmath.NewVector3(0, 0, 0);

				for k := 0; k < samples; k++ {
					var u float64 = (float64(i) + rand.Float64()) / width;
					var v float64 = (float64(j) + rand.Float64()) / height;
					color.Add(RaytraceScene(scene, camera.GetRay(u, v), depth));
				}

				color.DivideScalar(float64(samples));
			} else {
				var u float64;
				var v float64;

				if jitter {
					u = (float64(i) + rand.Float64()) / width;
					v = (float64(j) + rand.Float64()) / height;
				} else {
					u = float64(i) / width;
					v = float64(j) / height;
				}

				color = RaytraceScene(scene, camera.GetRay(u, v), depth);
			}

			//Apply gamma
			//color.DivideScalar(1.0);
			color.Sqrt();

			color.MulScalar(255);

			//Write to picture
			var index = picture.Index(pixel.Vec{X:float64(i), Y:float64(j)});
			picture.Pix[index].R = uint8(color.X);
			picture.Pix[index].G = uint8(color.Y);
			picture.Pix[index].B = uint8(color.Z);
		}
	}

	wg.Done();
}

// Render the scene to calculate the color for a ray.
// Receives the scene and the initial ray to be casted.
// It is called recursively until the ray does not hit anything, it is absorved of depth reaches 0.
func RaytraceScene(scene *hitable.HitableList, ray *vmath.Ray, depth int64) *vmath.Vector3 {
	var hitRecord = hitable.NewHitRecord();

	if scene.Hit(ray, 0.001, math.MaxFloat64, hitRecord) {

		var scattered *vmath.Ray = vmath.NewEmptyRay();
		var attenuation *vmath.Vector3 = vmath.NewVector3(0, 0, 0);

		if depth > 0 && hitRecord.Material.Scatter(ray, hitRecord, attenuation, scattered) {
			var color = attenuation.Clone();
			color.Mul(RaytraceScene(scene, scattered.Clone(), depth - 1));
			return color;
		} else {
			// Ray was absorved return black
			//return vmath.NewVector3(0, 0, 0);

			// The ray was absorved use the last value
			return attenuation.Clone();
		}

	} else {

		return BackgroundColor(ray);
	}
}

// Calculate the background color from ray.
// This method is used for multi threading.
func BackgroundColor(r *vmath.Ray) *vmath.Vector3 {
	var unitDirection = r.Direction.UnitVector();
	var t = 0.5 * (unitDirection.Y + 1.0);

	var a = vmath.NewVector3(1.0, 1.0, 1.0);
	a.MulScalar(1.0 - t);

	var b = vmath.NewVector3(0.5, 0.7, 1.0);
	b.MulScalar(t);

	a.Add(b);

	return a;
}

// Write the frame to a PPM file string
func WritePPM(picture *pixel.PictureData, fname string) {
	var size = picture.Rect.Size();

	var nx int = int(size.X);
	var ny int = int(size.Y);
	
	var file, err = os.Create("sky.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			//Write to file
			var index = picture.Index(pixel.Vec{X:float64(i), Y:float64(j)});
			file.WriteString(strconv.Itoa(int(picture.Pix[index].R)) + " " + strconv.Itoa(int(picture.Pix[index].G)) + " " + strconv.Itoa(int(picture.Pix[index].B)) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();
}

//CheckError an error	
func CheckError(e error) {
	if e != nil {
		panic(e);
	}
}
