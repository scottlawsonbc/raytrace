/*

// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.package phys
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Sphere represents a sphere with a center and radius.
type Sphere struct {
	Center r3.Point `json:"Center"`
	Radius Distance `json:"Radius"`
}

// Ensure Sphere satisfies the Shape interface.
var _ Shape = (*Sphere)(nil)

func (s Sphere) Validate() error {
	if s.Radius <= 0 {
		return fmt.Errorf("invalid Sphere radius: %v (has it been set?)", s.Radius)
	}
	return nil
}

// // Collide method implementation...
// func (s Sphere) Collide(r ray, tmin, tmax Distance) (bool, collision) {
// 	oc := r.origin.Sub(s.Center)
// 	a := r.direction.Dot(r.direction)
// 	b := oc.Dot(r.direction)
// 	c := oc.Dot(oc) - float64(s.Radius*s.Radius)
// 	discriminant := b*b - a*c
// 	if discriminant < 0 {
// 		return false, collision{}
// 	}
// 	sqrtD := math.Sqrt(discriminant)
// 	t := (-b - sqrtD) / a
// 	if t < float64(tmin) || t > float64(tmax) {
// 		t = (-b + sqrtD) / a
// 		if t < float64(tmin) || t > float64(tmax) {
// 			return false, collision{}
// 		}
// 	}
// 	at := r.at(Distance(t))
// 	normal := at.Sub(s.Center).Unit()
// 	// Compute UV coordinates.
// 	theta := math.Acos(normal.Y)          // from 0 to pi
// 	phi := math.Atan2(normal.Z, normal.X) // from -pi to pi
// 	if phi < 0 {
// 		phi += 2 * math.Pi
// 	}
// 	u := phi / (2 * math.Pi)
// 	v := theta / math.Pi // Or is it 1 - theta / math.Pi? SCOTT
// 	return true, collision{
// 		t:      Distance(t),
// 		at:     at,
// 		normal: at.Sub(s.Center).Unit(),
// 		uv:     r3.Vec{X: u, Y: v},
// 	}
// }

// Bounds method implementation...
func (s Sphere) Bounds() AABB {
	radius := float64(s.Radius)
	return AABB{
		Min: r3.Point{
			X: s.Center.X - radius,
			Y: s.Center.Y - radius,
			Z: s.Center.Z - radius,
		},
		Max: r3.Point{
			X: s.Center.X + radius,
			Y: s.Center.Y + radius,
			Z: s.Center.Z + radius,
		},
	}
}

func init() {
	RegisterInterfaceType(Sphere{})
}

// Collide method implementation with Box Mapping.
func (s Sphere) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	oc := r.origin.Sub(s.Center)
	a := r.direction.Dot(r.direction)
	b := oc.Dot(r.direction)
	c := oc.Dot(oc) - float64(s.Radius*s.Radius)
	discriminant := b*b - a*c
	if discriminant < 0 {
		return false, collision{}
	}
	sqrtD := math.Sqrt(discriminant)
	t := (-b - sqrtD) / a
	if t < float64(tmin) || t > float64(tmax) {
		t = (-b + sqrtD) / a
		if t < float64(tmin) || t > float64(tmax) {
			return false, collision{}
		}
	}
	at := r.at(Distance(t))
	normal := at.Sub(s.Center).Unit()
	// Compute UV coordinates using Box Mapping.
	uv := boxMapUV(normal)
	return true, collision{
		t:      Distance(t),
		at:     at,
		normal: normal,
		uv:     uv,
	}
}

// boxMapUV maps a normal vector to UV coordinates using box mapping.
func boxMapUV(normal r3.Vec) (uv r2.Point) {
	absX := math.Abs(normal.X)
	absY := math.Abs(normal.Y)
	absZ := math.Abs(normal.Z)

	var isXPositive, isYPositive, isZPositive bool
	isXPositive = normal.X > 0
	isYPositive = normal.Y > 0
	isZPositive = normal.Z > 0

	var maxAxis float64
	var uc, vc float64

	// Determine the major axis
	if absX >= absY && absX >= absZ {
		// Major axis is X
		maxAxis = absX
		if isXPositive {
			// +X face
			uc = -normal.Z
			vc = normal.Y
		} else {
			// -X face
			uc = normal.Z
			vc = normal.Y
		}
	} else if absY >= absX && absY >= absZ {
		// Major axis is Y
		maxAxis = absY
		if isYPositive {
			// +Y face
			uc = normal.X
			vc = -normal.Z
		} else {
			// -Y face
			uc = normal.X
			vc = normal.Z
		}
	} else {
		// Major axis is Z
		maxAxis = absZ
		if isZPositive {
			// +Z face
			uc = normal.X
			vc = normal.Y
		} else {
			// -Z face
			uc = -normal.X
			vc = normal.Y
		}
	}

	// Calculate u and v in [0,1]
	uv = r2.Point{
		X: (uc/math.Abs(maxAxis) + 1) / 2,
		Y: (vc/math.Abs(maxAxis) + 1) / 2,
	}.Clip(0, 1)
	return
}

*/
// Package phys provides physically motivated types and utilities for rendering.
// It includes shapes, textures, spectra, and helpers. Symbols document units
// explicitly and state zero-value and concurrency properties.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func init() {
	RegisterInterfaceType(Sphere{})
}

