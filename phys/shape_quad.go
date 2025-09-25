// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Quad represents a finite rectangular plane in 3D space.
// It is defined by a center point, a normal vector, and dimensions (width and height).
// Internally, it is represented as two triangles for collision detection.
type Quad struct {
	Center r3.Point // Center of the plane.
	Normal r3.Vec   // Normal vector of the plane (should be a unit vector).
	Width  Distance // Width of the plane.
	Height Distance // Height of the plane.
}

func (q Quad) Validate() error {
	if q.Width <= 0 {
		return fmt.Errorf("invalid Quad width: %v (has it been set?)", q.Width)
	}
	if q.Height <= 0 {
		return fmt.Errorf("invalid Quad height: %v (has it been set?)", q.Height)
	}
	if q.Normal.IsZero() {
		return fmt.Errorf("invalid Quad Normal: %v (has it been set?)", q.Normal)
	}
	if q.Normal.Length() != 1 {
		return fmt.Errorf("invalid Quad Normal should be a unit vector, got: %v", q.Normal)
	}
	return nil
}

// Collide checks for an intersection between a ray and the plane by checking collisions with the two triangles.
func (q Quad) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	normal := q.Normal.Unit()

	// Compute two orthogonal vectors (u and v) in the plane.
	var arbitrary r3.Vec
	if math.Abs(normal.X) < 0.9 {
		arbitrary = r3.Vec{X: 1, Y: 0, Z: 0}
	} else {
		arbitrary = r3.Vec{X: 0, Y: 1, Z: 0}
	}

	// Compute orthogonal vectors u and v.
	u := normal.Cross(arbitrary).Unit()
	v := normal.Cross(u).Unit()

	// Scale u and v by half the width and height.
	halfWidth := float64(q.Width) / 2
	halfHeight := float64(q.Height) / 2
	u = u.Muls(halfWidth)
	v = v.Muls(halfHeight)

	// Compute the four corner points of the plane.
	p0 := q.Center.Subv(u).Subv(v) // Bottom-left corner.
	p1 := q.Center.Add(u).Subv(v)  // Bottom-right corner.
	p2 := q.Center.Add(u).Add(v)   // Top-right corner.
	p3 := q.Center.Subv(u).Add(v)  // Top-left corner.

	// Create two triangles from the corner points.
	tri1 := Triangle{P0: p0, P1: p1, P2: p2}
	tri2 := Triangle{P0: p0, P1: p2, P2: p3}

	// Check for collisions with the two triangles.
	hit1, c1 := tri1.Collide(r, tmin, tmax)
	hit2, c2 := tri2.Collide(r, tmin, tmax)

	var hit bool
	var c collision

	if hit1 && (!hit2 || c1.t < c2.t) {
		hit = true
		c = c1
	} else if hit2 {
		hit = true
		c = c2
	}

	if hit {
		// Compute UV coordinates based on the hit point.
		// Map the collision point back to local plane coordinates.

		// Set local origin to p1 to align UV (0,0) at p1
		localOrigin := p1
		localU := p2.Sub(p1) // Vector along U axis (Width)
		localV := p0.Sub(p1) // Vector along V axis (Height)

		// Compute local coordinates (s, t)
		hitPoint := c.at.Sub(localOrigin)
		uCoord := hitPoint.Dot(localU) / localU.Dot(localU)
		vCoord := hitPoint.Dot(localV) / localV.Dot(localV)

		// Clamp UV coordinates to [0,1] to handle floating-point inaccuracies
		uCoord = math.Max(0, math.Min(1, uCoord))
		vCoord = math.Max(0, math.Min(1, vCoord))

		// SCOTT TODO IS THIS RIGHT?
		uCoord = 1 - uCoord
		vCoord = 1 - vCoord

		c.uv = r2.Point{X: uCoord, Y: vCoord}
		c.normal = normal // Ensure normal is set correctly

		// Debugging Statements
		// fmt.Printf("Hit Point: %+v, UV: %+v\n", c.at, c.uv)
	}

	return hit, c
}

// Bounds computes the axis-aligned bounding box of the plane.
func (q Quad) Bounds() AABB {
	normal := q.Normal.Unit()
	// Compute two orthogonal vectors (u and v) in the plane.
	// Choose an arbitrary vector that is not parallel to the normal.
	var arbitrary r3.Vec
	if math.Abs(normal.X) < 0.9 {
		arbitrary = r3.Vec{X: 1, Y: 0, Z: 0}
	} else {
		arbitrary = r3.Vec{X: 0, Y: 1, Z: 0}
	}
	// Compute orthogonal vectors u and v.
	u := normal.Cross(arbitrary).Unit()
	v := normal.Cross(u).Unit()
	// Scale u and v by half the width and height.
	halfWidth := float64(q.Width) / 2
	halfHeight := float64(q.Height) / 2
	u = u.Muls(halfWidth)
	v = v.Muls(halfHeight)
	// Compute the four corner points of the plane.
	p0 := q.Center.Subv(u).Subv(v) // Bottom-left corner.
	p1 := q.Center.Add(u).Subv(v)  // Bottom-right corner.
	p2 := q.Center.Add(u).Add(v)   // Top-right corner.
	p3 := q.Center.Subv(u).Add(v)  // Top-left corner.
	// Compute bounds from the four corner points.
	minX := math.Min(math.Min(p0.X, p1.X), math.Min(p2.X, p3.X))
	minY := math.Min(math.Min(p0.Y, p1.Y), math.Min(p2.Y, p3.Y))
	minZ := math.Min(math.Min(p0.Z, p1.Z), math.Min(p2.Z, p3.Z))
	maxX := math.Max(math.Max(p0.X, p1.X), math.Max(p2.X, p3.X))
	maxY := math.Max(math.Max(p0.Y, p1.Y), math.Max(p2.Y, p3.Y))
	maxZ := math.Max(math.Max(p0.Z, p1.Z), math.Max(p2.Z, p3.Z))
	return AABB{
		Min: r3.Point{X: minX, Y: minY, Z: minZ},
		Max: r3.Point{X: maxX, Y: maxY, Z: maxZ},
	}
}

func init() {
	RegisterInterfaceType(Quad{})
}
