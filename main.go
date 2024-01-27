package main

import (
	"bytes"
	"gotracer/camera"
	"gotracer/geometry"
	"gotracer/material"
	"gotracer/vmath"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/v2"
	"github.com/sheenobu/go-obj/obj"
	"golang.org/x/image/colornames"
)

// Render size
const Width float64 = 640.0
const Height float64 = 480.0
const Upscale float64 = 1.0

// Max raytracing recursive depth
const MaxDepth int64 = 50

// Minimum distance to be considerd for ray collision
const MinDistance float64 = 1e-5

// If true multiple rays are casted and blended for each pixel
const Antialiasing = false

// If true the last n Frames are blended
const TemporalFilter = true
const TemporalFilterSamples = 32

// If true splits the image generation into threads
const Multithreaded = true
const MultithreadedTheads = 4
const MultithreadDataCopies = false

// Temporal acomulation buffers
var Frames []*pixel.PictureData

// Scene and camera copies for threads
var SceneCopies []*geometry.Scene
var CameraCopies []*camera.CameraDefocus

func main() {
	//runtime.GOMAXPROCS(8)
	pixelgl.Run(run)
}

func run() {
	var bounds = pixel.R(0, 0, Width, Height)
	var windowBounds = pixel.R(0, 0, Width*Upscale, Height*Upscale)

	// Prepare the scene
	var scene = geometry.NewScene()
	scene.Add(geometry.NewSphere(500.0, vmath.NewVector3(0.0, -500.5, -1.0), material.NewLightMaterial(vmath.NewVector3(0.4, 0.7, 0.0))))
	scene.Add(geometry.NewSphere(0.5, vmath.NewVector3(-1.0, 0.0, -3.0), material.NewNormalMaterial()))
	scene.Add(geometry.NewSphere(1.5, vmath.NewVector3(5.0, 1.0, -6.0), material.NewDieletricMaterial(1.3, vmath.NewVector3(0.90, 0.90, 0.90))))
	scene.Add(geometry.NewSphere(1.5, vmath.NewVector3(-1.0, 1.0, -3.0), material.NewMetalMaterial(vmath.NewVector3(0.6, 0.6, 0.6), 0.1)))

	var min = 15.0
	var distance = 30.0

	//LoadOBJ(scene, "bunny.obj", material.NewLightMaterial(vmath.NewVector3(0.90, 0.9, 0.9)))

	// Place random sphere objects
	for i := 0; i < 40; i++ {
		var radius = 0.4 + rand.Float64()*0.2
		var position = vmath.NewVector3(rand.Float64()*distance-min, radius-0.5, rand.Float64()*distance-min)
		scene.Add(geometry.NewSphere(radius, position, material.NewLightMaterial(vmath.NewRandomVector3(0.1, 1))))

		radius = 0.4 + rand.Float64()*0.2
		position = vmath.NewVector3(rand.Float64()*distance-min, radius-0.5, rand.Float64()*distance-min)
		scene.Add(geometry.NewSphere(radius, position, material.NewMetalMaterial(vmath.NewRandomVector3(0.1, 1), rand.Float64())))

		radius = 0.4 + rand.Float64()*0.2
		position = vmath.NewVector3(rand.Float64()*distance-min, radius-0.5, rand.Float64()*distance-min)
		scene.Add(geometry.NewSphere(radius, position, material.NewDieletricMaterial(2.0*rand.Float64(), vmath.NewRandomVector3(0.95, 1.0))))
	}

	// Random triangles
	for i := 0; i < 0; i++ {
		var size float64 = 1.0
		var position = vmath.NewVector3(rand.Float64()*distance-min, size/2.0-0.5, rand.Float64()*distance-min)

		var a = position.Clone()
		a.Add(vmath.NewVector3(0.0, size, 0.0))
		var b = position.Clone()
		b.Add(vmath.NewVector3(-size/1.5, 0, 0.0))
		var c = position.Clone()
		c.Add(vmath.NewVector3(size/1.5, 0, 0.0))

		scene.Add(geometry.NewTriangle(a, b, c, material.NewLightMaterial(vmath.NewRandomVector3(0.1, 1))))
	}

	var halfSize = vmath.NewVector3(0.5, 0.5, 0.5)

	//Place random box objects
	for i := 0; i < 10; i++ {
		var position = vmath.NewVector3(rand.Float64()*distance-min, halfSize.Y-0.5, rand.Float64()*distance-min)
		var bmin = position.Clone()
		bmin.Sub(halfSize)
		var bmax = position.Clone()
		bmax.Add(halfSize)
		scene.Add(geometry.NewBox(bmin, bmax, material.NewLightMaterial(vmath.NewRandomVector3(0.1, 1))))

		position = vmath.NewVector3(rand.Float64()*distance-min, halfSize.Y-0.5, rand.Float64()*distance-min)
		bmin = position.Clone()
		bmin.Sub(halfSize)
		bmax = position.Clone()
		bmax.Add(halfSize)
		scene.Add(geometry.NewBox(bmin, bmax, material.NewMetalMaterial(vmath.NewRandomVector3(0.6, 1), 0.0)))
	}

	var camera = camera.NewCameraDefocusBounds(bounds)

	if Multithreaded && MultithreadDataCopies {
		for i := 0; i < MultithreadedTheads; i++ {
			SceneCopies = append(SceneCopies, scene.Clone())
			CameraCopies = append(CameraCopies, camera.Clone())
		}
	}

	var config = pixelgl.WindowConfig{
		Resizable:   false,
		Undecorated: false,
		VSync:       false,
		Title:       "Gotracer",
		Bounds:      windowBounds}

	var window, err = pixelgl.NewWindow(config)

	CheckError(err)

	var delta time.Duration

	for !window.Closed() {

		var start = time.Now()

		window.Clear(colornames.Black)

		var picture *pixel.PictureData = Render(bounds, scene, camera)
		var sprite *pixel.Sprite

		if TemporalFilter {

			// Add new frame to the list
			Frames = append(Frames, picture)
			if len(Frames) > TemporalFilterSamples {
				Frames = Frames[1:]
			}

			var final = pixel.MakePictureData(bounds)

			// Average the Frames in the list
			for i := 0; i < len(final.Pix); i++ {

				var r, g, b int

				for j := 0; j < len(Frames); j++ {
					r += (int)(Frames[j].Pix[i].R)
					g += (int)(Frames[j].Pix[i].G)
					b += (int)(Frames[j].Pix[i].B)
				}

				final.Pix[i].R = (uint8)(r / len(Frames))
				final.Pix[i].G = (uint8)(g / len(Frames))
				final.Pix[i].B = (uint8)(b / len(Frames))
			}

			sprite = pixel.NewSprite(final, final.Bounds())
		} else {
			sprite = pixel.NewSprite(picture, picture.Bounds())
		}
		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()).Scaled(window.Bounds().Center(), Upscale))

		delta = time.Since(start)
		log.Printf("Frame time %s", delta)

		var speed = 1.0 * delta.Seconds()

		//Keyboard input
		if window.Pressed(pixelgl.KeyRight) {
			camera.Position.X += speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyLeft) {
			camera.Position.X -= speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyUp) {
			camera.Position.Z -= speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyDown) {
			camera.Position.Z += speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyLeftControl) || window.Pressed(pixelgl.KeyRightControl) {
			camera.Position.Y -= speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeySpace) {
			camera.Position.Y += speed
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyW) {
			camera.Aperture += 0.1
			UpdateCamera(camera)
		}
		if window.Pressed(pixelgl.KeyS) {
			camera.Aperture -= 0.1
			UpdateCamera(camera)
		}

		window.Update()
	}
}

