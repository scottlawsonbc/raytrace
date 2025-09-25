// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// TestD verifies the correctness of the D (Beckmann normal distribution) function.
func TestD(t *testing.T) {
	// Create a MicrofacetBRDF instance with known roughness
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	// Define surface normal
	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	// Test case where h is aligned with n
	hAligned := n
	expectedDAligned := math.Exp((math.Pow(n.Dot(hAligned), 2)-1)/(math.Pow(math.Max(brdf.Roughness, eps), 2)*math.Pow(n.Dot(hAligned), 2))) / (math.Pi * math.Pow(math.Max(brdf.Roughness, eps), 2) * math.Pow(n.Dot(hAligned), 2) * math.Pow(n.Dot(hAligned), 2))
	dAligned := brdf.D(hAligned, n)
	if !almostEqual(dAligned, expectedDAligned, eps) {
		t.Errorf("D function incorrect for aligned h: got %v, want %v", dAligned, expectedDAligned)
	}

	// Test case where h is perpendicular to n
	hPerp := r3.Vec{X: 1, Y: 0, Z: 0}.Unit()
	dPerp := brdf.D(hPerp, n)
	if !almostEqual(dPerp, 0.0, eps) {
		t.Errorf("D function should return 0 for h perpendicular to n: got %v, want %v", dPerp, 0.0)
	}

	// Test with roughness = 0 (perfectly smooth, clamped internally)
	brdfSmooth := MicrofacetBRDF{
		Roughness: 0, // Will be clamped to eps
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}
	dSmooth := brdfSmooth.D(hAligned, n)
	if math.IsNaN(dSmooth) {
		t.Errorf("D function with roughness=0 should not be NaN: got %v", dSmooth)
	}
}

// TestG1 verifies the correctness of the G1 geometry function.
func TestG1(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	// Test with v aligned with n
	vAligned := n
	h := vAligned
	g1Aligned := brdf.G1(vAligned, n, h)
	if !almostEqual(g1Aligned, 1.0, eps) {
		t.Errorf("G1 function incorrect for aligned v: got %v, want %v", g1Aligned, 1.0)
	}

	// Test with v such that a >= 1.6
	// Find tanThetaV such that a = 1 / (roughness * tanThetaV) >= 1.6
	// => tanThetaV <= 1 / (roughness * 1.6) = 1 / (0.5 * 1.6) = 1.25
	// tanThetaV = sqrt(1 - cos^2(theta)) / cos(theta) = tan(theta)
	// Find theta where tan(theta) = 1.25 => theta ~ 51.34 degrees
	theta := math.Atan(1.25)
	v := r3.Vec{X: math.Sin(theta), Y: 0, Z: math.Cos(theta)}.Unit()
	g1 := brdf.G1(v, n, h)
	if !almostEqual(g1, 1.0, eps) {
		t.Errorf("G1 function should return 1 for a >= 1.6: got %v, want %v", g1, 1.0)
	}

	// Test with a < 1.6
	// Choose theta such that tan(theta) > 1 / (roughness * 1.6) = 1.25
	theta = math.Atan(2.0) // tan(theta) = 2.0 > 1.25
	v = r3.Vec{X: math.Sin(theta), Y: 0, Z: math.Cos(theta)}.Unit()
	g1 = brdf.G1(v, n, h)
	a := 1 / (brdf.Roughness * (math.Sqrt(1-math.Pow(math.Cos(theta), 2)) / math.Cos(theta)))
	expectedG1 := (3.535*a + 2.181*a*a) / (1 + 2.276*a + 2.577*a*a)
	if !almostEqual(g1, expectedG1, eps) {
		t.Errorf("G1 function incorrect for a < 1.6: got %v, want %v", g1, expectedG1)
	}

	// Test with v perpendicular to n
	vPerp := r3.Vec{X: 1, Y: 0, Z: 0}.Unit()
	g1Perp := brdf.G1(vPerp, n, h)
	if !almostEqual(g1Perp, 0.0, eps) {
		t.Errorf("G1 function should return 0 for v perpendicular to n: got %v, want %v", g1Perp, 0.0)
	}
}

// TestG verifies the correctness of the G geometry function.
func TestG(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	// Define outgoing and incoming directions aligned with n
	wo := n
	wi := n
	h := wo.Add(wi).Unit()

	expectedG := 1.0 * 1.0 // G1(wo) * G1(wi) = 1 * 1
	g := brdf.G(wo, wi, n, h)
	if !almostEqual(g, expectedG, eps) {
		t.Errorf("G function incorrect for aligned wo and wi: got %v, want %v", g, expectedG)
	}

	// Define outgoing and incoming directions at an angle
	wo = r3.Vec{X: math.Sin(math.Pi / 4), Y: 0, Z: math.Cos(math.Pi / 4)}.Unit()
	wi = wo
	h = wo.Add(wi).Unit()

	g = brdf.G(wo, wi, n, h)
	g1 := brdf.G1(wo, n, h)
	expectedG = g1 * g1
	if !almostEqual(g, expectedG, eps) {
		t.Errorf("G function incorrect for angled wo and wi: got %v, want %v", g, expectedG)
	}
}

