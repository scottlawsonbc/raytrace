// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// // Transform represents a 4x4 transformation matrix.
// type Transform struct {
// 	Matrix        Matrix4x4
// 	InverseMatrix Matrix4x4 // Precompute the inverse for efficiency.
// }

// // NewTransform creates a new Transform with an idnode matrix.
// func NewTransform() Transform {
// 	identity := IdentityMatrix()
// 	return Transform{
// 		Matrix:        identity,
// 		InverseMatrix: identity,
// 	}
// }

// // ApplyToPoint applies the transformation to a r3.Point.
// func (t Transform) ApplyToPoint(p r3.Point) r3.Point {
// 	return t.Matrix.Transformr3.Point(p)
// }

// // ApplyToVector applies the transformation to a r3.Vec.
// func (t Transform) ApplyToVector(v r3.Vec) r3.Vec {
// 	return t.Matrix.Transformr3.Vec(v)
// }

// // Inverse returns the inverse of the transformation.
// func (t Transform) Inverse() Transform {
// 	return Transform{
// 		Matrix:        t.InverseMatrix,
// 		InverseMatrix: t.Matrix,
// 	}
// }

// // Combine combines the current transform with another.
// func (t Transform) Combine(other Transform) Transform {
// 	combinedMatrix := t.Matrix.Multiply(other.Matrix)
// 	combinedInverse := other.InverseMatrix.Multiply(t.InverseMatrix)
// 	return Transform{
// 		Matrix:        combinedMatrix,
// 		InverseMatrix: combinedInverse,
// 	}
// }

// Transform represents a transformation in 3D space, including translation,
// rotation (as a matrix), and scaling.
type Transform struct {
	Translation r3.Vec
	Rotation    r3.Mat3x3
	Scale       r3.Vec
}

// NewTransform creates a new Transform with default values (idnode).
func NewTransform() Transform {
	return Transform{
		Translation: r3.Vec{X: 0, Y: 0, Z: 0},
		Rotation:    r3.IdentityMat3x3(),
		Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
	}
}

// ApplyToPoint applies the transformation to a r3.Point.
func (t Transform) ApplyToPoint(p r3.Point) r3.Point {
	// Scale, then rotate, then translate.
	scaled := r3.Vec{X: p.X * t.Scale.X, Y: p.Y * t.Scale.Y, Z: p.Z * t.Scale.Z}
	rotated := t.Rotation.MulVec(scaled)
	translated := rotated.Add(t.Translation)
	return r3.Point(translated)
}

// ApplyToVector applies the transformation to a r3.Vec (ignoring translation).
func (t Transform) ApplyToVector(v r3.Vec) r3.Vec {
	// Scale, then rotate
	scaled := r3.Vec{X: v.X * t.Scale.X, Y: v.Y * t.Scale.Y, Z: v.Z * t.Scale.Z}
	rotated := t.Rotation.MulVec(scaled)
	return rotated
}

// Inverse returns the inverse of the transformation.
func (t Transform) Inverse() Transform {
	// Invert scale
	invScale := r3.Vec{
		X: 1 / t.Scale.X,
		Y: 1 / t.Scale.Y,
		Z: 1 / t.Scale.Z,
	}
	// Invert rotation
	invRotation := t.Rotation.Transpose()
	// Invert translation
	invTranslation := invRotation.MulVec(t.Translation.Muls(-1)).Mul(invScale)
	return Transform{
		Translation: invTranslation,
		Rotation:    invRotation,
		Scale:       invScale,
	}
}
