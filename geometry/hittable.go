package geometry

import (
	"gotracer/material"
	"gotracer/vmath"
)

// Hitable interface indicates a object that can be raytraced.
type Hitable interface {
	// The hit method indicates if the object was intersected by the ray.
	// If true the result is stored on the hitrecord object provided.
	Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *material.HitRecord) bool

	// Clone object create a new object with the same properties.
	Clone() Hitable
}
