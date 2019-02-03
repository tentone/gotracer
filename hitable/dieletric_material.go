package hitable

import "gotracer/vmath";

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

	var outwardNormal *vmath.Vector3;
	var reflected *vmath.Vector3 = vmath.Reflect(ray.Direction, hitRecord.Normal);
	var refractionRatio float64;

	attenuation.Set(1.0, 1.0, 1.0);

	var refracted *vmath.Vector3;

	if vmath.Dot(ray.Direction, hitRecord.Normal) > 0 {
		outwardNormal.Copy(hitRecord.Normal);
		outwardNormal.MulScalar(-1.0);
		refractionRatio = m.RefractiveIndice;
	} else {
		outwardNormal.Copy(hitRecord.Normal);
		refractionRatio = AirRefractiveIndice / m.RefractiveIndice;
	}

	// If there is a refracted ray use it as scattered ray
	// TODO <SUPPORT MULTIPLE SCATERED RAYS>
	if vmath.Refract(ray.Direction, outwardNormal, refractionRatio, refracted) {
		scattered.Set(hitRecord.P, refracted);
	} else {
		scattered.Set(hitRecord.P, reflected);
		return false;
	}

	return true;
}
