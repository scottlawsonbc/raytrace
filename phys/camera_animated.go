// Copyright 2024 Scott Lawson
//
// Package phys implements a physically based 3D renderer.
// This file adds AnimatedCamera, a tiny adapter that turns a parameterized
// camera generator into a concrete [Camera]. The parameter u in [0,1)
// typically represents progress through an animation cycle (for example,
// an orbit turn). AnimatedCamera composes with calibrated intrinsics and
// can implement an orbit camera in a single line.
package phys

import (
	"fmt"
	"math"
	"time"
)

// Compile-time interface check. AnimatedCamera must implement [Camera].
var _ Camera = AnimatedCamera{}

// CameraFunc is a functional camera factory that returns a concrete [Camera]
// for the given normalized progress parameter u in [0,1).
//
// Instance purpose:
// CameraFunc decouples "how to generate a camera" from "when to sample it".
// AnimatedCamera holds a CameraFunc and a current u, and delegates to it.
//
// Concurrency guarantees:
// CameraFunc must be pure from the perspective of AnimatedCamera. In other
// words, given the same u, it must return an equivalent [Camera] without
// observable data races when called concurrently.
//
// Zero value:
// The zero value of a function variable is nil and not usable.
type CameraFunc func(u float64) Camera

// AnimatedCamera implements [Camera] by sampling a parameterized camera
// generator at a fixed progress value u.
//
// Instance purpose:
// AnimatedCamera is a minimal adapter for frame-by-frame animation. Your
// outer loop updates u (or uses WithU/WithTime/Advance), replaces the camera
// in the [Scene], and renders a single frame. Internally, AnimatedCamera
// builds a concrete camera for each ray by calling the CameraFunc.
//
// Concurrency guarantees:
// AnimatedCamera is an immutable value when used as intended. The type has
// no interior mutability. Share by value across goroutines. If different
// frames need different u, pass distinct copies produced with WithU.
//
// Zero value:
// The zero value of AnimatedCamera is not usable because Build is nil.
// Construct it with NewAnimatedCamera or a convenience helper such as
// NewAnimatedOrbit.
type AnimatedCamera struct {
	// Build is the parameterized camera factory. Build must be non-nil.
	Build CameraFunc

	// U is the normalized progress in [0,1). U==0 selects the first pose.
	// Values outside [0,1) are wrapped by Cast using U - floor(U).
	U float64

	// Period controls time mapping for WithTime and Advance. When Period is
	// zero, those helpers will return an error instead of guessing.
	Period time.Duration
}

// Cast generates a primary ray for normalized image coordinates (s, t).
// Cast wraps U into [0,1), obtains a concrete [Camera] by calling Build,
// and delegates ray generation to it. Cast has no side effects.
func (ac AnimatedCamera) Cast(s, t float64, rand *Rand) ray {
	u := ac.wrap01(ac.U)
	cam := ac.Build(u)
	return cam.Cast(s, t, rand)
}

// Validate reports whether the AnimatedCamera can generate rays.
// Validate checks that Build is non-nil and that the concrete [Camera]
// produced at the current U validates successfully.
func (ac AnimatedCamera) Validate() error {
	if ac.Build == nil {
		return fmt.Errorf("AnimatedCamera.Build is nil")
	}
	u := ac.wrap01(ac.U)
	cam := ac.Build(u)
	if cam == nil {
		return fmt.Errorf("AnimatedCamera.Build(%g) returned nil Camera", u)
	}
	return cam.Validate()
}

// WithU returns a copy of AnimatedCamera with progress set to u.
// WithU does not clamp or wrap. Cast wraps at call time.
func (ac AnimatedCamera) WithU(u float64) AnimatedCamera {
	ac.U = u
	return ac
}

// WithTime returns a copy of AnimatedCamera with progress mapped from t.
// WithTime divides t by Period and wraps into [0,1). WithTime returns an
// error if Period is zero.
func (ac AnimatedCamera) WithTime(t time.Duration) (AnimatedCamera, error) {
	if ac.Period == 0 {
		return AnimatedCamera{}, fmt.Errorf("AnimatedCamera.WithTime: Period is zero")
	}
	turns := float64(t) / float64(ac.Period)
	return ac.WithU(turns), nil
}

// Advance returns a copy with dt added to time progress relative to Period.
// Advance returns an error if Period is zero.
func (ac AnimatedCamera) Advance(dt time.Duration) (AnimatedCamera, error) {
	if ac.Period == 0 {
		return AnimatedCamera{}, fmt.Errorf("AnimatedCamera.Advance: Period is zero")
	}
	du := float64(dt) / float64(ac.Period)
	return ac.WithU(ac.U + du), nil
}

