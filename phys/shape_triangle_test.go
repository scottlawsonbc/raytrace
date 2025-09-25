// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// FuzzTriangleCollide performs fuzz testing on the Triangle.Collide method.
// It generates random triangles and rays, validates the triangles, and checks
// for correct collision detection and normal calculations.
// To run this, use command like `go test -fuzz=Fuzz -fuzztime=10s` while
// in the phys directory.
func FuzzTriangleCollide(f *testing.F) {
	// Seed corpus with some valid triangles and rays.
	seedTriangles := []Triangle{
		{
			P0: r3.Point{X: 0, Y: 0, Z: 0},
			P1: r3.Point{X: 1, Y: 0, Z: 0},
			P2: r3.Point{X: 0, Y: 1, Z: 0},
		},
		{
			P0: r3.Point{X: -1, Y: -1, Z: 0},
			P1: r3.Point{X: 1, Y: -1, Z: 0},
			P2: r3.Point{X: 0, Y: 1, Z: 0},
		},
		{
			P0: r3.Point{X: 1, Y: 2, Z: 3},
			P1: r3.Point{X: 4, Y: 5, Z: 6},
			P2: r3.Point{X: 7, Y: 8, Z: 10}, // Non-colinear, non-degenerate
		},
	}

	seedRays := []ray{
		{
			origin:    r3.Point{X: 0.25, Y: 0.25, Z: -1},
			direction: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			origin:    r3.Point{X: 1, Y: 0.5, Z: -1},
			direction: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			origin:    r3.Point{X: 0.5, Y: 0.5, Z: 1},
			direction: r3.Vec{X: 0, Y: 0, Z: -1},
		},
	}

	// Add seed corpus to the fuzzer.
	for _, tri := range seedTriangles {
		for _, r := range seedRays {
			f.Add(
				tri.P0.X, tri.P0.Y, tri.P0.Z,
				tri.P1.X, tri.P1.Y, tri.P1.Z,
				tri.P2.X, tri.P2.Y, tri.P2.Z,
				r.origin.X, r.origin.Y, r.origin.Z,
				r.direction.X, r.direction.Y, r.direction.Z,
			)
		}
	}

	// Define the fuzz function.
	f.Fuzz(func(t *testing.T,
		p0x, p0y, p0z float64,
		p1x, p1y, p1z float64,
		p2x, p2y, p2z float64,
		rayOx, rayOy, rayOz float64,
		rayDx, rayDy, rayDz float64,
	) {
		// Create Triangle from fuzzed points.
		triangle := Triangle{
			P0: r3.Point{X: p0x, Y: p0y, Z: p0z},
			P1: r3.Point{X: p1x, Y: p1y, Z: p1z},
			P2: r3.Point{X: p2x, Y: p2y, Z: p2z},
		}

		// Validate the triangle.
		if err := triangle.Validate(); err != nil {
			// Invalid triangle; skip this input.
			return
		}

		// Create Ray from fuzzed origin and direction.
		ray := ray{
			origin:    r3.Point{X: rayOx, Y: rayOy, Z: rayOz},
			direction: r3.Vec{X: rayDx, Y: rayDy, Z: rayDz},
		}

		// Normalize the ray direction.
		ray.direction = ray.direction.Unit()
		if ray.direction.IsZero() {
			// Degenerate ray direction; skip this input.
			return
		}

		// Perform collision detection.
		hit, coll := triangle.Collide(ray, 0, Distance(math.MaxFloat64))

		if hit {
			// Verify that the intersection point lies on the triangle's plane.
			edge1 := triangle.P1.Sub(triangle.P0)
			edge2 := triangle.P2.Sub(triangle.P0)
			expectedNormal := edge1.Cross(edge2).Unit()

			// Vector from P0 to intersection point.
			vec := coll.at.Sub(triangle.P0)

			// Dot product should be close to zero (point lies on the plane).
			dot := expectedNormal.Dot(r3.Vec{X: vec.X, Y: vec.Y, Z: vec.Z})
			if math.Abs(dot) > eps {
				t.Errorf("Intersection point is not on the triangle's plane: dot=%v", dot)
			}

			// Verify barycentric coordinates are within [0,1] and u + v <= 1.
			if coll.uv.X < -eps || coll.uv.X > 1.0+eps {
				t.Errorf("Barycentric coordinate u out of bounds: u=%v", coll.uv.X)
			}
			if coll.uv.Y < -eps || coll.uv.Y > 1.0+eps {
				t.Errorf("Barycentric coordinate v out of bounds: v=%v", coll.uv.Y)
			}
			if coll.uv.X+coll.uv.Y > 1.0+eps {
				t.Errorf("Barycentric coordinates out of bounds: u+v=%v", coll.uv.X+coll.uv.Y)
			}

			// Verify the normal is correctly computed.
			if !coll.normal.IsClose(expectedNormal, eps) {
				t.Errorf("Normal vector incorrect: expected=%v, got=%v", expectedNormal, coll.normal)
			}

			// Optionally, verify the intersection point using ray equation.
			expectedAt := ray.origin.Add(ray.direction.Muls(float64(coll.t)))
			if !coll.at.IsClose(expectedAt, eps) {
				t.Errorf("Intersection point mismatch: expected=%v, got=%v", expectedAt, coll.at)
			}
		}
	})
}