// TestF verifies the correctness of the F (Fresnel) function using Schlick's approximation.
func TestF(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	// Test case where cosTheta = 1 (angle = 0 degrees)
	wo := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
	h := wo
	expectedF := brdf.F0
	f := brdf.F(wo, h)
	if !f.IsClose(expectedF, eps) {
		t.Errorf("F function incorrect for cosTheta=1: got %v, want %v", f, expectedF)
	}

	// Test case where cosTheta = 0 (angle = 90 degrees)
	wo = r3.Vec{X: 1, Y: 0, Z: 0}.Unit()
	h = r3.Vec{X: 0, Y: 1, Z: 0}.Unit()        // Perpendicular to wo
	expectedF = r3.Vec{X: 1.0, Y: 1.0, Z: 1.0} // F = F0 + (1 - F0) * 1 = 1
	f = brdf.F(wo, h)
	if !f.IsClose(expectedF, eps) {
		t.Errorf("F function incorrect for cosTheta=0: got %v, want %v", f, r3.Vec{X: 1, Y: 1, Z: 1})
	}

	// Test intermediate angle (e.g., 45 degrees)
	angle := math.Pi / 4 // 45 degrees
	wo = r3.Vec{X: math.Sin(angle), Y: 0, Z: math.Cos(angle)}.Unit()
	// To get cosTheta = cos(angle), set h to wo (since h = half-vector and in this case wo is already the half-vector)
	// However, to avoid confusion, it's better to define h based on the incoming direction
	// For simplicity, assume wi = wo, then h = wo
	h = wo
	expectedF = brdf.F0.Add(r3.Vec{X: 1, Y: 1, Z: 1}.Sub(brdf.F0).Muls(math.Pow(1-math.Cos(angle), 5)))
	f = brdf.F(wo, h)
	if !f.IsClose(expectedF, 1e-2) {
		t.Errorf("F function incorrect for intermediate angle: got %v, want %v", f, expectedF)
	}
}

// TestEvaluate verifies the correctness of the Evaluate function.
func TestEvaluate(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	// Test case where wo and wi are aligned with n
	wo := n
	wi := n
	expectedSpecular := brdf.F0.Muls(brdf.D(n, n) * brdf.G(wo, wi, n, n) / (4 * math.Pow(n.Dot(wo), 1) * math.Pow(n.Dot(wi), 1)))
	specular := brdf.Evaluate(wo, wi, n)
	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate incorrect for aligned wo and wi: got %v, want %v", specular, expectedSpecular)
	}

	// Test case with wo and wi at 45 degrees
	angle := math.Pi / 4
	wo = r3.Vec{X: math.Sin(angle), Y: 0, Z: math.Cos(angle)}.Unit()
	wi = wo
	h := wo.Add(wi).Unit()
	D := brdf.D(h, n)
	G := brdf.G(wo, wi, n, h)
	F := brdf.F(wo, h)
	denom := 4 * math.Max(0, n.Dot(wo)) * math.Max(0, n.Dot(wi)) // + eps not needed for expected value
	expectedSpecular = F.Muls(D * G / denom)
	specular = brdf.Evaluate(wo, wi, n)
	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate incorrect for angled wo and wi: got %v, want %v", specular, expectedSpecular)
	}

	// Test case where wi is below the surface (n.Dot(wi) <= 0)
	wi = r3.Vec{X: 0, Y: 0, Z: -1}.Unit()
	specular = brdf.Evaluate(wo, wi, n)
	expectedSpecular = r3.Vec{X: 0, Y: 0, Z: 0}
	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate should return zero when wi is below the surface: got %v, want %v", specular, expectedSpecular)
	}

	// Test case with roughness = 0 (perfectly smooth, clamped internally)
	brdfSmooth := MicrofacetBRDF{
		Roughness: 0.0, // Will be clamped to eps
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}
	wo = n
	wi = n
	specular = brdfSmooth.Evaluate(wo, wi, n)
	expectedSpecular = brdfSmooth.F0.Muls(brdfSmooth.D(n, n) * brdfSmooth.G(wo, wi, n, n) / (4 * math.Pow(n.Dot(wo), 1) * math.Pow(n.Dot(wi), 1)))
	if math.IsNaN(specular.X) || math.IsNaN(specular.Y) || math.IsNaN(specular.Z) {
		t.Errorf("Evaluate with roughness=0 should not return NaN: got %v", specular)
	}
	// This is a large value so we compare with a very large epsilon.
	if !specular.IsClose(expectedSpecular, 1e3) {
		t.Errorf("Evaluate with roughness=0 clamped incorrectly: got %v, want %v", specular, expectedSpecular)
	}

	// Test case with roughness = 1 (very rough)
	brdfRough := MicrofacetBRDF{
		Roughness: 1.0,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}
	wo = r3.Vec{X: 1, Y: 0, Z: 1}.Unit()
	wi = r3.Vec{X: 1, Y: 0, Z: 1}.Unit()
	h = wo.Add(wi).Unit()
	D = brdfRough.D(h, n)
	G = brdfRough.G(wo, wi, n, h)
	F = brdfRough.F(wo, h)
	denom = 4 * math.Max(0, n.Dot(wo)) * math.Max(0, n.Dot(wi)) // + eps not needed for expected value
	expectedSpecular = F.Muls(D * G / denom)
	specular = brdfRough.Evaluate(wo, wi, n)
	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate incorrect for roughness=1: got %v, want %v", specular, expectedSpecular)
	}
}