// Update the camera viewport
func UpdateCamera(camera *camera.CameraDefocus) {

	if TemporalFilter {
		Frames = nil
	}

	if Multithreaded && MultithreadDataCopies {
		for i := 0; i < MultithreadedTheads; i++ {
			CameraCopies[i].Copy(camera)
			CameraCopies[i].UpdateViewport()
		}
	} else {
		camera.UpdateViewport()
	}
}

// Render image the image
//
//go:norace
func Render(bounds pixel.Rect, scene *geometry.Scene, camera *camera.CameraDefocus) *pixel.PictureData {
	var size = bounds.Size()
	var picture *pixel.PictureData = pixel.MakePictureData(bounds)
	var nx = int(size.X)
	var ny = int(size.Y)
	var wg sync.WaitGroup

	if Multithreaded {
		wg.Add(MultithreadedTheads)
		var wtx = nx / MultithreadedTheads
		var itx = 0

		if MultithreadDataCopies {
			for i := 0; i < MultithreadedTheads; i++ {
				go RaytraceThread(&wg, picture, SceneCopies[i], CameraCopies[i], MaxDepth, TemporalFilter, Antialiasing, size.X, size.Y, itx, 0, itx+wtx, ny)
				itx += wtx
			}
		} else {
			for i := 0; i < MultithreadedTheads; i++ {
				go RaytraceThread(&wg, picture, scene, camera, MaxDepth, TemporalFilter, Antialiasing, size.X, size.Y, itx, 0, itx+wtx, ny)
				itx += wtx
			}
		}

		wg.Wait()
	} else {
		wg.Add(1)
		RaytraceThread(&wg, picture, scene, camera, MaxDepth, TemporalFilter, Antialiasing, size.X, size.Y, 0, 0, nx, ny)
	}

	return picture
}

