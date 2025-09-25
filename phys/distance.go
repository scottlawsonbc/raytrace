// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import "fmt"

// Distance represents a physical distance in float64 nanometers.
type Distance float64

const (
	Nanometer  Distance = 1
	Micrometer          = 1000 * Nanometer
	Millimeter          = 1000 * Micrometer
	Meter               = 1000 * Millimeter
	NM                  = Nanometer
	UM                  = Micrometer
	MM                  = Millimeter
	M                   = Meter
)

func (d Distance) Nanometers() float64 {
	return float64(d)
}

func (d Distance) Micrometers() float64 {
	return float64(d) / float64(Micrometer)
}

func (d Distance) Millimeters() float64 {
	return float64(d) / float64(Millimeter)
}

func (d Distance) Meters() float64 {
	return float64(d) / float64(Meter)
}

func (d Distance) String() string {
	if d < Micrometer {
		return fmt.Sprintf("%f nm", d.Nanometers())
	}
	if d < Millimeter {
		return fmt.Sprintf("%f µm", d.Micrometers())
	}
	if d < Meter {
		return fmt.Sprintf("%f mm", d.Millimeters())
	}
	return fmt.Sprintf("%f m", d.Meters())
}
