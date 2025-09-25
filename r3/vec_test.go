package r3_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func ExampleVec_angleBetweenVectors() {
	// Calculating the angle between two vectors
	v1 := r3.Vec{X: 1, Y: 0, Z: 0}
	v2 := r3.Vec{X: 0, Y: 1, Z: 0}

	// Calculate the dot product
	dotProduct := v1.Dot(v2)

	// Calculate the magnitudes
	magnitudeV1 := v1.Length()
	magnitudeV2 := v2.Length()

	// Calculate the angle in radians
	angleRadians := math.Acos(dotProduct / (magnitudeV1 * magnitudeV2))

	// Convert to degrees
	angleDegrees := angleRadians * (180 / math.Pi)

	fmt.Printf("The angle between %v and %v is %.2f degrees\n", v1, v2, angleDegrees)
	// Output: The angle between (1, 0, 0) and (0, 1, 0) is 90.00 degrees
}

func ExampleVec_reflection() {
	// Reflecting a vector off a surface
	// Incoming vector (must be a unit vector)
	incoming := r3.Vec{X: 1, Y: -1, Z: 0}.Unit()

	// Normal vector of the surface (must be a unit vector)
	normal := r3.Vec{X: 0, Y: 1, Z: 0}

	// Calculate the reflection vector
	dotProduct := incoming.Dot(normal)
	reflection := incoming.Sub(normal.Muls(2 * dotProduct))

	fmt.Printf("Reflection vector: %v\n", reflection)
	// Output: Reflection vector: (0.7071067811865475, 0.7071067811865475, 0)
}

func ExamplePoint_projectOntoPlane() {
	// Projecting a point onto a plane
	// Point to project
	point := r3.Point{X: 5, Y: 5, Z: 5}

	// Point on the plane
	planePoint := r3.Point{X: 0, Y: 0, Z: 0}

	// Normal vector of the plane (assumed to be normalized)
	normal := r3.Vec{X: 0, Y: 1, Z: 0}

	// Vector from plane point to the point
	vector := point.Sub(planePoint)

	// Distance from the point to the plane along the normal
	distance := vector.Dot(normal)

	// Projected point calculation
	projectedPoint := point.Subv(normal.Muls(distance))

	fmt.Printf("Projected point: %v\n", projectedPoint)
	// Output: Projected point: (5, 0, 5)
}

func ExamplePoint_withinSphere() {
	// Checking if a point lies within a sphere
	// Center of the sphere
	center := r3.Point{X: 0, Y: 0, Z: 0}

	// Radius of the sphere
	radius := 5.0

	// Point to check
	point := r3.Point{X: 3, Y: 4, Z: 0}

	// Calculate the distance squared between the point and the center
	distanceSquared := point.Sub(center).Dot(point.Sub(center))

	// Compare with the radius squared
	if distanceSquared <= radius*radius {
		fmt.Printf("Point %v lies within the sphere.\n", point)
	} else {
		fmt.Printf("Point %v is outside the sphere.\n", point)
	}
	// Output: Point (3, 4, 0) lies within the sphere.
}

func ExamplePoint_shortestDistanceToLine() {
	// Finding the shortest distance between a point and a line
	// Line defined by points A and B
	A := r3.Point{X: 0, Y: 0, Z: 0}
	B := r3.Point{X: 10, Y: 0, Z: 0}

	// Point P
	P := r3.Point{X: 5, Y: 5, Z: 0}

	// Vector from A to B
	AB := B.Sub(A)

	// Vector from A to P
	AP := P.Sub(A)

	// Project AP onto AB
	t := AP.Dot(AB) / AB.Dot(AB)

	// Closest point on the line to P
	closestPoint := A.Add(AB.Muls(t))

	// Distance from P to the closest point
	distance := P.Sub(closestPoint).Length()

	fmt.Printf("The shortest distance from %v to the line is %v\n", P, distance)
	// Output: The shortest distance from (5, 5, 0) to the line is 5
}

