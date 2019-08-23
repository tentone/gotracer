package hitable

import (
	"gotracer/vmath"
)

// Emissive material materials are diffuse objects that don’t emit light merely take on the color of their surroundings.
// But they modulate that with their own intrinsic color. Light that reflects off a diffuse surface has its direction randomized.
// They also might be absorbed rather than reflected. The darker the surface, the more likely  absorption is.
type LambertMaterial struct {
	// Albedo represents the base color of the material.
	Albedo *vmath.Vector3
}

func NewLambertMaterial(albedo *vmath.Vector3) *LambertMaterial {
	var m = new(LambertMaterial)
	m.Albedo = albedo
	return m
}

func (m *LambertMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {
	var target = hitRecord.Normal.Clone()
	target.Add(vmath.RandomInUnitSphere())

	scattered.Set(hitRecord.P, target)
	attenuation.Copy(m.Albedo)

	return true
}

func (o *LambertMaterial) Clone() Material {
	var m = new(LambertMaterial)
	m.Albedo = o.Albedo.Clone()
	return m
}