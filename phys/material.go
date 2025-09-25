// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// The outgoing direction (wo) is the direction that light leaves the surface, heading toward the viewer or camera.
// The incoming direction (wi) is the direction from which light arrives at the surface point, coming from light sources or other surfaces.

type surfaceInteraction struct {
	incoming  ray       // Incoming ray.
	outgoing  r3.Vec    // Outgoing direction. Leaves the surface and goes into camera.
	collision collision // Surface collision context.
	node      Node      // The node that was hit.
}

// resolution represents the outcome of a material interaction.
// It contains the color of the surface and the rays that should be traced next.
type resolution struct {
	scattered []ray    // One or more rays reflecting, refracting, or scattering from the surface.
	emission  Spectrum // Direct emission from the surface.
}

type Material interface {
	Resolve(ctx context.Context, si surfaceInteraction) resolution
	ComputeDirectLighting(ctx context.Context, si surfaceInteraction, scene *Scene) Spectrum
	Validate() error
}
