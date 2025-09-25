// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"
	"math"
)

type Lambertian struct {
	Texture Texture
}

func (m Lambertian) Validate() error {
	return m.Texture.Validate()
}

func (m Lambertian) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	p := s.collision.at
	n := s.collision.normal.Unit()
	directIllumination := Spectrum{}
	albedo := m.Texture.At(s.collision.uv.X, s.collision.uv.Y)
	for _, light := range scene.Light {
		dirToLight, distanceToLight, radiantIntensity := light.Sample(p, s.incoming.rand)
		// Offset the origin slightly to prevent self-intersection.
		shadowRayOrigin := p.Add(n.Muls(eps))
		shadowRay := ray{
			origin:    shadowRayOrigin,
			direction: dirToLight,
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
			nDotL := math.Max(0, n.Dot(dirToLight))
			// Accumulate the contribution.
			radiance := albedo.Mul(Spectrum(radiantIntensity)).Muls(nDotL)
			directIllumination = directIllumination.Add(radiance)
		}
	}
	return Spectrum(directIllumination)
}

// Resolve computes the reflection for a Lambertian surface interaction.
// It generates a new ray direction using cosine-weighted hemisphere sampling
// to accurately model diffuse reflection.
func (m Lambertian) Resolve(ctx context.Context, s surfaceInteraction) resolution {
	// Extract the point of collision and the surface normal.
	p := s.collision.at
	n := s.collision.normal.Unit() // Ensure the normal is normalized.

	// Sample a new direction using cosine-weighted hemisphere sampling.
	scatteredDirection := s.incoming.rand.CosineWeightedHemisphere(n)

	albedo := m.Texture.At(s.collision.uv.X, s.collision.uv.Y)
	// Create the scattered ray originating from the collision point in the sampled direction
	newRay := ray{
		origin:    p,
		direction: scatteredDirection,
		depth:     s.incoming.depth + 1,
		radiance:  s.incoming.radiance.Mul(albedo), // Scale the incoming radiance by the material's albedo.
		rand:      s.incoming.rand,
		pixelX:    s.incoming.pixelX,
		pixelY:    s.incoming.pixelY,
	}

	// Return the resolution containing the indirect scattered ray.
	return resolution{scattered: []ray{newRay}}
}

func init() {
	RegisterInterfaceType(Lambertian{})
}
