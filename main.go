package main;

import "os";
import "strconv";
import "fmt";
//import "bufio"
//import "io/ioutil"

//import "github.com/faiface/pixel/pixelgl";

import "gotracer/math";


func main() {
	//Render();

	var a *math.Vector3 = math.NewVector3(1, 1, 1);
	var b *math.Vector3 = math.NewVector3(1, 2, 3);
	a.Add(b);

	fmt.Println(b.ToString());
	fmt.Println(a.ToString());
}

//Render the image
func Render() {
	var nx int = 640;
	var ny int = 480;

	var file, err = os.Create("output.ppm");
	CheckError(err);

	file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n");

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {

			//Calculate color
			var color = math.NewVector3(float64(i) / float64(nx), float64(j) / float64(ny), 0.2)

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