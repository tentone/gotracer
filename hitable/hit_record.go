package hitable;

import "gotracer/vmath";

// Hit record indicates the intersection of a ray with a surface, indicates where the ray has collided.
type HitRecord struct {
	T float64;
	P *vmath.Vector3;
	Normal *vmath.Vector3;
	Material Material;
}

// Create new hitable list
func NewHitRecord() *HitRecord {
	var hr = new(HitRecord);
	hr.T = 0.0;
	hr.P = vmath.NewVector3(0.0, 0.0, 0.0);
	hr.Normal = vmath.NewVector3(0.0, 0.0, 0.0);
	return hr;
}

// Copy the content of another hit record to this one.
func (a *HitRecord) Copy(b *HitRecord) {
	a.T = b.T;
	a.P.Copy(b.P);
	a.Normal.Copy(b.Normal);
	a.Material = b.Material;
}