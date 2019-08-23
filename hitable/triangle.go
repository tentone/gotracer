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
}

// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm
func (triangle *Triangle) Hit(ray *vmath.Ray, tmin float64, tmax float64, hitRecord *HitRecord) bool {
	var v0v1 *vmath.Vector3 = triangle.B.Clone()
	v0v1.Sub(triangle.A)

	var v0v2 *vmath.Vector3 = triangle.C.Clone()
	v0v2.Sub(triangle.A)

	var pvec *vmath.Vector3 = vmath.Cross(ray.Direction, v0v2)
	var det = vmath.Dot(v0v1, pvec)
	if det < 0.000001 {
		return false
	}

	var invDet = 1.0 / det
	var tvec *vmath.Vector3 = ray.Origin.Clone()
	tvec.Sub(triangle.A)

	var u = vmath.Dot(tvec, pvec) * invDet
	if u < 0 || u > 1 {
		return false
	}

	var qvec *vmath.Vector3 = vmath.Cross(tvec, v0v1)
	var v = vmath.Dot(ray.Direction, qvec) * invDet
	if v < 0 || u+v > 1 {
		return false
	}

	var t = vmath.Dot(v0v2, qvec) * invDet
	hitRecord.T = t
	hitRecord.P = ray.PointAtParameter(t)
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