// TestTriangleCollide tests the Triangle.Collide method across various scenarios.
// It ensures that ray-triangle intersections are accurately detected and that
// the collision data (intersection point and normal) is correctly computed.
func TestTriangleCollide(t *testing.T) {
	// Define a standard triangle in the XY-plane.
	triangle := Triangle{
		P0: r3.Point{X: 0, Y: 0, Z: 0},
		P1: r3.Point{X: 1, Y: 0, Z: 0},
		P2: r3.Point{X: 0, Y: 1, Z: 0},
	}

	// Define test cases.
	testCases := []struct {
		description    string
		ray            ray
		expectedHit    bool
		expectedPoint  r3.Point
		expectedNormal r3.Vec
	}{
		{
			description: "Ray intersects at center",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.25, Y: 0.25, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray misses triangle (outside)",
			ray: ray{
				origin:    r3.Point{X: 1, Y: 1, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    false,
			expectedPoint:  r3.Point{},
			expectedNormal: r3.Vec{},
		},
		{
			description: "Ray parallel to triangle",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0.5, Z: 1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    false,
			expectedPoint:  r3.Point{},
			expectedNormal: r3.Vec{},
		},
		{
			description: "Ray intersects at vertex P0",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0, Y: 0, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray intersects at edge P0-P1",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.5, Y: 0, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray intersects at edge P1-P2",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0.5, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.5, Y: 0.5, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},

		{
			description: "Ray intersects above triangle",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0.5, Z: 1},
				direction: r3.Vec{X: 0, Y: 0, Z: -1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.5, Y: 0.5, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray misses triangle (below)",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: -1},
				direction: r3.Vec{X: 0, Y: -1, Z: 0},
			},
			expectedHit:    false,
			expectedPoint:  r3.Point{},
			expectedNormal: r3.Vec{},
		},
		{
			description: "Ray intersects from behind",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: 1},
				direction: r3.Vec{X: 0, Y: 0, Z: -1},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0.25, Y: 0.25, Z: 0},
			// Following PBR book convention, we don't flip the normal for back-facing hits.
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},

		{
			description: "Intersecting front face at center",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.25, Y: 0.25, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Intersecting front face at edge P1-P2",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0.5, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.5, Y: 0.5, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Intersecting back face at center",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: 1},
				direction: r3.Vec{X: 0, Y: 0, Z: -1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0.25, Y: 0.25, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1}, // Normal remains the same
		},
		{
			description: "Ray misses the triangle",
			ray: ray{
				origin:    r3.Point{X: 1.5, Y: 1.5, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    false,
			expectedPoint:  r3.Point{},
			expectedNormal: r3.Vec{},
		},
		{
			description: "Ray intersects exactly at vertex P0",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0, Y: 0, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray intersects exactly at vertex P1",
			ray: ray{
				origin:    r3.Point{X: 1, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 1, Y: 0, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray intersects exactly at vertex P2",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 1, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:    true,
			expectedPoint:  r3.Point{X: 0, Y: 1, Z: 0},
			expectedNormal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			description: "Ray is parallel to the triangle plane",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: 0},
				direction: r3.Vec{X: 1, Y: 1, Z: 0},
			},
			expectedHit:    false,
			expectedPoint:  r3.Point{},
			expectedNormal: r3.Vec{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			hit, coll := triangle.Collide(tc.ray, 0, Distance(math.MaxFloat64))
			if hit != tc.expectedHit {
				t.Errorf("Expected hit: %v, got: %v", tc.expectedHit, hit)
			}
			if hit {
				// Compare intersection points.
				if !coll.at.IsClose(tc.expectedPoint, eps) {
					t.Errorf("Expected intersection point: %v, got: %v", tc.expectedPoint, coll.at)
				}

				// Compare normals.
				if !coll.normal.IsClose(tc.expectedNormal, eps) {
					t.Errorf("Expected normal: %v, got: %v", tc.expectedNormal, coll.normal)
				}
			}
		})
	}
}

