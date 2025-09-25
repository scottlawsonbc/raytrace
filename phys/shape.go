// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type collision struct {
	t      Distance // Distance along the incoming ray to the collision point.
	at     r3.Point // Collision point on the shape.
	uv     r2.Point // Texture coordinates at the collision point.
	normal r3.Vec   // Normal vector of the surface at the collision point.
}

// Shape represents an geometric object that can collide with rays.
type Shape interface {
	Collide(r ray, tmin Distance, tmax Distance) (bool, collision)
	Bounds() AABB
	Validate() error // Validate checks if the shape is valid.
}

// AABB represents an axis-aligned bounding box.
// AABB is not a Shape itself, but describes the bounds of a Shape.
type AABB struct {
	Min r3.Point
	Max r3.Point
}

func (b AABB) surfaceArea() float64 {
	dx := b.Max.X - b.Min.X
	dy := b.Max.Y - b.Min.Y
	dz := b.Max.Z - b.Min.Z
	return 2 * (dx*dy + dy*dz + dz*dx)
}

func (b AABB) LongestAxis() int {
	dx := b.Max.X - b.Min.X
	dy := b.Max.Y - b.Min.Y
	dz := b.Max.Z - b.Min.Z
	if dx > dy && dx > dz {
		return 0
	} else if dy > dz {
		return 1
	} else {
		return 2
	}
}

func (b AABB) Union(other AABB) AABB {
	return AABB{
		Min: r3.Point{
			X: math.Min(b.Min.X, other.Min.X),
			Y: math.Min(b.Min.Y, other.Min.Y),
			Z: math.Min(b.Min.Z, other.Min.Z),
		},
		Max: r3.Point{
			X: math.Max(b.Max.X, other.Max.X),
			Y: math.Max(b.Max.Y, other.Max.Y),
			Z: math.Max(b.Max.Z, other.Max.Z),
		},
	}
}

func (b AABB) center() r3.Point {
	return r3.Point{
		X: (b.Min.X + b.Max.X) * 0.5,
		Y: (b.Min.Y + b.Max.Y) * 0.5,
		Z: (b.Min.Z + b.Max.Z) * 0.5,
	}
}

func (b AABB) intersects(other AABB) bool {
	return b.Min.X <= other.Max.X && b.Max.X >= other.Min.X &&
		b.Min.Y <= other.Max.Y && b.Max.Y >= other.Min.Y &&
		b.Min.Z <= other.Max.Z && b.Max.Z >= other.Min.Z
}

func (b AABB) hit(r ray, tmin, tmax Distance) bool {
	for axis := 0; axis < 3; axis++ {
		invD := 1.0 / r.direction.Get(axis)
		t0 := (b.Min.Get(axis) - r.origin.Get(axis)) * invD
		t1 := (b.Max.Get(axis) - r.origin.Get(axis)) * invD
		if invD < 0.0 {
			t0, t1 = t1, t0
		}
		tmin = Distance(math.Max(float64(t0), float64(tmin)))
		tmax = Distance(math.Min(float64(t1), float64(tmax)))
		if tmax <= tmin {
			return false
		}
	}
	return true
}
