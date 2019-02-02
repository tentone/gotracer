package graphics;

import (
	"github.com/faiface/pixel"
	"gotracer/vmath"
);

// Camera object describes how the objects are projected into the screen.
// The camera object is used to get the rays that need to be casted for each screen UV coordinate.
type Camera struct {
	// Position of the camera in the world
	Origin *vmath.Vector3;

	//The Lower left corner of the camera relative to the center considering the vertical and horizontal sizes.
	LowerLeftCorner *vmath.Vector3;

	// Vertical size of the camera (usually only uses Y)
	Vertical *vmath.Vector3;

	// Horizontal size of the camera (usually only uses X)
	Horizontal *vmath.Vector3;
}

// Create a new camera with default values.
func NewCamera() *Camera {
	var c = new(Camera);
	c.LowerLeftCorner = vmath.NewVector3(-2.0, -1.0, -1.0);
	c.Horizontal = vmath.NewVector3(4.0, 0.0, 0.0);
	c.Vertical = vmath.NewVector3(0.0, 2.0, 0.0);
	c.Origin = vmath.NewVector3(0.0, 0.0, 0.0);
	return c;
}

// Create camera from bouding box
func NewCameraBounds(bounds pixel.Rect) *Camera {
	var c = new(Camera);

	var size = bounds.Size();
	var aspect = size.X / size.Y;
	var scale = 2.0;

	c.LowerLeftCorner = vmath.NewVector3(-scale/2.0*aspect, -1.0, -1.0);
	c.Vertical = vmath.NewVector3(0.0, scale, 0.0);
	c.Horizontal = vmath.NewVector3(scale*aspect, 0.0, 0.0);
	c.Origin = vmath.NewVector3(0.0, 0.0, 0.0);

	return c;
};

// Get a ray from this camera, from a normalized UV screen coordinate.
func (s *Camera) GetRay(u float64, v float64) *vmath.Ray {
	var hor = s.Horizontal.Clone();
	hor.MulScalar(u);

	var vert = s.Vertical.Clone();
	vert.MulScalar(v);

	var direction = s.LowerLeftCorner.Clone();
	direction.Add(hor);
	direction.Add(vert);

	return vmath.NewRay(s.Origin, direction);
}

