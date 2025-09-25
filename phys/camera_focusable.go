// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// FocusableCamera represents a camera with adjustable position and orientation,
// including depth of field effects via aperture and focus distance.
type FocusableCamera struct {
	LookFrom        r3.Point // Camera position in world space.
	LookAt          r3.Point // Point in world space the camera is looking at.
	VUp             r3.Vec   // Up vector defining the camera's orientation.
	FOVHeight       Distance // Field of view height at the focal distance.
	FOVWidth        Distance // Field of view width at the focal distance.
	Aperture        Distance // Aperture size controlling depth of field.
	WorkingDistance Distance // Distance from the camera to the focal plane.
}

// Cast generates a ray from the camera through the image plane at (s, t),
// incorporating depth of field by simulating a thin lens.
func (cam FocusableCamera) Cast(s, t float64, rand *Rand) ray {
	// Compute the camera's orthonormal basis vectors.
	w := cam.LookFrom.Sub(cam.LookAt).Unit() // Camera direction vector (pointing backwards)
	u := cam.VUp.Cross(w).Unit()             // Camera right vector
	v := w.Cross(u)                          // Camera up vector

	// Calculate the size of the image plane at the focal distance
	horizontal := u.Muls(float64(cam.FOVWidth * cam.WorkingDistance)) // Horizontal span
	vertical := v.Muls(float64(cam.FOVHeight * cam.WorkingDistance))  // Vertical span

	// Compute the lower-left corner of the image plane.
	lowerLeftCorner := cam.LookFrom.
		Subv(horizontal.Divs(2)).
		Subv(vertical.Divs(2)).
		Subv(w.Muls(float64(cam.WorkingDistance)))

	// Simulate depth of field by sampling a random point on the lens aperture.
	lensRadius := cam.Aperture / 2
	rd := rand.InUnitDisk().Muls(float64(lensRadius)) // Random point in unit disk scaled by lens radius.
	offset := u.Muls(rd.X).Add(v.Muls(rd.Y))          // Offset from the lens center.

	// Compute the ray's origin and direction.
	origin := cam.LookFrom.Add(offset) // Ray origin with lens offset.
	imagePoint := lowerLeftCorner.
		Add(horizontal.Muls(s)).
		Add(vertical.Muls(t))
	direction := imagePoint.Sub(origin).Unit() // Direction from origin to image point.

	return ray{
		origin:    origin,
		direction: direction,
		depth:     0,                 // Depth 0 for primary rays.
		radiance:  Spectrum{1, 1, 1}, // TODO: scott revist this in the future when doing spectral rendering.
		rand:      rand,
		pixelX:    0, // This is set by the renderer.
		pixelY:    0, // This is set by the renderer.
	}
}

// Validate checks the FocusableCamera parameters for validity.
func (cam FocusableCamera) Validate() error {
	if cam.FOVHeight <= 0 || cam.FOVWidth <= 0 {
		return fmt.Errorf("FocusableCamera FOVHeight and FOVWidth must be positive: %v", cam)
	}
	if cam.WorkingDistance <= 0 {
		return fmt.Errorf("FocusableCamera WorkingDistance must be positive: %v", cam)
	}
	if cam.Aperture < 0 {
		return fmt.Errorf("FocusableCamera Aperture cannot be negative: %v", cam)
	}
	if cam.LookFrom.IsNaN() || cam.LookAt.IsNaN() || cam.VUp.IsNaN() {
		return fmt.Errorf("FocusableCamera has NaN values: %+v", cam)
	}
	if cam.LookFrom.IsInf() || cam.LookAt.IsInf() || cam.VUp.IsInf() {
		return fmt.Errorf("FocusableCamera has Inf values: %+v", cam)
	}
	if cam.LookFrom == cam.LookAt {
		return fmt.Errorf("FocusableCamera LookFrom and LookAt points are the same: %+v", cam)
	}
	if cam.VUp.IsZero() {
		return fmt.Errorf("FocusableCamera VUp vector is zero: %+v", cam)
	}
	// Compute the camera's orthonormal basis vectors.
	w := cam.LookFrom.Sub(cam.LookAt).Unit()
	u := cam.VUp.Cross(w).Unit()
	v := w.Cross(u)
	if u.IsNaN() || v.IsNaN() || w.IsNaN() {
		return fmt.Errorf("FocusableCamera basis vectors are NaN: u=%v, v=%v, w=%v for %+v", u, v, w, cam)
	}
	if math.Abs(u.Dot(v)) > eps || math.Abs(u.Dot(w)) > eps || math.Abs(v.Dot(w)) > eps {
		return fmt.Errorf("FocusableCamera basis vectors are not orthogonal: u·v=%f, u·w=%f, v·w=%f", u.Dot(v), u.Dot(w), v.Dot(w))
	}
	return nil
}
