package r3

import (
	"fmt"
	"math"
)

// Vec represents a three-dimensional vector with X, Y, and Z components.
type Vec struct {
	X float64
	Y float64
	Z float64
}

// Add returns the vector addition of the current vector and v2.
func (v Vec) Add(v2 Vec) Vec {
	return Vec{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

// Sub returns the vector subtraction of v2 from the current vector.
func (v Vec) Sub(v2 Vec) Vec {
	return Vec{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

// Mul returns the component-wise multiplication of the current vector and v2.
func (v Vec) Mul(v2 Vec) Vec {
	return Vec{v.X * v2.X, v.Y * v2.Y, v.Z * v2.Z}
}

// Div returns the component-wise division of the current vector by v2.
func (v Vec) Div(v2 Vec) Vec {
	return Vec{v.X / v2.X, v.Y / v2.Y, v.Z / v2.Z}
}

// Muls returns the current vector multiplied by a scalar value s.
func (v Vec) Muls(s float64) Vec {
	return Vec{v.X * s, v.Y * s, v.Z * s}
}

// Divs returns the current vector divided by a scalar value s.
func (v Vec) Divs(s float64) Vec {
	return Vec{v.X / s, v.Y / s, v.Z / s}
}

// Dot computes the dot product of the current vector with v2.
// The dot product is a scalar representing the cosine of the angle
// between the vectors multiplied by their magnitudes (lengths).
func (v Vec) Dot(v2 Vec) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

// Cross computes the cross product of the current vector with v2.
// The cross product results in a new vector perpendicular to both input vectors.
func (v Vec) Cross(v2 Vec) Vec {
	return Vec{v.Y*v2.Z - v.Z*v2.Y, v.Z*v2.X - v.X*v2.Z, v.X*v2.Y - v.Y*v2.X}
}

// Lerp performs linear interpolation between the current vector and v2 by parameter t.
// The parameter t is clamped between 0 and 1.
// When t=0, it returns the current vector.
// When t=1, it returns v2.
func (v Vec) Lerp(v2 Vec, t float64) Vec {
	// Clamp t to the range [0, 1] to ensure valid interpolation
	t = math.Max(0, math.Min(1, t))
	return Vec{
		X: v.X + t*(v2.X-v.X),
		Y: v.Y + t*(v2.Y-v.Y),
		Z: v.Z + t*(v2.Z-v.Z),
	}
}

// Eq checks if the current vector is exactly equal to v2.
// It returns true only if all corresponding components are identical.
func (v Vec) Eq(v2 Vec) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

// Close checks if the current vector is approximately equal to v2 within an absolute tolerance.
// It returns true only if all corresponding components are approximately equal.
func (v Vec) IsClose(v2 Vec, atol float64) bool {
	return math.Abs(v.X-v2.X) < atol && math.Abs(v.Y-v2.Y) < atol && math.Abs(v.Z-v2.Z) < atol
}

// Length returns the Euclidean length (magnitude) of the vector.
func (v Vec) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

// Unit returns the unit vector (vector with length 1) in the direction of the current vector.
// If the vector is zero, it returns the zero vector to avoid division by zero.
func (v Vec) Unit() Vec {
	length := v.Length()
	if length == 0 {
		return Vec{0, 0, 0}
	}
	return v.Divs(length)
}

// Clip clamps each component of the vector between the specified min and max values.
// It ensures that X, Y, and Z are not less than min and not greater than max.
func (v Vec) Clip(min, max float64) Vec {
	return Vec{math.Min(math.Max(v.X, min), max), math.Min(math.Max(v.Y, min), max), math.Min(math.Max(v.Z, min), max)}
}

// IsNaN checks if any component of the vector is NaN (Not a Number).
func (v Vec) IsNaN() bool {
	return math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsNaN(v.Z)
}

// IsInf checks if any component of the vector is infinite.
func (v Vec) IsInf() bool {
	return math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || math.IsInf(v.Z, 0)
}

// IsZero reports whether the vector is the zero vector (0, 0, 0).
func (v Vec) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

// Get returns the value of the vector component at the specified index.
// Index 0 corresponds to X, 1 to Y, and 2 to Z.
// It panics if the index is out of bounds.
func (v Vec) Get(i int) float64 {
	switch i {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	}
	panic(fmt.Sprintf("invalid index %d", i))
}

func (v Vec) String() string {
	return fmt.Sprintf("(%v, %v, %v)", v.X, v.Y, v.Z)
}
