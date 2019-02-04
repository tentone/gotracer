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

//If true the last n frames are blended
var TemporalFilter bool = false;
var TemporalFilterSamples int = 64;

//If true splits the image generation into threads
var Multithreaded bool = false;
var MultithreadedTheads int = 4;

func run() {
	var width = 640;
	var height = 480;

	var bounds = pixel.R(0, 0, float64(width), float64(height));

	Camera = graphics.NewCameraBounds(bounds);

	var config = pixelgl.WindowConfig{
		Resizable: false,
		Undecorated: false,
		VSync: false,
		Title: "Gotracer",
		Bounds: bounds};

	var window, err = pixelgl.NewWindow(config);
	
	CheckError(err);

	var frames []*pixel.PictureData;

	for !window.Closed() {
		
		var start = time.Now();

		window.Clear(colornames.Black);

		var picture = RaytraceImage(bounds, false);

		if TemporalFilter {
			// Add new frame to the list
			frames = append(frames, picture);
			if len(frames) > TemporalFilterSamples {
				frames = frames[1:];
			}

			var final = pixel.MakePictureData(bounds);

			// Average the frames in the list
			for i := 0; i < len(final.Pix); i++ {

				var r, g, b int;

				for j := 0; j < len(frames); j++ {
					r += (int)(frames[j].Pix[i].R);
					g += (int)(frames[j].Pix[i].G);
					b += (int)(frames[j].Pix[i].B);
				}

				final.Pix[i].R = (uint8)(r / len(frames));
				final.Pix[i].G = (uint8)(g / len(frames));
				final.Pix[i].B = (uint8)(b / len(frames));
			}

			var sprite = pixel.NewSprite(final, final.Bounds());
			sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()));
		} else {
			var sprite = pixel.NewSprite(picture, picture.Bounds());
			sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()));
		}

		//Keyboard input
		if window.Pressed(pixelgl.KeyLeft) {
			Camera.Position.X += 0.1;
			Camera.UpdateViewport();
		}
		if window.Pressed(pixelgl.KeyRight) {
			Camera.Position.X -= 0.1;
			Camera.UpdateViewport();
		}
		if window.Pressed(pixelgl.KeyDown) {
			Camera.Position.Y -= 0.1;
			Camera.UpdateViewport();
		}
		if window.Pressed(pixelgl.KeyUp) {
			Camera.Position.Y += 0.1;
			Camera.UpdateViewport();
		}

		window.Update();

		var delta = time.Since(start);
		log.Printf("Frame time %s", delta);
	}
}

func main() {
	// Prepare the scene
	Scene.Add(hitable.NewSphere(100.0, vmath.NewVector3(0.0, -100.5, -1.0), hitable.NewLambertMaterial(vmath.NewVector3(0.8, 0.8, 0.0))));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(0.0, 0.0, 0.0), hitable.NewLambertMaterial(vmath.NewVector3(0.8, 0.3, 0.3))));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(1.0, 0.0, -1.0), hitable.NewMetalMaterial(vmath.NewVector3(0.8, 0.6, 0.2), 0.4)));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -2.0), hitable.NewMetalMaterial(vmath.NewVector3(0.8, 0.8, 0.8), 0.2)));
	Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -1.0), hitable.NewDieletricMaterial(1.5)));

	//Scene.Add(hitable.NewSphere(0.4, vmath.NewVector3(-1.0, 1.0, -1.0), hitable.NewDieletricMaterial(0.2)));
	//Scene.Add(hitable.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -1.0), hitable.NewNormalMaterial()));

	// Start the renderer
	pixelgl.Run(run)
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
			return vmath.NewVector3(0, 0, 0);
			//TODO <EXPERIMENT USING THE SCATTERED RAY VALUE>
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

//Render sky with raytrace
func RaytraceImage(bounds pixel.Rect, alialiasing bool) *pixel.PictureData {
	var size = bounds.Size();
	var picture = pixel.MakePictureData(bounds);

	var nx int = int(size.X);
	var ny int = int(size.Y);

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
