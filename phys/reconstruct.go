package phys

// This file adds separable reconstruction filters (box, tent, Mitchell–Netravali)
// that can be applied to rendered frames. Using these filters is physically
// meaningful: they define the pixel footprint / reconstruction kernel rather
// than "denoising" the image. Applying them does not compromise physical accuracy
// as long as you consider the chosen kernel to be your sensor’s pixel response.

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// ReconFilter describes a 1D, even, compact-support reconstruction kernel.
// The kernel is separable; 2D weights are w(x)*w(y). Radius is in pixels.
type ReconFilter struct {
	// Name is a human-readable kernel name.
	Name string
	// Radius is the half-width of support (e.g., Tent=1, Mitchell=2).
	Radius float64
	// Eval returns w(|x|) for x in pixels. Implementations MUST return 0 for
	// |x| >= Radius (compact support).
	Eval func(x float64) float64
}

// BoxFilter returns a box (nearest) kernel with radius 0.5 (pixel average).
// This corresponds to uniform integration over a pixel footprint.
// Radius: 0.5, w(x)=1 for |x|<0.5 else 0.
func BoxFilter() ReconFilter {
	return ReconFilter{
		Name:   "Box(0.5)",
		Radius: 0.5,
		Eval: func(x float64) float64 {
			if math.Abs(x) < 0.5 {
				return 1
			}
			return 0
		},
	}
}

// TentFilter returns a triangular (linear) kernel with radius 1.0.
// Radius: 1, w(x)=max(0, 1-|x|).
func TentFilter() ReconFilter {
	return ReconFilter{
		Name:   "Tent(1)",
		Radius: 1,
		Eval: func(x float64) float64 {
			ax := math.Abs(x)
			if ax >= 1 {
				return 0
			}
			return 1 - ax
		},
	}
}

// MitchellNetravaliFilter returns the cubic Mitchell–Netravali kernel with
// the common B=C=1/3 parameters and radius 2. It gives crisp yet well-behaved
// results for Monte-Carlo images.
// Radius: 2.
func MitchellNetravaliFilter() ReconFilter {
	const B = 1.0 / 3.0
	const C = 1.0 / 3.0
	return ReconFilter{
		Name:   "Mitchell-Netravali(B=1/3,C=1/3)",
		Radius: 2,
		Eval: func(x float64) float64 {
			x = math.Abs(x)
			if x >= 2 {
				return 0
			}
			x2 := x * x
			x3 := x2 * x
			if x < 1 {
				return ((12-9*B-6*C)*x3 + (-18+12*B+6*C)*x2 + (6 - 2*B)) / 6.0
			}
			return ((-B-6*C)*x3 + (6*B+30*C)*x2 + (-12*B-48*C)*x + (8*B + 24*C)) / 6.0
		},
	}
}

// ApplySeparableFilterRGBA applies a separable reconstruction filter to src
// (assumed to be linear RGB) and returns a new RGBA image. Edges are clamped.
// This is a postprocess equivalent to reconstructing with that kernel.
// For perfect unbiasedness relative to the kernel, accumulate samples by
// splatting into a “film” with the same kernel; this postpass is a practical
// approximation that works well for uniform per-pixel sampling.
func ApplySeparableFilterRGBA(src *image.RGBA, f ReconFilter) *image.RGBA {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	// Horizontal pass to an intermediate linear float buffer.
	type px struct{ r, g, b float64 }
	tmp := make([]px, w*h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// Support: [x-R, x+R]
			x0 := int(math.Floor(float64(x) - f.Radius))
			x1 := int(math.Ceil(float64(x) + f.Radius))
			if x0 < 0 {
				x0 = 0
			}
			if x1 > w-1 {
				x1 = w - 1
			}
			var wr, wg, wb, wsum float64
			for xi := x0; xi <= x1; xi++ {
				wk := f.Eval(float64(x) - float64(xi))
				r8, g8, b8, _ := src.At(b.Min.X+xi, b.Min.Y+y).RGBA()
				// RGBA() is 16-bit per channel in Go; convert to [0,1].
				r := float64(r8) / 65535.0
				g := float64(g8) / 65535.0
				bl := float64(b8) / 65535.0

				wr += wk * r
				wg += wk * g
				wb += wk * bl
				wsum += wk
			}
			if wsum > 0 {
				wr /= wsum
				wg /= wsum
				wb /= wsum
			}
			tmp[y*w+x] = px{wr, wg, wb}
		}
	}

	// Vertical pass from tmp to dst.
	dst := image.NewRGBA(b)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			y0 := int(math.Floor(float64(y) - f.Radius))
			y1 := int(math.Ceil(float64(y) + f.Radius))
			if y0 < 0 {
				y0 = 0
			}
			if y1 > h-1 {
				y1 = h - 1
			}
			var wr, wg, wb, wsum float64
			for yi := y0; yi <= y1; yi++ {
				wk := f.Eval(float64(y) - float64(yi))
				p := tmp[yi*w+x]
				wr += wk * p.r
				wg += wk * p.g
				wb += wk * p.b
				wsum += wk
			}
			if wsum > 0 {
				wr /= wsum
				wg /= wsum
				wb /= wsum
			}
			// Clamp and convert back to 8-bit.
			r8 := uint8(math.Max(0, math.Min(255, 255.0*wr)))
			g8 := uint8(math.Max(0, math.Min(255, 255.0*wg)))
			b8 := uint8(math.Max(0, math.Min(255, 255.0*wb)))
			dst.Set(b.Min.X+x, b.Min.Y+y, color.RGBA{R: r8, G: g8, B: b8, A: 255})
		}
	}
	return dst
}

// ApplySeparableFilter is a convenience that accepts any image and returns a new
// RGBA after first copying the source (so callers may pass paletted images).
func ApplySeparableFilter(src image.Image, f ReconFilter) *image.RGBA {
	rgba := image.NewRGBA(src.Bounds())
	draw.Draw(rgba, rgba.Bounds(), src, src.Bounds().Min, draw.Src)
	return ApplySeparableFilterRGBA(rgba, f)
}