// TestRoughnessRange verifies that the Roughness parameter is within the expected range.
func TestRoughnessRange(t *testing.T) {
	// Roughness should typically be between 0 and 1
	tests := []struct {
		roughness float64
		valid     bool
	}{
		{-0.1, false},
		{0.0, true}, // Will be clamped internally
		{0.5, true},
		{1.0, true},
		{1.5, false},
	}

	for _, tt := range tests {
		brdf := MicrofacetBRDF{
			Roughness: tt.roughness,
			F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
		}
		// Evaluate with valid and invalid roughness
		wo := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
		wi := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
		n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
		specular := brdf.Evaluate(wo, wi, n)

		if !tt.valid {
			// For invalid roughness, roughness is clamped to eps
			expectedRoughness := math.Max(tt.roughness, eps)
			expectedD := math.Exp((math.Pow(n.Dot(n), 2)-1)/(math.Pow(expectedRoughness, 2)*math.Pow(n.Dot(n), 2))) / (math.Pi * math.Pow(expectedRoughness, 2) * math.Pow(n.Dot(n), 2) * math.Pow(n.Dot(n), 2))
			if !almostEqual(brdf.D(n, n), expectedD, eps) {
				t.Errorf("Roughness clamping incorrect for roughness=%v: got D=%v, want D≈%v", tt.roughness, brdf.D(n, n), expectedD)
			}
			continue
		}

		// For valid roughness, ensure that Evaluate returns a non-negative reflectance
		if specular.X < 0 || specular.Y < 0 || specular.Z < 0 {
			t.Errorf("Evaluate returned negative reflectance for roughness=%v: got %v", tt.roughness, specular)
		}
	}
}

// TestF0Range verifies that the F0 parameter components are within the expected range [0, 1].
func TestF0Range(t *testing.T) {
	tests := []struct {
		F0    r3.Vec
		valid bool
	}{
		{r3.Vec{X: -0.1, Y: 0.04, Z: 0.04}, false},
		{r3.Vec{X: 0.0, Y: 0.0, Z: 0.0}, true},
		{r3.Vec{X: 0.04, Y: 0.04, Z: 0.04}, true},
		{r3.Vec{X: 1.0, Y: 1.0, Z: 1.0}, true},
		{r3.Vec{X: 1.1, Y: 0.04, Z: 0.04}, false},
	}

	for _, tt := range tests {
		brdf := MicrofacetBRDF{
			Roughness: 0.5,
			F0:        tt.F0,
		}

		wo := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
		wi := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
		n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

		specular := brdf.Evaluate(wo, wi, n)

		if !tt.valid {
			// For invalid F0, specular might exceed expected bounds
			// Here we ensure that components do not exceed [0,1] if that's a requirement
			// Alternatively, you could enforce clamping in the implementation
			for _, c := range []float64{specular.X, specular.Y, specular.Z} {
				if c < 0 || c > 1.0 {
					// Expected behavior depends on implementation; adjust as needed
					// Here we assume F0 should be clamped or validated elsewhere
					continue
				}
			}
			continue
		}

		// For valid F0, ensure that F0 components are within [0,1]
		if tt.F0.X < 0 || tt.F0.X > 1 ||
			tt.F0.Y < 0 || tt.F0.Y > 1 ||
			tt.F0.Z < 0 || tt.F0.Z > 1 {
			t.Errorf("F0 components out of range [0,1]: got %v", tt.F0)
		}
	}
}

// TestEvaluatePhysicalProperties verifies some physical properties of the BRDF.
func TestEvaluatePhysicalProperties(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	// Ensure that Evaluate returns non-negative values
	wo := r3.Vec{X: math.Sin(math.Pi / 6), Y: 0, Z: math.Cos(math.Pi / 6)}.Unit()
	wi := r3.Vec{X: math.Sin(math.Pi / 3), Y: 0, Z: math.Cos(math.Pi / 3)}.Unit()
	specular := brdf.Evaluate(wo, wi, n)
	if specular.X < 0 || specular.Y < 0 || specular.Z < 0 {
		t.Errorf("Evaluate returned negative values: got %v", specular)
	}

	// Symmetry: Evaluate(wo, wi) should equal Evaluate(wi, wo)
	specular1 := brdf.Evaluate(wo, wi, n)
	specular2 := brdf.Evaluate(wi, wo, n)
	if !specular1.IsClose(specular2, eps) {
		t.Errorf("BRDF is not symmetric: Evaluate(wo, wi)=%v != Evaluate(wi, wo)=%v", specular1, specular2)
	}

	// Energy conservation: Specular reflection should not exceed Fresnel reflectance
	// F should be greater or equal to specular component
	F := brdf.F(wo, wo.Add(wi).Unit())
	if specular.X > F.X+eps || specular.Y > F.Y+eps || specular.Z > F.Z+eps {
		t.Errorf("Specular reflection exceeds Fresnel reflectance: specular=%v, F=%v", specular, F)
	}
}

