package hitable;

import "gotracer/vmath";

type Hitable interface {
	hit(ray *vmath.Ray) bool
}