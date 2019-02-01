package hitable

import "gotracer/vmath";

// Metalic object type reflect the rays that hit the object surface.
type Metal struct {
	Albedo *vmath.Vector3;
}

func NewMetal(albedo *vmath.Vector3) *Metal {
	var m = new(Metal);
	m.Albedo = albedo;
	return m;
}

func (m *Metal) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	//TODO <ADD CODE HERE>

	return true;
}
