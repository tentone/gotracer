package graphics;

type Frame struct {
	Data [921600]int; //640 * 480 * 3
	Width int;
	Height int;
}

// Write the frame to a PPM file string
func (f *Frame) WritePPM() string {
	return "TODO";
}

// String description of the frame
func (f *Frame) ToString() string {
	return "Frame string";
}