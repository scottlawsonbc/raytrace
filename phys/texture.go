// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

type Texture interface {
	// At returns the color of the texture at the given UV coordinates.
	At(u, v float64) Spectrum
	// Validate checks if the texture is valid.
	Validate() error
}

// TODO: rework how assets are retrieved and make it work in a way that can
// translate to other formats like glTF. Whatever solution I end up with
// should end up being some version of an asset loader.

