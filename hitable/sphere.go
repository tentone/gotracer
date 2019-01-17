package hitable;

import "gotracer/vmath";

// Sphere is represented by a center point and a radius.
type Sphere struct {
	Radius float64;
	Center *vmath.Vector3;
}
