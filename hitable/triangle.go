package hitable

import (
	"gotracer/vmath"
	"math"
)

// Triangle is hittable object represented by three points.
type Triangle struct {
	A *vmath.Vector3
	B *vmath.Vector3
	C *vmath.Vector3

	// Normal direction of the triangle plane
	Normal *vmath.Vector3

	// Material used to render the sphere.
	Material Material
}

func NewTriangle(a *vmath.Vector3, b *vmath.Vector3, c *vmath.Vector3, material Material) *Triangle {
	var t = new(Triangle)
	t.A = a
	t.B = b
	t.C = c
	t.GetNormal()
	t.Material = material
	return t
}

func (triangle *Triangle) GetNormal() {

	var c = triangle.C.Clone()
	c.Sub(triangle.B)

	var a = triangle.A.Clone()
	a.Sub(triangle.B)

	var cross = vmath.Cross(c, a)

	var targetLengthSq = cross.SquaredLength()

	if targetLengthSq > 0 {
		cross.MulScalar(1 / math.Sqrt(targetLengthSq))
		triangle.Normal = cross.Clone()
	}

	triangle.Normal = vmath.NewEmptyVector3()
}

func (triangle *Triangle) Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *HitRecord) bool {

	// Check if ray and plane are parallel
	var NdotRayDirection float64 = vmath.Dot(triangle.Normal, ray.Direction)
	if math.Abs(NdotRayDirection) < 0.001 {
		return false
	}

	// Compute d parameter
	var d float64 = vmath.Dot(triangle.Normal, triangle.A)

	// Compute t (distance))
	var t float64 = (vmath.Dot(triangle.Normal, ray.Origin) + d) / NdotRayDirection

	// Check if the triangle is in behind the ray
	if t < 0 {
		return false
	}

	// Compute the intersection point using equation 1
	var P *vmath.Vector3 = ray.Direction.Clone()
	P.MulScalar(t)
	P.Add(ray.Origin)

	// Edge 0
	var edge0 *vmath.Vector3 = triangle.B.Clone()
	edge0.Sub(triangle.A)
	var vp0 *vmath.Vector3 = P.Clone()
	vp0.Sub(triangle.A)

	// Vector perpendicular to triangle's plane
	var C *vmath.Vector3 = vmath.Cross(edge0, vp0)
	if vmath.Dot(triangle.Normal, C) < 0 {
		return false
	}

	// Edge 1
	var edge1 *vmath.Vector3 = triangle.C.Clone()
	edge1.Sub(triangle.B)
	var vp1 *vmath.Vector3 = P.Clone()
	vp1.Sub(triangle.B)
	C = vmath.Cross(edge1, vp1)
	if vmath.Dot(triangle.Normal, C) < 0 {
		return false
	}

	// Edge 2
	var edge2 *vmath.Vector3 = triangle.A.Clone()
	edge2.Sub(triangle.C)
	var vp2 *vmath.Vector3 = P.Clone()
	vp2.Sub(triangle.C)
	C = vmath.Cross(edge2, vp2)
	if vmath.Dot(triangle.Normal, C) < 0 {
		return false
	}

	hitRecord.T = t
	hitRecord.P = P //ray.PointAtParameter(temp)
	hitRecord.Normal = triangle.Normal.Clone()
	hitRecord.Material = triangle.Material
	return true
}

func (triangle *Triangle) Clone() Hitable {
	var s = new(Triangle)
	s.A = triangle.A.Clone()
	s.B = triangle.B.Clone()
	s.C = triangle.C.Clone()
	s.Normal = triangle.Normal.Clone()
	s.Material = triangle.Material.Clone()
	return s
}