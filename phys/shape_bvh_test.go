package phys

import (
	"math"
	"math/rand"
	"testing"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// BenchmarkBVHConstructionSmall benchmarks BVH construction with a small number of shapes.
func BenchmarkBVHConstructionSmall(b *testing.B) {
	shapes := generateRandomShapes(1000) // 1,000 triangles
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBVH(shapes, 0)
	}
}

// BenchmarkBVHConstructionLarge benchmarks BVH construction with a large number of shapes.
func BenchmarkBVHConstructionLarge(b *testing.B) {
	shapes := generateRandomShapes(1000000) // 1,000,000 triangles
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewBVH(shapes, 0)
	}
}

// BenchmarkBVHCollisionDetectionSingleRay benchmarks collision detection with a single ray.
func BenchmarkBVHCollisionDetectionSingleRay(b *testing.B) {
	// Build BVH with a complex scene.
	shapes := generateRandomShapes(1000000) // 1,000,000 triangles
	bvh := NewBVH(shapes, 0)
	r := ray{
		origin:    r3.Point{X: 0, Y: 0, Z: -10},
		direction: r3.Vec{X: 0, Y: 0, Z: 1},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bvh.Collide(r, 0.001, Distance(math.MaxFloat64))
	}
}

// BenchmarkBVHCollisionDetectionMultipleRays benchmarks collision detection with multiple rays.
func BenchmarkBVHCollisionDetectionMultipleRays(b *testing.B) {
	// Build BVH with a complex scene.
	shapes := generateRandomShapes(1000000) // 1,000,000 triangles
	bvh := NewBVH(shapes, 0)
	rays := generateRandomRays(1000000) // 1,000,000 rays
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range rays {
			_, _ = bvh.Collide(r, 0.001, Distance(math.MaxFloat64))
		}
	}
}

// Helper function to generate a list of random triangles.
func generateRandomShapes(n int) []Shape {
	shapes := make([]Shape, n)
	for i := 0; i < n; i++ {
		face := createRandomFace()
		shapes[i] = face
	}
	return shapes
}

// Helper function to generate a random triangle (Face).
func createRandomFace() Face {
	// Generate random vertices within a certain range.
	v0 := Vertex{
		Position: r3.Point{
			X: randFloat()*100 - 50,
			Y: randFloat()*100 - 50,
			Z: randFloat()*100 - 50,
		},
		UV: r2.Point{X: randFloat(), Y: randFloat()},
	}
	v1 := Vertex{
		Position: r3.Point{
			X: randFloat()*100 - 50,
			Y: randFloat()*100 - 50,
			Z: randFloat()*100 - 50,
		},
		UV: r2.Point{X: randFloat(), Y: randFloat()},
	}
	v2 := Vertex{
		Position: r3.Point{
			X: randFloat()*100 - 50,
			Y: randFloat()*100 - 50,
			Z: randFloat()*100 - 50,
		},
		UV: r2.Point{X: randFloat(), Y: randFloat()},
	}
	face := Face{Vertex: [3]Vertex{v0, v1, v2}}
	return face
}

// Helper function to generate a list of random rays.
func generateRandomRays(n int) []ray {
	rays := make([]ray, n)
	for i := 0; i < n; i++ {
		r := ray{
			origin: r3.Point{
				X: randFloat()*200 - 100,
				Y: randFloat()*200 - 100,
				Z: randFloat()*200 - 100,
			},
			direction: randomUnitVector(),
		}
		rays[i] = r
	}
	return rays
}

// Helper function to generate a random float64 between 0 and 1.
func randFloat() float64 {
	return rand.Float64()
}

// Helper function to generate a random unit vector.
func randomUnitVector() r3.Vec {
	theta := randFloat() * 2 * math.Pi
	phi := math.Acos(2*randFloat() - 1)
	x := math.Sin(phi) * math.Cos(theta)
	y := math.Sin(phi) * math.Sin(theta)
	z := math.Cos(phi)
	return r3.Vec{X: x, Y: y, Z: z}
}
