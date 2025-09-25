// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func TestAABBIntersects(t *testing.T) {
	testcases := []struct {
		box1, box2 AABB
		intersect  bool
	}{
		{
			box1:      AABB{Min: r3.Point{X: 0, Y: 0, Z: 0}, Max: r3.Point{X: 1, Y: 1, Z: 1}},
			box2:      AABB{Min: r3.Point{X: 0, Y: 0, Z: 0}, Max: r3.Point{X: 1, Y: 1, Z: 1}},
			intersect: true,
		},
		{
			box1:      AABB{Min: r3.Point{X: 0, Y: 0, Z: 0}, Max: r3.Point{X: 1, Y: 1, Z: 1}},
			box2:      AABB{Min: r3.Point{X: -2, Y: -2, Z: -2}, Max: r3.Point{X: -1, Y: -1, Z: -1}},
			intersect: false,
		},
	}
	for _, tc := range testcases {
		if got := tc.box1.intersects(tc.box2); got != tc.intersect {
			t.Errorf("got %v, want %v", got, tc.intersect)
		}
	}
}
