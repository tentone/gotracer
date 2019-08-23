package geometry

import (
	"gotracer/material"
	"gotracer/vmath"
)

// A scene (hittable list) contains hittable objects to be ray traced.
// Works in the same way as a scene in game engines.
type Scene struct {
	List []Hitable
}

// Create new hittable list
func NewScene() *Scene {
	return new(Scene)
}

// Add a hittable element to the list
func (scene *Scene) Add(h Hitable) {
	scene.List = append(scene.List, h)
}

// Hit iterates and tests all hittable object in the list.
func (scene *Scene) Hit(r *vmath.Ray, tmin float64, tmax float64, rec *material.HitRecord) bool {

	var hitAnything = false
	var closestSoFar = tmax
	var tempRec = material.NewHitRecord()

	for i := 0; i < len(scene.List); i++ {
		if scene.List[i].Hit(r, tmin, closestSoFar, tempRec) {
			hitAnything = true
			closestSoFar = tempRec.T
			rec.Copy(tempRec)
		}
	}

	return hitAnything
}

// Clone the hittable list and the objects in the list
func (scene *Scene) Clone() *Scene {
	var l = NewScene()

	for i := 0; i < len(scene.List); i++ {
		l.Add(scene.List[i].Clone())
	}

	return l
}