// TestTriangleCollideEdgeCases tests specific edge cases for the Triangle.Collide method.
func TestTriangleCollideEdgeCases(t *testing.T) {
	triangle := Triangle{
		P0: r3.Point{X: 0, Y: 0, Z: 0},
		P1: r3.Point{X: 1, Y: 0, Z: 0},
		P2: r3.Point{X: 0, Y: 1, Z: 0},
	}

	// Define edge cases.
	edgeCases := []struct {
		name        string
		ray         ray
		expectHit   bool
		expectPoint r3.Point
		expectNorm  r3.Vec
	}{
		{
			name: "Ray intersects exactly at P2",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 1, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectHit:   true,
			expectPoint: r3.Point{X: 0, Y: 1, Z: 0},
			expectNorm:  r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			name: "Ray intersects exactly at P1",
			ray: ray{
				origin:    r3.Point{X: 1, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectHit:   true,
			expectPoint: r3.Point{X: 1, Y: 0, Z: 0},
			expectNorm:  r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			name: "Ray grazes the triangle (t == tmin)",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: 0.5, Z: 0},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectHit:   false,
			expectPoint: r3.Point{},
			expectNorm:  r3.Vec{},
		},
		{
			name: "Ray starts on the triangle and points away",
			ray: ray{
				origin:    r3.Point{X: 0.25, Y: 0.25, Z: 0},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectHit:   false,
			expectPoint: r3.Point{},
			expectNorm:  r3.Vec{},
		},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			hit, coll := triangle.Collide(ec.ray, 0, Distance(math.MaxFloat64))
			if hit != ec.expectHit {
				t.Errorf("Expected hit: %v, got: %v", ec.expectHit, hit)
			}
			if hit {
				if !coll.at.IsClose(ec.expectPoint, eps) {
					t.Errorf("Expected intersection point: %v, got: %v", ec.expectPoint, coll.at)
				}
				if !coll.normal.IsClose(ec.expectNorm, eps) {
					t.Errorf("Expected normal: %v, got: %v", ec.expectNorm, coll.normal)
				}
			}
		})
	}
}

// TestTriangleBounds verifies that the Bounds method correctly computes the AABB of the triangle.
func TestTriangleBounds(t *testing.T) {
	triangle := Triangle{
		P0: r3.Point{X: -1, Y: -1, Z: 0},
		P1: r3.Point{X: 1, Y: -1, Z: 0},
		P2: r3.Point{X: 0, Y: 1, Z: 0},
	}

	expectedMin := r3.Point{X: -1, Y: -1, Z: 0}
	expectedMax := r3.Point{X: 1, Y: 1, Z: 0}

	bounds := triangle.Bounds()

	if !bounds.Min.IsClose(expectedMin, eps) {
		t.Errorf("Expected Bounds.Min: %v, got: %v", expectedMin, bounds.Min)
	}

	if !bounds.Max.IsClose(expectedMax, eps) {
		t.Errorf("Expected Bounds.Max: %v, got: %v", expectedMax, bounds.Max)
	}
}