// Ray trace the picture in a thread and write it to the output object.
// The result is written to the picture object passed as argument.
// This method is intended to be called multiple threads.
//
//go:norace
func RaytraceThread(wg *sync.WaitGroup, picture *pixel.PictureData, scene *geometry.Scene, camera *camera.CameraDefocus, depth int64, jitter bool, antialiasing bool, width float64, height float64, ix int, iy int, nx int, ny int) {
	for j := iy; j < ny; j++ {
		for i := ix; i < nx; i++ {
			var color *vmath.Vector3

			//If using antialiasing jitter the UV and cast multiple rays
			if antialiasing {
				var samples = 4
				color = vmath.NewVector3(0, 0, 0)

				for k := 0; k < samples; k++ {
					var u = (float64(i) + rand.Float64()) / width
					var v = (float64(j) + rand.Float64()) / height
					color.Add(RaytraceScene(scene, camera.GetRay(u, v), depth))
				}

				color.DivideScalar(float64(samples))
			} else {
				var u float64
				var v float64

				if jitter {
					u = (float64(i) + rand.Float64()) / width
					v = (float64(j) + rand.Float64()) / height
				} else {
					u = float64(i) / width
					v = float64(j) / height
				}

				color = RaytraceScene(scene, camera.GetRay(u, v), depth)
			}

			//Apply gamma
			//color.DivideScalar(1.0);
			color.Sqrt()

			color.MulScalar(255)

			//Write to picture
			var index = picture.Index(pixel.Vec{X: float64(i), Y: float64(j)})
			picture.Pix[index].R = uint8(color.X)
			picture.Pix[index].G = uint8(color.Y)
			picture.Pix[index].B = uint8(color.Z)
		}
	}

	wg.Done()
}

// Render the scene to calculate the color for a ray.
// Receives the scene and the initial ray to be casted.
// It is called recursively until the ray does not hit anything, it is absorbed of depth reaches 0.
//
//go:norace
func RaytraceScene(scene *geometry.Scene, ray *vmath.Ray, depth int64) *vmath.Vector3 {
	var hitRecord = material.NewHitRecord()

	if scene.Hit(ray, MinDistance, math.MaxFloat64, hitRecord) {

		var scattered = vmath.NewEmptyRay()
		var attenuation = vmath.NewVector3(0, 0, 0)

		if depth > 0 && hitRecord.Material.Scatter(ray, hitRecord, attenuation, scattered) {
			var color = attenuation.Clone()
			color.Mul(RaytraceScene(scene, scattered.Clone(), depth-1))
			return color
		} else {
			// Ray was absorved return black
			//return vmath.NewVector3(0, 0, 0);

			// The ray was absorved use the last value
			return attenuation.Clone()
		}

	} else {

		return BackgroundColor(ray)
	}
}

// Calculate the background color from ray.
// This method is used for multi threading.
//
//go:norace
func BackgroundColor(r *vmath.Ray) *vmath.Vector3 {
	var unitDirection = r.Direction.UnitVector()
	var t = 0.5 * (unitDirection.Y + 1.0)

	var a = vmath.NewVector3(1.0, 1.0, 1.0)
	a.MulScalar(1.0 - t)

	var b = vmath.NewVector3(0.5, 0.7, 1.0)
	b.MulScalar(t)

	a.Add(b)

	return a
}

// Load obj file triangle into the scene.
//
//go:norace
func LoadOBJ(scene *geometry.Scene, fname string, material material.Material) {
	var file, _ = os.Open(fname)

	var data, _ = ioutil.ReadAll(file)

	var reader = obj.NewReader(bytes.NewBuffer(data))
	var object, _ = reader.Read()

	for i := 0; i < len(object.Faces); i++ {
		var points = object.Faces[i].Points
		var a = vmath.NewVector3(points[0].Vertex.X, points[0].Vertex.Y, points[0].Vertex.Z)
		var b = vmath.NewVector3(points[1].Vertex.X, points[1].Vertex.Y, points[1].Vertex.Z)
		var c = vmath.NewVector3(points[2].Vertex.X, points[2].Vertex.Y, points[2].Vertex.Z)
		scene.Add(geometry.NewTriangle(a, b, c, material))
	}
}

// Write the frame to a PPM file string.
//
//go:norace
func WritePPM(picture *pixel.PictureData, fname string) {
	var size = picture.Rect.Size()

	var nx = int(size.X)
	var ny = int(size.Y)

	var file, err = os.Create(fname)
	CheckError(err)

	_, _ = file.WriteString("P3\n" + strconv.Itoa(nx) + " " + strconv.Itoa(ny) + "\n255\n")

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			//Write to file
			var index = picture.Index(pixel.Vec{X: float64(i), Y: float64(j)})
			_, _ = file.WriteString(strconv.Itoa(int(picture.Pix[index].R)) + " " + strconv.Itoa(int(picture.Pix[index].G)) + " " + strconv.Itoa(int(picture.Pix[index].B)) + "\n")
		}
	}

	//Close file
	_ = file.Sync()
	_ = file.Close()
}

// CheckError an error.
//
//go:norace
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}
