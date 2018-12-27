package hitable;

import "gotracer/vmath";

type HitRecord struct {
	t float64;
	p vmath.Vector3;
	normal vmath.Vector3;
}

type Hitable interface {
	hit(ray *vmath.Ray) bool
}