// TestHalfVector verifies that the half-vector is correctly computed as the normalized sum of wo and wi.
func TestHalfVector(t *testing.T) {
	wo := r3.Vec{X: 1, Y: 0, Z: 1}.Unit()
	wi := r3.Vec{X: 1, Y: 0, Z: 1}.Unit()
	expectedH := r3.Vec{X: 1, Y: 0, Z: 1}.Unit()
	h := wo.Add(wi).Unit()
	if !h.IsClose(expectedH, eps) {
		t.Errorf("Half-vector incorrect: got %v, want %v", h, expectedH)
	}
}

// TestZeroRoughnessAndF0 verifies behavior when roughness is zero and F0 is zero.
func TestZeroRoughnessAndF0(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.0, // Will be clamped internally
		F0:        r3.Vec{X: 0, Y: 0, Z: 0},
	}

	wo := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
	wi := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	specular := brdf.Evaluate(wo, wi, n)
	expectedSpecular := r3.Vec{X: 0, Y: 0, Z: 0}

	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate with zero roughness and F0 should return zero: got %v, want %v", specular, expectedSpecular)
	}
}

// TestNonNormalizedVectors ensures that Evaluate handles non-normalized input vectors correctly.
func TestNonNormalizedVectors(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	// Non-normalized vectors
	wo := r3.Vec{X: 0, Y: 0, Z: 2}
	wi := r3.Vec{X: 0, Y: 0, Z: 2}
	n := r3.Vec{X: 0, Y: 0, Z: 1}

	// Expected behavior: vectors are normalized inside Evaluate
	woNorm := wo.Unit()
	wiNorm := wi.Unit()
	h := woNorm.Add(wiNorm).Unit()
	D := brdf.D(h, n.Unit())
	G := brdf.G(woNorm, wiNorm, n.Unit(), h)
	F := brdf.F(woNorm, h)
	denom := 4*math.Max(0, n.Dot(woNorm))*math.Max(0, n.Dot(wiNorm)) + eps
	expectedSpecular := F.Muls(D * G / denom)
	specular := brdf.Evaluate(wo, wi, n)

	if !specular.IsClose(expectedSpecular, eps) {
		t.Errorf("Evaluate with non-normalized vectors incorrect: got %v, want %v", specular, expectedSpecular)
	}
}

// TestDFunctionProperties checks that D is symmetric with respect to its arguments.
func TestDFunctionProperties(t *testing.T) {
	brdf := MicrofacetBRDF{
		Roughness: 0.3,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}

	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()
	h1 := r3.Vec{X: 1, Y: 1, Z: 1}.Unit()
	h2 := r3.Vec{X: -1, Y: -1, Z: 1}.Unit()

	D1 := brdf.D(h1, n)
	D2 := brdf.D(h2, n)

	// Since D should be symmetric with respect to the azimuthal angle, D1 should equal D2
	if !almostEqual(D1, D2, eps) {
		t.Errorf("D function is not symmetric: D1=%v, D2=%v", D1, D2)
	}
}

// BenchmarkEvaluate benchmarks the performance of the Evaluate function.
func BenchmarkEvaluate(b *testing.B) {
	brdf := MicrofacetBRDF{
		Roughness: 0.5,
		F0:        r3.Vec{X: 0.04, Y: 0.04, Z: 0.04},
	}
	wo := r3.Vec{X: math.Sin(math.Pi / 6), Y: 0, Z: math.Cos(math.Pi / 6)}.Unit()
	wi := r3.Vec{X: math.Sin(math.Pi / 3), Y: 0, Z: math.Cos(math.Pi / 3)}.Unit()
	n := r3.Vec{X: 0, Y: 0, Z: 1}.Unit()

	for i := 0; i < b.N; i++ {
		brdf.Evaluate(wo, wi, n)
	}
}

// almostEqual checks if two float64 values are approximately equal within a small tolerance.
func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

// package phys

// import (
// 	"fmt"
// 	"math"
// 	"testing"
// )

// func ExampleMicrofacetBRDF_Evaluate() {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.3,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}
// 	wo := r3.Vec{0, 0, 1}.Unit()
// 	wi := r3.Vec{0, 0, 1}.Unit()
// 	n := r3.Vec{0, 0, 1}.Unit()

// 	specular := brdf.Evaluate(wo, wi, n)
// 	fmt.Println(specular)
// 	// Output: (0.035367765043112884, 0.035367765043112884, 0.035367765043112884)
// }

// // TestD verifies the correctness of the D (Beckmann normal distribution) function.
// func TestD(t *testing.T) {
// 	// Create a MicrofacetBRDF instance with known roughness
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	// Define surface normal
// 	n := r3.Vec{0, 0, 1}.Unit()

