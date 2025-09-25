// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

// eps is a small value used to avoid floating point errors.
const eps = 1e-6

// import "math"

// // // This file contains extra math types that aren't actually used right now and
// // // are provided for reference in case they are needed in the future.
// // // These might get removed in the future.

// // Matrix4x4 represents a 4x4 transformation matrix used for performing
// // affine and projective transformations in three-dimensional space.
// //
// // The matrix operates on homogeneous coordinates, distinguishing between
// // points and vectors through the w component:
// //
// //   - **Points** are represented with w=1. Transformations applied to points
// //     include translation, rotation, and scaling.
// //
// //   - **Vectors** are represented with w=0. Transformations applied to vectors
// //     include rotation and scaling, but **exclude** translation to preserve
// //     direction and magnitude.
// //
// // This distinction ensures that points and vectors behave correctly under
// // various transformations, which is essential for accurate rendering and
// // lighting calculations.
// type Matrix4x4 [4][4]float64

// // IdentityMatrix returns an identity matrix, which leaves vectors and points
// // unchanged when applied as a transformation.
// func IdentityMatrix() Matrix4x4 {
// 	return Matrix4x4{
// 		{1, 0, 0, 0},
// 		{0, 1, 0, 0},
// 		{0, 0, 1, 0},
// 		{0, 0, 0, 1},
// 	}
// }

// // Multiply performs matrix multiplication between the current matrix and another Matrix4x4.
// // The result is a new Matrix4x4 representing the combined transformation.
// func (m Matrix4x4) Multiply(n Matrix4x4) Matrix4x4 {
// 	var result Matrix4x4
// 	for i := 0; i < 4; i++ {
// 		for j := 0; j < 4; j++ {
// 			for k := 0; k < 4; k++ {
// 				result[i][j] += m[i][k] * n[k][j]
// 			}
// 		}
// 	}
// 	return result
// }

// // TransformrVec applies the transformation matrix to a r3.Vec, treating it as a vector with w=0.
// // This means the vector is affected by rotation and scaling but not by translation.
// func (m Matrix4x4) TransformVec(v r3.Vec) r3.Vec {
// 	x := m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z
// 	y := m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z
// 	z := m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z
// 	return r3.Vec{X: x, Y: y, Z: z}
// }

// // TransformPoint applies the transformation matrix to a r3.Point, treating it as a point with w=1.
// // This means the point is affected by rotation, scaling, and translation.
// // If the resulting w component is not 1, the coordinates are normalized by w to maintain proper positioning.
// func (m Matrix4x4) TransformPoint(p r3.Point) r3.Point {
// 	x := m[0][0]*p.X + m[0][1]*p.Y + m[0][2]*p.Z + m[0][3]
// 	y := m[1][0]*p.X + m[1][1]*p.Y + m[1][2]*p.Z + m[1][3]
// 	z := m[2][0]*p.X + m[2][1]*p.Y + m[2][2]*p.Z + m[2][3]
// 	w := m[3][0]*p.X + m[3][1]*p.Y + m[3][2]*p.Z + m[3][3]
// 	if w != 0 && w != 1 {
// 		return r3.Point{X: x / w, Y: y / w, Z: z / w}
// 	}
// 	return r3.Point{X: x, Y: y, Z: z}
// }

// // TranslationMatrix creates a translation matrix.
// func TranslationMatrix(dx, dy, dz float64) Matrix4x4 {
// 	return Matrix4x4{
// 		{1, 0, 0, dx},
// 		{0, 1, 0, dy},
// 		{0, 0, 1, dz},
// 		{0, 0, 0, 1},
// 	}
// }

// // ScalingMatrix creates a scaling matrix.
// func ScalingMatrix(sx, sy, sz float64) Matrix4x4 {
// 	return Matrix4x4{
// 		{sx, 0, 0, 0},
// 		{0, sy, 0, 0},
// 		{0, 0, sz, 0},
// 		{0, 0, 0, 1},
// 	}
// }

// // RotationMatrixX creates a rotation matrix around the X-axis.
// func RotationMatrixX(angle float64) Matrix4x4 {
// 	c := math.Cos(angle)
// 	s := math.Sin(angle)
// 	return Matrix4x4{
// 		{1, 0, 0, 0},
// 		{0, c, -s, 0},
// 		{0, s, c, 0},
// 		{0, 0, 0, 1},
// 	}
// }

