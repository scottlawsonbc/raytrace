// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// TransformedShape wraps a Shape with a Transform.
type TransformedShape struct {
	Shape     Shape
	Transform Transform
}

func (ts TransformedShape) Validate() error {
	if ts.Shape == nil {
		return fmt.Errorf("TransformedShape: Shape is nil")
	}
	return ts.Shape.Validate()
}

// Collide transforms the ray into the local space, performs collision,
// and transforms the collision back to world space.
func (ts TransformedShape) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	// Transform the ray into the local space of the shape.
	invTransform := ts.Transform.Inverse()
	localOrigin := invTransform.ApplyToPoint(r.origin)
	localDirection := invTransform.ApplyToVector(r.direction)

	localRay := ray{
		origin:    localOrigin,
		direction: localDirection,
		depth:     r.depth,
		radiance:  r.radiance,
		rand:      r.rand,
		pixelX:    r.pixelX,
		pixelY:    r.pixelY,
	}

	hit, col := ts.Shape.Collide(localRay, tmin, tmax)
	if !hit {
		return false, collision{}
	}

	// Transform collision back to world space
	worldPoint := ts.Transform.ApplyToPoint(col.at)
	worldNormal := ts.Transform.ApplyToVector(col.normal).Unit()

	return true, collision{
		t:      col.t,
		at:     worldPoint,
		normal: worldNormal,
		uv:     col.uv,
	}
}

// Bounds transforms the bounding box of the shape.
func (ts TransformedShape) Bounds() AABB {
	bounds := ts.Shape.Bounds()

	// Transform all 8 corners of the AABB
	min := bounds.Min
	max := bounds.Max
	corners := []r3.Point{
		{X: min.X, Y: min.Y, Z: min.Z},
		{X: max.X, Y: min.Y, Z: min.Z},
		{X: min.X, Y: max.Y, Z: min.Z},
		{X: max.X, Y: max.Y, Z: min.Z},
		{X: min.X, Y: min.Y, Z: max.Z},
		{X: max.X, Y: min.Y, Z: max.Z},
		{X: min.X, Y: max.Y, Z: max.Z},
		{X: max.X, Y: max.Y, Z: max.Z},
	}

	var transformedCorners []r3.Point
	for _, corner := range corners {
		transformedCorners = append(transformedCorners, ts.Transform.ApplyToPoint(corner))
	}

	// Compute new AABB
	newMin := transformedCorners[0]
	newMax := transformedCorners[0]
	for _, p := range transformedCorners[1:] {
		newMin.X = math.Min(newMin.X, p.X)
		newMin.Y = math.Min(newMin.Y, p.Y)
		newMin.Z = math.Min(newMin.Z, p.Z)
		newMax.X = math.Max(newMax.X, p.X)
		newMax.Y = math.Max(newMax.Y, p.Y)
		newMax.Z = math.Max(newMax.Z, p.Z)
	}

	return AABB{Min: newMin, Max: newMax}
}

func init() {
	RegisterInterfaceType(TransformedShape{})
}
