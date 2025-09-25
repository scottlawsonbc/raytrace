package phys

import (
	"fmt"
	"math"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Vertex represents a single vertex in 3D space with texture coordinates.
type Vertex struct {
	Position r3.Point
	UV       r2.Point
	// Possibly other fields here in the future, like normals, tangents, etc.
}

// Face represents a triangular mesh face.
type Face struct {
	Vertex [3]Vertex
}

// Validate performs comprehensive validation checks on the Face instance.
func (face Face) Validate() error {
	const eps = 1e-8

	// Check every vertex (Point2 and Point3) for NaN and Inf.
	checkPoint2 := func(p r2.Point, label string) error {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) ||
			math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
			return fmt.Errorf("invalid Face: %s contains NaN or Inf: %+v", label, p)
		}
		return nil
	}
	checkPoint3 := func(p r3.Point, label string) error {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) || math.IsNaN(p.Z) ||
			math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) || math.IsInf(p.Z, 0) {
			return fmt.Errorf("invalid Face: %s contains NaN or Inf: %+v", label, p)
		}
		return nil
	}
	for i, v := range face.Vertex {
		if err := checkPoint3(v.Position, fmt.Sprintf("Vertex[%d].Position", i)); err != nil {
			return err
		}
		if err := checkPoint2(v.UV, fmt.Sprintf("Vertex[%d].UV", i)); err != nil {
			return err
		}
	}

	// Check for duplicate vertices.
	p0 := face.Vertex[0].Position
	p1 := face.Vertex[1].Position
	p2 := face.Vertex[2].Position
	if p0 == p1 || p0 == p2 || p1 == p2 {
		return fmt.Errorf("invalid Face: two or more vertices are identical")
	}

	// Compute edges.
	edge1 := p1.Sub(p0)
	edge2 := p2.Sub(p0)

	// Compute normal
	normal := edge1.Cross(edge2)
	if normal.IsZero() {
		return fmt.Errorf("invalid Face: normal is a zero vector")
	}
	normal = normal.Unit()

	// Check normal vector is unit length.
	normalLength := normal.Length()
	if normalLength < 1-eps || normalLength > 1+eps {
		return fmt.Errorf("invalid Face: normal is not a unit vector (length: %f)", normalLength)
	}

	// Check triangle area is not degenerate.
	crossProduct := edge1.Cross(edge2)
	area := 0.5 * crossProduct.Length()
	const epsilonArea = 1e-12
	if area < epsilonArea {
		return fmt.Errorf("invalid Face: triangle is degenerate (zero or near-zero area)")
	}

	// Check normal is orthogonal to edges.
	dot1 := normal.Dot(edge1)
	dot2 := normal.Dot(edge2)
	if math.Abs(dot1) > eps || math.Abs(dot2) > eps {
		return fmt.Errorf("invalid Face: normal is not orthogonal to triangle edges (dot1: %f, dot2: %f)", dot1, dot2)
	}
	return nil
}

// Collide determines whether a given ray intersects with the Face.
// It also interpolates the UV coordinates at the intersection point.
func (face Face) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	const eps = 1e-8

	// Möller-Trumbore intersection algorithm.
	p0 := face.Vertex[0].Position
	p1 := face.Vertex[1].Position
	p2 := face.Vertex[2].Position

	edge1 := p1.Sub(p0)
	edge2 := p2.Sub(p0)
	h := r.direction.Cross(edge2)
	a := edge1.Dot(h)
	if a > -eps && a < eps {
		return false, collision{}
	}
	f := 1 / a
	s := r.origin.Sub(p0)
	u := f * s.Dot(h)
	if u < -eps || u > 1.0+eps {
		return false, collision{}
	}
	q := s.Cross(edge1)
	v := f * r.direction.Dot(q)
	if v < -eps || u+v > 1.0+eps {
		return false, collision{}
	}
	t := f * edge2.Dot(q)
	if t < float64(tmin) || t > float64(tmax) {
		return false, collision{}
	}
	at := r.at(Distance(t))

	// Compute face normal.
	normal := edge1.Cross(edge2).Unit()

	// Interpolate UVs using barycentric coordinates.
	uv0 := face.Vertex[0].UV
	uv1 := face.Vertex[1].UV
	uv2 := face.Vertex[2].UV
	w := 1 - u - v
	interpolatedUV := uv0.Muls(w).Add(uv1.Muls(u)).Add(uv2.Muls(v))
	return true, collision{
		t:      Distance(t),
		at:     at,
		normal: normal,
		uv:     interpolatedUV,
	}
}

// Bounds computes the Axis-Aligned Bounding Box (AABB) of the Face.
func (f Face) Bounds() AABB {
	p0 := f.Vertex[0].Position
	p1 := f.Vertex[1].Position
	p2 := f.Vertex[2].Position

	min := r3.Point{
		X: math.Min(math.Min(p0.X, p1.X), p2.X),
		Y: math.Min(math.Min(p0.Y, p1.Y), p2.Y),
		Z: math.Min(math.Min(p0.Z, p1.Z), p2.Z),
	}
	max := r3.Point{
		X: math.Max(math.Max(p0.X, p1.X), p2.X),
		Y: math.Max(math.Max(p0.Y, p1.Y), p2.Y),
		Z: math.Max(math.Max(p0.Z, p1.Z), p2.Z),
	}
	return AABB{Min: min, Max: max}
}

// Mesh represents a collection of Faces forming a mesh.
type Mesh struct {
	Face []Face
	BVH  *BVH
}

// NewMesh creates a new Mesh and builds the BVH.
func NewMesh(faces []Face) (*Mesh, error) {
	m := &Mesh{Face: faces}
	var shapes []Shape
	for _, face := range faces {
		shapes = append(shapes, face)
	}
	m.BVH = NewBVH(shapes, 0)
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return m, nil
}

// Validate performs comprehensive validation checks on the Mesh instance.
func (m Mesh) Validate() error {
	if len(m.Face) == 0 {
		return fmt.Errorf("Mesh must contain at least one face")
	}
	for i, face := range m.Face {
		if err := face.Validate(); err != nil {
			return fmt.Errorf("Mesh face %d is invalid: %v", i, err)
		}
	}
	if m.BVH == nil {
		return fmt.Errorf("Mesh must have a BVH")
	}
	if err := m.BVH.Validate(); err != nil {
		return fmt.Errorf("Mesh BVH validation error: %v", err)
	}
	return nil
}

// Collide determines whether a given ray intersects with the Mesh.
// It delegates collision detection to the BVH.
func (m Mesh) Collide(r ray, tmin, tmax Distance) (bool, collision) {
	return m.BVH.Collide(r, tmin, tmax)
}

// Bounds computes the Axis-Aligned Bounding Box (AABB) of the Mesh.
func (m Mesh) Bounds() AABB {
	return m.BVH.Bounds()
}

// String returns a string representation of the Mesh.
func (m *Mesh) String() string {
	return fmt.Sprintf("Mesh{Faces: %d, BVH: %v}", len(m.Face), m.BVH)
}

func init() {
	RegisterInterfaceType(Mesh{})
	RegisterInterfaceType(Face{})
}