// // RotationMatrixY creates a rotation matrix around the Y-axis.
// func RotationMatrixY(angle float64) Matrix4x4 {
// 	c := math.Cos(angle)
// 	s := math.Sin(angle)
// 	return Matrix4x4{
// 		{c, 0, s, 0},
// 		{0, 1, 0, 0},
// 		{-s, 0, c, 0},
// 		{0, 0, 0, 1},
// 	}
// }

// // RotationMatrixZ creates a rotation matrix around the Z-axis.
// func RotationMatrixZ(angle float64) Matrix4x4 {
// 	c := math.Cos(angle)
// 	s := math.Sin(angle)
// 	return Matrix4x4{
// 		{c, -s, 0, 0},
// 		{s, c, 0, 0},
// 		{0, 0, 1, 0},
// 		{0, 0, 0, 1},
// 	}
// }

// // Quaternion represents a quaternion with X, Y, Z, and W components.
// type Quaternion struct {
// 	X, Y, Z, W float64
// }

// // Multiply multiplies the current quaternion with q2.
// func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
// 	return Quaternion{
// 		X: q.W*q2.X + q.X*q2.W + q.Y*q2.Z - q.Z*q2.Y,
// 		Y: q.W*q2.Y - q.X*q2.Z + q.Y*q2.W + q.Z*q2.X,
// 		Z: q.W*q2.Z + q.X*q2.Y - q.Y*q2.X + q.Z*q2.W,
// 		W: q.W*q2.W - q.X*q2.X - q.Y*q2.Y - q.Z*q2.Z,
// 	}
// }

// // ToRotationMatrix converts the quaternion to a rotation matrix.
// func (q Quaternion) ToRotationMatrix() Matrix4x4 {
// 	return Matrix4x4{
// 		{
// 			1 - 2*q.Y*q.Y - 2*q.Z*q.Z,
// 			2*q.X*q.Y - 2*q.Z*q.W,
// 			2*q.X*q.Z + 2*q.Y*q.W,
// 			0,
// 		},
// 		{
// 			2*q.X*q.Y + 2*q.Z*q.W,
// 			1 - 2*q.X*q.X - 2*q.Z*q.Z,
// 			2*q.Y*q.Z - 2*q.X*q.W,
// 			0,
// 		},
// 		{
// 			2*q.X*q.Z - 2*q.Y*q.W,
// 			2*q.Y*q.Z + 2*q.X*q.W,
// 			1 - 2*q.X*q.X - 2*q.Y*q.Y,
// 			0,
// 		},
// 		{0, 0, 0, 1},
// 	}
// }

// // Slerp performs spherical linear interpolation between two quaternions.
// func Slerp(q1, q2 Quaternion, t float64) Quaternion {
// 	// Compute the cosine of the angle between the quaternions.
// 	cosTheta := q1.X*q2.X + q1.Y*q2.Y + q1.Z*q2.Z + q1.W*q2.W

// 	// If cosTheta < 0, the interpolation will take the long way around the sphere.
// 	// To fix this, invert one quaternion.
// 	if cosTheta < 0 {
// 		q2 = Quaternion{-q2.X, -q2.Y, -q2.Z, -q2.W}
// 		cosTheta = -cosTheta
// 	}

// 	// If the quaternions are close, use linear interpolation.
// 	if cosTheta > 0.9995 {
// 		return Quaternion{
// 			X: q1.X + t*(q2.X-q1.X),
// 			Y: q1.Y + t*(q2.Y-q1.Y),
// 			Z: q1.Z + t*(q2.Z-q1.Z),
// 			W: q1.W + t*(q2.W-q1.W),
// 		}.Unit()
// 	}

// 	// Compute the angle between the quaternions.
// 	theta := math.Acos(cosTheta)
// 	sinTheta := math.Sin(theta)

// 	// Compute interpolation factors.
// 	a := math.Sin((1-t)*theta) / sinTheta
// 	b := math.Sin(t*theta) / sinTheta

// 	return Quaternion{
// 		X: a*q1.X + b*q2.X,
// 		Y: a*q1.Y + b*q2.Y,
// 		Z: a*q1.Z + b*q2.Z,
// 		W: a*q1.W + b*q2.W,
// 	}
// }

// // Unit normalizes the quaternion to unit length.
// func (q Quaternion) Unit() Quaternion {
// 	length := math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
// 	if length > eps {
// 		return Quaternion{q.X / length, q.Y / length, q.Z / length, q.W / length}
// 	}
// 	return Quaternion{0, 0, 0, 1}
// }
