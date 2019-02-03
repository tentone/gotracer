package hitable

import (
	"gotracer/vmath"
	"math/rand"
);

// Dielectric material allow light to pass trough them.
// When a light ray hits them, it splits into a reflected ray and a refracted (transmitted) ray.
type DieletricMaterial struct {
	RefractiveIndice float64;
}

func NewDieletricMaterial (refractiveIndice float64) *DieletricMaterial  {
	var m = new(DieletricMaterial);
	m.RefractiveIndice = refractiveIndice;
	return m;
}

// Refractive indice of the air.
var AirRefractiveIndice = 1.0;

func (m *DieletricMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	var outwardNormal *vmath.Vector3 = vmath.NewEmptyVector3();
	var refracted *vmath.Vector3 = vmath.NewEmptyVector3();
	var reflected *vmath.Vector3 = vmath.Reflect(ray.Direction, hitRecord.Normal);
	var refractionRatio float64;
	var reflectionProbe float64;
	var cosine float64;

	attenuation.Set(1.0, 1.0, 1.0);

	var dot float64 = vmath.Dot(ray.Direction, hitRecord.Normal);

	if dot > 0 {
		outwardNormal.Copy(hitRecord.Normal);
		outwardNormal.MulScalar(-1.0);
		refractionRatio = m.RefractiveIndice;
		cosine = m.RefractiveIndice * dot / ray.Direction.Length();
	} else {
		outwardNormal.Copy(hitRecord.Normal);
		refractionRatio = AirRefractiveIndice / m.RefractiveIndice;
		cosine = -dot / ray.Direction.Length();
	}

	if vmath.Refract(ray.Direction, outwardNormal, refractionRatio, refracted) {
		reflectionProbe = vmath.Schlick(cosine, m.RefractiveIndice);
	} else {
		reflectionProbe = 1.0;
		scattered.Set(hitRecord.P, reflected);
		return true;
	}

	// TODO <SUPPORT MULTIPLE SCATERED RAYS>
	// Return reflected of refracted randomly with reflection probe probability.
	if rand.Float64() < reflectionProbe {
		scattered.Set(hitRecord.P, reflected);
	} else {
		scattered.Set(hitRecord.P, refracted);
	}

	return true;
}
