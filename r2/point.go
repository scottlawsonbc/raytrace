package r2

import (
	"fmt"
	"math"
)

// Point represents a 2D coordinate with X and Y components.
// This is used for texture coordinates in the context of ray tracing.
type Point struct {
	X float64 // aka U.
	Y float64 // aka V.
}

// Add adds another Point to the current Point component-wise.
// It returns a new Point with the summed values.
func (p Point) Add(p2 Point) Point {
	return Point{
		X: p.X + p2.X,
		Y: p.Y + p2.Y,
	}
}

// Muls multiplies each component of the Point by the scalar s.
// It returns a new Point with the scaled values.
func (p Point) Muls(s float64) Point {
	return Point{
		X: p.X * s,
		Y: p.Y * s,
	}
}

// Lerp performs linear interpolation between the current Point and another Point.
// The parameter t is clamped between 0 and 1 to ensure valid interpolation.
// When t=0, it returns the current Point.
// When t=1, it returns the other Point.
// Values between 0 and 1 interpolate between the two points.
func (p Point) Lerp(p2 Point, t float64) Point {
	// Clamp t to the range [0, 1]
	t = math.Max(0, math.Min(1, t))
	return Point{
		X: (1-t)*p.X + t*p2.X,
		Y: (1-t)*p.Y + t*p2.Y,
	}
}

// Eq checks if the current Point is exactly equal to p2.
func (p Point) Eq(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}

// Close checks if the current Point is approximately equal to p2 within a small absolute tolerance.
func (p Point) IsClose(p2 Point, atol float64) bool {
	return math.Abs(p.X-p2.X) < atol && math.Abs(p.Y-p2.Y) < atol
}

// Sub subtracts another Point from the current Point, resulting in a Vec.
// This represents the vector from p2 to p.
func (p Point) Sub(p2 Point) Vec {
	return Vec{p.X - p2.X, p.Y - p2.Y}
}

// Clip clamps each component of the point between the specified min and max values.
func (p Point) Clip(min, max float64) Point {
	return Point{math.Min(math.Max(p.X, min), max), math.Min(math.Max(p.Y, min), max)}
}

// String returns a string representation of the point.
func (p Point) String() string {
	return fmt.Sprintf("(%v, %v)", p.X, p.Y)
}
