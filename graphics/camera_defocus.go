package graphics;

import (
	"github.com/faiface/pixel"
	"gotracer/vmath"
	"math"
);

// Camera defocus is a camera that has support for defocus blur.
type CameraDefocus struct {
	*Camera;

	// Lens radius affects how much the rays can drift from the center.
	LensRadius float64;

	U *vmath.Vector3;

	V *vmath.Vector3;

	W *vmath.Vector3;
}

// Create camera from bouding box
func NewCameraDefocus (bounds pixel.Rect, position *vmath.Vector3, lookAt *vmath.Vector3, up *vmath.Vector3, fov float64) *CameraDefocus {
	var c = new(CameraDefocus);
	var size = bounds.Size();

	c.Fov = fov;
	c.AspectRatio = size.X / size.Y;
	c.Position = position;
	c.LookAt = lookAt;
	c.Up = up;
	c.UpdateViewport();

	return c;
};

// Create camera from bouding box
func NewCameraDefocusBounds(bounds pixel.Rect) *CameraDefocus {
	var c = new(CameraDefocus);
	var size = bounds.Size();

	c.Fov = 90;
	c.AspectRatio = size.X / size.Y;
	c.Position = vmath.NewVector3(-2.0, 2.0, 1.0);
	c.LookAt = vmath.NewVector3(0.0, 0.0, -1.0);
	c.Up = vmath.NewVector3(0.0, 1.0, 0.0);
	c.UpdateViewport();

	return c;
};

// UpdateViewport camera projection properties.
func (c *CameraDefocus) UpdateViewport() {

	var fovRad float64 = c.Fov * (math.Pi / 180.0);

	var halfHeight float64 = math.Tan(fovRad / 2.0);
	var halfWidth float64 = c.AspectRatio * halfHeight;

	var direction *vmath.Vector3 = c.Position.Clone();
	direction.Sub(c.LookAt);

	var w *vmath.Vector3 = direction.UnitVector();
	var u *vmath.Vector3 = vmath.Cross(c.Up, w);
	var v *vmath.Vector3 = vmath.Cross(w, u);

	u.MulScalar(halfWidth);
	v.MulScalar(halfHeight);

	c.LowerLeftCorner = c.Position.Clone();
	c.LowerLeftCorner.Sub(u);
	c.LowerLeftCorner.Sub(v);
	c.LowerLeftCorner.Sub(w);

	c.Horizontal = u.Clone();
	c.Horizontal.MulScalar(2.0);

	c.Vertical = v.Clone();
	c.Vertical.MulScalar(2.0);
}

// Get a ray from this camera, from a normalized UV screen coordinate.
func (c *CameraDefocus) GetRay(u float64, v float64) *vmath.Ray {

	var rd = vmath.RandomInUnitDisk();
	rd.MulScalar(c.LensRadius);

	var hor = c.Horizontal.Clone();
	hor.MulScalar(u);

	var vert = c.Vertical.Clone();
	vert.MulScalar(v);

	var direction = c.LowerLeftCorner.Clone();
	direction.Add(hor);
	direction.Add(vert);
	direction.Sub(c.Position);

	return vmath.NewRay(c.Position, direction);
}

