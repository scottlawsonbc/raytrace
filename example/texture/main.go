// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
//
// This example shows how to render a 3D scanned model and its color texture
// from a Wavefront .obj file and associated .mtl file.

package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func translate(x, y, z float64, s phys.Shape) phys.TransformedShape {
	return phys.TransformedShape{
		Shape: s,
		Transform: phys.Transform{
			Translation: r3.Vec{X: x, Y: y, Z: z},
			Rotation:    r3.RotationMatrixZ(0),
			Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
		},
	}
}

func rotate(angleDegrees float64, s phys.Shape) phys.TransformedShape {
	angleRadians := angleDegrees * math.Pi / 180
	return phys.TransformedShape{
		Shape: s,
		Transform: phys.Transform{
			Translation: r3.Vec{X: 0, Y: 0, Z: 0},
			Rotation:    r3.RotationMatrixZ(angleRadians),
			Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
		},
	}
}

func main() {
	prefabTriangle := phys.TriangleUV{
		P0:     r3.Point{X: 0, Y: 0, Z: 0},
		P1:     r3.Point{X: .1, Y: 0, Z: 0},
		P2:     r3.Point{X: 0, Y: .1, Z: 0},
		UV0:    r2.Point{X: 0, Y: 0},
		UV1:    r2.Point{X: 1, Y: 0},
		UV2:    r2.Point{X: 0, Y: 1},
		Normal: r3.Vec{X: 0, Y: 0, Z: 1},
	}
	prefabQuad := phys.Quad{
		Center: r3.Point{X: 0, Y: 0, Z: 0},
		Width:  0.1,
		Height: 0.1,
		Normal: r3.Vec{X: 0, Y: 0, Z: 1},
	}
	prefabSphere := phys.Sphere{
		Center: r3.Point{X: 0, Y: 0, Z: 0},
		Radius: 0.05,
	}

	prefabEmitterUV := phys.Emitter{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0, Y: 0, Z: 1}}}
	prefabTextureUV := phys.Emitter{Texture: phys.MustNewTextureImage("../../../../3d/debug/texture_debug_uv.png", "", "")}

	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,
			MaxRayDepth:  10,
			Dx:           1024,
			Dy:           1024,
		},
		Light: []phys.Light{},
		Camera: []phys.Camera{
			phys.OrthographicCamera{
				LookFrom:  r3.Point{X: 0.5, Y: 0.5, Z: 2.0},
				LookAt:    r3.Point{X: 0.5, Y: 0.5, Z: 0},
				VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
				FOVHeight: 1,
				FOVWidth:  1,
			},
		},

		Node: []phys.Node{
			{
				Name:     "triangle 1",
				Shape:    translate(0.1, 0.8, 0, rotate(0, prefabTriangle)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "triangle 1 origin",
				Shape:    translate(0.1, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 2",
				Shape:    translate(0.3, 0.8, 0, rotate(22.5, prefabTriangle)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "triangle 2 origin",
				Shape:    translate(0.3, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 3",
				Shape:    translate(0.5, 0.8, 0, rotate(45, prefabTriangle)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "triangle 3 origin",
				Shape:    translate(0.5, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 4",
				Shape:    translate(0.7, 0.8, 0, rotate(90, prefabTriangle)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "triangle 4 origin",
				Shape:    translate(0.7, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 5",
				Shape:    translate(0.8, 0.8, 0, rotate(135, prefabTriangle)),
				Material: phys.DebugUV{},
	},
			{
				Name:     "triangle 5 origin",
				Shape:    translate(0.8, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 6",
				Shape:    translate(0.9, 0.8, 0, rotate(0, prefabTriangle)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "triangle 6 origin",
				Shape:    translate(0.9, 0.8, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 21",
				Shape:    translate(0.1, 0.6, 0, rotate(0, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 21 origin",
				Shape:    translate(0.1, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 22",
				Shape:    translate(0.3, 0.6, 0, rotate(22.5, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 22 origin",
				Shape:    translate(0.3, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 23",
				Shape:    translate(0.5, 0.6, 0, rotate(45, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 23 origin",
				Shape:    translate(0.5, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 24",
				Shape:    translate(0.7, 0.6, 0, rotate(90, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 24 origin",
				Shape:    translate(0.7, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 25",
				Shape:    translate(0.8, 0.6, 0, rotate(135, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 25 origin",
				Shape:    translate(0.8, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "triangle 26",
				Shape:    translate(0.9, 0.6, 0, rotate(0, prefabTriangle)),
				Material: prefabTextureUV,
			},
			{
				Name:     "triangle 26 origin",
				Shape:    translate(0.9, 0.6, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 31",
				Shape:    translate(0.1, 0.4, 0, rotate(0, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 31 origin",
				Shape:    translate(0.1, 0.4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 32",
				Shape:    translate(0.3, .4, 0, rotate(22.5, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 32 origin",
				Shape:    translate(0.3, .4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 33",
				Shape:    translate(0.5, .4, 0, rotate(45, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 33 origin",
				Shape:    translate(0.5, .4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 34",
				Shape:    translate(0.7, .4, 0, rotate(90, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 34 origin",
				Shape:    translate(0.7, .4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 35",
				Shape:    translate(0.8, .4, 0, rotate(135, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 35 origin",
				Shape:    translate(0.8, .4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "sphere 36",
				Shape:    translate(0.9, .4, 0, rotate(0, prefabSphere)),
				Material: prefabTextureUV,
			},
			{
				Name:     "sphere 36 origin",
				Shape:    translate(0.9, .4, 0, phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 0.008}),
				Material: prefabEmitterUV,
			},
			{
				Name:     "quad 1",
				Shape:    translate(0.1, 0.1, 0, rotate(0, prefabQuad)),
				Material: phys.DebugUV{},
			},
			{
				Name:     "quad 2",
				Shape:    translate(0.2, 0.1, 0, rotate(0, prefabQuad)),
				Material: prefabTextureUV,
			},
			{
				Name:     "quad debug",
				Shape:    translate(0.3, 0.1, 0, rotate(0, prefabQuad)),
				Material: prefabTextureUV,
			},
			{
				Name:     "quad rotate",
				Shape:    translate(0.4, 0.1, 0, rotate(45, prefabQuad)),
				Material: prefabTextureUV,
			},
		},
	}

	// Render the scene and save it.
	r, err := phys.Render(context.Background(), &scene)
	if err != nil {
		panic(err)
	}
	path := time.Now().Format("./out/out_20060102_150405.png")
	err = phys.SavePNG(path, r.Image)
	if err != nil {
		panic(err)
	}
	// Save another copy with the same filename so that for debugging the
	// image can be opened in one pane and automatically reloads when rendered.
	err = phys.SavePNG("./texture.png", r.Image)
	if err != nil {
		panic(err)
	}
	log.Printf("Saved to ./texture.png\n")
}
