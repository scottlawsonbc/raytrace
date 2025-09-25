package phys

import (
	"encoding/json"
	"fmt"
	"math"
)

func init() {
	RegisterInterfaceType(TextureCheckerboard{})
}

// TextureCheckerboard represents a 2D procedural checkerboard texture.
// The zero value is not usable.
// Callers should assign valid textures and a positive Frequency, then call Validate.
//
// Frequency is an angular spatial frequency in radians per UV unit.
// The UV domain is dimensionless; 1.0 corresponds to one UV unit along each axis.
// The checker alternates whenever sin(Frequency*u) or sin(Frequency*v) changes sign.
// The edge-to-edge size of a single square is π/Frequency UV units.
type TextureCheckerboard struct {
	// Odd is the texture sampled for squares where the sign product of
	// sin(Frequency*u) and sin(Frequency*v) is negative.
	// Odd must be non-nil for a valid instance.
	Odd Texture

	// Even is the texture sampled for squares where the sign product of
	// sin(Frequency*u) and sin(Frequency*v) is non-negative.
	// Even must be non-nil for a valid instance.
	Even Texture

	// Frequency is the angular spatial frequency in radians per UV unit.
	// Larger values produce smaller squares; the square size is π/Frequency.
	// Frequency must be strictly positive.
	Frequency float64
}

// Validate reports whether tex has usable parameters.
//
// Validate returns nil if Odd and Even are non-nil and Frequency is strictly
// positive. Otherwise it returns a descriptive error. Validate does not mutate
// the receiver or its fields.
func (tex TextureCheckerboard) Validate() error {
	if tex.Odd == nil {
		return fmt.Errorf("error TextureCheckerboard.Odd Texture is nil")
	}
	if tex.Even == nil {
		return fmt.Errorf("error TextureCheckerboard.Even Texture is nil")
	}
	if tex.Frequency <= 0 {
		return fmt.Errorf("error TextureChecker.Frequency is negative: %v", tex.Frequency)
	}
	return nil
}

// At returns the spectrum at UV coordinates (u, v).
//
// At interprets u and v as dimensionless UV units. It computes
// sin(Frequency*u) * sin(Frequency*v); a negative product selects Odd, and a
// non-negative product selects Even. Boundaries occur along sin(...) == 0
// lines (i.e., at integer multiples of π/Frequency). At makes no assumptions
// about the range of u or v and does not clamp or wrap them.
func (tex TextureCheckerboard) At(u, v float64) Spectrum {
	sines := math.Sin(tex.Frequency*u) * math.Sin(tex.Frequency*v)
	if sines < 0 {
		return tex.Odd.At(u, v)
	}
	return tex.Even.At(u, v)
}

// MarshalJSON encodes a TextureCheckerboard as JSON.
//
// MarshalJSON includes a "Type" discriminator, the serialized Odd and Even
// textures (using their concrete encodings via marshalInterface), and the
// Frequency value. MarshalJSON returns an error if Odd or Even cannot be
// marshaled.
func (tex *TextureCheckerboard) MarshalJSON() ([]byte, error) {
	type TextureCheckerboardData struct {
		Type      string          `json:"Type"`
		Odd       json.RawMessage `json:"Odd"`
		Even      json.RawMessage `json:"Even"`
		Frequency float64         `json:"Frequency"`
	}
	oddData, err := marshalInterface(tex.Odd)
	if err != nil {
		return nil, err
	}
	evenData, err := marshalInterface(tex.Even)
	if err != nil {
		return nil, err
	}
	data := TextureCheckerboardData{
		Type:      "TextureCheckerboard",
		Odd:       oddData,
		Even:      evenData,
		Frequency: tex.Frequency,
	}
	return json.Marshal(data)
}

