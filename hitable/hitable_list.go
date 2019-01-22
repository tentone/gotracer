package hitable;

import "gotracer/vmath";

// A hitable list contains hitable objects to be raytraced.
// Works in the same way as a scene in game engines.
type HitableList struct {
	List []Hitable;
}

// Create new hitable list
func NewHitableList() *HitableList {
	return new(HitableList);
}

// Hit iterates and tests all hitable object in the list.
func (hl *HitableList) Hit(r *vmath.Ray, tmin float64, tmax float64, rec *HitRecord) (bool) {

	var hitAnything bool = false;
	var closestSoFar float64 = tmax;
	var tempRec *HitRecord;

	/*
	for _, h := range hl.List {
		if hit, hr := h.Hit(r, tmin, closestSoFar); hit {
			hitAnything = true;
			rec = hr;
			closestSoFar = hr.t;
		}
	}
	*/
	
	return hitAnything;
}
