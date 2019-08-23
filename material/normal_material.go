package material

import (
	"gotracer/vmath"
)

// Material to preview/debug the normal direction of a hitable object.
type NormalMaterial struct {}

func NewNormalMaterial() *NormalMaterial {
	return new(NormalMaterial)
}

func (m *NormalMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	var target = hitRecord.Normal.Clone()
	target.Add(vmath.RandomInUnitSphere())

	var color = vmath.NewVector3(hitRecord.Normal.X + 1.0, hitRecord.Normal.Y + 1.0, hitRecord.Normal.Z + 1.0)
	color.MulScalar(0.5)

	scattered.Set(hitRecord.P, target)
	attenuation.Copy(color)

	return true
}

func (o *NormalMaterial) Clone() Material {
	return new(NormalMaterial)
}