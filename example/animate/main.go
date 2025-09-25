package main

import (
	"context"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// lerpPoint returns the linear interpolation between a and b by t in [0,1].
func lerpPoint(a, b r3.Point, t float64) r3.Point {
	return r3.Point{
		X: a.X + (b.X-a.X)*t,
		Y: a.Y + (b.Y-a.Y)*t,
		Z: a.Z + (b.Z-a.Z)*t,
	}
}

// buildLinearDolly returns a phys.CameraFunc that produces a calibrated camera
// whose LookFrom linearly interpolates from p0 to p1 as u goes from 0 to 1.
// The function is pure and has no side effects.
func buildLinearDolly(
	intr phys.CameraIntrinsics,
	p0 r3.Point,
	p1 r3.Point,
	lookAt r3.Point,
	vup r3.Vec,
) phys.CameraFunc {
	return func(u float64) phys.Camera {
		uWrapped := u - math.Floor(u) // wrap into [0,1)
		lookFrom := lerpPoint(p0, p1, uWrapped)
		return phys.CalibratedCamera{
			Intrinsics: intr,
			Extrinsics: phys.CameraExtrinsics{
				LookFrom: lookFrom,
				LookAt:   lookAt,
				VUp:      vup,
			},
		}
	}
}

// newScene constructs a minimal scene with a lit checkerboard ground plane,
// origin axes, and a reference sphere. The function returns a validated
// phys.Scene using the provided render dimensions.
func newScene(dx, dy int) phys.Scene {
	ground := phys.Node{
		Name: "Ground",
		Material: phys.Lambertian{
			Texture: phys.TextureCheckerboard{
				Odd:       phys.TextureUniform{Color: phys.Spectrum{X: 0.176, Y: 0.404, Z: 0.671}},
				Even:      phys.TextureUniform{Color: phys.Spectrum{X: 0.8, Y: 0.8, Z: 0.8}},
				Frequency: 20,
			},
		},
		Shape: phys.Quad{
			Center: r3.Point{X: 0, Y: 0, Z: 0},
			Width:  400 * phys.MM,
			Height: 400 * phys.MM,
			Normal: r3.Vec{X: 0, Y: 0, Z: 1},
		},
	}

	refSphere := phys.Node{
		Name: "RefSphere",
		Material: phys.Lambertian{
			Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.671, Y: 0.196, Z: 0.176}},
		},
		Shape: phys.Sphere{
			Center: r3.Point{X: 0, Y: 0, Z: float64(25 * phys.MM)},
			Radius: 25 * phys.MM,
		},
	}

	axes := phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, 0.75*phys.MM, 200*phys.MM, "WorldAxes")

	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 32, // Increase for less noise.
			MaxRayDepth:  6,
			Dx:           dx,
			Dy:           dy,
		},
		Light: []phys.Light{
			phys.PointLight{
				Position:         r3.Point{X: 150 * float64(phys.MM), Y: 150 * float64(phys.MM), Z: 250 * float64(phys.MM)},
				RadiantIntensity: r3.Vec{X: 1.2, Y: 1.2, Z: 1.2},
			},
			phys.PointLight{
				Position:         r3.Point{X: -150 * float64(phys.MM), Y: -150 * float64(phys.MM), Z: 50 * float64(phys.MM)},
				RadiantIntensity: r3.Vec{X: .2, Y: .2, Z: .2},
			},
		},
		Camera: []phys.Camera{}, // filled per-frame
		Node:   append([]phys.Node{ground, refSphere}, axes...),
	}
	return scene
}

// palettize converts an RGBA frame to a Paletted frame suitable for GIF.
// The function uses Floyd–Steinberg dithering and the WebSafe palette.
func palettize(src *image.RGBA) *image.Paletted {
	dst := image.NewPaletted(src.Bounds(), palette.WebSafe)
	draw.FloydSteinberg.Draw(dst, dst.Rect, src, image.Point{})
	return dst
}

// main sets up an AnimatedCamera for a linear dolly and renders a short sequence.
// Each frame is saved as a before/after PNG and also appended to an animation GIF.
func main() {
	// Render settings.
	const (
		dx      = 1440 / 5
		dy      = 1080 / 5
		nFrames = 60
		fps     = 60 // frames per second in the output GIF
	)

	// Choose intrinsics (dimensions do not need to match dx,dy for this demo).
	intr := phys.IntrinsicsFireflyDLGeneric6mm

	// Define the dolly path and view.
	p0 := r3.Point{X: -250 * float64(phys.MM), Y: -200 * float64(phys.MM), Z: 150 * float64(phys.MM)}
	p1 := r3.Point{X: +250 * float64(phys.MM), Y: -200 * float64(phys.MM), Z: 150 * float64(phys.MM)}
	lookAt := r3.Point{X: 0, Y: 0, Z: 25 * float64(phys.MM)}
	vup := r3.Vec{X: 0, Y: 0, Z: -1}

	build := buildLinearDolly(intr, p0, p1, lookAt, vup)
	ac := phys.NewAnimatedCamera(build, 0, 2*time.Second) // 2 s cycle (helpers use this)

	// Prepare scene and output directory.
	scene := newScene(dx, dy)
	if err := os.MkdirAll("./out", 0o755); err != nil {
		log.Fatalf("failed to create out directory: %v", err)
	}

	// Prepare GIF container (side-by-side width).
	delayCS := int(math.Round(100.0 / float64(fps))) // delay in 1/100 s
	anim := &gif.GIF{
		Image:     make([]*image.Paletted, 0, nFrames),
		Delay:     make([]int, 0, nFrames),
		LoopCount: 0, // loop forever
	}

	ctx := context.Background()
	for i := 0; i < nFrames; i++ {
		u := float64(i) / float64(nFrames) // sample [0,1)
		scene.Camera = []phys.Camera{ac.WithU(u)}

		artifact, err := phys.Render(ctx, &scene)
		if err != nil {
			log.Fatalf("render failed at frame %d: %v", i, err)
		}

		// Append to GIF.
		anim.Image = append(anim.Image, palettize(artifact.Image))
		anim.Delay = append(anim.Delay, delayCS)

		fmt.Printf("rendered frame %d/%d\n", i+1, nFrames)
	}

	// Make animation loop by mirroring frames.
	for i := nFrames - 1; i >= 0; i-- {
		anim.Image = append(anim.Image, anim.Image[i])
		anim.Delay = append(anim.Delay, delayCS)
	}

	// Write animated GIF.
	outPath := "./out/animation.gif"
	f, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("failed to create %s: %v", outPath, err)
	}
	defer f.Close()

	if err := gif.EncodeAll(f, anim); err != nil {
		log.Fatalf("gif encode failed: %v", err)
	}
	fmt.Printf("✅ wrote %s (%d frames at ~%d fps)\n", filepath.ToSlash(outPath), len(anim.Image), fps)
}
