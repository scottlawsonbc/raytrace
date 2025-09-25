package obj

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"testing/fstest"
)

// TestParseOBJ_Basic tests parsing a simple OBJ file with vertices and faces.
func TestParseOBJ_Basic(t *testing.T) {
	objData := `
# Simple triangle
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.0 1.0 0.0
f 1 2 3
`

	fsys := fstest.MapFS{
		"triangle.obj": {Data: []byte(objData)},
	}

	file, err := fsys.Open("triangle.obj")
	if err != nil {
		t.Fatalf("Failed to open OBJ file: %v", err)
	}
	defer file.Close()

	obj, err := ParseFS(fsys, "triangle.obj")
	if err != nil {
		t.Fatalf("Failed to parse OBJ file: %v", err)
	}

	// Validate parsed data
	if len(obj.Vertices) != 3 {
		t.Errorf("Expected 3 vertices, got %d", len(obj.Vertices))
	}
	if len(obj.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(obj.Faces))
	}
}

// TestParseOBJ_Materials tests parsing an OBJ file with materials and textures.
func TestParseOBJ_Materials(t *testing.T) {
	objData := `
# Cube with materials
v -1.0 -1.0 -1.0
v 1.0 -1.0 -1.0
v 1.0 1.0 -1.0
v -1.0 1.0 -1.0
v -1.0 -1.0 1.0
v 1.0 -1.0 1.0
v 1.0 1.0 1.0
v -1.0 1.0 1.0

usemtl Material001
mtllib cube.mtl

f 1 2 3 4
f 5 6 7 8
`

	mtlData := `
# Material definition
newmtl Material001
Kd 0.8 0.8 0.8
map_Kd texture.jpg
`

	fsys := fstest.MapFS{
		"cube.obj":    {Data: []byte(objData)},
		"cube.mtl":    {Data: []byte(mtlData)},
		"texture.jpg": {Data: []byte("fake image data")},
	}

	obj, err := ParseFS(fsys, "cube.obj")
	if err != nil {
		t.Fatalf("Failed to parse OBJ file: %v", err)
	}

	// Validate materials
	if len(obj.Materials) != 1 {
		t.Errorf("Expected 1 material, got %d", len(obj.Materials))
	}
	mat, ok := obj.Materials["Material001"]
	if !ok {
		t.Errorf("Material 'Material001' not found")
	}
	expectedKd := [3]float64{0.8, 0.8, 0.8}
	if mat.Diffuse != expectedKd {
		t.Errorf("Expected Diffuse %v, got %v", expectedKd, mat.Diffuse)
	}
	if mat.Texture != "texture.jpg" {
		t.Errorf("Expected Texture 'texture.jpg', got '%s'", mat.Texture)
	}
}

// TestParseOBJ_NegativeIndices tests parsing faces with negative indices.
func TestParseOBJ_NegativeIndices(t *testing.T) {
	objData := `
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
v 0.0 1.0 0.0
f -4 -3 -2 -1
`

	fsys := fstest.MapFS{
		"quad.obj": {Data: []byte(objData)},
	}

	obj, err := ParseFS(fsys, "quad.obj")
	if err != nil {
		t.Fatalf("Failed to parse OBJ file: %v", err)
	}

	if len(obj.Faces) != 1 {
		t.Errorf("Expected 1 face, got %d", len(obj.Faces))
	}

	// Check that negative indices are resolved correctly
	indices := obj.Faces[0].Indices
	expectedIndices := []int{1, 2, 3, 4}
	for i, idx := range indices {
		if idx.Vertex != expectedIndices[i] {
			t.Errorf("Expected vertex index %d, got %d", expectedIndices[i], idx.Vertex)
		}
	}
}

// TestParseOBJ_InvalidSyntax tests the parser's handling of invalid syntax.
func TestParseOBJ_InvalidSyntax(t *testing.T) {
	objData := `
v 0.0 0.0
f 1 2
`

	fsys := fstest.MapFS{
		"invalid.obj": {Data: []byte(objData)},
	}
	_, err := ParseFS(fsys, "invalid.obj")
	if err == nil {
		t.Fatal("Expected error for invalid OBJ data, got nil")
	}
}

// TestParseOBJ_EmptyFile tests parsing an empty OBJ file.
func TestParseOBJ_EmptyFile(t *testing.T) {
	fsys := fstest.MapFS{
		"empty.obj": {Data: []byte("")},
	}

	obj, err := ParseFS(fsys, "empty.obj")
	if err != nil {
		t.Fatalf("Failed to parse empty OBJ file: %v", err)
	}

	// Verify that the object has no data
	if len(obj.Vertices) != 0 {
		t.Errorf("Expected 0 vertices, got %d", len(obj.Vertices))
	}
	if len(obj.Faces) != 0 {
		t.Errorf("Expected 0 faces, got %d", len(obj.Faces))
	}
}

// ExampleParse demonstrates how to parse a Wavefront .obj file.
func ExampleParseFS() {
	objData := `
# Example OBJ data
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
f 1 2 3
`
	fsys := fstest.MapFS{
		"example.obj": {Data: []byte(objData)},
	}
	obj, err := ParseFS(fsys, "example.obj")
	if err != nil {
		fmt.Printf("Error parsing OBJ file: %v\n", err)
		return
	}

	fmt.Printf("Parsed OBJ with %d vertices and %d faces\n", len(obj.Vertices), len(obj.Faces))
	// Output: Parsed OBJ with 3 vertices and 1 faces
}