// UVMapKind selects how a shape converts a surface point (or normal) to UV.
//
// Zero value:
//
//	The zero value is UVMapEquirect and is usable.
//
// Semantics:
//
//	Equirectangular mappings use spherical coordinates. Box mapping selects a
//	face of the axis-aligned cube that encloses the unit sphere and projects
//	onto that face.
//
// Units:
//
//	Inputs are dimensionless directions (typically unit surface normals).
//	Outputs are UV in [0,1]x[0,1].
type UVMapKind uint8

const (
	// UVMapEquirect maps using longitude/latitude with +Y as the north pole.
	// u = (atan2(z, x) + π) / (2π), v = 1 - acos(y)/π.
	UVMapEquirect UVMapKind = iota

	// UVMapEquirectSouthUp maps using longitude/latitude with +Y at v = 0.
	// u = (atan2(z, x) + π) / (2π), v = acos(y)/π.
	UVMapEquirectSouthUp

	// UVMapBox maps using axis-major "cube" (a.k.a. box) projection.
	// The dominant axis selects the face, then coordinates are normalized to [0,1].
	UVMapBox
)

// Sphere represents a sphere with a center and radius.
//
// Zero value:
//
//	The zero value is not usable because Radius is 0. Callers must set a
//	positive Radius before using the sphere.
//
// Units and semantics:
//
//	Center is in scene length units.
//	Radius is in scene length units.
//	UVMap chooses how the sphere provides UV coordinates at hit points.
//
// Concurrency:
//
//	Sphere is immutable after construction. Concurrent calls to Collide and
//	Bounds are safe provided the referenced types (e.g., ray) are used safely.
type Sphere struct {
	// Center is the sphere center in scene units.
	Center r3.Point `json:"Center"`

	// Radius is the sphere radius in scene units. Radius must be > 0 for a
	// valid instance.
	Radius Distance `json:"Radius"`

	// UVMap selects the UV parameterization used by Collide. The zero value
	// is UVMapEquirect.
	UVMap UVMapKind `json:"UVMap,omitempty"`
}

// Ensure Sphere satisfies the Shape interface.
var _ Shape = (*Sphere)(nil)

// Validate reports whether s has a positive radius.
//
// Validate returns nil when Radius > 0. It does not mutate the receiver.
func (s Sphere) Validate() error {
	if s.Radius <= 0 {
		return fmt.Errorf("invalid Sphere radius: %v (has it been set?)", s.Radius)
	}
	return nil
}

// Bounds returns the axis-aligned bounding box of the sphere in scene units.
//
// Bounds assumes Radius >= 0 and does not mutate the receiver.
func (s Sphere) Bounds() AABB {
	r := float64(s.Radius)
	return AABB{
		Min: r3.Point{X: s.Center.X - r, Y: s.Center.Y - r, Z: s.Center.Z - r},
		Max: r3.Point{X: s.Center.X + r, Y: s.Center.Y + r, Z: s.Center.Z + r},
	}
}

