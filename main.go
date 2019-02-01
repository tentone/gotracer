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

var world hitable.HitableList;
var camera *graphics.Camera;

func run() {
	var width = 640;
	var height = 480;

	var bounds = pixel.R(0, 0, float64(width), float64(height));

	camera = graphics.NewCameraBounds(bounds);

	var config = pixelgl.WindowConfig{
		Resizable: false,
		Undecorated: false,
		VSync: false,
		Title: "Gotracer",
		Bounds: bounds};

	var window, err = pixelgl.NewWindow(config);
	
	CheckError(err);

	for !window.Closed() {
		
		var start = time.Now();

		window.Clear(colornames.Black);

		var picture = Raytrace(bounds, false);

		var sprite = pixel.NewSprite(picture, picture.Bounds());
		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()));	

		window.Update();

		var delta = time.Since(start);
		log.Printf("Frame time %s", delta);
	}
}

func main() {
	// Prepare the scene
	world.Add(hitable.NewSphere(0.5, vmath.NewVector3(0.0, 0.0, -1.0), hitable.NewLambertMaterial(vmath.NewVector3(1.0, 1.0, 1.0))));
	world.Add(hitable.NewSphere(100.0, vmath.NewVector3(0.0, -100.5, -1.0), hitable.NewLambertMaterial(vmath.NewVector3(1.0, 1.0, 1.0))));

	// Start the renderer
	pixelgl.Run(run)
}

// Calculate a random unitary vector in the surface of a sphere.
func RandomInUnitSphere() *vmath.Vector3 {
	var p *vmath.Vector3 = vmath.NewVector3(0, 0, 0);

	for {
		p.Set(rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0);

		if p.SquaredLength() < 1.0 {
			break
		}
	}

	return p;
}

// Raytrace the scene to calculate the color for a ray.
func CalculateColor(ray *vmath.Ray) *vmath.Vector3 {
	var rec = hitable.NewHitRecord();

	if world.Hit(ray, 0.001, math.MaxFloat64, rec) {

		//var scattered *vmath.Ray;
		//var attenuation *vmath.Vector3;

		var target *vmath.Vector3 = vmath.NewVector3(0, 0, 0);
		target.Add(rec.Normal);
		target.Add(RandomInUnitSphere());

		var color *vmath.Vector3 = CalculateColor(vmath.NewRay(rec.P, target));
		color.MulScalar(0.5);
		return color;

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
func Raytrace(bounds pixel.Rect, alialiasing bool) *pixel.PictureData {
	var size = bounds.Size();
	var picture = pixel.MakePictureData(bounds);

	var nx int = int(size.X);
	var ny int = int(size.Y);

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			var color *vmath.Vector3;

			//If using antialiasing jitter the UV and cast multiple rays
			if alialiasing {
				var samples = 16;
				color = vmath.NewVector3(0, 0, 0);

				for k := 0; k < samples; k++ {
					var u = (float64(i) + rand.Float64()) / size.X;
					var v = (float64(j) + rand.Float64()) / size.Y;
					var ray = camera.GetRay(u, v);
					color.Add(CalculateColor(ray));
				}

				color.DivideScalar(float64(samples));
			} else {
				var u = float64(i) / size.X;
				var v = float64(j) / size.Y;
				var ray = camera.GetRay(u, v);

				color = CalculateColor(ray);
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
