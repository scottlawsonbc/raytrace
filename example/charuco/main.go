// Copyright 2024 Scott Lawson. All rights reserved.
package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"runtime"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/gl"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func main() {
	// All window/GL calls on main OS thread.
	runtime.LockOSThread()

	const winW, winH = 1440 / 2, 1080 / 2
	opts := gl.Options{
		Name:   "charuco",
		Width:  winW,
		Height: winH,
		X:      800,
		Y:      100,
	}
	win, err := gl.New(opts)
	if err != nil {
		panic(err)
	}
	defer win.Close()

	// -------- Scene constants --------
	var colormap = []phys.Spectrum{
		{X: 0.5, Y: 0.5, Z: 0.5},

		{X: 227.0 / 255.0, Y: 26.0 / 255.0, Z: 28.0 / 255.0},   // dark red
		{X: 251.0 / 255.0, Y: 154.0 / 255.0, Z: 153.0 / 255.0}, // light red

		{X: 51.0 / 255.0, Y: 160.0 / 255.0, Z: 44.0 / 255.0},   // dark green
		{X: 178.0 / 255.0, Y: 223.0 / 255.0, Z: 138.0 / 255.0}, // light green

		{X: 31.0 / 255.0, Y: 120.0 / 255.0, Z: 180.0 / 255.0},  // dark blue
		{X: 166.0 / 255.0, Y: 206.0 / 255.0, Z: 227.0 / 255.0}, // light blue

		{X: 1, Y: 127.0 / 255.0, Z: 0.0 / 255.0},               // dark orange
		{X: 253.0 / 255.0, Y: 191.0 / 255.0, Z: 111.0 / 255.0}, // light orange
	}
	// 203, 136, 206
	P := phys.Spectrum{X: 203.0 / 255.0, Y: 136.0 / 255.0, Z: 206.0 / 255.0} // purple

	// Box dimensions
	const W = phys.MM * 150
	const H = phys.MM * 150
	const D = phys.MM * 150

	colorX := phys.Spectrum{X: 1, Y: 0, Z: 157.0 / 255.0}
	colorY := phys.Spectrum{X: 157 / 255.0, Y: 1, Z: 0}
	colorZ := phys.Spectrum{X: 0 / 255.0, Y: 57.0 / 255.0, Z: 1}

	// Base scene (camera updated per frame)
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,    // keep realtime-ish; increase for quality
			MaxRayDepth:  6,    // modest recursion
			Dx:           winW, // match window size
			Dy:           winH,
		},
		Light: []phys.Light{
			phys.PointLight{
				Position: r3.Point{X: float64(500 * phys.MM), Y: float64(500 * phys.MM), Z: float64(500 * phys.MM)},
				RadiantIntensity: r3.Vec{
					X: 0.3, Y: 0.3, Z: 0.3,
				},
			},
		},
		Camera: []phys.Camera{}, // set each frame
		Node: []phys.Node{
			node("Origin", phys.Emitter{Texture: phys.TextureUniform{Color: colormap[0]}}, 0, 0, 0),
			node("P1", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, -float64(W)/2, -float64(H)/2, float64(D)),
			node("P2", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, float64(W)/2, -float64(H)/2, float64(D)),
			node("P3", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, float64(W)/2, float64(H)/2, float64(D)),
			node("P4", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, -float64(W)/2, float64(H)/2, float64(D)),

			node("P5", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, -float64(W)/2, -float64(H)/2, 0),
			node("P6", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, float64(W)/2, -float64(H)/2, 0),
			node("P7", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, float64(W)/2, float64(H)/2, 0),
			node("P8", phys.Emitter{Texture: phys.TextureUniform{Color: P}}, -float64(W)/2, float64(H)/2, 0),

			// Tiny axes at origin
			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, 0.5*phys.MM, 20*phys.MM, "")[0],
			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, 0.5*phys.MM, 20*phys.MM, "")[1],
			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, 0.5*phys.MM, 20*phys.MM, "")[2],

			// 12 rods outlining the box edges
			{
				Name:     "AxisX1Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorX}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: -float64(H) / 2, Z: float64(D)},
					Direction: r3.Vec{X: 1, Y: 0, Z: 0},
					Radius:    1 * phys.MM,
					Height:    W,
				},
			},
			{
				Name:     "AxisX2Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorX}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: float64(H) / 2, Z: float64(D)},
					Direction: r3.Vec{X: 1, Y: 0, Z: 0},
					Radius:    1 * phys.MM,
					Height:    W,
				},
			},
			{
				Name:     "AxisY1Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorY}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: -float64(H) / 2, Z: float64(D)},
					Direction: r3.Vec{X: 0, Y: 1, Z: 0},
					Radius:    1 * phys.MM,
					Height:    H,
				},
			},
			{
				Name:     "AxisY2Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorY}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: float64(W) / 2, Y: -float64(H) / 2, Z: float64(D)},
					Direction: r3.Vec{X: 0, Y: 1, Z: 0},
					Radius:    1 * phys.MM,
					Height:    H,
				},
			},
			{
				Name:     "AxisX1Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorX}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: -float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 1, Y: 0, Z: 0},
					Radius:    1 * phys.MM,
					Height:    W,
				},
			},
			{
				Name:     "AxisX2Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorX}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 1, Y: 0, Z: 0},
					Radius:    1 * phys.MM,
					Height:    W,
				},
			},
			{
				Name:     "AxisZ1Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorZ}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 0, Z: 1},
					Radius:    1 * phys.MM,
					Height:    D,
				},
			},
			{
				Name:     "AxisZ2Top",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorZ}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: float64(W) / 2, Y: float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 0, Z: 1},
					Radius:    1 * phys.MM,
					Height:    D,
				},
			},
			{
				Name:     "AxisZ1Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorZ}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: -float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 0, Z: 1},
					Radius:    1 * phys.MM,
					Height:    D,
				},
			},
			{
				Name:     "AxisZ2Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorZ}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: float64(W) / 2, Y: -float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 0, Z: 1},
					Radius:    1 * phys.MM,
					Height:    D,
				},
			},
			{
				Name:     "AxisY1Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorY}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: -float64(W) / 2, Y: -float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 1, Z: 0},
					Radius:    1 * phys.MM,
					Height:    H,
				},
			},
			{
				Name:     "AxisY2Bottom",
				Material: phys.Emitter{Texture: phys.TextureUniform{Color: colorY}},
				Shape: phys.Cylinder{
					Origin:    r3.Point{X: float64(W) / 2, Y: -float64(H) / 2, Z: 0},
					Direction: r3.Vec{X: 0, Y: 1, Z: 0},
					Radius:    1 * phys.MM,
					Height:    H,
				},
			},

			{
				Name:     "CharucoSquare150MM",
				Material: phys.Emitter{Texture: phys.MustNewTextureImage("./asset/Square150MM.png", "", "")},
				Shape: phys.Quad{
					Center: r3.Point{X: 0, Y: 0, Z: 0},
					Width:  150 * phys.MM,
					Height: 150 * phys.MM,
					Normal: r3.Vec{X: 0, Y: 0, Z: 1},
				},
			},
		},
	}

	fmt.Printf("bbox: %v\n", scene.Node[0].Shape.Bounds())

	// -------- Orbit parameters & cache --------
	const fps = 30
	const secondsPerOrbit = 8.0
	const frameDur = time.Second / fps
	framesPerOrbit := int(math.Round(secondsPerOrbit * fps))

	const orbitRadius = 400 * phys.MM // distance in XY plane
	const orbitHeight = 320 * phys.MM // Z height

	// Frame cache (filled once during first orbit)
	cache := make([]*image.RGBA, framesPerOrbit)

	// -------- Main loop --------
	ctx := context.Background()
	frameIndex := 0
	firstOrbitDone := false

	var fpsCounter int
	lastFPS := time.Now()

	first := true

	for !win.ShouldClose() {
		frameStart := time.Now()
		win.PollEvents()

		idx := frameIndex % framesPerOrbit
		var img image.Image

		if firstOrbitDone && cache[idx] != nil {
			// Playback from cache
			img = cache[idx]
		} else {
			// Compute the camera pose for this frame index (deterministic orbit)
			theta := 2 * math.Pi * (float64(idx) / float64(framesPerOrbit))
			lookFrom := r3.Point{
				X: float64(orbitRadius) * math.Cos(theta),
				Y: float64(orbitRadius) * math.Sin(theta),
				Z: float64(orbitHeight),
			}
			lookAt := r3.Point{X: 0, Y: 0, Z: float64(D) / 2}
			vup := r3.Vec{X: 0, Y: 0, Z: -1}

			scene.Camera = []phys.Camera{
				phys.CalibratedCamera{
					Intrinsics: phys.IntrinsicsFireflyDLGeneric6mm,
					Extrinsics: phys.CameraExtrinsics{
						LookFrom: lookFrom,
						LookAt:   lookAt,
						VUp:      vup,
					},
				},
			}

			// Render this frame
			res, err := phys.Render(ctx, &scene)
			if err != nil {
				fmt.Println("render error:", err)
				return
			}
			if first {
				path := time.Now().Format("./out/out_20060102_150405.png")
				err = phys.SavePNG(path, res.Image)
				if err != nil {
					log.Println("first frame save error:", err)
					panic(err)
				}
				log.Printf("✅ Saved first frame: %s\n", path)
				first = false
			}
			img = res.Image

			// Cache a copy on the first orbit
			if !firstOrbitDone {
				cache[idx] = cloneRGBA(img)
				if idx == framesPerOrbit-1 {
					firstOrbitDone = true
					fmt.Printf("✅ Cache complete: %d frames\n", framesPerOrbit)
				}
			}
		}

		// Draw the current frame
		win.Draw(img)

		// Simple FPS-ish report
		fpsCounter++
		if time.Since(lastFPS) >= time.Second {
			fmt.Printf("fps≈%d  mode=%s  frame=%d/%d\n",
				fpsCounter,
				map[bool]string{true: "playback", false: "render"}[firstOrbitDone],
				idx+1, framesPerOrbit)
			fpsCounter = 0
			lastFPS = time.Now()
		}

		// Pace to 30 fps (only matters if rendering/playback < 30fps)
		if dt := time.Since(frameStart); dt < frameDur {
			time.Sleep(frameDur - dt)
		}

		frameIndex++
	}
}

func node(name string, mat phys.Material, x, y, z float64) phys.Node {
	// Create a sphere with a radius of 4mm.
	radius := 4 * phys.MM
	return phys.Node{
		Name:     name,
		Shape:    phys.Sphere{Center: r3.Point{X: x, Y: y, Z: z}, Radius: radius},
		Material: mat,
	}
}

// cloneRGBA makes a deep *image.RGBA copy of an image.Image.
func cloneRGBA(src image.Image) *image.RGBA {
	r := src.Bounds()
	dst := image.NewRGBA(r)
	draw.Draw(dst, r, src, r.Min, draw.Src)
	return dst
}