// Frames returns n [Camera] values sampled uniformly over one full cycle.
// If n <= 0, Frames returns nil.
func (ac AnimatedCamera) Frames(n int) []Camera {
	if n <= 0 || ac.Build == nil {
		return nil
	}
	out := make([]Camera, n)
	for i := 0; i < n; i++ {
		u := float64(i) / float64(n)
		out[i] = ac.Build(u)
	}
	return out
}

// NewAnimatedCamera constructs an AnimatedCamera from a CameraFunc, an
// initial progress u, and a period. The function returns an AnimatedCamera
// that implements [Camera].
func NewAnimatedCamera(build CameraFunc, u float64, period time.Duration) AnimatedCamera {
	return AnimatedCamera{
		Build:  build,
		U:      u,
		Period: period,
	}
}

// wrap01 wraps x into [0,1) using x - floor(x).
func (ac AnimatedCamera) wrap01(x float64) float64 {
	return x - math.Floor(x)
}

// // NewAnimatedOrbit returns an AnimatedCamera that implements a simple orbit
// // over an [OrbitSpec]. The animation parameter u in [0,1) maps to the orbit
// // azimuth angle theta = 2πu. The returned AnimatedCamera uses the provided
// // period for time mapping helpers.
// //
// // Behavior:
// //   - The spec is validated once. If spec is invalid, the function panics.
// //     Use NewAnimatedOrbitSafe to receive an error instead of a panic.
// //   - The underlying camera for each sample is an [OrbitCamera] with the
// //     same spec but a per-sample angle.
// func NewAnimatedOrbit(spec OrbitSpec, period time.Duration) AnimatedCamera {
// 	if err := spec.Validate(); err != nil {
// 		panic(fmt.Errorf("NewAnimatedOrbit: invalid spec: %v", err))
// 	}
// 	// Capture a base OrbitCamera by value and vary only the angle per sample.
// 	base := OrbitCamera{Spec: spec}
// 	build := func(u float64) Camera {
// 		return base.WithAngle(2 * math.Pi * (u - math.Floor(u)))
// 	}
// 	return NewAnimatedCamera(build, 0, period)
// }

// // NewAnimatedOrbitSafe is the error-returning variant of NewAnimatedOrbit.
// // The function reports whether spec is valid and returns a constructed
// // AnimatedCamera if so.
// func NewAnimatedOrbitSafe(spec OrbitSpec, period time.Duration) (AnimatedCamera, error) {
// 	if err := spec.Validate(); err != nil {
// 		return AnimatedCamera{}, err
// 	}
// 	base := OrbitCamera{Spec: spec}
// 	build := func(u float64) Camera {
// 		return base.WithAngle(2 * math.Pi * (u - math.Floor(u)))
// 	}
// 	return NewAnimatedCamera(build, 0, period), nil
// }

// // NewAnimatedXYOrbit is a convenience for a circular turntable in the XY
// // plane with a height offset along +Z. It constructs an [OrbitSpec] from
// // axes and returns an AnimatedCamera. The function panics on invalid input.
// // Use NewAnimatedXYOrbitSafe for error returns.
// func NewAnimatedXYOrbit(
// 	intr CameraIntrinsics,
// 	center r3.Point,
// 	radius Distance,
// 	heightOffset Distance,
// 	vup r3.Vec,
// 	period time.Duration,
// ) AnimatedCamera {
// 	spec, err := NewOrbitSpecFromAxes(
// 		intr,
// 		center,
// 		r3.Vec{X: 1, Y: 0, Z: 0},
// 		r3.Vec{X: 0, Y: 1, Z: 0},
// 		radius,
// 		radius,
// 		heightOffset,
// 		0,
// 		vup,
// 	)
// 	if err != nil {
// 		panic(fmt.Errorf("NewAnimatedXYOrbit: %v", err))
// 	}
// 	return NewAnimatedOrbit(spec, period)
// }

// // NewAnimatedXYOrbitSafe is the error-returning variant of NewAnimatedXYOrbit.
// // The function reports whether the inputs yield a valid orbit.
// func NewAnimatedXYOrbitSafe(
// 	intr CameraIntrinsics,
// 	center r3.Point,
// 	radius Distance,
// 	heightOffset Distance,
// 	vup r3.Vec,
// 	period time.Duration,
// ) (AnimatedCamera, error) {
// 	spec, err := NewOrbitSpecFromAxes(
// 		intr,
// 		center,
// 		r3.Vec{X: 1, Y: 0, Z: 0},
// 		r3.Vec{X: 0, Y: 1, Z: 0},
// 		radius,
// 		radius,
// 		heightOffset,
// 		0,
// 		vup,
// 	)
// 	if err != nil {
// 		return AnimatedCamera{}, err
// 	}
// 	return NewAnimatedOrbitSafe(spec, period)
// }
