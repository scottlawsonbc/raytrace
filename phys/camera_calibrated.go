// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

// Package phys implements a physically based 3D renderer.
// This file adds a calibrated, lens-distorted perspective camera that uses
// OpenCV-style intrinsics and distortion coefficients, with intrinsics and
// extrinsics represented as separate types.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func init() {
	RegisterInterfaceType(CalibratedCamera{})
}

// CalibratedCamera is a perspective camera that applies OpenCV distortion.
//
// Instance purpose:
// CalibratedCamera casts rays that reproduce the projection of a pinhole model
// followed by OpenCV radial and tangential distortion. The Intrinsics and
// Extrinsics are specified separately for clarity.
//
// Concurrency guarantees:
// CalibratedCamera is an immutable value. It is safe to copy by value and share.
//
// Zero value:
// The zero value of CalibratedCamera is not useful. Callers must populate fields.
type CalibratedCamera struct {
	// Intrinsics holds the image geometry and distortion model.
	Intrinsics CameraIntrinsics
	// Extrinsics holds the camera pose and orientation.
	Extrinsics CameraExtrinsics
}

// Cast generates a primary ray for the normalized sample position (s, t).
// The function returns a ray that starts at LookFrom and passes through the
// pixel center corresponding to (s, t) under the distorted projection.
func (cam CalibratedCamera) Cast(s, t float64, rand *Rand) ray {
	ci := cam.Intrinsics
	ce := cam.Extrinsics

	// Camera orthonormal basis.
	w := ce.LookFrom.Sub(ce.LookAt).Unit() // Backward
	u := ce.VUp.Cross(w).Unit()            // Right
	v := w.Cross(u)                        // Up

	// Convert normalized sample to pixel coordinates (top-left origin).
	uPix := s * float64(ci.Width)
	vPix := t * float64(ci.Height)

	// Distorted normalized image coordinates.
	xd := (uPix - ci.Cx) / ci.Fx
	yd := (vPix - ci.Cy) / ci.Fy

	// Undistort to ideal normalized coordinates.
	x, y := ci.undistortNormalized(xd, yd)

	// Camera-space direction. Note: image y grows downward; camera +Y is up.
	dirCam := r3.Vec{X: x, Y: -y, Z: -1.0}.Unit()

	// World-space direction.
	dirWorld :=
		u.Muls(dirCam.X).
			Add(v.Muls(dirCam.Y)).
			Add(w.Muls(dirCam.Z)).
			Unit()

	return ray{
		origin:    ce.LookFrom,
		direction: dirWorld,
		depth:     0,
		radiance:  Spectrum{1, 1, 1},
		rand:      rand,
	}
}

// Validate reports whether the camera can be used to render a scene.
func (cam CalibratedCamera) Validate() error {
	if err := cam.Intrinsics.Validate(); err != nil {
		return fmt.Errorf("CalibratedCamera intrinsics invalid: %v", err)
	}
	if err := cam.Extrinsics.Validate(); err != nil {
		return fmt.Errorf("CalibratedCamera extrinsics invalid: %v", err)
	}
	return nil
}

// NewCalibratedCamera constructs a CalibratedCamera from separate intrinsics and extrinsics.
// The function returns a CalibratedCamera that implements Camera.
func NewCalibratedCamera(intr CameraIntrinsics, extr CameraExtrinsics) CalibratedCamera {
	return CalibratedCamera{Intrinsics: intr, Extrinsics: extr}
}

// CameraIntrinsics stores OpenCV-style intrinsic parameters and image size.
//
// Instance purpose:
// CameraIntrinsics represents the pixel geometry of the image formation
// (focal lengths, principal point), the distortion model, and the resolution.
//
// Concurrency guarantees:
// CameraIntrinsics is an immutable value once constructed. It is safe to copy
// by value and to share across goroutines without additional synchronization.
//
// Zero value:
// The zero value of CameraIntrinsics is not useful. Callers must populate all fields.
type CameraIntrinsics struct {
	// Width is the image width in pixels for which the intrinsics are defined.
	Width int
	// Height is the image height in pixels for which the intrinsics are defined.
	Height int

	// Fx is the focal length in pixels along the x axis.
	Fx float64
	// Fy is the focal length in pixels along the y axis.
	Fy float64
	// Cx is the principal point x coordinate in pixels.
	Cx float64
	// Cy is the principal point y coordinate in pixels.
	Cy float64

	// Distortion parameters follow OpenCV ordering.
	// The standard model uses K1, K2, P1, P2, K3.
	// The rational model additionally uses K4, K5, K6.
	K1 float64
	K2 float64
	P1 float64
	P2 float64
	K3 float64
	K4 float64
	K5 float64
	K6 float64
}