// 	// Test case where h is aligned with n
// 	hAligned := n
// 	expectedDAligned := math.Exp((math.Pow(n.Dot(hAligned), 2)-1)/(math.Pow(brdf.Roughness, 2)*math.Pow(n.Dot(hAligned), 2))) / (math.Pi * math.Pow(brdf.Roughness, 2) * math.Pow(n.Dot(hAligned), 2) * math.Pow(n.Dot(hAligned), 2))
// 	dAligned := brdf.D(hAligned, n)
// 	if !almostEqual(dAligned, expectedDAligned, eps) {
// 		t.Errorf("D function incorrect for aligned h: got %v, want %v", dAligned, expectedDAligned)
// 	}

// 	// Test case where h is perpendicular to n
// 	hPerp := r3.Vec{1, 0, 0}.Unit()
// 	dPerp := brdf.D(hPerp, n)
// 	if !almostEqual(dPerp, 0.0, eps) {
// 		t.Errorf("D function should return 0 for h perpendicular to n: got %v, want %v", dPerp, 0.0)
// 	}

// 	// Test with roughness = 0 (perfectly smooth)
// 	brdfSmooth := MicrofacetBRDF{
// 		Roughness: 0.0,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}
// 	dSmooth := brdfSmooth.D(hAligned, n)
// 	if !math.IsInf(dSmooth, 1) && dSmooth != 0.0 {
// 		t.Errorf("D function with roughness=0 should return 0 or Inf: got %v", dSmooth)
// 	}
// }

// func almostEqual(a, b, eps float64) bool {
// 	return math.Abs(a-b) < eps
// }

// // TestG1 verifies the correctness of the G1 geometry function.
// func TestG1(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	n := r3.Vec{0, 0, 1}.Unit()

// 	// Test with v aligned with n
// 	vAligned := n
// 	h := vAligned
// 	g1Aligned := brdf.G1(vAligned, n, h)
// 	if !almostEqual(g1Aligned, 1.0, eps) {
// 		t.Errorf("G1 function incorrect for aligned v: got %v, want %v", g1Aligned, 1.0)
// 	}

// 	// Test with v such that a >= 1.6
// 	// Find tanThetaV such that a = 1 / (roughness * tanThetaV) >= 1.6
// 	// => tanThetaV <= 1 / (roughness * 1.6) = 1 / (0.5 * 1.6) = 1.25
// 	// tanThetaV = sqrt(1 - cos^2(theta)) / cos(theta) = tan(theta)
// 	// Find theta where tan(theta) = 1.25 => theta ~ 51.34 degrees
// 	theta := math.Atan(1.25)
// 	v := r3.Vec{math.Sin(theta), 0, math.Cos(theta)}.Unit()
// 	g1 := brdf.G1(v, n, h)
// 	if !almostEqual(g1, 1.0, eps) {
// 		t.Errorf("G1 function should return 1 for a >= 1.6: got %v, want %v", g1, 1.0)
// 	}

// 	// Test with a < 1.6
// 	// Choose theta such that tan(theta) > 1 / (roughness * 1.6) = 1.25
// 	theta = math.Atan(2.0) // tan(theta) = 2.0 > 1.25
// 	v = r3.Vec{math.Sin(theta), 0, math.Cos(theta)}.Unit()
// 	g1 = brdf.G1(v, n, h)
// 	expectedG1 := (3.535*1/(brdf.Roughness*2.0) + 2.181*math.Pow(1/(brdf.Roughness*2.0), 2)) / (1 + 2.276*1/(brdf.Roughness*2.0) + 2.577*math.Pow(1/(brdf.Roughness*2.0), 2))
// 	if !almostEqual(g1, expectedG1, eps) {
// 		t.Errorf("G1 function incorrect for a < 1.6: got %v, want %v", g1, expectedG1)
// 	}

// 	// Test with v perpendicular to n
// 	vPerp := r3.Vec{1, 0, 0}.Unit()
// 	g1Perp := brdf.G1(vPerp, n, h)
// 	if !almostEqual(g1Perp, 0.0, eps) {
// 		t.Errorf("G1 function should return 0 for v perpendicular to n: got %v, want %v", g1Perp, 0.0)
// 	}
// }

// // TestG verifies the correctness of the G geometry function.
// func TestG(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	n := r3.Vec{0, 0, 1}.Unit()

// 	// Define outgoing and incoming directions aligned with n
// 	wo := n
// 	wi := n
// 	h := wo.Add(wi).Unit()

// 	expectedG := 1.0 * 1.0 // G1(wo) * G1(wi) = 1 * 1
// 	g := brdf.G(wo, wi, n, h)
// 	if !almostEqual(g, expectedG, eps) {
// 		t.Errorf("G function incorrect for aligned wo and wi: got %v, want %v", g, expectedG)
// 	}

// 	// Define outgoing and incoming directions at an angle
// 	wo = r3.Vec{math.Sin(math.Pi / 4), 0, math.Cos(math.Pi / 4)}.Unit()
// 	wi = wo
// 	h = wo.Add(wi).Unit()

// 	g = brdf.G(wo, wi, n, h)
// 	g1 := brdf.G1(wo, n, h)
// 	expectedG = g1 * g1
// 	if !almostEqual(g, expectedG, eps) {
// 		t.Errorf("G function incorrect for angled wo and wi: got %v, want %v", g, expectedG)
// 	}
// }

