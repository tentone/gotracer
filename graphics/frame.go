package graphics;

import "strconv";

type Frame struct {
	Data [921600]int; //640 * 480 * 3
	Width int;
	Height int;
}

func NewFrame(width int, height int) *Frame {
	var f = new(Frame)
	f.Width = width;
	f.Height = height;
	return f;
}

// Write the frame to a PPM file string
func (f *Frame) WritePPM() string {
	return "TODO";
}

// String description of the frame
func (f *Frame) ToString() string {
	return "Width:" + strconv.Itoa(f.Width) + ", Height:" + strconv.Itoa(f.Height);
}