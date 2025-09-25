// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/obj"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Node represents a physical object in the scene.
// It combines a geometric shape with a material that interacts with light.
// Analogous to the "primitive" concept in some ray tracing systems.
type Node struct {
	Name      string
	Transform Transform
	Shape     Shape
	Material  Material
}

func (n Node) Validate() error {
	if n.Name == "" {
		return fmt.Errorf("Node must have a name")
	}
	if n.Shape == nil {
		return fmt.Errorf("Node %q: missing Shape", n.Name)
	}
	if n.Material == nil {
		return fmt.Errorf("Node %q: missing Material", n.Name)
	}
	if err := n.Shape.Validate(); err != nil {
		return fmt.Errorf("Shape %q: %v", n.Name, err)
	}
	if err := n.Material.Validate(); err != nil {
		return fmt.Errorf("Material %q: %v", n.Name, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Node.
func (n Node) MarshalJSON() ([]byte, error) {
	shapeJSON, err := marshalInterface(n.Shape)
	if err != nil {
		return nil, err
	}
	materialJSON, err := marshalInterface(n.Material)
	if err != nil {
		return nil, err
	}
	wrapped := map[string]interface{}{
		"Name":     n.Name,
		"Shape":    shapeJSON,
		"Material": materialJSON,
	}
	return json.Marshal(wrapped)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Node.
func (n *Node) UnmarshalJSON(data []byte) error {
	var wrapper struct {
		Name     string          `json:"Name"`
		Shape    json.RawMessage `json:"Shape"`
		Material json.RawMessage `json:"Material"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	// Unmarshal Shape.
	iface, err := unmarshalInterface(wrapper.Shape)
	if err != nil {
		return err
	}
	shape, ok := iface.(Shape)
	if !ok {
		return fmt.Errorf("expected Shape, got %T", iface)
	}
	// Unmarshal Material.
	iface, err = unmarshalInterface(wrapper.Material)
	if err != nil {
		return err
	}
	material, ok := iface.(Material)
	if !ok {
		return fmt.Errorf("expected Material, got %T", iface)
	}
	n.Name = wrapper.Name
	n.Shape = shape
	n.Material = material
	return nil
}

// ConvertObjectToNodes converts an obj.Object into a slice of phys.Node.
// Each node corresponds to a mesh with a unique material.
func ConvertObjectToNodes(src *obj.Object, assetFS fs.FS) ([]Node, error) {
	// Convert materials.
	materials, err := ConvertObjectToMaterial(src, assetFS)
	if err != nil {
		return nil, err
	}

	// Group faces by material.
	materialToFaces := make(map[string][]obj.Face)
	for _, face := range src.Faces {
		mat := face.Material
		if mat == "" {
			mat = "default"
		}
		materialToFaces[mat] = append(materialToFaces[mat], face)
	}

	var nodes []Node
	for matName, faces := range materialToFaces {
		var meshFaces []Face
		for _, face := range faces {
			// Triangulate faces if necessary.
			if len(face.Indices) < 3 {
				continue // Skip degenerate faces.
			}
			// Triangulate polygonal faces using fan triangulation.
			for i := 1; i < len(face.Indices)-1; i++ {
				v0Index := face.Indices[0]
				v1Index := face.Indices[i]
				v2Index := face.Indices[i+1]
				v0, err := getVertexFromIndex(src, v0Index)
				if err != nil {
					return nil, err
				}
				v1, err := getVertexFromIndex(src, v1Index)
				if err != nil {
					return nil, err
				}
				v2, err := getVertexFromIndex(src, v2Index)
				if err != nil {
					return nil, err
				}
				meshFace := Face{
					Vertex: [3]Vertex{v0, v1, v2},
				}
				if err := meshFace.Validate(); err != nil {
					log.Printf("invalid face: %v", err)
					continue
				}
				meshFaces = append(meshFaces, meshFace)
			}
		}

		// Create a new Mesh with BVH.
		mesh, err := NewMesh(meshFaces)
		if err != nil {
			return nil, fmt.Errorf("failed to create mesh for material '%s': %v", matName, err)
		}

		// Get the corresponding phys.Material.
		material, exists := materials[matName]
		if !exists {
			material = materials["default"]
		}

		// Create Node.
		node := Node{Name: matName, Shape: mesh, Material: material}
		if err := node.Validate(); err != nil {
			return nil, fmt.Errorf("invalid node for material '%s': %v", matName, err)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// Helper function to get Vertex from obj.Index
func getVertexFromIndex(src *obj.Object, idx obj.Index) (Vertex, error) {
	vertexIndex := idx.Vertex - 1
	if vertexIndex < 0 || vertexIndex >= len(src.Vertices) {
		return Vertex{}, fmt.Errorf("vertex index out of range")
	}
	position := src.Vertices[vertexIndex]

	var uv r2.Point
	if idx.TexCoord > 0 && idx.TexCoord-1 < len(src.TexCoords) {
		texCoord := src.TexCoords[idx.TexCoord-1]
		uv = r2.Point{X: texCoord.U, Y: texCoord.V}
	} else {
		uv = r2.Point{X: 0, Y: 0}
	}

	return Vertex{
		Position: r3.Point{X: position.X, Y: position.Y, Z: position.Z},
		UV:       uv,
	}, nil
}

func ConvertObjectToShape(src *obj.Object) (mesh *Mesh, err error) {
	var meshFaces []Face
	for _, face := range src.Faces {
		// Triangulate faces if necessary.
		if len(face.Indices) < 3 {
			continue // Skip degenerate faces.
		}
		// Triangulate polygonal faces using fan triangulation.
		for i := 1; i < len(face.Indices)-1; i++ {
			// Create a Face.
			v0Index := face.Indices[0]
			v1Index := face.Indices[i]
			v2Index := face.Indices[i+1]
			v0, err := getVertexFromIndex(src, v0Index)
			if err != nil {
				return nil, err
			}
			v1, err := getVertexFromIndex(src, v1Index)
			if err != nil {
				return nil, err
			}
			v2, err := getVertexFromIndex(src, v2Index)
			if err != nil {
				return nil, err
			}
			meshFace := Face{Vertex: [3]Vertex{v0, v1, v2}}
			if err := meshFace.Validate(); err != nil {
				log.Printf("invalid face: %v", err)
				continue
			}
			meshFaces = append(meshFaces, meshFace)
		}
	}
	mesh = &Mesh{Face: meshFaces}
	return mesh, nil
}

func walk(fsys fs.FS, msg string) {
	fmt.Printf("Walking %v msg=%s\n", fsys, msg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			st, err := fs.Stat(fsys, p)
			if err != nil {
				return err
			}
			r, err := fsys.Open(p)
			if err != nil {
				return err
			}
			defer r.Close()

			// Read prefix
			var buf [md5.Size]byte
			n, _ := io.ReadFull(r, buf[:])

			// Hash remainder
			h := md5.New()
			_, err = io.Copy(h, r)
			if err != nil {
				return err
			}
			s := h.Sum(nil)
			log.Printf("| %s %d %x %x\n", p, st.Size(), buf[:n], s)
		}
		return nil
	})
}

// ConvertObjectToMaterial converts the materials defined in obj.Object into phys.Material instances.
// It returns a map from material names to phys.Material. This allows associating different
// parts of the geometry with their respective materials.
func ConvertObjectToMaterial(src *obj.Object, assetFS fs.FS) (map[string]Material, error) {
	materialMap := make(map[string]Material)
	for name, mat := range src.Materials {
		var texture Texture
		var err error
		if mat.Texture != "" {
			log.Printf("loading texture %s", mat.Texture)
			texturePath := mat.Texture
			texture, err = NewTextureImageFS(assetFS, texturePath, "bilinear", "repeat")
			if err != nil {
				walk(assetFS, "phys.ConvertObjectToMaterial.assetFS")
				return nil, fmt.Errorf("failed to load texture '%s' for material '%s': %v", texturePath, name, err)
			}
		} else {
			// Assign a uniform color if no texture is defined.
			r := mat.Diffuse[0]
			g := mat.Diffuse[1]
			b := mat.Diffuse[2]
			texture = TextureUniform{Color: Spectrum{X: r, Y: g, Z: b}}
		}
		m := Emitter{Texture: texture}
		if err := m.Validate(); err != nil {
			return nil, fmt.Errorf("invalid material '%s': %v", name, err)
		}

		materialMap[name] = m
	}

	// Handle the case where no materials are defined.
	if len(materialMap) == 0 {
		// Create a default Emitter material
		defaultMaterial := Emitter{
			Texture: TextureUniform{
				Color: Spectrum{X: 0.8, Y: 0.8, Z: 0.8}, // Default gray color.
			},
		}
		materialMap["default"] = defaultMaterial
	}

	return materialMap, nil
}

func (n Node) String() string {
	return fmt.Sprintf("Node{Name: %q, Shape: %v, Material: %v}", n.Name, n.Shape, n.Material)
}