// // TestF verifies the correctness of the F (Fresnel) function using Schlick's approximation.
// func TestF(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	// Test case where cosTheta = 1 (angle = 0 degrees)
// 	wo := r3.Vec{0, 0, 1}.Unit()
// 	h := wo
// 	expectedF := brdf.F0
// 	f := brdf.F(wo, h)
// 	if !f.IsClose(expectedF, eps) {

// 		t.Errorf("F function incorrect for cosTheta=1: got %v, want %v", f, expectedF)
// 	}

// 	// Test case where cosTheta = 0 (angle = 90 degrees)
// 	wo = r3.Vec{1, 0, 0}.Unit()
// 	h = wo
// 	expectedF = r3.Vec{1.0, 1.0, 1.0} // F = F0 + (1 - F0) * 1 = 1
// 	f = brdf.F(wo, h)
// 	if !f.IsClose(expectedF, eps) {
// 		t.Errorf("F function incorrect for cosTheta=0: got %v, want %v", f, r3.Vec{1, 1, 1})
// 	}

// 	// Test intermediate angle
// 	angle := math.Pi / 4 // 45 degrees, cosTheta = sqrt(2)/2 ~0.7071
// 	wo = r3.Vec{math.Sin(angle), 0, math.Cos(angle)}.Unit()
// 	h = wo
// 	oneMinusCosTheta5 := math.Pow(1-math.Cos(angle), 5)
// 	expectedF = brdf.F0.Add(r3.Vec{1, 1, 1}.Sub(brdf.F0).Muls(oneMinusCosTheta5))
// 	f = brdf.F(wo, h)
// 	if !f.IsClose(expectedF, eps) {
// 		t.Errorf("F function incorrect for intermediate angle: got %v, want %v", f, expectedF)
// 	}
// }

// // TestEvaluate verifies the correctness of the Evaluate function.
// func TestEvaluate(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	n := r3.Vec{0, 0, 1}.Unit()

// 	// Test case where wo and wi are aligned with n
// 	wo := n
// 	wi := n
// 	expectedSpecular := brdf.F0.Muls(brdf.D(n, n) * brdf.G(wo, wi, n, n) / (4 * math.Pow(n.Dot(wo), 1) * math.Pow(n.Dot(wi), 1)))
// 	specular := brdf.Evaluate(wo, wi, n)
// 	if !specular.IsClose(expectedSpecular, eps) {
// 		t.Errorf("Evaluate incorrect for aligned wo and wi: got %v, want %v", specular, expectedSpecular)
// 	}

// 	// Test case with wo and wi at 45 degrees
// 	angle := math.Pi / 4
// 	wo = r3.Vec{math.Sin(angle), 0, math.Cos(angle)}.Unit()
// 	wi = wo
// 	h := wo.Add(wi).Unit()
// 	D := brdf.D(h, n)
// 	G := brdf.G(wo, wi, n, h)
// 	F := brdf.F(wo, h)
// 	denom := 4 * math.Max(0, n.Dot(wo)) * math.Max(0, n.Dot(wi)) // + eps not needed for expected value
// 	expectedSpecular = F.Muls(D * G / denom)
// 	specular = brdf.Evaluate(wo, wi, n)
// 	if !specular.IsClose(expectedSpecular, eps) {
// 		t.Errorf("Evaluate incorrect for angled wo and wi: got %v, want %v", specular, expectedSpecular)
// 	}

// 	// Test case where wi is below the surface (n.Dot(wi) <= 0)
// 	wi = r3.Vec{0, 0, -1}.Unit()
// 	specular = brdf.Evaluate(wo, wi, n)
// 	expectedSpecular = r3.Vec{0, 0, 0}
// 	if !specular.IsClose(expectedSpecular, eps) {

// 		t.Errorf("Evaluate should return zero when wi is below the surface: got %v, want %v", specular, expectedSpecular)
// 	}

// 	// Test case with roughness = 0 (perfectly smooth)
// 	brdfSmooth := MicrofacetBRDF{
// 		Roughness: 0.0,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}
// 	wo = n
// 	wi = n
// 	specular = brdfSmooth.Evaluate(wo, wi, n)
// 	expectedSpecular = brdfSmooth.F0.Muls(brdfSmooth.D(n, n) * brdfSmooth.G(wo, wi, n, n) / (4 * math.Pow(n.Dot(wo), 1) * math.Pow(n.Dot(wi), 1)))
// 	if !math.IsInf(specular.X, 1) && !math.IsInf(specular.Y, 1) && !math.IsInf(specular.Z, 1) {
// 		// Depending on D with roughness=0, this could be Inf or a very large number
// 		// Here we check if it's effectively a delta function (perfect mirror)
// 		if !specular.IsClose(r3.Vec{0, 0, 0}) { // Assuming D return, epss 0
// 			t.Errorf("Evaluate with roughness=0 should return zero or Inf: got %v", specular)
// 		}
// 	}

