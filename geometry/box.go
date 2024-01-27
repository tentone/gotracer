package geometry

import (
	"gotracer/material"
	"gotracer/vmath"
)

// Box hitable object.
type Box struct {
	// Origin corner of the box
	Min *vmath.Vector3

	// The opposite maximum corner of the box
	Max *vmath.Vector3

	// Material used to render the box.
	Material material.Material
}

func NewBox(min *vmath.Vector3, max *vmath.Vector3, material material.Material) *Box {
	var box = new(Box)
	box.Min = min
	box.Max = max
	box.Material = material
	return box
}

func (box *Box) Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *material.HitRecord) bool {

	var txmin = (box.Min.X - ray.Origin.X) / ray.Direction.X
	var txmax = (box.Max.X - ray.Origin.X) / ray.Direction.X

	var normal = vmath.NewVector3(0, 1.0, 0)
	var signal = -1.0

	if txmin > txmax {
		var temp = txmin
		txmin = txmax
		txmax = temp
		signal = 1.0
	}

	if (tmin > txmax) || (txmin > tmax) {
		return false
	}

	if txmin > tmin {
		normal.Set(signal, 0.0, 0.0)
		tmin = txmin
	}
	if txmax < tmax {
		tmax = txmax
	}

	var tymin = (box.Min.Y - ray.Origin.Y) / ray.Direction.Y
	var tymax = (box.Max.Y - ray.Origin.Y) / ray.Direction.Y
	signal = -1.0

	if tymin > tymax {
		var temp = tymin
		tymin = tymax
		tymax = temp
		signal = 1.0
	}

	if (tmin > tymax) || (tymin > tmax) {
		return false
	}

	if tymin > tmin {
		normal.Set(0.0, signal, 0.0)
		tmin = tymin
	}
	if tymax < tmax {
		tmax = tymax
	}

	var tzmin = (box.Min.Z - ray.Origin.Z) / ray.Direction.Z
	var tzmax = (box.Max.Z - ray.Origin.Z) / ray.Direction.Z
	signal = -1.0

	if tzmin > tzmax {
		var temp = tzmin
		tzmin = tzmax
		tzmax = temp
		signal = 1.0
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return false
	}

	if tzmin > tmin {
		normal.Set(0.0, 0.0, signal)
		tmin = tzmin
	}
	if tzmax < tmax {
		tmax = tzmax
	}

	if tmin > tmax {
		var temp = tmin
		tmin = tmax
		tmax = temp
	}

	hitRecord.Material = box.Material
	hitRecord.T = tmin
	hitRecord.P = ray.PointAtParameter(hitRecord.T)
	hitRecord.Normal = normal

	return true
}

func (o *Box) Clone() Hitable {
	var box = new(Box)
	box.Min = o.Min.Clone()
	box.Max = o.Max.Clone()
	box.Material = o.Material.Clone()
	return box
}
