package r3

import "math"

// Mat3x3 represents a 3x3 matrix.
type Mat3x3 struct {
	M [3][3]float64
}

// IdentityMat3x3 returns an identity matrix.
func IdentityMat3x3() Mat3x3 {
	return Mat3x3{
		M: [3][3]float64{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
	}
}

// MulVec multiplies the matrix by a vector.
func (m Mat3x3) MulVec(v Vec) Vec {
	return Vec{
		X: m.M[0][0]*v.X + m.M[0][1]*v.Y + m.M[0][2]*v.Z,
		Y: m.M[1][0]*v.X + m.M[1][1]*v.Y + m.M[1][2]*v.Z,
		Z: m.M[2][0]*v.X + m.M[2][1]*v.Y + m.M[2][2]*v.Z,
	}
}

// Multiply multiplies the current matrix with another Mat3x3.
func (m Mat3x3) Mul(n Mat3x3) Mat3x3 {
	var result Mat3x3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := 0.0
			for k := 0; k < 3; k++ {
				sum += m.M[i][k] * n.M[k][j]
			}
			result.M[i][j] = sum
		}
	}
	return result
}

// Transpose returns the transpose of the matrix.
func (m Mat3x3) Transpose() Mat3x3 {
	return Mat3x3{
		M: [3][3]float64{
			{m.M[0][0], m.M[1][0], m.M[2][0]},
			{m.M[0][1], m.M[1][1], m.M[2][1]},
			{m.M[0][2], m.M[1][2], m.M[2][2]},
		},
	}
}

// Rotation matrices around X, Y, and Z axes.
func RotationMatrixX(angle float64) Mat3x3 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Mat3x3{
		M: [3][3]float64{
			{1, 0, 0},
			{0, c, -s},
			{0, s, c},
		},
	}
}

func RotationMatrixY(angle float64) Mat3x3 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Mat3x3{
		M: [3][3]float64{
			{c, 0, s},
			{0, 1, 0},
			{-s, 0, c},
		},
	}
}

// RotationMatrixZ returns the rotation matrix about the Z axis for the radian argument angle.
func RotationMatrixZ(angle float64) Mat3x3 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	return Mat3x3{
		M: [3][3]float64{
			{c, -s, 0},
			{s, c, 0},
			{0, 0, 1},
		},
	}
}
