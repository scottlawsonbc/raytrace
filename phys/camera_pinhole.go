// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"fmt"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// PinholeCamera represents a basic perspective camera model.
type PinholeCamera struct {
	LowerLeftCorner r3.Point // Lower-left corner of the image plane in world space.
	Origin          r3.Point // Camera's origin point in world space.
	Horizontal      r3.Vec   // Horizontal span of the image plane in world space.
	Vertical        r3.Vec   // Vertical span of the image plane in world space.
}

// Cast generates a ray from the camera origin through the image plane at (s, t).
func (c PinholeCamera) Cast(s, t float64, rand *Rand) ray {
	// Compute the point on the image plane corresponding to (s, t)
	h := c.Horizontal.Muls(s)
	v := c.Vertical.Muls(t)
	imagePoint := c.LowerLeftCorner.Add(h).Add(v)
	// Compute the direction from the camera origin to the image point.
	direction := imagePoint.Sub(c.Origin).Unit()
	return ray{
		origin:    c.Origin,
		direction: direction,
		depth:     0,
		radiance:  Spectrum{1, 1, 1},
		rand:      rand,
	}
}

// Validate checks the PinholeCamera parameters for validity.
func (cam PinholeCamera) Validate() error {
	if cam.LowerLeftCorner.IsNaN() || cam.Origin.IsNaN() {
		return fmt.Errorf("PinholeCamera has NaN values: %+v", cam)
	}
	if cam.LowerLeftCorner.IsInf() || cam.Origin.IsInf() {
		return fmt.Errorf("PinholeCamera has Inf values: %+v", cam)
	}
	if cam.Horizontal.IsZero() {
		return fmt.Errorf("PinholeCamera Horizontal vector is zero: %+v", cam)
	}
	if cam.Vertical.IsZero() {
		return fmt.Errorf("PinholeCamera Vertical vector is zero: %+v", cam)
	}
	if cam.Horizontal.IsNaN() || cam.Vertical.IsNaN() {
		return fmt.Errorf("PinholeCamera has NaN in vectors: Horizontal=%v, Vertical=%v", cam.Horizontal, cam.Vertical)
	}
	if cam.Horizontal.IsInf() || cam.Vertical.IsInf() {
		return fmt.Errorf("PinholeCamera has Inf in vectors: Horizontal=%v, Vertical=%v", cam.Horizontal, cam.Vertical)
	}
	if cam.Horizontal.Cross(cam.Vertical).IsZero() {
		return fmt.Errorf("PinholeCamera Horizontal and Vertical vectors are colinear: Horizontal=%v, Vertical=%v", cam.Horizontal, cam.Vertical)
	}
	return nil
}
