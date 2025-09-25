// Package gltf provides types and functions for working with glTF 2.0 assets.
package gltf

// Work in progress. Not yet implemented.

import (
	"encoding/json"
	"os"
)

// Asset represents the root object for a glTF asset.
type Asset struct {
	Asset              AssetInfo              `json:"asset"`
	Scene              *uint32                `json:"scene,omitempty"`
	Scenes             []*Scene               `json:"scenes,omitempty"`
	Nodes              []*Node                `json:"nodes,omitempty"`
	Meshes             []*Mesh                `json:"meshes,omitempty"`
	Accessors          []*Accessor            `json:"accessors,omitempty"`
	Animations         []*Animation           `json:"animations,omitempty"`
	Buffers            []*Buffer              `json:"buffers,omitempty"`
	BufferViews        []*BufferView          `json:"bufferViews,omitempty"`
	Cameras            []*Camera              `json:"cameras,omitempty"`
	Images             []*Image               `json:"images,omitempty"`
	Materials          []*Material            `json:"materials,omitempty"`
	Samplers           []*Sampler             `json:"samplers,omitempty"`
	Skins              []*Skin                `json:"skins,omitempty"`
	Textures           []*Texture             `json:"textures,omitempty"`
	ExtensionsUsed     []string               `json:"extensionsUsed,omitempty"`
	ExtensionsRequired []string               `json:"extensionsRequired,omitempty"`
	Extensions         map[string]interface{} `json:"extensions,omitempty"`
	Extras             interface{}            `json:"extras,omitempty"`
}

