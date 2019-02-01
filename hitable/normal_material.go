package hitable

import "gotracer/vmath";

// Material to preview/debug the normal direction of a hitable object.
type NormalMaterial struct {}

func NewNormalMaterial() *NormalMaterial {
	return new(NormalMaterial);
}

func (m *NormalMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	var color = vmath.NewVector3(hitRecord.Normal.X + 1.0, hitRecord.Normal.Y + 1.0, hitRecord.Normal.Z + 1.0);
	color.MulScalar(0.5);

	attenuation.Copy(color);

	return false;
}
