// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func TestQuadCollide(t *testing.T) {
	// Define the quad's parameters.
	quad := Quad{
		Center: r3.Point{X: 0, Y: 0, Z: 0},
		Normal: r3.Vec{X: 0, Y: 1, Z: 0},
		Width:  2.0,
		Height: 2.0,
	}

	// Define test cases.
	testCases := []struct {
		name          string
		ray           ray
		expectedHit   bool
		expectedPoint r3.Point
		expectedNorm  r3.Vec
		expectedUV    r2.Point
	}{
		{
			name: "Ray hits quad from front",
			ray: ray{
				origin:    r3.Point{X: 0, Y: -1, Z: 0},
				direction: r3.Vec{X: 0, Y: 1, Z: 0},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: 0, Y: 0, Z: 0},
			expectedNorm:  r3.Vec{X: 0, Y: 1, Z: 0},
			expectedUV:    r2.Point{X: 0.5, Y: 0.5}, // Center of the quad.
		},
		{
			name: "Ray misses quad",
			ray: ray{
				origin:    r3.Point{X: 0, Y: -1, Z: 2},
				direction: r3.Vec{X: 0, Y: 1, Z: 0},
			},
			expectedHit: false,
		},
		{
			name: "Ray hits quad at corner",
			ray: ray{
				origin:    r3.Point{X: -1, Y: -1, Z: -1},
				direction: r3.Vec{X: 0, Y: 1, Z: 0},
			},
			expectedHit:   true,
			expectedPoint: r3.Point{X: -1, Y: 0, Z: -1},
			expectedNorm:  r3.Vec{X: 0, Y: 1, Z: 0},
			expectedUV:    r2.Point{X: 0.0, Y: 0.0}, // Bottom-left corner.
		},
		{
			name: "Ray parallel to quad",
			ray: ray{
				origin:    r3.Point{X: 0, Y: 0, Z: 0},
				direction: r3.Vec{X: 1, Y: 0, Z: 0},
			},
			expectedHit: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hit, collision := quad.Collide(tc.ray, eps, Distance(1000))
			if hit != tc.expectedHit {
				t.Errorf("Expected hit: %v, got: %v", tc.expectedHit, hit)
			}
			if hit {
				if !collision.at.IsClose(tc.expectedPoint, eps) {
					t.Errorf("Expected collision point: %v, got: %v", tc.expectedPoint, collision.at)
				}
				if !collision.normal.IsClose(tc.expectedNorm, eps) {
					t.Errorf("Expected normal: %v, got: %v", tc.expectedNorm, collision.normal)
				}
				if !collision.uv.IsClose(tc.expectedUV, eps) {
					t.Errorf("Expected UV: %v, got: %v", tc.expectedUV, collision.uv)
				}
			}
		})
	}
}
