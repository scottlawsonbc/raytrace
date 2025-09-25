// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
// Props useful for creating scenes.
package phys

import (
	"fmt"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func PropSkySphere(radius Distance, mat Material) Node {
	return Node{
		Name:     "SkySphere",
		Material: mat,
		Shape: Sphere{
			Center: r3.Point{},
			Radius: radius,
		}}
}

func PropAxes(origin r3.Point, radius Distance, len Distance, prefix string) []Node {
	// Sanity checks:
	// - radius should be > 0 and << len
	// - len should be > 0
	// - prefix can be "" or a short string
	// In case of error, always panic with human readable message to find errors quickly.
	// At this stage of development we don't want to let errors go unnoticed.
	if radius <= 0 {
		panic(fmt.Sprintf("PropAxes: invalid radius %v (must be > 0)", radius))
	}
	if len <= 0 {
		panic(fmt.Sprintf("PropAxes: invalid len %v (must be > 0)", len))
	}
	if len <= radius {
		panic(fmt.Sprintf("PropAxes: invalid radius %v (must be << len %v)", radius, len))
	}
	return []Node{
		{
			Name: prefix + "AxisX",
			Material: Emitter{
				Texture: TextureUniform{Color: Spectrum{X: 1, Y: 0, Z: 0}},
			},
			Shape: Cylinder{
				Origin:    origin,
				Direction: r3.Vec{X: 1, Y: 0, Z: 0},
				Radius:    radius,
				Height:    len,
			},
		},
		{
			Name: prefix + "AxisY",
			Material: Emitter{
				Texture: TextureUniform{Color: Spectrum{X: 0, Y: 1, Z: 0}},
			},
			Shape: Cylinder{
				Origin:    origin,
				Direction: r3.Vec{X: 0, Y: 1, Z: 0},
				Radius:    radius,
				Height:    len,
			},
		},
		{
			Name: prefix + "AxisZ",
			Material: Emitter{
				Texture: TextureUniform{Color: Spectrum{X: 0, Y: 0, Z: 1}},
			},
			Shape: Cylinder{
				Origin:    origin,
				Direction: r3.Vec{X: 0, Y: 0, Z: 1},
				Radius:    radius,
				Height:    len,
			},
		},
	}
}
