package hitable;

import "gotracer/vmath";

// Hit record indicates the intersection of a ray with a surface, indicates where the ray has colided.
type HitRecord struct {
	T float64;
	P *vmath.Vector3;
	Normal *vmath.Vector3;
}
