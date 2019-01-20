package hitable;

import "gotracer/vmath";

// A hitable list contains hitable objects to be raytraced.
// Works in the same way as a scene in game engines.
type HitableList struct {
	List []Hitable;
}

// Hit iterates and tests all hitable object in the list.
func (hl *HitableList) Hit(r *Ray, tmin float64, tmax float64, rec *HitRecord) (bool) {

	var hitAnything = false;
	var closestSoFar = tmax;

	for _, h := range hl.List {
		if hit, hr := h.Hit(r, tmin, closestSoFar); hit {
			hitAnything = true;
			rec = hr;
			closestSoFar = hr.t;
		}
	}

	return hitAnything;
}