// Validate reports whether the intrinsics are self-consistent and usable.
func (ci CameraIntrinsics) Validate() error {
	if ci.Width <= 0 || ci.Height <= 0 {
		return fmt.Errorf("CameraIntrinsics bad image size: %dx%d", ci.Width, ci.Height)
	}
	if !(ci.Fx > 0 && ci.Fy > 0) {
		return fmt.Errorf("CameraIntrinsics bad focal lengths: Fx=%g Fy=%g", ci.Fx, ci.Fy)
	}
	if math.IsNaN(ci.Cx) || math.IsNaN(ci.Cy) {
		return fmt.Errorf("CameraIntrinsics NaN principal point: Cx=%g Cy=%g", ci.Cx, ci.Cy)
	}
	return nil
}

// K returns the 3x3 pinhole matrix corresponding to the intrinsics.
func (ci CameraIntrinsics) K() [3][3]float64 {
	return [3][3]float64{
		{ci.Fx, 0, ci.Cx},
		{0, ci.Fy, ci.Cy},
		{0, 0, 1},
	}
}

// D returns the distortion vector in OpenCV ordering.
// The function returns a slice of length 5 or 8 depending on whether any K4..K6 are non-zero.
func (ci CameraIntrinsics) D() []float64 {
	if ci.K4 == 0 && ci.K5 == 0 && ci.K6 == 0 {
		return []float64{ci.K1, ci.K2, ci.P1, ci.P2, ci.K3}
	}
	return []float64{ci.K1, ci.K2, ci.P1, ci.P2, ci.K3, ci.K4, ci.K5, ci.K6}
}

// undistortNormalized inverts OpenCV distortion for a single normalized point.
// xd, yd are distorted normalized image coordinates after division by Fx, Fy.
func (ci CameraIntrinsics) undistortNormalized(xd, yd float64) (x, y float64) {
	// Initial guess assumes small distortion.
	x = xd
	y = yd
	const iters = 8
	for i := 0; i < iters; i++ {
		r2 := x*x + y*y
		r4 := r2 * r2
		r6 := r4 * r2

		// Radial term. Use rational model if K4..K6 are provided.
		num := 1.0 + ci.K1*r2 + ci.K2*r4 + ci.K3*r6
		den := 1.0 + ci.K4*r2 + ci.K5*r4 + ci.K6*r6
		if den == 0 {
			den = 1
		}
		radial := num / den

		// Tangential term.
		dx := 2.0*ci.P1*x*y + ci.P2*(r2+2.0*x*x)
		dy := ci.P1*(r2+2.0*y*y) + 2.0*ci.P2*x*y

		// Update by inverting the forward mapping.
		x = (xd - dx) / radial
		y = (yd - dy) / radial
	}
	return x, y
}

// NewCameraIntrinsicsFromKAndD constructs CameraIntrinsics from K and D.
// K is the 3x3 matrix. D is 5 or 8 coefficients in OpenCV order.
func NewCameraIntrinsicsFromKAndD(
	width int,
	height int,
	K [3][3]float64,
	D []float64,
) CameraIntrinsics {
	ci := CameraIntrinsics{
		Width:  width,
		Height: height,
		Fx:     K[0][0],
		Fy:     K[1][1],
		Cx:     K[0][2],
		Cy:     K[1][2],
	}
	if len(D) >= 5 {
		ci.K1, ci.K2, ci.P1, ci.P2, ci.K3 = D[0], D[1], D[2], D[3], D[4]
	}
	if len(D) >= 8 {
		ci.K4, ci.K5, ci.K6 = D[5], D[6], D[7]
	}
	return ci
}

// CameraExtrinsics stores the camera pose and orientation basis.
//
// Instance purpose:
// CameraExtrinsics represents the rigid-body transform of the camera in world
// coordinates via LookFrom, LookAt, and VUp, matching the rest of the package.
//
// Concurrency guarantees:
// CameraExtrinsics is immutable once constructed. It is safe to copy by value.
//
// Zero value:
// The zero value of CameraExtrinsics is not useful. Callers must populate fields.
type CameraExtrinsics struct {
	// LookFrom is the camera origin in world space.
	LookFrom r3.Point
	// LookAt is the point in world space the camera aims at.
	LookAt r3.Point
	// VUp is the camera up direction.
	VUp r3.Vec
}

// Validate reports whether the extrinsics define a proper camera frame.
func (ce CameraExtrinsics) Validate() error {
	if ce.LookFrom == ce.LookAt {
		return fmt.Errorf("CameraExtrinsics LookFrom and LookAt are identical")
	}
	if ce.VUp.IsZero() {
		return fmt.Errorf("CameraExtrinsics VUp is zero")
	}
	// Check orthogonality.
	w := ce.LookFrom.Sub(ce.LookAt).Unit()
	u := ce.VUp.Cross(w).Unit()
	v := w.Cross(u)
	if u.IsNaN() || v.IsNaN() || w.IsNaN() {
		return fmt.Errorf("CameraExtrinsics basis has NaN")
	}
	if math.Abs(u.Dot(v)) > eps || math.Abs(u.Dot(w)) > eps || math.Abs(v.Dot(w)) > eps {
		return fmt.Errorf("CameraExtrinsics basis vectors are not orthogonal: u·v=%g u·w=%g v·w=%g",
			u.Dot(v), u.Dot(w), v.Dot(w))
	}
	return nil
}
