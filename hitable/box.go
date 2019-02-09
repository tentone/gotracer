package hitable;

import (
	"gotracer/vmath"
);

// Box hitable object.
type Box struct {
	// Origin corner of the box
	Min *vmath.Vector3;

	// The opposite maximum corner of the box
	Max *vmath.Vector3;

	// Material used to render the box.
	Material Material;
}

func NewBox(min *vmath.Vector3, max *vmath.Vector3, material Material) *Box {
	var box = new(Box);
	box.Min = min;
	box.Max = max;
	box.Material = material;
	return box;
}

func (box *Box) Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *HitRecord) bool {

	var normal *vmath.Vector3 = vmath.NewVector3(0, 0, 1);
	var txmin float64 = (box.Min.X - ray.Origin.X) / ray.Direction.X;
	var txmax  float64 = (box.Max.X - ray.Origin.X) / ray.Direction.X;

	if txmin > txmax {
		var temp float64 = txmin;
		txmin = txmax;
		txmax = temp;
	}

	if (tmin > txmax) || (txmin > tmax) {
		return false;
	}

	if txmin > tmin {
		tmin = txmin;
	}
	if txmax < tmax {
		tmax = txmax;
	}

	var tymin float64 = (box.Min.Y - ray.Origin.Y) / ray.Direction.Y;
	var tymax float64 = (box.Max.Y - ray.Origin.Y) / ray.Direction.Y;

	if tymin > tymax {
		var temp float64 = tymin;
		tymin = tymax;
		tymax = temp;
	}

	if (tmin > tymax) || (tymin > tmax) {
		return false;
	}

	if tymin > tmin {
		tmin = tymin;
	}
	if tymax < tmax {
		tmax = tymax;
	}

	var tzmin float64 = (box.Min.Z - ray.Origin.Z) / ray.Direction.Z;
	var tzmax float64 = (box.Max.Z - ray.Origin.Z) / ray.Direction.Z;

	if tzmin > tzmax {
		var temp float64 = tzmin;
		tzmin = tzmax;
		tzmax = temp;
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return false;
	}

	if tzmin > tmin {
		tmin = tzmin;
	}
	if tzmax < tmax {
		tmax = tzmax;
	}

	hitRecord.Material = box.Material;
	hitRecord.T = tmin;
	hitRecord.P = ray.PointAtParameter(hitRecord.T);
	hitRecord.Normal = normal;

	return true;
}

func (o *Box) Clone() Hitable {
	var box = new(Box);
	box.Min = o.Min.Clone();
	box.Max = o.Max.Clone();
	box.Material = o.Material.Clone();
	return box;
}