// TestTriangleNormals verifies that the computed normals are correct based on vertex order.
func TestTriangleNormals(t *testing.T) {
	// Triangle in the XY-plane with counter-clockwise winding.
	triangleCCW := Triangle{
		P0: r3.Point{X: 0, Y: 0, Z: 0},
		P1: r3.Point{X: 1, Y: 0, Z: 0},
		P2: r3.Point{X: 0, Y: 1, Z: 0},
	}

	// Triangle in the XY-plane with clockwise winding.
	triangleCW := Triangle{
		P0: r3.Point{X: 0, Y: 0, Z: 0},
		P1: r3.Point{X: 0, Y: 1, Z: 0},
		P2: r3.Point{X: 1, Y: 0, Z: 0},
	}

	// Expected normals.
	expectedNormalCCW := r3.Vec{X: 0, Y: 0, Z: 1}
	expectedNormalCW := r3.Vec{X: 0, Y: 0, Z: -1}

	// Test CCW winding.
	hit, coll := triangleCCW.Collide(ray{
		origin:    r3.Point{X: 0.25, Y: 0.25, Z: -1},
		direction: r3.Vec{X: 0, Y: 0, Z: 1},
	}, 0, Distance(math.MaxFloat64))
	if hit {
		if !coll.normal.IsClose(expectedNormalCCW, eps) {
			t.Errorf("Expected normal: %v, got: %v", expectedNormalCCW, coll.normal)
		}
	} else {
		t.Errorf("Expected hit, but got no hit for CCW triangle.")
	}

	// Test CW winding.
	hit, coll = triangleCW.Collide(ray{
		origin:    r3.Point{X: 0.25, Y: 0.25, Z: 1},
		direction: r3.Vec{X: 0, Y: 0, Z: -1},
	}, 0, Distance(math.MaxFloat64))
	if hit {
		if !coll.normal.IsClose(expectedNormalCW, eps) {
			t.Errorf("Expected normal: %v, got: %v", expectedNormalCW, coll.normal)
		}
	} else {
		t.Errorf("Expected hit, but got no hit for CW triangle.")
	}
}

// TestTriangleCollideMultipleRays tests multiple rays against the same triangle to ensure consistency.
func TestTriangleCollideMultipleRays(t *testing.T) {
	triangle := Triangle{
		P0: r3.Point{X: -1, Y: -1, Z: 0},
		P1: r3.Point{X: 1, Y: -1, Z: 0},
		P2: r3.Point{X: 0, Y: 1, Z: 0},
	}

	rays := []struct {
		name          string
		ray           ray
		expectedHit   bool
		expectedPoint r3.Point
		expectedNorm  r3.Vec
	}{
		{
			name: "Intersect near P0",
			ray: ray{
				origin:    r3.Point{X: -1, Y: -2, Z: -1},
				direction: r3.Vec{X: 0, Y: 1, Z: 1},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: -1, Y: -1, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			name: "Intersect near P1",
			ray: ray{
				origin:    r3.Point{X: 1, Y: -2, Z: -1},
				direction: r3.Vec{X: 0, Y: 1, Z: 1},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 1, Y: -1, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			name: "Intersect above center",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 0, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0, Y: 0, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 0, Z: 1},
		},
		{
			name: "Miss due to direction",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 0, Z: -1},
				direction: r3.Vec{X: 1, Y: 1, Z: 0},
			},
			expectedHit:   false,
			expectedPoint: r3.Point{},
			expectedNorm:  r3.Vec{},
		},
		{
			name: "Miss due to offset",
			ray: ray{
				origin:    r3.Point{X: 2, Y: 2, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1},
			},
			expectedHit:   false,
			expectedPoint: r3.Point{},
			expectedNorm:  r3.Vec{},
		},
	}

	for _, r := range rays {
		t.Run(r.name, func(t *testing.T) {
			hit, coll := triangle.Collide(r.ray, 0, Distance(math.MaxFloat64))
			if hit != r.expectedHit {
				t.Errorf("Expected hit: %v, got: %v", r.expectedHit, hit)
			}
			if hit {
				if !coll.at.IsClose(r.expectedPoint, eps) {
					t.Errorf("Expected intersection point: %v, got: %v", r.expectedPoint, coll.at)
				}
				if !coll.normal.IsClose(r.expectedNorm, eps) {
					t.Errorf("Expected normal: %v, got: %v", r.expectedNorm, coll.normal)
				}
			}
		})
	}
}

