// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// OrthographicCamera represents an orthographic camera model where rays are parallel
// and orthogonal to the image plane.
type OrthographicCamera struct {
	LookFrom  r3.Point // Camera position in world space.
	LookAt    r3.Point // Point in world space the camera is looking at.
	VUp       r3.Vec   // Camera relative up vector.
	FOVHeight Distance // Field of view height in world units.
	FOVWidth  Distance // Field of view width in world units.
}

// Cast generates a parallel ray from a point on the image plane at (s, t).
// s and t are normalized coordinates across the image plane.
func (cam OrthographicCamera) Cast(s, t float64, rand *Rand) ray {
	// Compute the camera's orthonormal basis vectors.
	w := cam.LookFrom.Sub(cam.LookAt).Unit() // Camera direction vector (pointing backwards).
	u := cam.VUp.Cross(w).Unit()             // Camera right vector.
	v := w.Cross(u)                          // Camera up vector.

	// Compute the origin point on the image plane corresponding to (s, t).
	origin := cam.LookFrom.
		Add(u.Muls(float64(cam.FOVWidth) * (s - 0.5))). // Offset along the right vector.
		Add(v.Muls(float64(cam.FOVHeight) * (t - 0.5))) // Offset along the up vector.

	// The direction is constant for all rays in an orthographic projection.
	direction := cam.LookAt.Sub(cam.LookFrom).Unit() // Direction from camera to lookat point.

	r := ray{
		origin:    origin,
		direction: direction,
		depth:     0,
		radiance:  Spectrum{1, 1, 1},
		rand:      rand,
	}
	return r
}

func (cam OrthographicCamera) Validate() error {
	if cam.FOVHeight <= 0 || cam.FOVWidth <= 0 {
		return fmt.Errorf("Cast FOVHeight and FOVWidth must be positive: %v", cam)
	}
	if cam.LookFrom.IsNaN() || cam.LookAt.IsNaN() || cam.VUp.IsNaN() {
		return fmt.Errorf("Camera has has NaN values: %v", cam)
	}
	if cam.LookFrom.IsInf() || cam.LookAt.IsInf() || cam.VUp.IsInf() {
		return fmt.Errorf("Camera has Inf values: %v", cam)
	}
	if cam.LookFrom == cam.LookAt {
		return fmt.Errorf("Camera LookFrom and LookAt points are the same: %v", cam)
	}
	if cam.VUp.IsZero() {
		return fmt.Errorf("Camera VUp vector is zero: %v", cam)
	}
	w := cam.LookFrom.Sub(cam.LookAt).Unit() // Camera direction vector (pointing backwards).
	u := cam.VUp.Cross(w).Unit()             // Camera right vector.
	v := w.Cross(u)                          // Camera up vector.
	if u.IsNaN() || v.IsNaN() || w.IsNaN() {
		return fmt.Errorf("Camera basis vectors are NaN: u=%v, v=%v, w=%v for %+v", u, v, w, cam)
	}
	if math.Abs(u.Dot(v)) > eps || math.Abs(u.Dot(w)) > eps || math.Abs(v.Dot(w)) > eps {
		return fmt.Errorf("Camera basis vectors are not orthogonal: u·v=%f, u·w=%f, v·w=%f", u.Dot(v), u.Dot(w), v.Dot(w))
	}
	return nil
}

func init() {
	RegisterInterfaceType(OrthographicCamera{})
}
