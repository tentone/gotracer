package hitable;

import (
	"gotracer/vmath"
);

// EmissiveMaterial materials are diffuse objects that emit light ray that collide with them get more light.
type EmissiveMaterial struct {
	// Albedo represents the base color of the material.
	Albedo *vmath.Vector3;

	// Intesity of the emissive material
	Intensity float64;
}

func NewEmissiveMaterial(albedo *vmath.Vector3, intensity float64) *EmissiveMaterial {
	var m = new(EmissiveMaterial);
	m.Albedo = albedo;
	m.Intensity = intensity;
	return m;
}

func (m *EmissiveMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {
	var target *vmath.Vector3 = hitRecord.Normal.Clone();
	target.Add(vmath.RandomInUnitSphere());

	scattered.Set(hitRecord.P, target);
	attenuation.Copy(m.Albedo);
	attenuation.MulScalar(m.Intensity);

	return true;
}

func (o *EmissiveMaterial) Clone() Material {
	var m = new(EmissiveMaterial);
	m.Albedo = o.Albedo.Clone();
	m.Intensity = o.Intensity;
	return m;
}