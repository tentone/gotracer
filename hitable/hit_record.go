package hitable;

import "gotracer/vmath";


type HitRecord struct {
	t float64;
	p vmath.Vector3;
	normal vmath.Vector3;
}