func ExampleVec_rotateAroundAxis() {
	// Rotating a vector around an axis
	// Vector to rotate
	v := r3.Vec{X: 1, Y: 0, Z: 0}

	// Axis of rotation (must be a unit vector)
	axis := r3.Vec{X: 0, Y: 0, Z: 1}

	// Angle in degrees
	angleDegrees := 90.0
	angleRadians := angleDegrees * (math.Pi / 180)

	// Rodrigues' rotation formula
	cosTheta := math.Cos(angleRadians)
	sinTheta := math.Sin(angleRadians)

	vRotated := v.Muls(cosTheta).
		Add(axis.Cross(v).Muls(sinTheta)).
		Add(axis.Muls(axis.Dot(v) * (1 - cosTheta)))

	fmt.Printf("Rotated vector: %v\n", vRotated)
	// Output: Rotated vector: (6.123233995736757e-17, 1, 0)
}

func ExampleVec_scaleToLength() {
	// Scaling a vector to a desired length
	// Original vector
	v := r3.Vec{X: 3, Y: 4, Z: 0}

	// Desired length
	newLength := 10.0

	// Scale the vector
	vScaled := v.Unit().Muls(newLength)

	fmt.Printf("Scaled vector: %v\n", vScaled)
	// Output: Scaled vector: (6, 8, 0)
}

func TestVecSub(t *testing.T) {
	v1 := r3.Vec{1, 2, 3}
	v2 := r3.Vec{4, 5, 6}
	expected := r3.Vec{-3, -3, -3}
	result := v1.Sub(v2)
	if result != expected {
		t.Errorf("Sub: expected %v, got %v", expected, result)
	}
}

func TestVecMul(t *testing.T) {
	v1 := r3.Vec{1, 2, 3}
	v2 := r3.Vec{4, 5, 6}
	expected := r3.Vec{4, 10, 18}
	result := v1.Mul(v2)
	if result != expected {
		t.Errorf("Mul: expected %v, got %v", expected, result)
	}
}

func TestVecDiv(t *testing.T) {
	v1 := r3.Vec{4, 9, 16}
	v2 := r3.Vec{2, 3, 4}
	expected := r3.Vec{2, 3, 4}
	result := v1.Div(v2)
	if result != expected {
		t.Errorf("Div: expected %v, got %v", expected, result)
	}

	// Test division by zero
	vZero := r3.Vec{1, 1, 1}
	vDivZero := r3.Vec{0, 1, 1}
	result = vZero.Div(vDivZero)
	if !math.IsInf(result.X, 1) {
		t.Errorf("Div by zero: expected Inf, got %v", result.X)
	}
}

func TestVecMuls(t *testing.T) {
	v := r3.Vec{1, 2, 3}
	s := 2.0
	expected := r3.Vec{2, 4, 6}
	result := v.Muls(s)
	if result != expected {
		t.Errorf("Muls: expected %v, got %v", expected, result)
	}
}

func TestVecDivs(t *testing.T) {
	v := r3.Vec{2, 4, 6}
	s := 2.0
	expected := r3.Vec{1, 2, 3}
	result := v.Divs(s)
	if result != expected {
		t.Errorf("Divs: expected %v, got %v", expected, result)
	}

	// Test division by zero scalar
	result = v.Divs(0)
	if !math.IsInf(result.X, 1) || !math.IsInf(result.Y, 1) || !math.IsInf(result.Z, 1) {
		t.Errorf("Divs by zero: expected Inf, got %v", result)
	}
}

func TestVecDot(t *testing.T) {
	v1 := r3.Vec{1, 2, 3}
	v2 := r3.Vec{4, 5, 6}
	expected := 32.0
	result := v1.Dot(v2)
	if result != expected {
		t.Errorf("Dot: expected %v, got %v", expected, result)
	}
}

func TestVecCross(t *testing.T) {
	v1 := r3.Vec{1, 2, 3}
	v2 := r3.Vec{4, 5, 6}
	expected := r3.Vec{-3, 6, -3}
	result := v1.Cross(v2)
	if result != expected {
		t.Errorf("Cross: expected %v, got %v", expected, result)
	}
}

func TestVecLerp(t *testing.T) {
	v1 := r3.Vec{0, 0, 0}
	v2 := r3.Vec{10, 10, 10}
	tests := []struct {
		tParam   float64
		expected r3.Vec
	}{
		{0, r3.Vec{0, 0, 0}},
		{0.5, r3.Vec{5, 5, 5}},
		{1, r3.Vec{10, 10, 10}},
		{-0.5, r3.Vec{0, 0, 0}},   // Clamped to 0
		{1.5, r3.Vec{10, 10, 10}}, // Clamped to 1
	}

	for _, test := range tests {
		result := v1.Lerp(v2, test.tParam)
		if result != test.expected {
			t.Errorf("Lerp(%v): expected %v, got %v", test.tParam, test.expected, result)
		}
	}
}

