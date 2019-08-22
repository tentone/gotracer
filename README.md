# Gotracer
 - Software Raytracer written in golang.
 - Images can be previewed directly on the window as they are rendered.
 - Interaction can be done using keys from the keyboard, to control the camera.
 
## Screenshots
![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/a.png)![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/b.png)![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/c.png)

## Performance
 - To improve performance multi-threading was added to the renderer.
 - Data had to be reordered for go to actually scale properly with multiple threads.
 - Calling multiple goroutines that use that same data cause them to lock on each other to access data and the performance gains are minimal.
 - Tests were performed with 154 objects in the render scene.
 - Test platform was a Core i7 3537u running Go 1.12.8 w/ 4 threads
 - Base performance number for a single thread.
    - 320x240 ~590ms per frame
    - 640x480 ~2250ms per frame
 - Splitting the work by multiple goroutines
    - 320x240 ~ms per frame (x faster)
    - 640x480 ~ms per frame (x faster)
 - Explicit object copies for each thread (4 threads)
    - 320x240 ~300ms per frame (x faster)
    - 640x480 ~ms per frame (x faster)
 - Using the //go:norace code annotation to skip the data race condition analysis.
    - 320x240 ~285ms per frame (x faster)
    - 640x480 ~ms per frame (x faster)

## Features
 - Geometries (Sphere, Box).
 - Materials (Dieletrics, Lambert, Metal, Normal).
 - Camera defocus.
 - Filtering
    - Antialiased image from ray jittering.
    - Temporal accomulation from single ray raytraced images.

## Build
 - Install golang development tools
 - If youre running on Windows install the GCC compiler
    - http://win-builds.org/doku.php.
    - Add the winbuils/bin folder to the path.
    - Alternatively you can use MSYS2
       - https://github.com/faiface/pixel/wiki/Building-Pixel-on-Windows
 - Run go get and go build.
 - Run the executable

## Libraries
 - PixelGL
    - Used to create windows, output image and get key inputs
    - https://godoc.org/github.com/faiface/pixel/pixelgl
    - https://github.com/faiface/pixel/wiki/Creating-a-Window
    - https://godoc.org/github.com/faiface/pixel#PictureData

## References
 - Raytracer in a Weekend (Peter Shirley)
 - An efficient and robust ray-box intersection algorithm (2003) (Amy Williams , Steve Barrus , R. Keith , Morley Peter Shirley)