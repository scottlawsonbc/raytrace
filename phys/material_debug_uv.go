// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"
	"math"
)

// DebugUV visualizes the UV coordinates as colors.
type DebugUV struct{}

func (m DebugUV) Validate() error {
	return nil
}

// Resolve computes the UV-based color for a surface interaction.
func (m DebugUV) Resolve(ctx context.Context, s surfaceInteraction) resolution {
	// Clamp UVs to [0,1] for visualization.
	if s.collision.uv.X < 0.0 || s.collision.uv.X > 1.0 {
		return resolution{emission: Spectrum{X: 1.0, Y: 0.0, Z: 0.0}}
	}
	u := math.Min(math.Max(s.collision.uv.X, 0.0), 1.0)
	v := math.Min(math.Max(s.collision.uv.Y, 0.0), 1.0)
	// Map U to Red, V to Green, and set Blue to 0.5 for visibility.
	color := Spectrum{X: u, Y: v, Z: 0.5}
	return resolution{emission: color}
}

func (m DebugUV) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	return Spectrum{} // No direct lighting for UV shader.
}

func init() {
	RegisterInterfaceType(DebugUV{})
}
