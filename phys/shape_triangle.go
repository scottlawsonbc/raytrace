// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Triangle represents a triangle defined by three vertices in 3D space.
// The vertices P0, P1, and P2 should be defined in a counter-clockwise order
// when looking from the front side of the triangle. This winding order ensures
// that the computed normal points in the correct direction.
//
// Visual Representation:
//
//	      P2
//	     /  \
//	    /    \
//	   /      \
//	  /________\
//	P0          P1
//	      |
//	      | Normal (N)
//	      V
//
// Normal Calculation:
//  1. Compute two edge vectors:
//     - Edge1 = P1 - P0
//     - Edge2 = P2 - P0
//  2. Calculate the cross product of Edge1 and Edge2:
//     - Cross = Edge1 × Edge2
//  3. Normalize the cross product to obtain the unit normal vector:
//     - Normal = Cross.Unit()
//
// In this configuration, the normal vector points outward from the front face
// of the triangle, perpendicular to its surface.
//
// The Triangle type is fundamental for collision detection and rendering within
// the ray tracer, serving as a basic geometric primitive.
type Triangle struct {
	P0, P1, P2 r3.Point // Vertices of the triangle in counter-clockwise order.
}

// Validate performs comprehensive validation checks on the Triangle instance.
// It ensures that:
// 1. All three vertices are distinct.
// 2. The triangle is non-degenerate (has a non-zero area).
// 3. The triangle's vertices are not colinear.
//
// Returns:
//   - error: An error describing the validation failure, or nil if the triangle is valid.
func (tri Triangle) Validate() error {
	// 1. Check that all three vertices are distinct.
	if tri.P0 == tri.P1 || tri.P0 == tri.P2 || tri.P1 == tri.P2 {
		return fmt.Errorf("invalid Triangle: two or more vertices are identical")
	}

	// 2. Compute the vectors representing two edges of the triangle.
	edge1 := tri.P1.Sub(tri.P0)
	edge2 := tri.P2.Sub(tri.P0)

	// 3. Compute the cross product of the edge vectors to determine if the triangle is degenerate.
	crossProduct := edge1.Cross(edge2)

	// 4. Compute the area of the triangle. If the area is zero (or nearly zero), the triangle is degenerate.
	// Area = 0.5 * |edge1 x edge2|
	area := 0.5 * crossProduct.Length()

	// Define a small epsilon to account for floating-point precision errors.
	const epsilonArea = 1e-12

	if area < epsilonArea {
		return fmt.Errorf("invalid Triangle: triangle is degenerate (zero or near-zero area)")
	}

	// 5. (Optional) Check for colinearity: if the cross product is zero, the points are colinear.
	// This is already implied by the area check, but can be explicitly stated if desired.
	if crossProduct.IsZero() {
		return fmt.Errorf("invalid Triangle: vertices are colinear")
	}

	return nil
}

// Collide determines whether a given ray intersects with the triangle.
// It implements the Möller–Trumbore intersection algorithm, which is efficient
// for detecting ray-triangle intersections.
//
// The algorithm computes if and where the ray intersects the triangle within
// the bounds [tmin, tmax]. If an intersection occurs, it returns true along
// with the collision details, including the intersection point and the
// triangle's normal at that point.
//
// Parameters:
//   - r ray: The ray to test for intersection.
//   - tmin Distance: The minimum distance along the ray to consider for intersections.
//   - tmax Distance: The maximum distance along the ray to consider for intersections.
//
// Returns:
//   - bool: True if the ray intersects the triangle within [tmin, tmax], otherwise false.
//   - collision: The collision data containing the intersection point, normal, and other relevant information.
func (tri Triangle) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	edge1 := tri.P1.Sub(tri.P0)
	edge2 := tri.P2.Sub(tri.P0)
	h := r.direction.Cross(edge2)
	a := edge1.Dot(h)
	if a > -eps && a < eps {
		return false, collision{}
	}
	f := 1 / a
	s := r.origin.Sub(tri.P0)
	u := f * s.Dot(h)
	// if u < 0 || u > 1 {
	// 	return false, collision{}
	// }
	// By including eps, allow endpoints to be hit.
	if u < -eps || u > 1.0+eps {
		return false, collision{}
	}
	q := s.Cross(edge1)
	v := f * r.direction.Dot(q)
	if v < -eps || u+v > 1.0+eps {
		return false, collision{}
	}
	t := f * edge2.Dot(q)
	if t < float64(tmin) || t > float64(tmax) {
		return false, collision{}
	}
	at := r.at(Distance(t))
	return true, collision{
		t:      Distance(t),
		at:     at,
		normal: edge1.Cross(edge2).Unit(),
		uv:     r2.Point{X: u, Y: v},
	}
}

// Bounds computes the Axis-Aligned Bounding Box (AABB) of the triangle.
// The AABB is the smallest box aligned with the coordinate axes that completely
// contains the triangle. This is useful for broad-phase collision detection.
//
// Returns:
//   - AABB: The axis-aligned bounding box of the triangle.
func (tri Triangle) Bounds() AABB {
	min := r3.Point{
		X: math.Min(math.Min(tri.P0.X, tri.P1.X), tri.P2.X),
		Y: math.Min(math.Min(tri.P0.Y, tri.P1.Y), tri.P2.Y),
		Z: math.Min(math.Min(tri.P0.Z, tri.P1.Z), tri.P2.Z),
	}
	max := r3.Point{
		X: math.Max(math.Max(tri.P0.X, tri.P1.X), tri.P2.X),
		Y: math.Max(math.Max(tri.P0.Y, tri.P1.Y), tri.P2.Y),
		Z: math.Max(math.Max(tri.P0.Z, tri.P1.Z), tri.P2.Z),
	}
	return AABB{Min: min, Max: max}
}
