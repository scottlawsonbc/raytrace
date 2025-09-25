package phys

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"math"
	"os"
)

func init() {
	RegisterInterfaceType(TextureImage{})
}

// TextureImage represents a texture loaded from an image file.
type TextureImage struct {
	Image    image.Image
	FilePath string // Path to the image file (for serialization)
	Interp   string // Interpolation method: "nearest" or "bilinear"
	WrapMode string // Wrapping mode: "repeat" or "clamp"
}

func (it TextureImage) Validate() error {
	if it.Image == nil {
		return fmt.Errorf("image texture is nil")
	}
	return nil
}

// NewTextureImage loads an image from a file within the provided filesystem and returns a TextureImage.
func NewTextureImageFS(fsys fs.FS, filePath string, interp string, wrapMode string) (*TextureImage, error) {
	file, err := fsys.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return &TextureImage{
		Image:    img,
		FilePath: filePath,
		Interp:   interp,
		WrapMode: wrapMode,
	}, nil
}

func MustNewTextureImageFS(fsys fs.FS, filePath string, interp string, wrapMode string) *TextureImage {
	tex, err := NewTextureImageFS(fsys, filePath, interp, wrapMode)
	if err != nil {
		panic(err)
	}
	return tex
}

// NewTextureImage loads an image from a file and returns an ImageTexture.
func NewTextureImage(filePath string, interp string, wrapMode string) (*TextureImage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	return &TextureImage{
		Image:    img,
		FilePath: filePath,
		Interp:   interp,
		WrapMode: wrapMode,
	}, nil
}

func MustNewTextureImage(filePath string, interp string, wrapMode string) *TextureImage {
	tex, err := NewTextureImage(filePath, interp, wrapMode)
	if err != nil {
		panic(err)
	}
	return tex
}

// At returns the color value at the given UV coordinates.
func (it TextureImage) At(u, v float64) Spectrum {
	if it.Image == nil {
		// Return a default color if the image failed to load.
		return Spectrum{X: 1, Y: 0, Z: 1} // Magenta indicates missing texture.
	}
	// Handle wrapping modes.
	switch it.WrapMode {
	case "repeat":
		u = u - math.Floor(u)
		v = v - math.Floor(v)
	case "clamp":
		u = math.Min(math.Max(u, 0.0), 1.0)
		v = math.Min(math.Max(v, 0.0), 1.0)
	default:
		// Default to repeat
		u = u - math.Floor(u)
		v = v - math.Floor(v)
	}

	// Flip V coordinate to match image coordinate system.
	// TODO: scott what is the name of this conversion, uv to screen space?
	v = 1.0 - v

	// Convert UV coordinates to image coordinates
	width := it.Image.Bounds().Dx()
	height := it.Image.Bounds().Dy()
	x := u * float64(width-1)
	y := v * float64(height-1)

	var c color.Color

	switch it.Interp {
	case "bilinear":
		c = bilinearSample(it.Image, x, y)
	default:
		// Default to nearest neighbor
		ix := int(math.Round(x))
		iy := int(math.Round(y))
		c = it.Image.At(ix, iy)
	}

	// Convert color.Color to r3.Vec
	r, g, b, _ := c.RGBA()
	// Normalize to [0,1]
	return Spectrum{
		X: float64(r) / 65535.0,
		Y: float64(g) / 65535.0,
		Z: float64(b) / 65535.0,
	}
}

// bilinearSample performs bilinear interpolation on the image at (x, y).
func bilinearSample(img image.Image, x, y float64) color.Color {
	x0 := int(math.Floor(x))
	x1 := x0 + 1
	y0 := int(math.Floor(y))
	y1 := y0 + 1

	fx := x - float64(x0)
	fy := y - float64(y0)

	c00 := img.At(clamp(x0, 0, img.Bounds().Dx()-1), clamp(y0, 0, img.Bounds().Dy()-1))
	c10 := img.At(clamp(x1, 0, img.Bounds().Dx()-1), clamp(y0, 0, img.Bounds().Dy()-1))
	c01 := img.At(clamp(x0, 0, img.Bounds().Dx()-1), clamp(y1, 0, img.Bounds().Dy()-1))
	c11 := img.At(clamp(x1, 0, img.Bounds().Dx()-1), clamp(y1, 0, img.Bounds().Dy()-1))

	r := lerpColorComponent(c00, c10, c01, c11, fx, fy, 'R')
	g := lerpColorComponent(c00, c10, c01, c11, fx, fy, 'G')
	b := lerpColorComponent(c00, c10, c01, c11, fx, fy, 'B')

	return color.NRGBA64{
		R: uint16(r),
		G: uint16(g),
		B: uint16(b),
		A: 65535,
	}
}

func lerpColorComponent(c00, c10, c01, c11 color.Color, fx, fy float64, component rune) float64 {
	var c00c, c10c, c01c, c11c uint32
	switch component {
	case 'R':
		c00c, _, _, _ = c00.RGBA()
		c10c, _, _, _ = c10.RGBA()
		c01c, _, _, _ = c01.RGBA()
		c11c, _, _, _ = c11.RGBA()
	case 'G':
		_, c00c, _, _ = c00.RGBA()
		_, c10c, _, _ = c10.RGBA()
		_, c01c, _, _ = c01.RGBA()
		_, c11c, _, _ = c11.RGBA()
	case 'B':
		_, _, c00c, _ = c00.RGBA()
		_, _, c10c, _ = c10.RGBA()
		_, _, c01c, _ = c01.RGBA()
		_, _, c11c, _ = c11.RGBA()
	}

	c0 := float64(c00c)*(1-fx) + float64(c10c)*fx
	c1 := float64(c01c)*(1-fx) + float64(c11c)*fx
	return c0*(1-fy) + c1*fy
}
