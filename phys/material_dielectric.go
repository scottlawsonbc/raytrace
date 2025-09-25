// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type Dielectric struct {
	RefractiveIndexInterior float64
	RefractiveIndexExterior float64
	Roughness               float64 // Value between 0 (smooth) and 1 (rough). Usually 0.2 is really high already.
}

func (m Dielectric) Validate() error {
	if m.RefractiveIndexInterior < 1 || m.RefractiveIndexExterior < 1 {
		return fmt.Errorf("invalid Dielectric refractive index: %v", m)
	}
	if m.Roughness < 0 || m.Roughness > 1 {
		return fmt.Errorf("invalid Dielectric roughness: %v", m)
	}
	return nil
}

func (m Dielectric) Resolve(ctx context.Context, s surfaceInteraction) resolution {
	var outwardNormal r3.Vec
	var niOverNt float64
	var cosine float64
	var n1, n2 float64
	var rays []ray
	rand := s.incoming.rand

	// Determine if the ray is entering or exiting the material.
	if s.incoming.direction.Dot(s.collision.normal) > 0 {
		// Ray is exiting the material. Going from interior to exterior.
		outwardNormal = s.collision.normal.Muls(-1)
		niOverNt = m.RefractiveIndexInterior / m.RefractiveIndexExterior
		n1 = m.RefractiveIndexInterior
		n2 = m.RefractiveIndexExterior
		cosine = s.incoming.direction.Dot(s.collision.normal) / s.incoming.direction.Length()
		// Adjust cosine for total internal reflection
		cosine = math.Sqrt(1 - niOverNt*niOverNt*(1-cosine*cosine))
	} else {
		// Ray is entering the dielectric. Going from exterior to interior.
		outwardNormal = s.collision.normal
		niOverNt = m.RefractiveIndexExterior / m.RefractiveIndexInterior
		n1 = m.RefractiveIndexExterior
		n2 = m.RefractiveIndexInterior
		cosine = -s.incoming.direction.Dot(s.collision.normal) / s.incoming.direction.Length()
	}

	refracted, ok := refract(s.incoming.direction, outwardNormal, niOverNt)
	reflectProb := 1.0

	if ok {
		// Use Schlick's approximation for reflectance
		reflectProb = reflectance(cosine, n1, n2)

		// Add roughness to the refracted ray.
		if m.Roughness > 0 {
			refracted = refracted.Add(rand.InUnitSphere().Muls(m.Roughness)).Unit()
		}

		transmitted := ray{
			origin:    s.collision.at,
			direction: refracted,
			depth:     s.incoming.depth + 1,
			radiance:  s.incoming.radiance.Muls(1 - reflectProb),
			rand:      rand,
			pixelX:    s.incoming.pixelX,
			pixelY:    s.incoming.pixelY,
		}
		rays = append(rays, transmitted)
	}

	reflected := reflectRay(s.incoming.direction, s.collision.normal)

	// Add roughness to the reflected ray, scattering the direction slightly.
	if m.Roughness > 0 {
		reflected = reflected.Add(rand.InUnitSphere().Muls(m.Roughness)).Unit()
	}
	reflectedRay := ray{
		origin:    s.collision.at,
		direction: reflected,
		depth:     s.incoming.depth + 1,
		radiance:  s.incoming.radiance.Muls(reflectProb),
		rand:      rand,
		pixelX:    s.incoming.pixelX,
		pixelY:    s.incoming.pixelY,
	}
	rays = append(rays, reflectedRay)
	return resolution{scattered: rays}
}

// func (m Dielectric) ComputeDirectLighting(s surfaceInteraction, scene *Scene) r3.Vec {
// 	// Dielectrics (like glass) don't have direct diffuse contribution.
// 	// TODO: implement a small fuzz contribution
// 	return r3.Vec{}
// }

func (m Dielectric) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	p := s.collision.at
	n := s.collision.normal.Unit()
	wo := s.outgoing.Unit()
	directIllumination := Spectrum{}

	// Determine if the ray is entering or exiting the material.
	outside := wo.Dot(n) > 0
	etaI := m.RefractiveIndexExterior
	etaT := m.RefractiveIndexInterior
	normal := n
	if !outside {
		// Inside the material
		etaI, etaT = etaT, etaI
		normal = n.Muls(-1)
	}

	// Initialize the microfacet BRDF for reflection
	brdf := MicrofacetBRDF{
		Roughness: m.Roughness,
		F0:        r3.Vec{X: 1, Y: 1, Z: 1}, // Assuming dielectric with total internal reflection
	}

	for _, light := range scene.Light {
		dirToLight, distanceToLight, radiantIntensity := light.Sample(p, s.incoming.rand)
		wi := dirToLight.Unit()

		// Compute Fresnel term
		cosThetaI := math.Max(0, wi.Dot(normal))
		fresnel := reflectance(cosThetaI, etaI, etaT)

		// Offset the origin slightly to prevent self-intersection.
		shadowRayOrigin := p.Add(normal.Muls(eps))
		shadowRay := ray{
			origin:    shadowRayOrigin,
			direction: wi,
			depth:     s.incoming.depth + 1,
			radiance:  Spectrum{1, 1, 1},
			rand:      s.incoming.rand,
		}

		// Check for occlusion.
		occluded := false
		for _, node := range scene.Node {
			if node.Shape == s.node.Shape {
				continue // Skip self.
			}
			hit, _ := node.Shape.Collide(shadowRay, eps, distanceToLight)
			if hit {
				occluded = true
				break
			}
		}

		if !occluded {
			// Evaluate the BRDF
			brdfValue := brdf.Evaluate(wo, wi, normal).Muls(fresnel)

			// Compute the cosine term
			cosTheta := math.Max(0, normal.Dot(wi))

			// Accumulate the contribution.
			radiance := radiantIntensity.Mul(brdfValue).Muls(cosTheta)
			directIllumination = directIllumination.Add(Spectrum(radiance))
		}
	}
	return directIllumination
}

// reflectance computes the reflection coefficient using Schlick's approximation.
func reflectance(cosTheta, ni, nt float64) float64 {
	r0 := (ni - nt) / (ni + nt)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosTheta, 5)
}

// reflectRay computes the reflected ray direction given an incident vector and normal.
func reflectRay(v, n r3.Vec) r3.Vec {
	return v.Sub(n.Muls(2 * v.Dot(n))).Unit()
}

// refract computes the refracted ray direction given an incident vector, normal, and ratio of indices of refraction.
func refract(v, n r3.Vec, niOverNt float64) (r3.Vec, bool) {
	uv := v.Unit()
	dt := uv.Dot(n)
	discriminant := 1 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		refracted := uv.Sub(n.Muls(dt)).Muls(niOverNt).Sub(n.Muls(math.Sqrt(discriminant)))
		return refracted, true
	}
	return r3.Vec{}, false // Total internal reflection
}

func init() {
	RegisterInterfaceType(Dielectric{})
}
