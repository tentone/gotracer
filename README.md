# GoTracer
 - Software Raytracer written in golang.
 - Images can be previewed directly on the window as they are rendered.
 - Interaction can be done using keys from the keyboard, to control the camera.



## Screenshots
![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/a.png)![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/b.png)4![alt tag](https://raw.githubusercontent.com/tentone/gotracer/master/c.png)



## Performance

 - To improve performance multi-threading was added to the renderer.
 - Data had to be reordered for go to actually scale properly with multiple threads.
 - Calling multiple goroutines that use that same data cause them to lock on each other to access data and the performance gains are minimal.
 - Tests were performed with 154 objects in the render scene.
 - Tests performed on a Core i5 6500 (4 Core) with DDR3 memory w/ 4 threads.
 - Tested base performance number for a single thread.
 - Tested splitting the work by multiple goroutines each routine processes one portion of the output.
 - Tested explicit object copies for each thread was able to reduce time further by skipping the synchronization points added by the go race condition detector.

| Mode               | Resolution | Time p/ frame | Speedup |
| ------------------ | ---------- | ------------- | ------- |
| Single-thread      | 320x240    | ~430ms        | 1x      |
| Single-thread      | 640x480    | ~1750ms       | 1x      |
| GoRoutine          | 320x240    | ~220ms        | 1.95x   |
| GoRoutine          | 640x480    | ~650ms        | 2.69x   |
| Explicit Data Copy | 320x240    | ~120ms        | 3.58x   |
| Explicit Data Copy | 640x480    | ~620ms        | 3.64x   |
| go:norace          | 320x240    | TODO          | TODO    |
| go:norace          | 640x480    | TODO          | TODO    |

 - Performed new tests on new Go 1.12.8, also tested the go:norace flag, will have to repeat on a CPU with more cores.
 - Using the //go:norace code annotation to skip the data race condition analysis.
 - Test platform was a Core i7 3537u (2 Core HT) running Go 1.12.8 w/ 4 threads.

| Mode               | Resolution | Time p/ frame |
| ------------------ | ---------- | ------------- |
| Single-thread      | 320x240    | ~590ms        |
| Single-thread      | 640x480    | ~2250ms       |
| GoRoutine          | 320x240    | ~295ms        |
| GoRoutine          | 640x480    | ~1100ms       |
| Explicit Data Copy | 320x240    | ~280ms        |
| Explicit Data Copy | 640x480    | ~1200ms       |
| go:norace          | 320x240    | ~285ms        |
| go:norace          | 640x480    | ~1090ms       |




## Features
 - Geometries (Sphere, Box, Triangles).
 - Materials (Dieletrics, Lambert, Metal, Normal).
 - Camera defocus.
 - Filtering
    - Antialiased image from ray jittering.
    - Temporal accomulation from single ray raytraced images.
 - File loaders (.obj)



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