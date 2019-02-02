package hitable

import "gotracer/vmath";

// Metalic object type reflect the rays that hit the object surface.
type MetalMaterial struct {
	// Albedo represents the base color of the material.
	Albedo *vmath.Vector3;

	// Fuzz indicates the roughness of the metallic surface.
	// The more fuzz there is the more the ray are reflected with an offset applied.
	Fuzz float64;

}

func NewMetalMaterial(albedo *vmath.Vector3) *MetalMaterial {
	var m = new(MetalMaterial);
	m.Albedo = albedo;
	return m;
}

func (m *MetalMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	var unit *vmath.Vector3 = ray.Direction.UnitVector();
	var reflected = vmath.Reflect(unit, hitRecord.Normal);

	scattered.Set(hitRecord.P, reflected);
	attenuation.Copy(m.Albedo);

	return vmath.Dot(scattered.Direction, hitRecord.Normal) > 0;
}
