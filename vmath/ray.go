package vmath;

//Ray is represented by an origin point A and a normalized direction vector B
type Ray struct {
	Origin *Vector3;
	Direction *Vector3;
}

// Create new ray from origin point and direction
func NewRay(origin *Vector3, direction *Vector3) *Ray {
	var r = new(Ray);
	r.Origin = origin;
	r.Direction = direction;
	return r;
}

//Get the point at a certain distance in the direction of the vector from the origin
func (r *Ray) PointAtParameter(t float64) *Vector3 {
	var v = r.Direction.Clone();
	v.MulScalar(t);
	v.Add(r.Origin);

	return v;
}

