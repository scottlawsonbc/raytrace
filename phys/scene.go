// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"encoding/json"
	"fmt"
)

type Scene struct {
	Camera        []Camera
	Node          []Node
	Light         []Light
	RenderOptions RenderOptions
}

func (s *Scene) Add(e ...Node) {
	s.Node = append(s.Node, e...)
}

func (s *Scene) Bounds() AABB {
	bounds := AABB{}
	for _, e := range s.Node {
		bounds = bounds.Union(e.Shape.Bounds())
	}
	return bounds
}

// Validate returns nil if the scene is capable of rendering, or an error if not.
func (s *Scene) Validate() error {
	// Verify valid RenderOptions.
	err := s.RenderOptions.Validate()
	if err != nil {
		return fmt.Errorf("bad RenderOptions=%v err=%v", s.RenderOptions, err)
	}
	// Verify at least one camera.
	if len(s.Camera) == 0 {
		return fmt.Errorf("no cameras in the scene")
	}
	// Verify nodes have unique names.
	names := make(map[string]bool)
	for _, e := range s.Node {
		err := e.Validate()
		if err != nil {
			return fmt.Errorf("node %s: %v", e.Name, err)
		}
		if names[e.Name] {
			return fmt.Errorf("duplicate node name: %s", e.Name)
		}
		names[e.Name] = true
	}
	// Verify nodes all have a material and a shape.
	for i, e := range s.Node {
		if e.Material == nil {
			return fmt.Errorf("node %d has no material", i)
		}
		if e.Shape == nil {
			return fmt.Errorf("node %d has no shape", i)
		}
	}
	// Verify at least one camera.
	if len(s.Camera) == 0 {
		return fmt.Errorf("no cameras in the scene")
	}
	// Verify all provided cameras as valid.
	for i, c := range s.Camera {
		err := c.Validate()
		if err != nil {
			return fmt.Errorf("bad camera %d: %v", i, err)
		}
	}
	// Verify all provided lights as valid.
	for i, l := range s.Light {
		err := l.Validate()
		if err != nil {
			return fmt.Errorf("bad light %d: %v", i, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface for Scene.
func (s Scene) MarshalJSON() ([]byte, error) {
	camera := make([]json.RawMessage, len(s.Camera))
	for i, c := range s.Camera {
		cameraJSON, err := marshalInterface(c)
		if err != nil {
			return nil, err
		}
		camera[i] = cameraJSON
	}
	node := make([]json.RawMessage, len(s.Node))
	for i, e := range s.Node {
		nodeJSON, err := e.MarshalJSON()
		if err != nil {
			return nil, err
		}
		node[i] = nodeJSON
	}
	light := make([]json.RawMessage, len(s.Light))
	for i, l := range s.Light {
		lightJSON, err := marshalInterface(l)
		if err != nil {
			return nil, err
		}
		light[i] = lightJSON
	}
	wrapped := map[string]interface{}{
		"Camera":        camera,
		"Node":          node,
		"Light":         light,
		"RenderOptions": s.RenderOptions,
	}
	return json.Marshal(wrapped)
}

// UnmarshalJSON implements the json.Unmarshaler interface for Scene.
func (s *Scene) UnmarshalJSON(data []byte) error {
	var wrapper struct {
		Camera        []json.RawMessage `json:"Camera"`
		Node          []json.RawMessage `json:"Node"`
		Light         []json.RawMessage `json:"Light"`
		RenderOptions RenderOptions     `json:"RenderOptions"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}
	for _, cRaw := range wrapper.Camera {
		iface, err := unmarshalInterface(cRaw)
		if err != nil {
			return err
		}
		cam, ok := iface.(Camera)
		if !ok {
			return err
		}
		s.Camera = append(s.Camera, cam)
	}
	for _, eRaw := range wrapper.Node {
		var e Node
		err := e.UnmarshalJSON(eRaw)
		if err != nil {
			return err
		}
		s.Node = append(s.Node, e)
	}
	for _, lRaw := range wrapper.Light {
		iface, err := unmarshalInterface(lRaw)
		if err != nil {
			return err
		}
		light, ok := iface.(Light)
		if !ok {
			return err
		}
		s.Light = append(s.Light, light)
	}
	s.RenderOptions = wrapper.RenderOptions
	return nil
}
