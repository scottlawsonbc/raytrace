package r3

import (
	"fmt"
	"math"
)

// Point represents a point in three-dimensional space with X, Y, and Z coordinates.
type Point struct {
	X float64
	Y float64
	Z float64
}

// Set returns a new Point with the specified component set to v.
// Index 0 corresponds to X, 1 to Y, and 2 to Z.
// It panics if the index is out of bounds.
func (p Point) Set(i int, v float64) Point {
	switch i {
	case 0:
		return Point{v, p.Y, p.Z}
	case 1:
		return Point{p.X, v, p.Z}
	case 2:
		return Point{p.X, p.Y, v}
	}
	panic(fmt.Sprintf("invalid index `%d` for Point", i))
}

// Get returns the value of the point component at the specified index.
// Index 0 corresponds to X, 1 to Y, and 2 to Z.
// It panics if the index is out of bounds.
func (p Point) Get(i int) float64 {
	switch i {
	case 0:
		return p.X
	case 1:
		return p.Y
	case 2:
		return p.Z
	}
	panic(fmt.Sprintf("invalid index `%d` for Point", i))
}

// Sub subtracts another Point from the current Point, resulting in a Vec.
// This represents the vector from p2 to p.
func (p Point) Sub(p2 Point) Vec {
	return Vec{p.X - p2.X, p.Y - p2.Y, p.Z - p2.Z}
}

// Add adds a Vec to the current Point, resulting in a new Point.
// This translates the point by the given vector.
func (p Point) Add(v Vec) Point {
	return Point{p.X + v.X, p.Y + v.Y, p.Z + v.Z}
}

// Subv subtracts a Vec from the current Point, resulting in a new Point.
// This translates the point by the negative of the given vector.
func (p Point) Subv(v Vec) Point {
	return Point{p.X - v.X, p.Y - v.Y, p.Z - v.Z}
}

// Lerp performs linear interpolation between the current point and p2 by parameter t.
// The parameter t is clamped between 0 and 1.
// When t=0, it returns the current point.
// When t=1, it returns p2.
func (p Point) Lerp(p2 Point, t float64) Point {
	// Clamp t to the range [0, 1] to ensure valid interpolation
	t = math.Max(0, math.Min(1, t))
	return Point{
		X: p.X + t*(p2.X-p.X),
		Y: p.Y + t*(p2.Y-p.Y),
		Z: p.Z + t*(p2.Z-p.Z),
	}
}

// Eq checks if the current Point is exactly equal to p2.
// It returns true only if all corresponding coordinates are identical.
func (p Point) Eq(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y && p.Z == p2.Z
}

// Close checks if the current Point is approximately equal to p2 within a small epsilon.
// It returns true only if all corresponding coordinates are approximately equal.
func (p Point) IsClose(p2 Point, atol float64) bool {
	return math.Abs(p.X-p2.X) < atol && math.Abs(p.Y-p2.Y) < atol && math.Abs(p.Z-p2.Z) < atol
}

// IsNaN checks if any coordinate of the point is NaN (Not a Number).
func (p Point) IsNaN() bool {
	return math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsNaN(p.Z)
}

// IsInf checks if any coordinate of the point is infinite.
func (p Point) IsInf() bool {
	return math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) || math.IsInf(p.Z, 0)
}

// IsZero reports whether the point is the origin (0, 0, 0).
func (p Point) IsZero() bool {
	return p.X == 0 && p.Y == 0 && p.Z == 0
}

// String returns a string representation of the point.
func (p Point) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p.X, p.Y, p.Z)
}
