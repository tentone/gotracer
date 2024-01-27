package camera

import (
	"github.com/gopxl/pixel/v2"
	"gotracer/vmath"
	"math"
)

// CameraDefocus object describes how the objects are projected into the screen.
// The camera object is used to get the rays that need to be casted for each screen UV coordinate.
type Camera struct {
	// Aspect ratio of the camera viewport (X / Y)
	AspectRatio float64

	// Field of view of the camera in degrees.
	Fov float64

	// World position of the camera
	Position *vmath.Vector3

	// Point where the camera is looking at
	LookAt *vmath.Vector3

	// Up direction to calculate the camera look direction
	Up *vmath.Vector3

	// The Lower left corner of the camera relative to the center considering the vertical and horizontal sizes.
	// Calculated by the UpdateViewport method.
	LowerLeftCorner *vmath.Vector3

	// Vertical size of the camera (usually only uses Y).
	// Calculated by the UpdateViewport method.
	Vertical *vmath.Vector3

	// Horizontal size of the camera (usually only uses X).
	// Calculated by the UpdateViewport method.
	Horizontal *vmath.Vector3
}

// Create camera from bouding box
func NewCamera(bounds pixel.Rect, position *vmath.Vector3, lookAt *vmath.Vector3, up *vmath.Vector3, fov float64) *Camera {
	var c = new(Camera)
	var size = bounds.Size()

	c.Fov = fov
	c.AspectRatio = size.X / size.Y
	c.Position = position
	c.LookAt = lookAt
	c.Up = up
	c.UpdateViewport()

	return c
}

// Create camera from bouding box
func NewCameraBounds(bounds pixel.Rect) *Camera {
	var c = new(Camera)
	var size = bounds.Size()

	c.Fov = 70
	c.AspectRatio = size.X / size.Y
	c.Position = vmath.NewVector3(-2.0, 2.0, 1.0)
	c.LookAt = vmath.NewVector3(0.0, 0.0, -1.0)
	c.Up = vmath.NewVector3(0.0, 1.0, 0.0)
	c.UpdateViewport()

	return c
}

// UpdateViewport camera projection properties.
func (c *Camera) UpdateViewport() {

	var fovRad = c.Fov * (math.Pi / 180.0)

	var halfHeight = math.Tan(fovRad / 2.0)
	var halfWidth = c.AspectRatio * halfHeight

	var direction = c.Position.Clone()
	direction.Sub(c.LookAt)
	var w = direction.UnitVector()

	var u = vmath.Cross(c.Up, w)
	var v = vmath.Cross(w, u)

	u.MulScalar(halfWidth)
	v.MulScalar(halfHeight)

	c.LowerLeftCorner = c.Position.Clone()
	c.LowerLeftCorner.Sub(u)
	c.LowerLeftCorner.Sub(v)
	c.LowerLeftCorner.Sub(w)

	c.Horizontal = u.Clone()
	c.Horizontal.MulScalar(2.0)

	c.Vertical = v.Clone()
	c.Vertical.MulScalar(2.0)
}

// Get a ray from this camera, from a normalized UV screen coordinate.
func (c *Camera) GetRay(u float64, v float64) *vmath.Ray {
	var hor = c.Horizontal.Clone()
	hor.MulScalar(u)

	var vert = c.Vertical.Clone()
	vert.MulScalar(v)

	var direction = c.LowerLeftCorner.Clone()
	direction.Add(hor)
	direction.Add(vert)
	direction.Sub(c.Position)

	return vmath.NewRay(c.Position, direction)
}

// Copy data from another camera object
func (c *Camera) Copy(o *Camera) {
	c.Fov = o.Fov
	c.AspectRatio = o.AspectRatio
	c.Position.Copy(o.Position)
	c.LookAt.Copy(o.LookAt)
	c.Up.Copy(o.Up)
}

// Clone the camera object
func (o *Camera) Clone() *Camera {
	var c = new(Camera)
	c.Fov = o.Fov
	c.AspectRatio = o.AspectRatio
	c.Position = o.Position.Clone()
	c.LookAt = o.LookAt.Clone()
	c.Up = o.Up.Clone()
	c.UpdateViewport()
	return c
}
