package main;

import "os";
import "strconv";
import "fmt";
//import "bufio"
//import "io/ioutil"

//import "github.com/faiface/pixel/pixelgl";

import "gotracer/math";
import "gotracer/graphics";

func main() {
	//Create window
	//var config = new(pixelgl.WindowConfig);
	//config.Resizable = true;
	//config.Undecorated = false;
	//config.VSync = false;
	//config.Title = "Gotracer";

	//var bounds = new(pixelgl);
	//config.Bounds = ;

	//First raytrace
	RenderSky();

	//Gradient
	//RenderGradient();

	var a *math.Vector3 = math.NewVector3(1, 1, 1);
	var b *math.Vector3 = math.NewVector3(1, 2, 3);
	a.Add(b);

	var frame = new(graphics.Frame);
	fmt.Println(frame.ToString());

	fmt.Println(b.ToString());
	fmt.Println(a.ToString());
}

func SkyColor(r *math.Ray) *math.Vector3 {
	var unitDirection = r.Direction.UnitVector();
	var t = 0.5 * (unitDirection.Y + 1.0);

	var a = math.NewVector3(1.0, 1.0, 1.0);
	a.MulScalar(1.0 - t);

	var b = math.NewVector3(0.5, 0.7, 1.0);
	b.MulScalar(t);

	a.Add(b);
	
	return a;
}

//Render sky with raytrace
func RenderSky() {
	var nx int = 200;
	var ny int = 100;

	var lowerLeftCorner = math.NewVector3(-2.0, -1.0, -1.0);
	//var horizontal = math.NewVector3(4.0, 0.0, 0.0);
	//var vertical = math.NewVector3(0.0, 2.0, 0.0);
	var origin = math.NewVector3(0.0, 0.0, 0.0);

	var file, err = os.Create("sky.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			//var u = float64(i) / float64(nx);
			//var v = float64(j) / float64(ny);

			var direction = lowerLeftCorner.Clone();

			var ray = math.NewRay(origin, direction);

			//Calculate color
			var color = SkyColor(ray);

			var ir int = int(256 * color.X);
			var ig int = int(256 * color.Y);
			var ib int = int(256 * color.Z);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
		}
	}

	//Close file
	file.Sync();
	file.Close();
}

//RenderGradient the image
func RenderGradient() {
	var nx int = 640;
	var ny int = 480;

	var file, err = os.Create("gradient.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			//Calculate color
			var color = math.NewVector3(float64(i) / float64(nx), float64(j) / float64(ny), 0.2);

			var ir int = int(256 * color.X);
			var ig int = int(256 * color.Y);
			var ib int = int(256 * color.Z);

			file.WriteString(strconv.Itoa(ir) + " " + strconv.Itoa(ig) + " " + strconv.Itoa(ib) + "\n");
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