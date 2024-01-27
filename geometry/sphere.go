package geometry

import (
	"gotracer/material"
	"gotracer/vmath"
)
import "math"

// Sphere is hitable object represented by a center point and a radius.
// The sphere object has a material attached to it.
type Sphere struct {
	// Radius of the sphere
	Radius float64

	// Center position of the sphere
	Center *vmath.Vector3

	// Material used to render the sphere.
	Material material.Material
}

func NewSphere(radius float64, center *vmath.Vector3, material material.Material) *Sphere {
	var s = new(Sphere)
	s.Radius = radius
	s.Center = center
	s.Material = material
	return s
}

func (s *Sphere) Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *material.HitRecord) bool {

	var oc = ray.Origin.Clone()
	oc.Sub(s.Center)

	var a = vmath.Dot(ray.Direction, ray.Direction)
	var b = vmath.Dot(oc, ray.Direction)
	var c = vmath.Dot(oc, oc) - s.Radius*s.Radius
	var discriminant = b*b - a*c

	if discriminant > 0 {

		//First root result
		var temp = (-b - math.Sqrt(discriminant)) / a

		if temp < tmax && temp > tmin {
			hitRecord.T = temp
			hitRecord.P = ray.PointAtParameter(temp)
			hitRecord.Normal = hitRecord.P.Clone()
			hitRecord.Normal.Sub(s.Center)
			hitRecord.Normal.DivideScalar(s.Radius)
			hitRecord.Material = s.Material
			return true
		}

		//Second root result
		temp = (-b + math.Sqrt(discriminant)) / a

		if temp < tmax && temp > tmin {
			hitRecord.T = temp
			hitRecord.P = ray.PointAtParameter(temp)
			hitRecord.Normal = hitRecord.P.Clone()
			hitRecord.Normal.Sub(s.Center)
			hitRecord.Normal.DivideScalar(s.Radius)
			hitRecord.Material = s.Material
			return true
		}

	}

	return false
}

func (o *Sphere) Clone() Hitable {
	var s = new(Sphere)
	s.Radius = o.Radius
	s.Center = o.Center.Clone()
	s.Material = o.Material.Clone()
	return s
}
