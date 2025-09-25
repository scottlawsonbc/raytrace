// ./shape_triangle_uv.go
package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// TriangleUV represents a triangle with per-vertex texture coordinates.
type TriangleUV struct {
	P0, P1, P2    r3.Point
	UV0, UV1, UV2 r2.Point
	Normal        r3.Vec // Precomputed normal for efficiency.
}

// Validate performs comprehensive validation checks on the TriangleUV instance.
func (tri TriangleUV) Validate() error {
	// Helper function to check for NaN or Inf in r3.Point
	checkPoint3 := func(p r3.Point, label string) error {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsNaN(p.Z) ||
			math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) || math.IsInf(p.Z, 0) {
			return fmt.Errorf("invalid TriangleUV: %s contains NaN or Inf: %+v", label, p)
		}
		return nil
	}

	// Helper function to check for NaN or Inf in r2.Point
	checkPoint2 := func(p r2.Point, label string) error {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) ||
			math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return fmt.Errorf("invalid TriangleUV: %s contains NaN or Inf: %+v", label, p)
		}
		return nil
	}

	// Check all r3.Point coordinates
	if err := checkPoint3(tri.P0, "P0"); err != nil {
		return err
	}
	if err := checkPoint3(tri.P1, "P1"); err != nil {
		return err
	}
	if err := checkPoint3(tri.P2, "P2"); err != nil {
		return err
	}

	// Check all r2.Point UV coordinates
	if err := checkPoint2(tri.UV0, "UV0"); err != nil {
		return err
	}
	if err := checkPoint2(tri.UV1, "UV1"); err != nil {
		return err
	}
	if err := checkPoint2(tri.UV2, "UV2"); err != nil {
		return err
	}

	// Check for duplicate vertices
	if tri.P0 == tri.P1 || tri.P0 == tri.P2 || tri.P1 == tri.P2 {
		return fmt.Errorf("invalid TriangleUV: two or more vertices are identical: %+v", tri)
	}

	// Check normal vector is not zero
	if tri.Normal.IsZero() {
		return fmt.Errorf("invalid TriangleUV: normal is a zero vector: %+v", tri)
	}

	// Check normal vector is unit length
	normalLength := tri.Normal.Length()
	if normalLength < 1-eps || normalLength > 1+eps {
		return fmt.Errorf("invalid TriangleUV: normal is not a unit vector (length: %f): %+v", normalLength, tri)
	}

	// Check triangle area is not degenerate
	edge1 := tri.P1.Sub(tri.P0)
	edge2 := tri.P2.Sub(tri.P0)
	crossProduct := edge1.Cross(edge2)
	area := 0.5 * crossProduct.Length()
	const epsilonArea = 1e-12
	if area < epsilonArea {
		return fmt.Errorf("invalid TriangleUV: triangle is degenerate (zero or near-zero area): %+v", tri)
	}

	// Check normal is orthogonal to edges
	dot1 := tri.Normal.Dot(edge1)
	dot2 := tri.Normal.Dot(edge2)
	if math.Abs(dot1) > eps || math.Abs(dot2) > eps {
		return fmt.Errorf("invalid TriangleUV: normal is not orthogonal to triangle edges (dot1: %f, dot2: %f): %+v", dot1, dot2, tri)
	}

	// Verify normal direction consistency with cross product
	derivedNormal := crossProduct.Unit()
	dotNormal := tri.Normal.Dot(derivedNormal)
	if dotNormal < 1-eps {
		return fmt.Errorf("invalid TriangleUV: normal does not match the direction of the cross product (dot: %f): %+v", dotNormal, tri)
	}

	// Optionally, check UV mapping is non-degenerate (optional based on use-case)
	// Example: Ensure UVs do not form a degenerate triangle
	uvEdge1 := tri.UV1.Sub(tri.UV0)
	uvEdge2 := tri.UV2.Sub(tri.UV0)
	uvCrossProduct := uvEdge1.Cross(uvEdge2)
	uvArea := 0.5 * uvCrossProduct
	const epsilonUVArea = 1e-12
	if uvArea < epsilonUVArea {
		return fmt.Errorf("invalid TriangleUV: UV coordinates form a degenerate triangle: %+v", tri)
	}

	return nil
}

// Collide determines whether a given ray intersects with the TriangleUV.
// It also interpolates the UV coordinates at the intersection point.
func (tri TriangleUV) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	// Möller-Trumbore intersection algorithm.
	edge1 := tri.P1.Sub(tri.P0)
	edge2 := tri.P2.Sub(tri.P0)
	h := r.direction.Cross(edge2)
	a := edge1.Dot(h)
	if a > -eps && a < eps {
		return false, collision{}
	}
	f := 1 / a
	s := r.origin.Sub(tri.P0)
	// Barycentric coordinates.
	u := f * s.Dot(h)
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
	interpolatedUV := tri.UV0.Lerp(tri.UV1, u).Add(tri.UV2.Muls(v))
	return true, collision{
		t:      Distance(t),
		at:     at,
		normal: tri.Normal.Unit(),
		uv:     interpolatedUV,
	}
}

// Bounds computes the Axis-Aligned Bounding Box (AABB) of the TriangleUV.
func (tri TriangleUV) Bounds() AABB {
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

func (tri TriangleUV) String() string {
	return fmt.Sprintf("TriangleUV{P0: %v, P1: %v, P2: %v, UV0: %v, UV1: %v, UV2: %v, Normal: %v}",
		tri.P0, tri.P1, tri.P2, tri.UV0, tri.UV1, tri.UV2, tri.Normal)
}

func init() {
	RegisterInterfaceType(TriangleUV{})
}
