// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

// Package phys implements physically based materials and utility shaders used by
// the raytracer. Package phys follows the Go standard library conventions and
// avoids hidden global state.
package phys

import (
	"context"
	"math"
)

// DebugNormal draws the surface normal as a false-color visualization.
//
// The zero value is ready for use. DebugNormal has no internal state and
// therefore provides no concurrency hazards; values are safe for concurrent
// use by multiple goroutines.
//
// Exported fields: none.
type DebugNormal struct{}

// Validate reports whether the material configuration is valid.
//
// DebugNormal requires no configuration and therefore always reports that the
// material is valid. The method returns a nil error to indicate success and
// has no side effects.
func (m DebugNormal) Validate() error {
	return nil
}

// Resolve a surface interaction for this material.
//
// Resolve returns a purely emissive color that encodes the unit surface normal
// in RGB as:
//
//	R = (nx + 1) / 2
//	G = (ny + 1) / 2
//	B = (nz + 1) / 2
//
// The method normalizes the input normal defensively and maps components from
// [-1, 1] into [0, 1]. Resolve ignores the incoming radiance; the purpose of
// this material is to *debug* normals, not to participate in lighting, so the
// result should not be modulated by illumination. The method returns the
// resulting color as emission and performs no other side effects.
func (m DebugNormal) Resolve(ctx context.Context, c surfaceInteraction) resolution {
	n := c.collision.normal

	// Normalize defensively.
	length := math.Sqrt(n.X*n.X + n.Y*n.Y + n.Z*n.Z)
	if !(length > 0) { // catches 0 and NaN
		length = 1
	}
	nx := n.X / length
	ny := n.Y / length
	nz := n.Z / length

	// Map from [-1, 1] to [0, 1].
	r := 0.5 * (nx + 1.0)
	g := 0.5 * (ny + 1.0)
	b := 0.5 * (nz + 1.0)

	// Clamp to [0, 1] to be robust against tiny numeric excursions.
	if r < 0 {
		r = 0
	} else if r > 1 {
		r = 1
	}
	if g < 0 {
		g = 0
	} else if g > 1 {
		g = 1
	}
	if b < 0 {
		b = 0
	} else if b > 1 {
		b = 1
	}

	s := Spectrum{X: r, Y: g, Z: b}
	return resolution{emission: s}
}

// ComputeDirectLighting reports whether and how much direct lighting should be
// added for this material.
//
// DebugNormal is an unlit diagnostic material and therefore contributes no
// direct lighting. The method always returns the zero [Spectrum] and has no
// side effects.
func (m DebugNormal) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	return Spectrum{}
}

func init() {
	RegisterInterfaceType(DebugNormal{})
}

// // Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

// package phys

// import (
// 	"context"
// 	"math"
// )

// // DebugNormal draws the normal vector of the surface as a color.
// type DebugNormal struct{}

// func (m DebugNormal) Validate() error {
// 	return nil
// }

// func (m DebugNormal) Resolve(ctx context.Context, c surfaceInteraction) resolution {
// 	// Normalize the vector to unit length.
// 	n := c.collision.normal
// 	length := math.Sqrt(n.X*n.X + n.Y*n.Y + n.Z*n.Z)
// 	if length == 0 {
// 		length = 1 // Avoid division by zero.
// 	}
// 	nx := n.X / length
// 	ny := n.Y / length
// 	nz := n.Z / length
// 	// Map components from [-1, 1] to [0, 1].
// 	r := (nx + 1) * 0.5
// 	g := (ny + 1) * 0.5
// 	b := (nz + 1) * 0.5
// 	s := Spectrum{X: r, Y: g, Z: b}.Mul(c.incoming.radiance)
// 	return resolution{emission: s}
// }

// func (m DebugNormal) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
// 	// This shader is for debugging; no direct lighting.
// 	return Spectrum{}
// }

// func init() {
// 	RegisterInterfaceType(DebugNormal{})
// }
