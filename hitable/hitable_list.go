package hitable

import "gotracer/vmath"

// A hitable list contains hitable objects to be raytraced.
// Works in the same way as a scene in game engines.
type HitableList struct {
	List []Hitable
}

// Create new hitable list
func NewHitableList() *HitableList {
	return new(HitableList)
}

// Add a hitable element to the list
func (hl *HitableList) Add(h Hitable) {
	hl.List = append(hl.List, h)
}

// Hit iterates and tests all hitable object in the list.
func (hl *HitableList) Hit(r *vmath.Ray, tmin float64, tmax float64, rec *HitRecord) bool {

	var hitAnything = false
	var closestSoFar = tmax
	var tempRec = NewHitRecord()

	for i := 0; i < len(hl.List); i++ {
		if hl.List[i].Hit(r, tmin, closestSoFar, tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			rec.Copy(tempRec)
		}
	}

	return hitAnything
}

// Clone the hitable list and the objects in the list
func (hl *HitableList) Clone() *HitableList {
	var l = NewHitableList()

	for i := 0; i < len(hl.List); i++ {
		l.Add(hl.List[i].Clone())
	}

	return l
}
