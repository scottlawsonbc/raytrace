// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
//
// This example shows how to render a 3D scanned model and its color texture
// from a Wavefront .obj file and associated .mtl file.

package main

import (
	"context"
	"log"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// gotraceui command:  go run -ldflags="-H windowsgui" honnef.co/go/gotraceui/cmd/gotraceui@latest

func main() {
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
				Name: "triangle",
				Shape: phys.TriangleUV{
					P0:     r3.Point{X: 0, Y: 0, Z: 0},
					P1:     r3.Point{X: 1, Y: 0, Z: 0},
					P2:     r3.Point{X: 0, Y: 1, Z: 0},
					UV0:    r2.Point{X: 0, Y: 0},
					UV1:    r2.Point{X: 1, Y: 0},
					UV2:    r2.Point{X: 0, Y: 1},
					Normal: r3.Vec{X: 0, Y: 0, Z: 1},
				},
				Material: phys.DebugUV{},
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
	err = phys.SavePNG("./triangle.png", r.Image)
	if err != nil {
		panic(err)
	}
	log.Printf("Saved to ./triangle.png\n")
}
