// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// Emitter represents a material that emits light of a specific color.
type Emitter struct {
	Texture Texture
}

func (m Emitter) Validate() error {
	return m.Texture.Validate()
}

func (m Emitter) Resolve(ctx context.Context, c surfaceInteraction) resolution {
	e := m.Texture.At(c.collision.uv.X, c.collision.uv.Y)
	return resolution{emission: Spectrum(r3.Vec(e).Mul(r3.Vec(c.incoming.radiance)))}
}

func (m Emitter) ComputeDirectLighting(ctx context.Context, s surfaceInteraction, scene *Scene) Spectrum {
	// Emitters emit light but don't receive direct lighting.
	return Spectrum{}
}

// Implement custom JSON marshalling for Emitter
func (e *Emitter) MarshalJSON() ([]byte, error) {
	type EmitterData struct {
		Type    string          `json:"Type"`
		Texture json.RawMessage `json:"Texture"`
	}
	textureData, err := marshalInterface(e.Texture)
	if err != nil {
		return nil, err
	}
	data := EmitterData{
		Type:    "Emitter",
		Texture: textureData,
	}
	return json.Marshal(data)
}

// Implement custom JSON unmarshalling for Emitter
func (e *Emitter) UnmarshalJSON(data []byte) error {
	type EmitterData struct {
		Type    string          `json:"Type"`
		Texture json.RawMessage `json:"Texture"`
	}
	var temp EmitterData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Type != "Emitter" {
		return fmt.Errorf("invalid type: expected Emitter, got %s", temp.Type)
	}
	texture, err := unmarshalInterface(temp.Texture)
	if err != nil {
		return err
	}
	e.Texture = texture.(Texture)
	return nil
}

func init() {
	RegisterInterfaceType(Emitter{})
}