func TestVecEq(t *testing.T) {
	v1 := r3.Vec{1, 2, 3}
	v2 := r3.Vec{1, 2, 3}
	v3 := r3.Vec{4, 5, 6}

	if !v1.Eq(v2) {
		t.Errorf("Eq: expected %v to equal %v", v1, v2)
	}
	if v1.Eq(v3) {
		t.Errorf("Eq: expected %v not to equal %v", v1, v3)
	}
}

func TestVecIsClose(t *testing.T) {
	v1 := r3.Vec{1.0000001, 2.0000001, 3.0000001}
	v2 := r3.Vec{1.0000002, 2.0000002, 3.0000002}
	v3 := r3.Vec{1.1, 2.1, 3.1}
	atol := 1e-6

	if !v1.IsClose(v2, atol) {
		t.Errorf("IsClose: expected %v to be close to %v within %v", v1, v2, atol)
	}
	if v1.IsClose(v3, atol) {
		t.Errorf("IsClose: expected %v not to be close to %v within %v", v1, v3, atol)
	}
}

func TestVecLength(t *testing.T) {
	v := r3.Vec{3, 4, 0}
	expected := 5.0
	result := v.Length()
	if result != expected {
		t.Errorf("Length: expected %v, got %v", expected, result)
	}
}

func TestVecUnit(t *testing.T) {
	v := r3.Vec{3, 4, 0}
	expected := r3.Vec{0.6, 0.8, 0}
	result := v.Unit()
	if !result.IsClose(expected, 1e-6) {
		t.Errorf("Unit: expected %v, got %v", expected, result)
	}

	// Test zero vector
	vZero := r3.Vec{0, 0, 0}
	expectedZero := r3.Vec{0, 0, 0}
	resultZero := vZero.Unit()
	if resultZero != expectedZero {
		t.Errorf("Unit of zero vector: expected %v, got %v", expectedZero, resultZero)
	}
}

func TestVecClip(t *testing.T) {
	v := r3.Vec{-2, 0, 2}
	min, max := -1.0, 1.0
	expected := r3.Vec{-1, 0, 1}
	result := v.Clip(min, max)
	if result != expected {
		t.Errorf("Clip: expected %v, got %v", expected, result)
	}
}

func TestVecIsNaN(t *testing.T) {
	vNaN := r3.Vec{math.NaN(), 0, 0}
	if !vNaN.IsNaN() {
		t.Errorf("IsNaN: expected %v to be NaN", vNaN)
	}

	vValid := r3.Vec{0, 0, 0}
	if vValid.IsNaN() {
		t.Errorf("IsNaN: expected %v not to be NaN", vValid)
	}
}

func TestVecIsInf(t *testing.T) {
	vInf := r3.Vec{math.Inf(1), 0, 0}
	if !vInf.IsInf() {
		t.Errorf("IsInf: expected %v to be Inf", vInf)
	}

	vValid := r3.Vec{0, 0, 0}
	if vValid.IsInf() {
		t.Errorf("IsInf: expected %v not to be Inf", vValid)
	}
}

func TestVecIsZero(t *testing.T) {
	vZero := r3.Vec{0, 0, 0}
	if !vZero.IsZero() {
		t.Errorf("IsZero: expected %v to be zero", vZero)
	}

	vNonZero := r3.Vec{1e-9, 0, 0}
	if vNonZero.IsZero() {
		t.Errorf("IsZero: expected %v not to be zero", vNonZero)
	}
}

func TestVecGet(t *testing.T) {
	v := r3.Vec{1, 2, 3}
	tests := []struct {
		index    int
		expected float64
	}{
		{0, 1},
		{1, 2},
		{2, 3},
	}

	for _, test := range tests {
		result := v.Get(test.index)
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
			v.Get(idx)
		}(index)
	}
}

func TestVecString(t *testing.T) {
	v := r3.Vec{1.1, 2.2, 3.3}
	expected := "(1.1, 2.2, 3.3)"
	result := v.String()
	if result != expected {
		t.Errorf("String: expected %v, got %v", expected, result)
	}
}
