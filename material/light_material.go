package material

import (
	"gotracer/vmath"
)

// Light material emits light. The color of the object is the solid color of the light.
type LightMaterial struct {
	// Color of the light.
	Color *vmath.Vector3
}

func NewLightMaterial(color *vmath.Vector3) *LightMaterial {
	var m = new(LightMaterial)
	m.Color = color
	return m
}

func (m *LightMaterial) Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool {

	scattered.Set(hitRecord.P, hitRecord.Normal.Clone())
	attenuation.Add(m.Color)

	return true
}

func (o *LightMaterial) Clone() Material {
	var m = new(LightMaterial)
	m.Color = o.Color.Clone()
	return m
}
