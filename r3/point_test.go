// Copyright Scott Lawson 2024. All rights reserverd.

package r3_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func ExamplePoint_distance() {
	// Calculating the distance between two points
	p1 := r3.Point{X: 1, Y: 2, Z: 3}
	p2 := r3.Point{X: 4, Y: 5, Z: 6}

	// Calculate the vector difference
	vec := p1.Sub(p2)

	// Compute the distance using the length of the vector
	distance := vec.Length()

	fmt.Printf("The distance between %v and %v is %v\n", p1, p2, distance)
	// Output: The distance between (1, 2, 3) and (4, 5, 6) is 5.196152422706632
}

func ExamplePoint_movingObject() {
	// Moving an object along a direction
	position := r3.Point{X: 0, Y: 0, Z: 0}

	// Direction vector (must be a unit vector)
	direction := r3.Vec{X: 1, Y: 1, Z: 0}.Unit()

	// Speed and time
	speed := 10.0    // units per second
	deltaTime := 0.5 // seconds
	distance := speed * deltaTime

	// Calculate displacement
	displacement := direction.Muls(distance)

	// Update position
	newPosition := position.Add(displacement)

	fmt.Printf("New position of the object is %v\n", newPosition)
	// Output: New position of the object is (3.5355339059327373, 3.5355339059327373, 0)
}

func ExamplePoint_interpolation() {
	// Interpolating between two points
	start := r3.Point{X: 0, Y: 0, Z: 0}
	end := r3.Point{X: 10, Y: 10, Z: 10}

	// Parameter t ranges from 0.0 to 1.0
	t := 0.5
	// Interpolate between start and end
	position := start.Lerp(end, t)
	fmt.Printf("Position at t=%.2f: %v\n", t, position)
	// Output: Position at t=0.50: (5, 5, 5)
}

func TestPointSet(t *testing.T) {
	p := r3.Point{1, 2, 3}
	tests := []struct {
		index    int
		value    float64
		expected r3.Point
	}{
		{0, 10, r3.Point{10, 2, 3}},
		{1, 20, r3.Point{1, 20, 3}},
		{2, 30, r3.Point{1, 2, 30}},
	}

	for _, test := range tests {
		result := p.Set(test.index, test.value)
		if result != test.expected {
			t.Errorf("Set(%d, %v): expected %v, got %v", test.index, test.value, test.expected, result)
		}
	}

	// Test for panic on invalid index
	invalidIndices := []int{-1, 3}
	for _, index := range invalidIndices {
		func(idx int) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Set did not panic on invalid index %d", idx)
				}
			}()
			p.Set(idx, 0)
		}(index)
	}
}

func TestPointGet(t *testing.T) {
	p := r3.Point{1, 2, 3}
	tests := []struct {
		index    int
		expected float64
	}{
		{0, 1},
		{1, 2},
		{2, 3},
	}

	for _, test := range tests {
		result := p.Get(test.index)
		if result != test.expected {
			t.Errorf("Get(%d): expected %v, got %v", test.index, test.expected, result)
		}
	}

	// Test for panic on invalid index
	invalidIndices := []int{-1, 3}
	for _, index := range invalidIndices {
		func(idx int) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Get did not panic on invalid index %d", idx)
				}
			}()
			p.Get(idx)
		}(index)
	}
}

func TestPointSub(t *testing.T) {
	p1 := r3.Point{1, 2, 3}
	p2 := r3.Point{4, 5, 6}
	expected := r3.Vec{-3, -3, -3}
	result := p1.Sub(p2)
	if result != expected {
		t.Errorf("Sub: expected %v, got %v", expected, result)
	}
}

func TestPointAdd(t *testing.T) {
	p := r3.Point{1, 2, 3}
	v := r3.Vec{4, 5, 6}
	expected := r3.Point{5, 7, 9}
	result := p.Add(v)
	if result != expected {
		t.Errorf("Add: expected %v, got %v", expected, result)
	}
}

func TestPointSubv(t *testing.T) {
	p := r3.Point{1, 2, 3}
	v := r3.Vec{4, 5, 6}
	expected := r3.Point{-3, -3, -3}
	result := p.Subv(v)
	if result != expected {
		t.Errorf("Subv: expected %v, got %v", expected, result)
	}
}

func TestPointLerp(t *testing.T) {
	p1 := r3.Point{0, 0, 0}
	p2 := r3.Point{10, 10, 10}
	tests := []struct {
		tParam   float64
		expected r3.Point
	}{
		{0, r3.Point{0, 0, 0}},
		{0.5, r3.Point{5, 5, 5}},
		{1, r3.Point{10, 10, 10}},
		{-0.5, r3.Point{0, 0, 0}},   // Clamped to 0
		{1.5, r3.Point{10, 10, 10}}, // Clamped to 1
	}

	for _, test := range tests {
		result := p1.Lerp(p2, test.tParam)
		if result != test.expected {
			t.Errorf("Lerp(%v): expected %v, got %v", test.tParam, test.expected, result)
		}
	}
}

func TestPointEq(t *testing.T) {
	p1 := r3.Point{1, 2, 3}
	p2 := r3.Point{1, 2, 3}
	p3 := r3.Point{4, 5, 6}

	if !p1.Eq(p2) {
		t.Errorf("Eq: expected %v to equal %v", p1, p2)
	}
	if p1.Eq(p3) {
		t.Errorf("Eq: expected %v not to equal %v", p1, p3)
	}
}

func TestPointIsClose(t *testing.T) {
	p1 := r3.Point{1.0000001, 2.0000001, 3.0000001}
	p2 := r3.Point{1.0000002, 2.0000002, 3.0000002}
	p3 := r3.Point{1.1, 2.1, 3.1}
	atol := 1e-6

	if !p1.IsClose(p2, atol) {
		t.Errorf("IsClose: expected %v to be close to %v within %v", p1, p2, atol)
	}
	if p1.IsClose(p3, atol) {
		t.Errorf("IsClose: expected %v not to be close to %v within %v", p1, p3, atol)
	}
}

func TestPointIsNaN(t *testing.T) {
	pNaN := r3.Point{math.NaN(), 0, 0}
	if !pNaN.IsNaN() {
		t.Errorf("IsNaN: expected %v to be NaN", pNaN)
	}

	pValid := r3.Point{0, 0, 0}
	if pValid.IsNaN() {
		t.Errorf("IsNaN: expected %v not to be NaN", pValid)
	}
}

func TestPointIsInf(t *testing.T) {
	pInf := r3.Point{math.Inf(1), 0, 0}
	if !pInf.IsInf() {
		t.Errorf("IsInf: expected %v to be Inf", pInf)
	}

	pValid := r3.Point{0, 0, 0}
	if pValid.IsInf() {
		t.Errorf("IsInf: expected %v not to be Inf", pValid)
	}
}

func TestPointIsZero(t *testing.T) {
	pZero := r3.Point{0, 0, 0}
	if !pZero.IsZero() {
		t.Errorf("IsZero: expected %v to be zero", pZero)
	}

	pNonZero := r3.Point{1e-9, 0, 0}
	if pNonZero.IsZero() {
		t.Errorf("IsZero: expected %v not to be zero", pNonZero)
	}
}

func TestPointString(t *testing.T) {
	p := r3.Point{1.1, 2.2, 3.3}
	expected := "(1.1, 2.2, 3.3)"
	result := p.String()
	if result != expected {
		t.Errorf("String: expected %v, got %v", expected, result)
	}
}