// AssetInfo represents metadata about the glTF asset.
/*
3.2. AssetInfo (Required)
Each glTF asset MUST have an asset property.
The asset object MUST contain a version property that specifies the target
glTF version of the asset. Additionally, an optional minVersion property
MAY be used to specify the minimum glTF version support required to load
the asset. The minVersion property allows asset creators to specify a
minimum version that a client implementation MUST support in order to load
the asset. This is very similar to the extensionsRequired concept described
in Section 3.12, where an asset SHOULD NOT be loaded if the client does not
support the specified extension. Additional metadata MAY be stored in
optional properties such as generator or copyright. For example,
*/
type AssetInfo struct {
	Version    string                 `json:"version"`
	Generator  string                 `json:"generator,omitempty"`
	MinVersion string                 `json:"minVersion,omitempty"`
	Copyright  string                 `json:"copyright,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

/*
3.5. Scenes
3.5.1. Overview
glTF 2.0 is used everywhere here.
Asset MAY contain zero or more scenes, the set of visual objects to render.
Scenes are defined in a scenes array.
All nodes listed in scene.nodes array MUST be root nodes, i.e., they MUST NOT be listed in a node.children array of any node. The same root node MAY appear in multiple scenes.

An additional root-level property, scene (note singular),
identifies which of the scenes in the array SHOULD be displayed at load time.
When scene is undefined, client implementations MAY delay rendering until a particular scene is requested.

A glTF asset that does not contain any scenes SHOULD be treated as a library
of individual nodes such as materials or meshes.

The following example defines a glTF asset with a single scene that contains a single node.

{
    "nodes": [
        {
            "name": "singleNode"
        }
    ],
    "scenes": [
        {
            "name": "singleScene",
            "nodes": [
                0
            ]
        }
    ],
    "scene": 0
}
*/
// Scene represents the root nodes of a scene.
type Scene struct {
	Nodes      []uint32               `json:"nodes,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

/*
3.5.2. Nodes and Hierarchy
glTF assets MAY define nodes, that is, the objects comprising the scene to render.

Nodes MAY have transform properties, as described later.

Nodes are organized in a parent-child hierarchy known informally as the node hierarchy. A node is called a root node when it doesn’t have a parent.

The node hierarchy MUST be a set of disjoint strict trees. That is node hierarchy MUST NOT contain cycles and each node MUST have zero or one parent node.

The node hierarchy is defined using a node’s children property, as in the following example:

{
    "nodes": [
        {
            "name": "Car",
            "children": [1, 2, 3, 4]
        },
        {
            "name": "wheel_1"
        },
        {
            "name": "wheel_2"
        },
        {
            "name": "wheel_3"
        },
        {
            "name": "wheel_4"
        }
    ]
}
*/
// Node represents a node in the node hierarchy.
type Node struct {
	Camera      *uint32                `json:"camera,omitempty"`
	Children    []uint32               `json:"children,omitempty"`
	Skin        *uint32                `json:"skin,omitempty"`
	Matrix      []float64              `json:"matrix,omitempty"` // Length must be 16
	Mesh        *uint32                `json:"mesh,omitempty"`
	Rotation    []float64              `json:"rotation,omitempty"`    // Length must be 4
	Scale       []float64              `json:"scale,omitempty"`       // Length must be 3
	Translation []float64              `json:"translation,omitempty"` // Length must be 3
	Weights     []float64              `json:"weights,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
	Extras      interface{}            `json:"extras,omitempty"`
}

// Mesh represents a set of primitives to be rendered.
type Mesh struct {
	Primitives []*MeshPrimitive       `json:"primitives"`
	Weights    []float64              `json:"weights,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// MeshPrimitive represents geometry to be rendered with the given material.
type MeshPrimitive struct {
	Attributes map[string]uint32      `json:"attributes"`
	Indices    *uint32                `json:"indices,omitempty"`
	Material   *uint32                `json:"material,omitempty"`
	Mode       *uint32                `json:"mode,omitempty"`
	Targets    []map[string]uint32    `json:"targets,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Accessor represents a typed view into a bufferView.
type Accessor struct {
	BufferView    *uint32                `json:"bufferView,omitempty"`
	ByteOffset    uint32                 `json:"byteOffset,omitempty"`
	ComponentType uint32                 `json:"componentType"`
	Normalized    bool                   `json:"normalized,omitempty"`
	Count         uint32                 `json:"count"`
	Type          string                 `json:"type"`
	Max           []float64              `json:"max,omitempty"`
	Min           []float64              `json:"min,omitempty"`
	Sparse        *AccessorSparse        `json:"sparse,omitempty"`
	Name          string                 `json:"name,omitempty"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
	Extras        interface{}            `json:"extras,omitempty"`
}

// AccessorSparse represents sparse storage of accessor values that deviate from their initialization value.
type AccessorSparse struct {
	Count      uint32                 `json:"count"`
	Indices    AccessorSparseIndices  `json:"indices"`
	Values     AccessorSparseValues   `json:"values"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// AccessorSparseIndices represents indices of accessor sparse values.
type AccessorSparseIndices struct {
	BufferView    uint32                 `json:"bufferView"`
	ByteOffset    uint32                 `json:"byteOffset,omitempty"`
	ComponentType uint32                 `json:"componentType"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
	Extras        interface{}            `json:"extras,omitempty"`
}

// AccessorSparseValues represents the accessor sparse values.
type AccessorSparseValues struct {
	BufferView uint32                 `json:"bufferView"`
	ByteOffset uint32                 `json:"byteOffset,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Animation represents a keyframe animation.
type Animation struct {
	Channels   []*AnimationChannel    `json:"channels"`
	Samplers   []*AnimationSampler    `json:"samplers"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// AnimationChannel combines an animation sampler with a target property being animated.
type AnimationChannel struct {
	Sampler    uint32                 `json:"sampler"`
	Target     AnimationChannelTarget `json:"target"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// AnimationChannelTarget describes the animated property.
type AnimationChannelTarget struct {
	Node       *uint32                `json:"node,omitempty"`
	Path       string                 `json:"path"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// AnimationSampler combines timestamps with output values and defines an interpolation algorithm.
type AnimationSampler struct {
	Input         uint32                 `json:"input"`
	Interpolation string                 `json:"interpolation,omitempty"`
	Output        uint32                 `json:"output"`
	Extensions    map[string]interface{} `json:"extensions,omitempty"`
	Extras        interface{}            `json:"extras,omitempty"`
}

// Buffer points to binary geometry, animation, or skins.
type Buffer struct {
	URI        string                 `json:"uri,omitempty"`
	ByteLength uint32                 `json:"byteLength"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// BufferView represents a view into a buffer generally representing a subset of the buffer.
type BufferView struct {
	Buffer     uint32                 `json:"buffer"`
	ByteOffset uint32                 `json:"byteOffset,omitempty"`
	ByteLength uint32                 `json:"byteLength"`
	ByteStride uint32                 `json:"byteStride,omitempty"`
	Target     uint32                 `json:"target,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Camera represents a camera's projection. A node may reference a camera to apply a transform.
type Camera struct {
	Orthographic *CameraOrthographic    `json:"orthographic,omitempty"`
	Perspective  *CameraPerspective     `json:"perspective,omitempty"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty"`
	Extras       interface{}            `json:"extras,omitempty"`
}

// CameraOrthographic contains properties to create an orthographic projection matrix.
type CameraOrthographic struct {
	XMag       float64                `json:"xmag"`
	YMag       float64                `json:"ymag"`
	ZFar       float64                `json:"zfar"`
	ZNear      float64                `json:"znear"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// CameraPerspective contains properties to create a perspective projection matrix.
type CameraPerspective struct {
	AspectRatio *float64               `json:"aspectRatio,omitempty"`
	YFov        float64                `json:"yfov"`
	ZFar        *float64               `json:"zfar,omitempty"`
	ZNear       float64                `json:"znear"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
	Extras      interface{}            `json:"extras,omitempty"`
}

// Image represents image data used to create a texture.
type Image struct {
	URI        string                 `json:"uri,omitempty"`
	MimeType   string                 `json:"mimeType,omitempty"`
	BufferView *uint32                `json:"bufferView,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Material defines the material appearance of a primitive.
type Material struct {
	Name                 string                 `json:"name,omitempty"`
	PBRMetallicRoughness *PBRMetallicRoughness  `json:"pbrMetallicRoughness,omitempty"`
	NormalTexture        *NormalTextureInfo     `json:"normalTexture,omitempty"`
	OcclusionTexture     *OcclusionTextureInfo  `json:"occlusionTexture,omitempty"`
	EmissiveTexture      *TextureInfo           `json:"emissiveTexture,omitempty"`
	EmissiveFactor       []float64              `json:"emissiveFactor,omitempty"` // Length must be 3
	AlphaMode            string                 `json:"alphaMode,omitempty"`
	AlphaCutoff          float64                `json:"alphaCutoff,omitempty"`
	DoubleSided          bool                   `json:"doubleSided,omitempty"`
	Extensions           map[string]interface{} `json:"extensions,omitempty"`
	Extras               interface{}            `json:"extras,omitempty"`
}

// PBRMetallicRoughness defines the metallic-roughness material model.
type PBRMetallicRoughness struct {
	BaseColorFactor          []float64              `json:"baseColorFactor,omitempty"` // Length must be 4
	BaseColorTexture         *TextureInfo           `json:"baseColorTexture,omitempty"`
	MetallicFactor           float64                `json:"metallicFactor,omitempty"`
	RoughnessFactor          float64                `json:"roughnessFactor,omitempty"`
	MetallicRoughnessTexture *TextureInfo           `json:"metallicRoughnessTexture,omitempty"`
	Extensions               map[string]interface{} `json:"extensions,omitempty"`
	Extras                   interface{}            `json:"extras,omitempty"`
}

// TextureInfo references a texture.
type TextureInfo struct {
	Index      uint32                 `json:"index"`
	TexCoord   uint32                 `json:"texCoord,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// NormalTextureInfo is the normal texture info.
type NormalTextureInfo struct {
	TextureInfo
	Scale float64 `json:"scale,omitempty"`
}

// OcclusionTextureInfo is the occlusion texture info.
type OcclusionTextureInfo struct {
	TextureInfo
	Strength float64 `json:"strength,omitempty"`
}

// Sampler contains properties for texture filtering and wrapping modes.
type Sampler struct {
	MagFilter  *uint32                `json:"magFilter,omitempty"`
	MinFilter  *uint32                `json:"minFilter,omitempty"`
	WrapS      uint32                 `json:"wrapS,omitempty"`
	WrapT      uint32                 `json:"wrapT,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Skin defines joints and matrices.
type Skin struct {
	InverseBindMatrices *uint32                `json:"inverseBindMatrices,omitempty"`
	Skeleton            *uint32                `json:"skeleton,omitempty"`
	Joints              []uint32               `json:"joints"`
	Name                string                 `json:"name,omitempty"`
	Extensions          map[string]interface{} `json:"extensions,omitempty"`
	Extras              interface{}            `json:"extras,omitempty"`
}

// Texture and its sampler.
type Texture struct {
	Sampler    *uint32                `json:"sampler,omitempty"`
	Source     *uint32                `json:"source,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Extras     interface{}            `json:"extras,omitempty"`
}

// Load reads a glTF asset from a file.
func Load(path string) (*Asset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var doc Asset
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// Save writes the glTF asset to a file.
func (d *Asset) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(d)
}
