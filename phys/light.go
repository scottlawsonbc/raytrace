// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"fmt"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type Light interface {
	Sample(p r3.Point, rand *Rand) (direction r3.Vec, distance Distance, radiance r3.Vec)
	Validate() error
}

type PointLight struct {
	Position         r3.Point
	RadiantIntensity r3.Vec // Radiant intensity (color and strength) (W/sr)
}

func (pl PointLight) Validate() error {
	if pl.RadiantIntensity.X < 0 || pl.RadiantIntensity.Y < 0 || pl.RadiantIntensity.Z < 0 {
		return fmt.Errorf("invalid PointLight RadiantIntensity: %v (should be non-negative)", pl.RadiantIntensity)
	}
	return nil
}

// TODO: what are the physical units? What is intensity? Is it radiance?
func (pl PointLight) Sample(p r3.Point, rand *Rand) (direction r3.Vec, distance Distance, radiantIntensity r3.Vec) {
	dir := pl.Position.Sub(p)
	dist := dir.Length()
	dir = dir.Divs(dist) // Normalize direction
	// Intensity does not attenuate in a path tracer (the attenuation is handled by the rendering equation)
	return dir, Distance(dist), pl.RadiantIntensity
}

func init() {
	RegisterInterfaceType(PointLight{})
}