// BenchmarkParseOBJ_LargeWithMaterials benchmarks the Parse function with a large OBJ file
// that includes materials, using an in-memory filesystem.
func BenchmarkParseOBJ_LargeWithMaterials(b *testing.B) {
	// Number of vertices and faces for the large test example
	numVertices := 1000000 // Adjust as needed to make the test large
	numFaces := 3333333    // Assuming each face is a triangle

	// Generate the large OBJ data with material references
	objData, mtlData := generateLargeOBJWithMaterials(numVertices, numFaces)

	// Create an in-memory filesystem with the OBJ and MTL data
	fsys := fstest.MapFS{
		"large.obj":    &fstest.MapFile{Data: []byte(objData)},
		"material.mtl": &fstest.MapFile{Data: []byte(mtlData)},
		"texture.jpg":  &fstest.MapFile{Data: []byte("fake image data")}, // Placeholder texture
	}

	// Record the number of bytes processed per iteration for reporting
	b.SetBytes(int64(len(objData)))

	// Reset the benchmark timer after setup
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := ParseFS(fsys, "large.obj")
		if err != nil {
			b.Fatalf("Failed to parse OBJ data: %v", err)
		}
	}
}

// generateLargeOBJWithMaterials generates a large OBJ data string with the specified number
// of vertices and faces, including material references, and a corresponding MTL data string.
func generateLargeOBJWithMaterials(numVertices, numFaces int) (string, string) {
	var objBuilder strings.Builder
	var mtlBuilder strings.Builder

	// Seed the random number generator for reproducibility
	rand.Seed(1)

	// Define a simple material in the MTL file
	mtlName := "Material001"
	mtlBuilder.WriteString(fmt.Sprintf("newmtl %s\n", mtlName))
	mtlBuilder.WriteString("Kd 0.8 0.8 0.8\n")
	mtlBuilder.WriteString("Ka 0.2 0.2 0.2\n")
	mtlBuilder.WriteString("Ks 1.0 1.0 1.0\n")
	mtlBuilder.WriteString("Ns 100.0\n")
	mtlBuilder.WriteString("map_Kd texture.jpg\n")

	// Start the OBJ file with material references
	objBuilder.WriteString("mtllib material.mtl\n")
	objBuilder.WriteString(fmt.Sprintf("usemtl %s\n", mtlName))

	// Write vertices
	for i := 0; i < numVertices; i++ {
		x := rand.Float64()
		y := rand.Float64()
		z := rand.Float64()
		objBuilder.WriteString(fmt.Sprintf("v %f %f %f\n", x, y, z))
	}

	// Write faces using consecutive vertices
	for i := 1; i+2 <= numVertices && (i-1)/3 < numFaces; i += 3 {
		objBuilder.WriteString(fmt.Sprintf("f %d %d %d\n", i, i+1, i+2))
	}

	return objBuilder.String(), mtlBuilder.String()
}

// FuzzParseFS provides random inputs to the parser to check for panics or crashes.
func FuzzParseFS(f *testing.F) {
	// Seed corpus with a minimal valid OBJ file
	// and some edge-case or slightly malformed files if desired.
	f.Add(`
# Minimal valid OBJ file
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 0.0 1.0 0.0
f 1 2 3
`)

	f.Add(`
# OBJ file with a material
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
usemtl Material001
mtllib materials.mtl
f 1 2 3
`)

	f.Add(`
# Malformed OBJ file with incomplete vertex
v 0.0 0.0
f 1 2 3
`)

	f.Add(`
# OBJ file with negative indices
v 0.0 0.0 0.0
v 1.0 0.0 0.0
v 1.0 1.0 0.0
v 0.0 1.0 0.0
f -4 -3 -2 -1
`)

	f.Fuzz(func(t *testing.T, objContent string) {
		// Create an in-memory filesystem with the fuzzed OBJ content
		fsys := fstest.MapFS{
			"fuzzed.obj": {Data: []byte(objContent)},
		}

		// Attempt to parse the fuzzed OBJ file
		obj, err := ParseFS(fsys, "fuzzed.obj")

		// If an error occurs, it might be expected for malformed inputs.
		if err != nil {
			// Optionally, check if the error is recoverable or expected.
			t.Logf("ParseFS returned an error: %v", err)
			return
		}

		// Perform sanity checks on the parsed object
		// Example: Ensure that face indices are within bounds
		for _, face := range obj.Faces {
			for _, idx := range face.Indices {
				if idx.Vertex < 1 || idx.Vertex > len(obj.Vertices) {
					t.Errorf("Vertex index %d out of bounds (1 to %d)", idx.Vertex, len(obj.Vertices))
				}
				if idx.TexCoord > 0 && (idx.TexCoord < 1 || idx.TexCoord > len(obj.TexCoords)) {
					t.Errorf("Texture coordinate index %d out of bounds (1 to %d)", idx.TexCoord, len(obj.TexCoords))
				}
				if idx.Normal > 0 && (idx.Normal < 1 || idx.Normal > len(obj.Normals)) {
					t.Errorf("Normal index %d out of bounds (1 to %d)", idx.Normal, len(obj.Normals))
				}
			}
		}

	})
}
