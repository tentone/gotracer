package math;

import "strconv";
import "math";

// Vector 3 is represented by a x,y,z values.
type Vector3 struct {
	X float64;
	Y float64;
	Z float64;
}

// Create new vector3 with values.
func NewVector3(x float64, y float64, z float64) *Vector3 {
	var v = new(Vector3)
	v.X = x;
	v.Y = y;
	v.Z = z;
	return v;
}

// Add vectors
func (v *Vector3) Add(b *Vector3) {
	v.X += b.X;
	v.Y += b.Y;
	v.Z += b.Z;
}

// Subtract vectors
func (v *Vector3) Sub(b *Vector3) {
	v.X -= b.X;
	v.Y -= b.Y;
	v.Z -= b.Z;
}

// Multiply vectors
func (v *Vector3) Mul(b *Vector3) {
	v.X *= b.X;
	v.Y *= b.Y;
	v.Z *= b.Z;
}

// Multiply vectors
func (v *Vector3) Divide(b *Vector3) {
	v.X *= b.X;
	v.Y *= b.Y;
	v.Z *= b.Z;
}

// Multiply vector by scalar
func (v *Vector3) MulScalar(b float64) {
	v.X *= b;
	v.Y *= b;
	v.Z *= b;
}

// Divide vector by scalar
func (v *Vector3) DivideScalar(b float64) {
	v.X /= b;
	v.Y /= b;
	v.Z /= b;
}

// Return a copy of the vector
func (v *Vector3) Clone() *Vector3 {
	return NewVector3(v.X, v.Y, v.Z);
};

// Create a new copy vector with a unit length vector with the same direction as this one.
func (v Vector3) UnitVector() Vector3 {
	v.DivideScalar(v.Length());
	return v;
}

// Length of the vector
func (v *Vector3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z);
}

// Squared length of the vector (useful for comparisons, avoids the squaredroot calc).
func (v *Vector3) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z;
}

// Generate a string with the vector values
func (v *Vector3) ToString() string {
	return "(" + strconv.FormatFloat(v.X, 'f', -1, 64) + ", " + strconv.FormatFloat(v.Y, 'f', -1, 64) + ", " + strconv.FormatFloat(v.Z, 'f', -1, 64) + ")";
}
