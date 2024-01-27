package material

import (
	"gotracer/vmath"
)

// Material class can be used to calculate how the light rays are affected by the hitable objects surface.
// Should be used to implement multiple types of materials.
type Material interface {
	// Calculate a scattered ray based on the input ray (that hit the surface).
	// Produce a scattered ray (or say it absorbed the incident ray), if scattered, say how much the ray should be attenuated.
	// The return value indicates if if the ray was scatered, if returned false we assume that the ray was absorved.
	Scatter(ray *vmath.Ray, hitRecord *HitRecord, attenuation *vmath.Vector3, scattered *vmath.Ray) bool

	// Clone object create a new object with the same properties.
	Clone() Material
}
