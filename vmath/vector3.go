package vmath;

import (
	"math"
	"math/rand"
	"strconv"
);

// Vector 3 is represented by a x,y,z values.
type Vector3 struct {
	X float64;
	Y float64;
	Z float64;
}

// Create new vector3 with values.
func NewVector3(x float64, y float64, z float64) *Vector3 {
	var v = new(Vector3);
	v.X = x;
	v.Y = y;
	v.Z = z;
	return v;
}

// Create new empty vector3 with values.
func NewEmptyVector3() *Vector3 {
	return new(Vector3);
}


// Set value of the vector.
func (v *Vector3) Set(x float64, y float64, z float64) {
	v.X = x;
	v.Y = y;
	v.Z = z;
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

// Apply sqrt to the individual components of the vector
func (v *Vector3) Sqrt() {
	v.X = math.Sqrt(v.X);
	v.Y = math.Sqrt(v.Y);
	v.Z = math.Sqrt(v.Z);
}

// Normalize this vector
func (v *Vector3) Normalize() {
	v.DivideScalar(v.Length());
}

// Calculate the reflection of a vector relative to a normal vector.
func Reflect(v *Vector3, n *Vector3) *Vector3 {
	var normal = n.Clone();
	normal.MulScalar(2.0 * Dot(v, n));

	var reflected = v.Clone();
	reflected.Sub(normal);

	return reflected;
}


// Calculate the refracted vector of a vector relative to a normal vector.
// This calculation is done using the snells law. Ni is the initial refractive indice and No is the out refraction indice.
// The refractionRatio parameters is calculated from Ni/No.
func Refract(v *Vector3, normal *Vector3, refractionRatio float64, refracted *Vector3) bool {

	var uv *Vector3 = v.UnitVector();
	var dt float64 = Dot(uv, normal);
	var discriminant float64 = 1.0 - math.Pow(refractionRatio, 2) * (1 - math.Pow(dt, 2));

	if discriminant > 0 {

		var normalDt = normal.Clone();
		normalDt.MulScalar(dt);
		uv.Sub(normalDt);
		uv.MulScalar(refractionRatio);

		var normalDisc = normal.Clone();
		normalDisc.MulScalar(math.Sqrt(discriminant));

		uv.Sub(normalDisc);

		refracted.Copy(uv);

		return true;
	}

	return false;
}

// Real glass has reflectivity that varies with angle look at a window at a steep angle and it becomes a mirror.
// The behavior can be approximated by Christophe Schlick polynomial aproximation.
func Schlick(cosine float64, reflectiveIndex float64) float64 {
	var r = math.Pow((1 - reflectiveIndex) / (1 + reflectiveIndex), 2);
	return r + (1 - r) * math.Pow(1 - cosine, 5);
}

// Calculate a random unitary vector in the surface of a sphere.
func RandomInUnitSphere() *Vector3 {
	var p *Vector3 = NewVector3(0, 0, 0);

	for {
		p.Set(rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0, rand.Float64() * 2.0 - 1.0);

		if p.SquaredLength() < 1.0 {
			break
		}
	}

	return p;
}

// Dot product between two vectors
func Dot(a *Vector3, b *Vector3) float64 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z;
}

// Cross product between two vectors
func Cross(result *Vector3, a *Vector3, b *Vector3) {
	result.X = a.Y * b.Z - a.Z * b.Y;
	result.Y = a.Z * b.X - a.X * b.Z;
	result.Z = a.X * b.Y - a.Y * b.X;
}

// Return a copy of the vector
func (v *Vector3) Clone() *Vector3 {
	return NewVector3(v.X, v.Y, v.Z);
};

// Copy the context of another vector to this one
func (v *Vector3) Copy(b *Vector3) {
	v.X = b.X;
	v.Y = b.Y;
	v.Z = b.Z;
}

// Create a new copy vector with a unit length vector with the same direction as this one.
func (v *Vector3) UnitVector() *Vector3 {
	var unit = v.Clone();
	unit.DivideScalar(v.Length());
	return unit;
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
