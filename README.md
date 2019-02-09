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
 - Here are the base performance number for a single thread on my i5 6500 with DDR3 memory.
    - 320x240 averaged 430ms per frame
    - 640x480 averaged 1750ms per frame
 - Just by splitting the work by multiple goroutines (4 threads)
    - 320x240 averaged 220ms per frame (1.95x faster)
    - 640x480 averaged 650ms per frame (2.69x faster)
 - Just by making explicit object copies for each thread (4 threads)
    - 320x240 averaged 120ms per frame (3.58x faster)
    - 640x480 averaged 620ms per frame (3.64x faster)
    
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