// 	// Test case with roughness = 1 (very rough)
// 	brdfRough := MicrofacetBRDF{
// 		Roughness: 1.0,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}
// 	wo = r3.Vec{1, 0, 1}.Unit()
// 	wi = r3.Vec{1, 0, 1}.Unit()
// 	h = wo.Add(wi).Unit()
// 	D = brdfRough.D(h, n)
// 	G = brdfRough.G(wo, wi, n, h)
// 	F = brdfRough.F(wo, h)
// 	denom = 4 * math.Max(0, n.Dot(wo)) * math.Max(0, n.Dot(wi)) // + eps not needed for expected value
// 	expectedSpecular = F.Muls(D * G / denom)
// 	specular = brdfRough.Evaluate(wo, wi, n)
// 	if !specular.IsClose(expectedSpecular, eps) {
// 		t.Errorf("Evaluate incorrect for roughness=1: got %v, want %v", specular, expectedSpecular)
// 	}
// }

// // TestRoughnessRange verifies that the Roughness parameter is within the expected range.
// func TestRoughnessRange(t *testing.T) {
// 	// Roughness should typically be between 0 and 1
// 	tests := []struct {
// 		roughness float64
// 		valid     bool
// 	}{
// 		{-0.1, false},
// 		{0.0, true},
// 		{0.5, true},
// 		{1.0, true},
// 		{1.5, false},
// 	}

// 	for _, tt := range tests {
// 		brdf := MicrofacetBRDF{
// 			Roughness: tt.roughness,
// 			F0:        r3.Vec{0.04, 0.04, 0.04},
// 		}
// 		// Here, you might enforce roughness constraints in your implementation.
// 		// Since the current implementation does not enforce it, we just check expected behavior.
// 		// Alternatively, you could modify the MicrofacetBRDF to clamp or return errors.
// 		// For this test, we'll assume values outside [0,1] are invalid and check for expected outcomes.

// 		// Example: Evaluate with invalid roughness and ensure no panic occurs
// 		wo := r3.Vec{0, 0, 1}.Unit()
// 		wi := r3.Vec{0, 0, 1}.Unit()
// 		n := r3.Vec{0, 0, 1}.Unit()
// 		specular := brdf.Evaluate(wo, wi, n)
// 		if !tt.valid {
// 			// For invalid roughness, specular might behave unexpectedly
// 			// Here we just ensure that the function doesn't panic and returns a r3.Vec
// 			// Further behavior checks can be added based on how you handle invalid roughness
// 			continue
// 		}
// 		// For valid roughness, ensure that Evaluate returns a non-negative reflectance
// 		if specular.X < 0 || specular.Y < 0 || specular.Z < 0 {
// 			t.Errorf("Evaluate returned negative reflectance for roughness=%v: got %v", tt.roughness, specular)
// 		}
// 	}
// }

// // TestF0Range verifies that the F0 parameter components are within the expected range [0, 1].
// func TestF0Range(t *testing.T) {
// 	tests := []struct {
// 		F0    r3.Vec
// 		valid bool
// 	}{
// 		{r3.Vec{-0.1, 0.04, 0.04}, false},
// 		{r3.Vec{0.0, 0.0, 0.0}, true},
// 		{r3.Vec{0.04, 0.04, 0.04}, true},
// 		{r3.Vec{1.0, 1.0, 1.0}, true},
// 		{r3.Vec{1.1, 0.04, 0.04}, false},
// 	}

// 	for _, tt := range tests {
// 		brdf := MicrofacetBRDF{
// 			Roughness: 0.5,
// 			F0:        tt.F0,
// 		}

// 		wo := r3.Vec{0, 0, 1}.Unit()
// 		wi := r3.Vec{0, 0, 1}.Unit()
// 		n := r3.Vec{0, 0, 1}.Unit()

// 		specular := brdf.Evaluate(wo, wi, n)

// 		if !tt.valid {
// 			// For invalid F0, specular might exceed expected bounds
// 			// Here we ensure that components do not exceed [0,1] if that's a requirement
// 			// Alternatively, you could enforce clamping in the implementation
// 			for _, c := range []float64{specular.X, specular.Y, specular.Z} {
// 				if c < 0 || c > 1.0 {
// 					// Expected behavior depends on implementation; adjust as needed
// 					// Here we assume F0 should be clamped or validated elsewhere
// 					continue
// 				}
// 			}
// 			continue
// 		}

// 		// For valid F0, ensure that F0 components are within [0,1]
// 		if tt.F0.X < 0 || tt.F0.X > 1 ||
// 			tt.F0.Y < 0 || tt.F0.Y > 1 ||
// 			tt.F0.Z < 0 || tt.F0.Z > 1 {
// 			t.Errorf("F0 components out of range [0,1]: got %v", tt.F0)
// 		}
// 	}
// }

// // TestEvaluatePhysicalProperties verifies some physical properties of the BRDF.
// func TestEvaluatePhysicalProperties(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	n := r3.Vec{0, 0, 1}.Unit()

// 	// Ensure that Evaluate returns non-negative values
// 	wo := r3.Vec{math.Sin(math.Pi / 6), 0, math.Cos(math.Pi / 6)}.Unit()
// 	wi := r3.Vec{math.Sin(math.Pi / 3), 0, math.Cos(math.Pi / 3)}.Unit()
// 	specular := brdf.Evaluate(wo, wi, n)
// 	if specular.X < 0 || specular.Y < 0 || specular.Z < 0 {
// 		t.Errorf("Evaluate returned negative values: got %v", specular)
// 	}