// Collide reports whether the ray r intersects the sphere within [tmin, tmax].
//
// On success it returns true and a populated collision:
//   - t is the parametric distance along r in scene length units.
//   - at is the hit position in scene units.
//   - normal is the outward unit surface normal at the hit.
//   - uv is computed according to s.UVMap, with u, v in [0,1].
//
// On failure it returns false and a zero collision. Collide does not mutate r.
func (s Sphere) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	oc := r.origin.Sub(s.Center)

	// Quadratic coefficients for |o + t*d - c|^2 == R^2.
	a := r.direction.Dot(r.direction)            // = |d|^2
	b := oc.Dot(r.direction)                     // = (o-c)·d
	c := oc.Dot(oc) - float64(s.Radius*s.Radius) // = |o-c|^2 - R^2
	discriminant := b*b - a*c                    // = (b^2 - a*c)
	if discriminant < 0 {
		return false, collision{}
	}

	sqrtD := math.Sqrt(discriminant)

	// Try nearer root first.
	t := (-b - sqrtD) / a
	if t < float64(tmin) || t > float64(tmax) {
		t = (-b + sqrtD) / a
		if t < float64(tmin) || t > float64(tmax) {
			return false, collision{}
		}
	}

	hitT := Distance(t)
	at := r.at(hitT)

	// Geometric normal.
	normal := at.Sub(s.Center).Unit()

	// UV parameterization.
	var uv r2.Point
	switch s.UVMap {
	case UVMapEquirect:
		uv = equirectUV(normal, true)
	case UVMapEquirectSouthUp:
		uv = equirectUV(normal, false)
	case UVMapBox:
		uv = boxMapUV(normal)
	default:
		uv = equirectUV(normal, true)
	}

	return true, collision{
		t:      hitT,
		at:     at,
		normal: normal,
		uv:     uv,
	}
}

// equirectUV returns longitude/latitude UV for a unit direction n.
//
// Semantics:
//
//	u increases with longitude from -π..π (X toward +Z), remapped to [0,1].
//	v increases from south to north if northUp is false (v = θ/π),
//	or from north to south if northUp is true (v = 1 - θ/π).
//
// Inputs/Outputs:
//
//	n is expected to be unit length (method does not re-normalize).
//	Return value uv is clamped to [0,1].
func equirectUV(n r3.Vec, northUp bool) r2.Point {
	phi := math.Atan2(n.Z, n.X) // [-π, π]
	if phi < 0 {
		phi += 2 * math.Pi
	}
	u := phi / (2 * math.Pi)

	theta := math.Acos(max(-1, min(1, n.Y))) // [0, π]
	var v float64
	if northUp {
		// v = 1 at north pole (+Y), 0 at south pole (-Y).
		v = 1 - theta/math.Pi
	} else {
		// v = 0 at north pole (+Y), 1 at south pole (-Y).
		v = theta / math.Pi
	}
	return r2.Point{X: u, Y: v}.Clip(0, 1)
}

// boxMapUV maps a unit direction n to UV using axis-major cube projection.
//
// The dominant component selects a face:
//
//	+X, -X, +Y, -Y, +Z, -Z.
//
// The face-local coordinates are normalized to [0,1]. The mapping is continuous
// per face but not across edges. Input n is expected to be unit length.
func boxMapUV(n r3.Vec) r2.Point {
	ax := math.Abs(n.X)
	ay := math.Abs(n.Y)
	az := math.Abs(n.Z)

	var maxAxis, uc, vc float64

	switch {
	case ax >= ay && ax >= az:
		maxAxis = ax
		if n.X >= 0 {
			// +X
			uc = -n.Z
			vc = n.Y
		} else {
			// -X
			uc = n.Z
			vc = n.Y
		}
	case ay >= ax && ay >= az:
		maxAxis = ay
		if n.Y >= 0 {
			// +Y
			uc = n.X
			vc = -n.Z
		} else {
			// -Y
			uc = n.X
			vc = n.Z
		}
	default:
		maxAxis = az
		if n.Z >= 0 {
			// +Z
			uc = n.X
			vc = n.Y
		} else {
			// -Z
			uc = -n.X
			vc = n.Y
		}
	}

	u := (uc/math.Abs(maxAxis) + 1) / 2
	v := (vc/math.Abs(maxAxis) + 1) / 2
	return r2.Point{X: u, Y: v}.Clip(0, 1)
}
