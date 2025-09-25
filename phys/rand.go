// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"
	"math/rand"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Rand wraps a rand.Rand instance to provide additional random generators
// useful for physically based rendering. Each tile worker uses a different one.
type Rand struct {
	*rand.Rand
}

// NewRand creates a new Rand instance with a given seed.
// This allows for reproducible random number sequences.
func NewRand(seed int64) *Rand {
	return &Rand{rand.New(rand.NewSource(seed))}
}

// InUnitSphere returns a random vector uniformly distributed within a unit sphere.
// Useful for volumetric scattering and diffuse reflections.
// Length of the vector is guaranteed to be less than 1.
func (r *Rand) InUnitSphere() r3.Vec {
	for {
		// Generate a random point in the cube [-1,1] x [-1,1] x [-1,1].
		p := r3.Vec{
			X: r.Float64(),
			Y: r.Float64(),
			Z: r.Float64(),
		}.Muls(2).Sub(r3.Vec{X: 1, Y: 1, Z: 1})
		// If the point is inside the unit sphere, return it.
		if p.Length() < 1.0 {
			return p
		}
	}
}

// UnitVector returns a random unit vector uniformly distributed on the surface of a unit sphere.
// Essential for specular reflections and directional sampling.
// The vector is guaranteed to have a length of 1.
func (r *Rand) UnitVector() r3.Vec {
	// Random azimuthal angle between 0 and 2π.
	azimuth := r.Float64() * 2 * math.Pi
	// Random elevation angle (cosine of polar angle) between -1 and 1.
	z := r.Float64()*2 - 1
	// Radius at given z (since it's a unit sphere).
	radius := math.Sqrt(1 - z*z)

	return r3.Vec{
		X: radius * math.Cos(azimuth),
		Y: radius * math.Sin(azimuth),
		Z: z,
	}
}

// InUnitDisk returns a random vector inside a unit disk (circle) in the
// XY-plane centered at the origin. It uses the rejection sampling method
// to ensure uniform distribution.
func (r *Rand) InUnitDisk() r3.Vec {
	for {
		// Generate a random point in the square [-1,1] x [-1,1] at z=0.
		p := r3.Vec{
			X: r.Float64(),
			Y: r.Float64(),
			Z: 0,
		}.Muls(2).Sub(r3.Vec{X: 1, Y: 1, Z: 0})

		if p.Dot(p) < 1.0 {
			return p
		}
	}
}

// CosineWeightedHemisphere samples a random direction in the hemisphere
// with a cosine-weighted distribution. Samples are aligned to provided normal.
func (r *Rand) CosineWeightedHemisphere(normal r3.Vec) r3.Vec {
	// Generate two random numbers in [0,1)
	u1 := r.Float64()
	u2 := r.Float64()

	// Compute polar coordinates for cosine-weighted sampling.
	r1 := math.Sqrt(u1)
	theta := 2 * math.Pi * u2

	// Convert to Cartesian coordinates in the local space (z-up hemisphere).
	x := r1 * math.Cos(theta)
	y := r1 * math.Sin(theta)
	z := math.Sqrt(1 - u1)

	// Create an orthonormal basis (u, v, w) with w = normal.
	var tangent r3.Vec
	if math.Abs(normal.X) > math.Abs(normal.Y) {
		tangent = r3.Vec{X: -normal.Z, Y: 0, Z: normal.X}.Unit()
	} else {
		tangent = r3.Vec{X: 0, Y: normal.Z, Z: -normal.Y}.Unit()
	}
	bitangent := normal.Cross(tangent)

	// Transform the sampled direction from local space to world space.
	sampledDirection := tangent.Muls(x).Add(bitangent.Muls(y)).Add(normal.Muls(z)).Unit()
	return sampledDirection
}
