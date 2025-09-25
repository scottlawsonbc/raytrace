package phys

import (
	"encoding/json"
	"fmt"
)

func init() {
	RegisterInterfaceType(TextureUniform{})
}

// TextureUniform represents a texture with a single color.
type TextureUniform struct {
	Color Spectrum
}

// At returns the color of the texture at the given UV coordinates.
func (tex TextureUniform) At(u, v float64) Spectrum {
	return tex.Color
}

func (tex TextureUniform) Validate() error {
	return nil
}

// Implement custom JSON marshalling for TextureUniform
func (tu TextureUniform) MarshalJSON() ([]byte, error) {
	type TextureUniformData struct {
		Type  string   `json:"Type"`
		Color Spectrum `json:"Color"`
	}
	data := TextureUniformData{
		Type:  "TextureUniform",
		Color: tu.Color,
	}
	return json.Marshal(data)
}

// Implement custom JSON unmarshalling for TextureUniform
func (tu *TextureUniform) UnmarshalJSON(data []byte) error {
	type TextureUniformData struct {
		Type  string   `json:"Type"`
		Color Spectrum `json:"Color"`
	}
	var temp TextureUniformData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Type != "TextureUniform" {
		return fmt.Errorf("invalid type: expected TextureUniform, got %s", temp.Type)
	}
	tu.Color = temp.Color
	return nil
}
