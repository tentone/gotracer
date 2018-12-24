package math;

import "strconv";
import "math";

// Vector 3 is represented by a x,y,z values.
type Vector3 struct {
	x float64
	y float64
	z float64
}

// Create new vector3 with values.
func NewVector3(x float64, y float64, z float64) *Vector3 {
	var v = new(Vector3)
	v.x = x;
	v.y = y;
	v.z = z;
	return v;
}

// Add vectors
func (v *Vector3) Add(b *Vector3) {
	v.x += b.x;
	v.y += b.y;
	v.z += b.z;
}

// Subtract vectors
func (v *Vector3) Sub(b *Vector3) {
	v.x -= b.x;
	v.y -= b.y;
	v.z -= b.z;
}

// Multiply vectors
func (v *Vector3) Mul(b *Vector3) {
	v.x *= b.x;
	v.y *= b.y;
	v.z *= b.z;
}

// Multiply vectors
func (v *Vector3) Divide(b *Vector3) {
	v.x *= b.x;
	v.y *= b.y;
	v.z *= b.z;
}

// Length of the vector
func (v *Vector3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z);
}

// Squared length of the vector (useful for comparisons, avoids the squaredroot calc).
func (v *Vector3) SquaredLength() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z;
}

// Generate a string with the vector values
func (v *Vector3) ToString() string {
	return "(" + strconv.FormatFloat(v.x, 'f', -1, 64) + ", " + strconv.FormatFloat(v.y, 'f', -1, 64) + ", " + strconv.FormatFloat(v.z, 'f', -1, 64) + ")";
}
