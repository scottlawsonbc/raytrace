// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// MicrofacetBRDF represents the Cook-Torrance microfacet Bidirectional Reflectance Distribution Function.
// It models the reflection of light on rough surfaces by accounting for micro-scale surface variations.
// The key components of the model include the normal distribution function (D),
// the geometry (shadowing-masking) function (G), and the Fresnel term (F).
//
// Fields:
//   - Roughness: Controls the surface roughness. Typically ranges between 0 (perfectly smooth) and 1 (very rough).
//   - F0: The base reflectivity at normal incidence, represented as a r3.Vec. Each component usually ranges from 0 to 1.
type MicrofacetBRDF struct {
	Roughness float64 // If zero, will be clamped to a small epsilon min.
	F0        r3.Vec  // Base reflectivity at normal incidence.
}

// D calculates the Beckmann normal distribution function.
// It describes the distribution of microfacet normals on the surface.
//
// h is the half-vector between outgoing and incoming directions.
// n is the surface normal.
func (brdf *MicrofacetBRDF) D(h, n r3.Vec) float64 {
	// Clamp Roughness to prevent division by zero.
	roughness := math.Max(brdf.Roughness, eps)
	cosTheta := n.Dot(h)
	if cosTheta <= 0 {
		return 0
	}
	m2 := roughness * roughness
	cosTheta2 := cosTheta * cosTheta
	exponent := (cosTheta2 - 1) / (m2 * cosTheta2)
	return math.Exp(exponent) / (math.Pi * m2 * cosTheta2 * cosTheta2)
}

// G computes the geometry (shadowing-masking) function.
// It accounts for the shadowing and masking of microfacets.
//
// wo is the outgoing direction.
// wi is the incoming direction.
// n is the surface normal.
// h is the half-vector between wo and wi.
// Returns the value of the geometry function G.
func (brdf *MicrofacetBRDF) G(wo, wi, n, h r3.Vec) float64 {
	return brdf.G1(wo, n, h) * brdf.G1(wi, n, h)
}

// G1 computes the geometry function for a single direction.
// It estimates the shadowing-masking for one of the directions.
//
// v is the direction vector (either incoming or outgoing).
// n is the surface normal.
// h is the half-vector between the directions.
// Returns the value of the geometry function G1 for the given direction.
func (brdf *MicrofacetBRDF) G1(v, n, h r3.Vec) float64 {
	cosThetaV := math.Max(0, n.Dot(v))
	cosThetaH := math.Max(0, h.Dot(v))
	if cosThetaV <= 0 || cosThetaH <= 0 {
		return 0
	}
	tanThetaV := math.Sqrt(1-cosThetaV*cosThetaV) / cosThetaV
	a := 1 / (brdf.Roughness * tanThetaV)
	if a >= 1.6 {
		return 1
	}
	return (3.535*a + 2.181*a*a) / (1 + 2.276*a + 2.577*a*a)
}

// F computes the Fresnel term using Schlick's approximation.
// It models how light reflects off the surface based on the viewing angle.
//
// wo is the outgoing direction.
// h is the half-vector between wo and the incoming direction.
func (brdf *MicrofacetBRDF) F(wo, h r3.Vec) r3.Vec {
	cosTheta := math.Max(0, h.Dot(wo))
	oneMinusCosTheta5 := math.Pow(1-cosTheta, 5)
	return brdf.F0.Add(r3.Vec{X: 1, Y: 1, Z: 1}.Sub(brdf.F0).Muls(oneMinusCosTheta5))
}

// Evaluate computes the BRDF value given the outgoing and incoming directions.
// It combines the normal distribution, geometry, and Fresnel terms to produce the reflected radiance.
//
// wo is the outgoing direction.
// wi is the incoming direction.
// n is the surface normal.
func (brdf *MicrofacetBRDF) Evaluate(wo, wi, n r3.Vec) r3.Vec {
	wo = wo.Unit()
	wi = wi.Unit()
	n = n.Unit()
	h := wo.Add(wi).Unit()
	D := brdf.D(h, n)
	G := brdf.G(wo, wi, n, h)
	F := brdf.F(wo, h)
	// Prevent division by zero by adding a small epsilon
	denom := 4*math.Max(0, n.Dot(wo))*math.Max(0, n.Dot(wi)) + eps
	specular := F.Muls(D * G / denom)
	return specular
}