// TestTriangleDegenerate tests degenerate triangles (zero area).
func TestTriangleDegenerate(t *testing.T) {
	// Define a degenerate triangle where all points are the same.
	degenerateTriangle := Triangle{
		P0: r3.Point{X: 1, Y: 1, Z: 1},
		P1: r3.Point{X: 1, Y: 1, Z: 1},
		P2: r3.Point{X: 1, Y: 1, Z: 1},
	}

	r := ray{
		origin:    r3.Point{X: 1, Y: 1, Z: 0},
		direction: r3.Vec{X: 0, Y: 0, Z: 1},
	}

	err := degenerateTriangle.Validate()
	if err == nil {
		t.Errorf("Expected validation error for degenerate triangle with identical points, but got none")
	}

	hit, coll := degenerateTriangle.Collide(r, 0, Distance(math.MaxFloat64))
	if hit {
		t.Errorf("Expected no hit for degenerate triangle, but got hit at %v", coll.at)
	}

	// Define a degenerate triangle where all points lie on a line.
	linearTriangle := Triangle{
		P0: r3.Point{X: 0, Y: 0, Z: 0},
		P1: r3.Point{X: 1, Y: 1, Z: 1},
		P2: r3.Point{X: 2, Y: 2, Z: 2},
	}

	r = ray{
		origin:    r3.Point{X: 1, Y: 1, Z: -1},
		direction: r3.Vec{X: 0, Y: 0, Z: 1},
	}

	err = linearTriangle.Validate()
	if err == nil {
		t.Errorf("Expected validation error for degenerate triangle with colinear points, but got none")
	}

	hit, coll = linearTriangle.Collide(r, 0, Distance(math.MaxFloat64))
	if hit {
		t.Errorf("Expected no hit for linear triangle, but got hit at %v", coll.at)
	}
}

// TestTriangleValidate tests the Validate method for the Triangle type.
func TestTriangleValidate(t *testing.T) {
	// Define valid triangles.
	validTriangles := []struct {
		name     string
		triangle Triangle
	}{
		{
			name: "Standard triangle",
			triangle: Triangle{
				P0: r3.Point{X: 0, Y: 0, Z: 0},
				P1: r3.Point{X: 1, Y: 0, Z: 0},
				P2: r3.Point{X: 0, Y: 1, Z: 0},
			},
		},
		{
			name: "Non-axis-aligned triangle",
			triangle: Triangle{
				P0: r3.Point{X: 1, Y: 2, Z: 3},
				P1: r3.Point{X: 4, Y: 5, Z: 6},
				P2: r3.Point{X: 7, Y: 8, Z: 10},
			},
		},
	}

	for _, vt := range validTriangles {
		t.Run(vt.name, func(t *testing.T) {
			err := vt.triangle.Validate()
			if err != nil {
				t.Errorf("Expected valid triangle, but got error: %v", err)
			}
		})
	}

	// Define invalid triangles.
	invalidTriangles := []struct {
		name     string
		triangle Triangle
	}{
		{
			name: "Duplicate vertices P0=P1",
			triangle: Triangle{
				P0: r3.Point{X: 0, Y: 0, Z: 0},
				P1: r3.Point{X: 0, Y: 0, Z: 0},
				P2: r3.Point{X: 0, Y: 1, Z: 0},
			},
		},
		{
			name: "Duplicate vertices P0=P2",
			triangle: Triangle{
				P0: r3.Point{X: 1, Y: 1, Z: 1},
				P1: r3.Point{X: 2, Y: 2, Z: 2},
				P2: r3.Point{X: 1, Y: 1, Z: 1},
			},
		},
		{
			name: "Duplicate vertices P1=P2",
			triangle: Triangle{
				P0: r3.Point{X: -1, Y: -1, Z: -1},
				P1: r3.Point{X: 0, Y: 0, Z: 0},
				P2: r3.Point{X: 0, Y: 0, Z: 0},
			},
		},
		{
			name: "Colinear points",
			triangle: Triangle{
				P0: r3.Point{X: 0, Y: 0, Z: 0},
				P1: r3.Point{X: 1, Y: 1, Z: 1},
				P2: r3.Point{X: 2, Y: 2, Z: 2},
			},
		},
		{
			name: "Zero area triangle",
			triangle: Triangle{
				P0: r3.Point{X: 1, Y: 2, Z: 3},
				P1: r3.Point{X: 4, Y: 5, Z: 6},
				P2: r3.Point{X: 7, Y: 8, Z: 9},
			},
		},
	}

	for _, it := range invalidTriangles {
		t.Run(it.name, func(t *testing.T) {
			err := it.triangle.Validate()
			if err == nil {
				t.Errorf("Expected validation error for invalid triangle, but got none")
			} else {
				t.Logf("Received expected error: %v", err)
			}
		})
	}
}
