// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type Metal struct {
	Albedo r3.Vec
	Fuzz   float64
}

func (m Metal) Validate() error {
	if m.Albedo.X < 0 || m.Albedo.Y < 0 || m.Albedo.Z < 0 {
		return fmt.Errorf("invalid Metal albedo must be positive")
	}
	if m.Fuzz < 0 || m.Fuzz > 1 {
		return fmt.Errorf("invalid Metal fuzz must be in the range [0, 1]")
	}
	return nil
}

func (m Metal) Resolve(ctx context.Context, s surfaceInteraction) resolution {
	// TODO: scott should this actually return a resolution with multiple rays?
	reflected := reflectRay(s.incoming.direction.Unit(), s.collision.normal)
	scatteredDirection := reflected.Add(s.incoming.rand.InUnitSphere().Muls(m.Fuzz))
	if scatteredDirection.Dot(s.collision.normal) > 0 {
		newRay := ray{
			origin:    s.collision.at,
			direction: scatteredDirection.Unit(),
			depth:     s.incoming.depth + 1,
			radiance:  s.incoming.radiance.Mul(Spectrum(m.Albedo)),
			rand:      s.incoming.rand,
			pixelX:    s.incoming.pixelX,
			pixelY:    s.incoming.pixelY,
		}
		return resolution{scattered: []ray{newRay}}
	}
	// Absorb the ray (no outgoing rays).
	// TODO: scott should this ever be reached?
	// fmt.Println("absorbing ray")
	return resolution{emission: Spectrum{}}
}

// func (m Metal) ComputeDirectLighting(s surfaceInteraction, scene *Scene) r3.Vec {
// 	// Metals reflect light but don't have direct diffuse contribution.
// 	// TODO: implement a small fuzz contribution
// 	return r3.Vec{}
// }

func (m Metal) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	p := s.collision.at
	n := s.collision.normal.Unit()
	wo := s.outgoing.Unit()
	directIllumination := Spectrum{}

	// Initialize the microfacet BRDF with material properties
	brdf := MicrofacetBRDF{
		Roughness: m.Fuzz,
		F0:        m.Albedo, // Base reflectivity
	}

	for _, light := range scene.Light {
		dirToLight, distanceToLight, radiantIntensity := light.Sample(p, s.incoming.rand)
		wi := dirToLight.Unit()

		// Offset the origin slightly to prevent self-intersection.
		shadowRayOrigin := p.Add(n.Muls(eps))
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
			brdfValue := brdf.Evaluate(wo, wi, n)

			// Compute the cosine term
			cosTheta := math.Max(0, n.Dot(wi))

			// Accumulate the contribution
			contribution := radiantIntensity.Mul(brdfValue).Muls(cosTheta)
			directIllumination = directIllumination.Add(Spectrum(contribution))
		}
	}
	return directIllumination
}

func init() {
	RegisterInterfaceType(Metal{})
}
