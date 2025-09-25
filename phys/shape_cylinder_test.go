// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func TestCylinderCollide(t *testing.T) {
	// Define a standard cylinder for testing.
	cylinder := Cylinder{
		Origin:    r3.Point{X: 0, Y: 0, Z: 0},
		Direction: r3.Vec{X: 0, Y: 1, Z: 0}, // Y-axis
		Radius:    1.0,
		Height:    2.0,
	}

	// Define a set of test cases.
	testCases := []struct {
		name          string
		ray           ray
		expectedHit   bool
		expectedPoint r3.Point
		expectedNorm  r3.Vec
	}{
		{
			name: "Ray intersects cylindrical surface",
			ray: ray{
				origin:    r3.Point{X: 2, Y: 1, Z: 0},
				direction: r3.Vec{X: -1, Y: 0, Z: 0}.Unit(),
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 1, Y: 1, Z: 0},
			expectedNorm:  r3.Vec{X: 1, Y: 0, Z: 0},
		},
		{
			name: "Ray misses the cylinder",
			ray: ray{
				origin:    r3.Point{X: 2, Y: 3, Z: 0},
				direction: r3.Vec{X: 1, Y: 0, Z: 0}.Unit(),
			},
			expectedHit: false,
		},
		{
			name: "Ray intersects top cap",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 3, Z: 0},
				direction: r3.Vec{X: 0, Y: -1, Z: 0}.Unit(),
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0, Y: 2, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 1, Z: 0},
		},
		{
			name: "Ray intersects bottom cap",
			ray: ray{
				origin:    r3.Point{X: 0, Y: -1, Z: 0},
				direction: r3.Vec{X: 0, Y: 1, Z: 0}.Unit(),
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0, Y: 0, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: -1, Z: 0},
		},
		{
			name: "Ray parallel to cylinder axis and outside",
			ray: ray{
				origin:    r3.Point{X: 2, Y: 0, Z: 0},
				direction: r3.Vec{X: 0, Y: 1, Z: 0}.Unit(),
			},
			expectedHit: false,
		},
		{
			name: "Ray parallel to cylinder axis and inside",
			ray: ray{
				origin:    r3.Point{X: 0.5, Y: -1, Z: 0},
				direction: r3.Vec{X: 0, Y: 1, Z: 0}.Unit(),
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0.5, Y: 0, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: -1, Z: 0},
		},
		{
			name: "Ray tangent to cylinder",
			ray: ray{
				origin:    r3.Point{X: 1, Y: 1, Z: -1},
				direction: r3.Vec{X: 0, Y: 0, Z: 1}.Unit(),
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 1, Y: 1, Z: 0},
			expectedNorm:  r3.Vec{X: 1, Y: 0, Z: 0},
		},
		{
			name: "Ray intersects cylindrical surface and caps",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 3, Z: -1},
				direction: r3.Vec{X: 0, Y: -1, Z: 1}.Unit(),
			},
			expectedHit: true,
			// Expected to hit the top cap first
			expectedPoint: r3.Point{X: 0, Y: 2, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 1, Z: 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hit, collision := cylinder.Collide(tc.ray, 0.0001, math.MaxFloat64)
			if hit != tc.expectedHit {
				t.Errorf("Expected hit: %v, got: %v", tc.expectedHit, hit)
			}
			if hit {
				// Check intersection point
				if !collision.at.IsClose(tc.expectedPoint, eps) {
					t.Errorf("Expected intersection point: %v, got: %v", tc.expectedPoint, collision.at)
				}
				// Check normal
				if !collision.normal.IsClose(tc.expectedNorm, eps) {
					t.Errorf("Expected normal: %v, got: %v", tc.expectedNorm, collision.normal)
				}
			}
		})
	}
}

func TestCylinderBounds(t *testing.T) {
	// Define test cases for Bounds
	testCases := []struct {
		name        string
		cylinder    Cylinder
		expectedMin r3.Point
		expectedMax r3.Point
	}{
		{
			name: "Standard cylinder aligned with Y-axis",
			cylinder: Cylinder{
				Origin:    r3.Point{X: 0, Y: 0, Z: 0},
				Direction: r3.Vec{X: 0, Y: 1, Z: 0},
				Radius:    1.0,
				Height:    2.0,
			},
			expectedMin: r3.Point{X: -1, Y: 0, Z: -1},
			expectedMax: r3.Point{X: 1, Y: 2, Z: 1},
		},
		{
			name: "Cylinder offset from origin",
			cylinder: Cylinder{
				Origin:    r3.Point{X: 2, Y: 3, Z: 4},
				Direction: r3.Vec{X: 0, Y: 1, Z: 0},
				Radius:    0.5,
				Height:    1.5,
			},
			expectedMin: r3.Point{X: 1.5, Y: 3, Z: 3.5},
			expectedMax: r3.Point{X: 2.5, Y: 4.5, Z: 4.5},
		},
		{
			name: "Cylinder aligned with Z-axis",
			cylinder: Cylinder{
				Origin:    r3.Point{X: -1, Y: -1, Z: -1},
				Direction: r3.Vec{X: 0, Y: 0, Z: 1},
				Radius:    2.0,
				Height:    3.0,
			},
			expectedMin: r3.Point{X: -3, Y: -3, Z: -1},
			expectedMax: r3.Point{X: 1, Y: 1, Z: 2},
		},
		{
			name: "Cylinder with non-unit axis vector",
			cylinder: Cylinder{
				Origin:    r3.Point{X: 0, Y: 0, Z: 0},
				Direction: r3.Vec{X: 0, Y: 2, Z: 0}, // Not unit vector
				Radius:    1.0,
				Height:    2.0,
			},
			expectedMin: r3.Point{X: -1, Y: 0, Z: -1},
			expectedMax: r3.Point{X: 1, Y: 2, Z: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bounds := tc.cylinder.Bounds()
			if !bounds.Min.IsClose(tc.expectedMin, eps) {
				t.Errorf("Expected Min: %v, got: %v", tc.expectedMin, bounds.Min)
			}
			if !bounds.Max.IsClose(tc.expectedMax, eps) {
				t.Errorf("Expected Max: %v, got: %v", tc.expectedMax, bounds.Max)
			}
		})
	}
}