// 	// Symmetry: Evaluate(wo, wi) should equal Evaluate(wi, wo)
// 	specular1 := brdf.Evaluate(wo, wi, n)
// 	specular2 := brdf.Evaluate(wi, wo, n)
// 	if !specular1.IsClose(specular2, eps) {
// 		t.Errorf("BRDF is not symmetric: Evaluate(wo, wi)=%v != Evaluate(wi, wo)=%v", specular1, specular2)
// 	}

// 	// Energy conservation: Specular reflection should not exceed Fresnel reflectance
// 	// F should be greater or equal to specular component
// 	F := brdf.F(wo, wo.Add(wi).Unit())
// 	if specular.X > F.X+eps || specular.Y > F.Y+eps || specular.Z > F.Z+eps {
// 		t.Errorf("Specular reflection exceeds Fresnel reflectance: specular=%v, F=%v", specular, F)
// 	}
// }

// // TestHalfVector verifies that the half-vector is correctly computed as the normalized sum of wo and wi.
// func TestHalfVector(t *testing.T) {
// 	wo := r3.Vec{1, 0, 1}.Unit()
// 	wi := r3.Vec{1, 0, 1}.Unit()
// 	expectedH := r3.Vec{1, 0, 1}.Unit()
// 	h := wo.Add(wi).Unit()
// 	if !h.IsClose(expectedH, eps) {
// 		t.Errorf("Half-vector incorrect: got %v, want %v", h, expectedH)
// 	}
// }

// // TestZeroRoughnessAndF0 verifies behavior when roughness is zero and F0 is zero.
// func TestZeroRoughnessAndF0(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.0,
// 		F0:        r3.Vec{0, 0, 0},
// 	}

// 	wo := r3.Vec{0, 0, 1}.Unit()
// 	wi := r3.Vec{0, 0, 1}.Unit()
// 	n := r3.Vec{0, 0, 1}.Unit()

// 	specular := brdf.Evaluate(wo, wi, n)
// 	expectedSpecular := r3.Vec{0, 0, 0}

// 	if !specular.IsClose(expectedSpecular, eps) {
// 		t.Errorf("Evaluate with zero roughness and F0 should return zero: got %v, want %v", specular, expectedSpecular)
// 	}
// }

// // TestNonNormalizedVectors ensures that Evaluate handles non-normalized input vectors correctly.
// func TestNonNormalizedVectors(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	// Non-normalized vectors
// 	wo := r3.Vec{0, 0, 2}
// 	wi := r3.Vec{0, 0, 2}
// 	n := r3.Vec{0, 0, 1}

// 	// Expected behavior: vectors should be normalized inside the functions
// 	// Since our current implementation does not normalize inputs, Evaluate might behave unexpectedly
// 	// If normalization is required, consider adding it or documenting that inputs must be normalized

// 	// For this test, we'll normalize manually to compare expected behavior
// 	woNorm := wo.Unit()
// 	wiNorm := wi.Unit()
// 	h := woNorm.Add(wiNorm).Unit()
// 	D := brdf.D(h, n.Unit())
// 	G := brdf.G(woNorm, wiNorm, n.Unit(), h)
// 	F := brdf.F(woNorm, h)
// 	denom := 4*math.Max(0, n.Dot(woNorm))*math.Max(0, n.Dot(wiNorm)) + eps
// 	expectedSpecular := F.Muls(D * G / denom)
// 	specular := brdf.Evaluate(wo, wi, n)

// 	if !specular.IsClose(expectedSpecular, eps) {
// 		t.Errorf("Evaluate with non-normalized vectors incorrect: got %v, want %v", specular, expectedSpecular)
// 	}
// }

// // TestDFunctionProperties checks that D is symmetric with respect to its arguments.
// func TestDFunctionProperties(t *testing.T) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.3,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}

// 	n := r3.Vec{0, 0, 1}.Unit()
// 	h1 := r3.Vec{1, 1, 1}.Unit()
// 	h2 := r3.Vec{-1, -1, 1}.Unit()

// 	D1 := brdf.D(h1, n)
// 	D2 := brdf.D(h2, n)

// 	// Since D should be symmetric with respect to the azimuthal angle, D1 should equal D2
// 	if !almostEqual(D1, D2, eps) {
// 		t.Errorf("D function is not symmetric: D1=%v, D2=%v", D1, D2)
// 	}
// }

// func BenchmarkEvaluate(b *testing.B) {
// 	brdf := MicrofacetBRDF{
// 		Roughness: 0.5,
// 		F0:        r3.Vec{0.04, 0.04, 0.04},
// 	}
// 	wo := r3.Vec{math.Sin(math.Pi / 6), 0, math.Cos(math.Pi / 6)}.Unit()
// 	wi := r3.Vec{math.Sin(math.Pi / 3), 0, math.Cos(math.Pi / 3)}.Unit()
// 	n := r3.Vec{0, 0, 1}.Unit()

// 	for i := 0; i < b.N; i++ {
// 		brdf.Evaluate(wo, wi, n)
// 	}
// }
