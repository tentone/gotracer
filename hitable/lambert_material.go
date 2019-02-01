package hitable;

import (
	"gotracer/vmath"
	"math/rand"
);

// LambertMaterial materials are diffuse objects that don’t emit light merely take on the color of their surroundings.
// But they  modulate that with their own intrinsic color. Light that reflects off a diffuse surface has its direction randomized.
// They also might be absorbed rather than reflected. The darker the surface, the more likely  absorption is.
type LambertMaterial struct {
	Albedo *vmath.Vector3;
}

func NewLambertMaterial(albedo *vmath.Vector3) *LambertMaterial {
	var m = new(LambertMaterial);
	m.Albedo = albedo
	return m;
}

func (m *LambertMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {
	var target *vmath.Vector3 = vmath.NewVector3(0, 0, 0);
	target.Add(hitRecord.Normal);
	target.Add(RandomInUnitSphere());

	scattered = vmath.NewRay(hitRecord.P, target);
	attenuation.Copy(m.Albedo);

	return true;
}

// Calculate a random unitary vector in the surface of a sphere.
func RandomInUnitSphere() *vmath.Vector3 {
	var p *vmath.Vector3 = vmath.NewVector3(0, 0, 0);

	for {
		p.Set(rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0);

		if p.SquaredLength() < 1.0 {
			break
		}
	}

	return p;
}