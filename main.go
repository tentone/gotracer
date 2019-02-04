package main;

import (
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"gotracer/graphics"
	"gotracer/hitable"
	"gotracer/vmath"
);

var Scene hitable.HitableList;
var Camera *graphics.Camera;

// Max raytracing recursive depth
var MaxDepth int64 = 50;

//If true the last n Frames are blended
var TemporalFilter bool = true;
var TemporalFilterSamples int = 64;
var Frames []*pixel.PictureData;

//If true splits the image generation into threads
var Multithreaded bool = false;
var MultithreadedTheads int = 4;

func run() {
	var width float64 = 320.0;
	var height float64 = 240.0;
	var upscale float64 = 2.0;

	var bounds = pixel.R(0, 0, width, height);
	var windowBounds = pixel.R(0, 0, width * upscale, height * upscale);

	Camera = graphics.NewCameraBounds(bounds);

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

		var picture *pixel.PictureData = RaytraceImage(bounds, false);
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
			Camera.Position.X += 0.1;
			UpdateCamera();
		}
		if window.Pressed(pixelgl.KeyLeft) {
			Camera.Position.X -= 0.1;
			UpdateCamera();
		}
		if window.Pressed(pixelgl.KeyDown) {
			Camera.Position.Y -= 0.1;
			UpdateCamera();
		}
		if window.Pressed(pixelgl.KeyUp) {
			Camera.Position.Y += 0.1;
			UpdateCamera();
		}

		window.Update();
		delta = time.Since(start);

		log.Printf("Frame time %s", delta);
	}
}

// Update the camera viewport
func UpdateCamera(){

	if TemporalFilter {
		Frames = nil;
	}

	Camera.UpdateViewport();
}

func main() {
	// Prepare the scene
	Scene.Add(hitable.NewSphere(100.0, vmath.NewVector3(0.0, -100.5, -1.0), hitable.NewLambertMaterial(vmath.NewVector3(0.4, 0.7, 0.0))));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(0.0, 0.0, 0.0), hitable.NewLambertMaterial(vmath.NewVector3(0.3, 0.2, 0.9))));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(1.0, 0.0, -1.0), hitable.NewMetalMaterial(vmath.NewVector3(0.8, 0.6, 0.2), 0.0)));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -2.0), hitable.NewMetalMaterial(vmath.NewVector3(0.8, 0.8, 0.8), 0.5)));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -1.0), hitable.NewDieletricMaterial(1.5)));
	Scene.Add(hitable.NewSphere(0.4, vmath.NewVector3(-1.0, 1.0, -3.0), hitable.NewNormalMaterial()));
	Scene.Add(hitable.NewSphere(0.3, vmath.NewVector3(-2.0, 2.0, -1.0), hitable.NewDieletricMaterial(0.2)));

	// Start the renderer
	pixelgl.Run(run);
}

// RaytraceImage the scene to calculate the color for a ray.
func RaytraceColor(ray *vmath.Ray, depth int64) *vmath.Vector3 {
	var hitRecord = hitable.NewHitRecord();

	if Scene.Hit(ray, 0.001, math.MaxFloat64, hitRecord) {

		var scattered *vmath.Ray = vmath.NewEmptyRay();
		var attenuation *vmath.Vector3 = vmath.NewVector3(0, 0, 0);

		if depth < MaxDepth && hitRecord.Material.Scatter(ray, hitRecord, attenuation, scattered) {
			var color = attenuation.Clone();
			color.Mul(RaytraceColor(scattered.Clone(), depth + 1));
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

// Raytrace the picure in a thread and write it to the output object.
func RaytraceThread(output *pixel.PictureData, u int, v int, width int, height int) {
	// TODO <ADD CODE HERE>
}

//Render sky with raytrace
func RaytraceImage(bounds pixel.Rect, alialiasing bool) *pixel.PictureData {
	var size = bounds.Size();
	var picture *pixel.PictureData = pixel.MakePictureData(bounds);

	var nx int = int(size.X);
	var ny int = int(size.Y);

	if Multithreaded {
		//TODO <CALCULATE THREAD RANGE>
		//TODO <CHECK HOW TO STORE RESULT>
		go RaytraceThread(picture, 0, 0, nx, ny);
	} else {
		// Single threaded
		for j := 0; j < ny; j++ {
			for i := 0; i < nx; i++ {
				var color *vmath.Vector3;

				//If using antialiasing jitter the UV and cast multiple rays
				if alialiasing {
					var samples int = 16;
					color = vmath.NewVector3(0, 0, 0);

					for k := 0; k < samples; k++ {
						var u float64 = (float64(i) + rand.Float64()) / size.X;
						var v float64 = (float64(j) + rand.Float64()) / size.Y;
						color.Add(RaytraceColor(Camera.GetRay(u, v), 0));
					}

					color.DivideScalar(float64(samples));
				} else {
					var u float64;
					var v float64;

					if TemporalFilter {
						u = (float64(i) + rand.Float64()) / size.X;
						v = (float64(j) + rand.Float64()) / size.Y;
					} else {
						u = float64(i) / size.X;
						v = float64(j) / size.Y;
					}

					color = RaytraceColor(Camera.GetRay(u, v), 0);
				}

				//Apply gamma
				color.DivideScalar(1.0);
				color.Sqrt();

				color.MulScalar(255);

				//Write to picture
				var index = picture.Index(pixel.Vec{X:float64(i), Y:float64(j)});
				picture.Pix[index].R = uint8(color.X);
				picture.Pix[index].G = uint8(color.Y);
				picture.Pix[index].B = uint8(color.Z);
			}
		}
	}



	return picture;
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
