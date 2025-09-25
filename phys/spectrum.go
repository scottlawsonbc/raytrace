// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"image/color"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Spectrum represents a sampled spectrum of light with discrete bands.
// The spectrum is discretely sampled and stored as a slice of values.
// The underlying type may change as this type evolves.
// For convenience, has method to convert to color.Color for image display.
type Spectrum r3.Vec

// Add returns the sum of two spectra.
func (s Spectrum) Add(other Spectrum) Spectrum {
	return Spectrum(r3.Vec(s).Add(r3.Vec(other)))
}

// Mul returns the element-wise product of two spectra.
func (s Spectrum) Mul(other Spectrum) Spectrum {
	return Spectrum(r3.Vec(s).Mul(r3.Vec(other)))
}

// Muls returns the spectrum multiplied by a scalar.
func (s Spectrum) Muls(t float64) Spectrum {
	return Spectrum(r3.Vec(s).Muls(t))
}

// Divs returns the spectrum divided by a scalar.
func (s Spectrum) Divs(t float64) Spectrum {
	return Spectrum(r3.Vec(s).Divs(t))
}

// Clip returns the spectrum with each component clipped to the range [min, max].
func (s Spectrum) Clip(min, max float64) Spectrum {
	return Spectrum(r3.Vec(s).Clip(min, max))
}

// ToColor converts the spectrum to a color.Color.
func (s Spectrum) ToColor() color.Color {
	c := s.Clip(0, 1)
	return color.RGBA{
		R: uint8(c.X * 255),
		G: uint8(c.Y * 255),
		B: uint8(c.Z * 255),
		A: 255,
	}
}

// String returns a string representation of the spectrum.
func (s Spectrum) String() string {
	return r3.Vec(s).String()
}
