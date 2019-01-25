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

// Add a hitable element to the list
func (hl *HitableList) Add(h Hitable) {
	hl.List = append(hl.List, h);
}

// Hit iterates and tests all hitable object in the list.
func (hl *HitableList) Hit(r *vmath.Ray, tmin float64, tmax float64, rec *HitRecord) (bool) {

	var hitAnything bool = false;
	var closestSoFar float64 = tmax;
	var tempRec *HitRecord = NewHitRecord();
	var i int = 0;

	for i < len(hl.List) {
		if hl.List[i].Hit(r, tmin, closestSoFar, tempRec){
			hitAnything = true;
			closestSoFar = tempRec.T;
			rec = tempRec;
		}
		i++;
	}

	return hitAnything;
}
