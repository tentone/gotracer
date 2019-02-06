package graphics;

import "C"
import (
	"github.com/faiface/pixel"
	"gotracer/vmath"
	"math"
);

// Camera defocus is a camera that has support for defocus blur.
type CameraDefocus struct {
	Camera;

	// Lens radius affects how much the rays can drift from the center.
	LensRadius float64;

	// Lens aperture.
	Aperture float64;

	// Distance to be in perfect focus of the camera.
	FocusDistance float64;

	U *vmath.Vector3;
	V *vmath.Vector3;
	W *vmath.Vector3;
}

// Create camera from bouding box
func NewCameraDefocus (bounds pixel.Rect, position *vmath.Vector3, lookAt *vmath.Vector3, up *vmath.Vector3, fov float64, aperture float64, focusDistance float64) *CameraDefocus {
	var c = new(CameraDefocus);
	var size = bounds.Size();

	c.Aperture = aperture;
	c.FocusDistance = focusDistance;
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
	c.Aperture = 0.0;

	var direction = c.Position.Clone()
	direction.Sub(c.LookAt);
	c.FocusDistance = direction.Length();
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

	c.LensRadius = c.Aperture / 2.0;
	c.W = direction.UnitVector();
	c.U = vmath.Cross(c.Up, c.W).UnitVector();
	c.V = vmath.Cross(c.W, c.U);

	var u *vmath.Vector3 = c.U.Clone();
	var v *vmath.Vector3 = c.V.Clone();
	var w *vmath.Vector3 = c.W.Clone();

	u.MulScalar(halfWidth * c.FocusDistance);
	v.MulScalar(halfHeight * c.FocusDistance);
	w.MulScalar(c.FocusDistance);

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

	var rd *vmath.Vector3 = vmath.RandomInUnitDisk();
	rd.MulScalar(c.LensRadius);

	var offset *vmath.Vector3 = c.U.Clone();
	offset.MulScalar(rd.X);

	var vclone *vmath.Vector3 = c.V.Clone();
	vclone.MulScalar(rd.Y);
	offset.Add(vclone);

	var hor *vmath.Vector3 = c.Horizontal.Clone();
	hor.MulScalar(u);

	var vert *vmath.Vector3 = c.Vertical.Clone();
	vert.MulScalar(v);

	var direction *vmath.Vector3 = c.LowerLeftCorner.Clone();
	direction.Add(hor);
	direction.Add(vert);
	direction.Sub(c.Position);
	direction.Sub(offset);

	offset.Add(c.Position);

	return vmath.NewRay(offset, direction);
}

