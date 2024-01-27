package vmath

// Ray is represented by an origin point A and a normalized direction vector B
type Ray struct {
	// Origin of the ray
	Origin *Vector3

	// Normalized direction of the ray
	Direction *Vector3
}

// Create new ray from origin point and direction
func NewRay(origin *Vector3, direction *Vector3) *Ray {
	var r = new(Ray)
	r.Origin = origin
	r.Direction = direction
	return r
}

// Create new empty ray
func NewEmptyRay() *Ray {
	var r = new(Ray)
	r.Origin = NewVector3(0.0, 0.0, 0.0)
	r.Direction = NewVector3(0.0, 0.0, 1.0)
	return r
}

// Get the point at a certain distance in the direction of the vector from the origin
func (r *Ray) PointAtParameter(t float64) *Vector3 {
	var v = r.Direction.Clone()
	v.MulScalar(t)
	v.Add(r.Origin)
	return v
}

// Set the values of this array.
// Internally copies the value of the vectors passed as paramters.
func (r *Ray) Set(origin *Vector3, direction *Vector3) {
	r.Origin.Copy(origin)
	r.Direction.Copy(direction)
}

// Clone this array into a new array.
func (r *Ray) Clone() *Ray {
	var nr = new(Ray)
	nr.Origin = r.Origin.Clone()
	nr.Direction = r.Direction.Clone()
	return nr
}