// UnmarshalJSON decodes a TextureCheckerboard from JSON.
//
// UnmarshalJSON expects a "Type" field with value "TextureCheckerboard",
// JSON-encoded Odd and Even textures (decoded via unmarshalInterface), and
// a Frequency value. UnmarshalJSON sets the receiver's fields accordingly.
// It returns an error if the type discriminator is wrong, if nested textures
// cannot be decoded, or if JSON is malformed. Callers may invoke Validate
// after unmarshaling to check semantic validity (e.g., positive Frequency).
func (tex *TextureCheckerboard) UnmarshalJSON(data []byte) error {
	type TextureCheckerboardData struct {
		Type      string          `json:"Type"`
		Odd       json.RawMessage `json:"Odd"`
		Even      json.RawMessage `json:"Even"`
		Frequency float64         `json:"Frequency"`
	}
	var temp TextureCheckerboardData
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if temp.Type != "TextureCheckerboard" {
		return fmt.Errorf("invalid type: expected TextureCheckerboard, got %s", temp.Type)
	}
	oddTexture, err := unmarshalInterface(temp.Odd)
	if err != nil {
		return err
	}
	evenTexture, err := unmarshalInterface(temp.Even)
	if err != nil {
		return err
	}
	tex.Odd = oddTexture.(Texture)
	tex.Even = evenTexture.(Texture)
	tex.Frequency = temp.Frequency
	return nil
}

// package phys

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// )

// // TextureCheckerboard represents a texture with a checkerboard pattern.
// type TextureCheckerboard struct {
// 	Odd       Texture
// 	Even      Texture
// 	Frequency float64
// }

// func (tex TextureCheckerboard) Validate() error {
// 	if tex.Odd == nil {
// 		return fmt.Errorf("error TextureCheckerboard.Odd Texture is nil")
// 	}
// 	if tex.Even == nil {
// 		return fmt.Errorf("error TextureCheckerboard.Even Texture is nil")
// 	}
// 	if tex.Frequency <= 0 {
// 		return fmt.Errorf("error TextureChecker.Frequency is negative: %v", tex.Frequency)
// 	}
// 	return nil
// }

// // At returns the color of the texture at the given UV coordinates.
// func (tex TextureCheckerboard) At(u, v float64) Spectrum {
// 	sines := math.Sin(tex.Frequency*u) * math.Sin(tex.Frequency*v)
// 	if sines < 0 {
// 		return tex.Odd.At(u, v)
// 	}
// 	return tex.Even.At(u, v)
// }

// // Implement custom JSON marshalling for TextureCheckerboard
// func (tex *TextureCheckerboard) MarshalJSON() ([]byte, error) {
// 	type TextureCheckerboardData struct {
// 		Type      string          `json:"Type"`
// 		Odd       json.RawMessage `json:"Odd"`
// 		Even      json.RawMessage `json:"Even"`
// 		Frequency float64         `json:"Frequency"`
// 	}
// 	oddData, err := marshalInterface(tex.Odd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	evenData, err := marshalInterface(tex.Even)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data := TextureCheckerboardData{
// 		Type:      "TextureCheckerboard",
// 		Odd:       oddData,
// 		Even:      evenData,
// 		Frequency: tex.Frequency,
// 	}
// 	return json.Marshal(data)
// }

// // Implement custom JSON unmarshalling for TextureCheckerboard
// func (tex *TextureCheckerboard) UnmarshalJSON(data []byte) error {
// 	type TextureCheckerboardData struct {
// 		Type      string          `json:"Type"`
// 		Odd       json.RawMessage `json:"Odd"`
// 		Even      json.RawMessage `json:"Even"`
// 		Frequency float64         `json:"Frequency"`
// 	}
// 	var temp TextureCheckerboardData
// 	if err := json.Unmarshal(data, &temp); err != nil {
// 		return err
// 	}
// 	if temp.Type != "TextureCheckerboard" {
// 		return fmt.Errorf("invalid type: expected TextureCheckerboard, got %s", temp.Type)
// 	}
// 	oddTexture, err := unmarshalInterface(temp.Odd)
// 	if err != nil {
// 		return err
// 	}
// 	evenTexture, err := unmarshalInterface(temp.Even)
// 	if err != nil {
// 		return err
// 	}
// 	tex.Odd = oddTexture.(Texture)
// 	tex.Even = evenTexture.(Texture)
// 	tex.Frequency = temp.Frequency
// 	return nil
// }
