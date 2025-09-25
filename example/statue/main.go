// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

func main() {

}

// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
// )

// func main() {
// 	scene := phys.Scene{
// 		RenderOptions: phys.RenderOptions{
// 			Seed:         0,
// 			RaysPerPixel: 5,
// 			MaxRayDepth:     20,
// 			Dx:           1024,
// 			Dy:           1024,
// 		},
// 		Light: []phys.Light{
// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(0 * phys.NM),
// 					Y: float64(0 * phys.NM),
// 					Z: float64(3 * phys.NM)},
// 				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			},
// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(0 * phys.NM),
// 					Y: float64(-1 * phys.NM),
// 					Z: float64(3 * phys.NM)},
// 				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			},
// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(1 * phys.NM),
// 					Y: float64(-1 * phys.NM),
// 					Z: float64(3 * phys.NM)},
// 				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			},
// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(1 * phys.NM),
// 					Y: float64(-1 * phys.NM),
// 					Z: float64(-3 * phys.NM)},
// 				RadiantIntensity: r3.Vec{X: 0.2, Y: 0.2, Z: 0.2},
// 			},

// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(-3 * phys.NM),
// 					Y: float64(0 * phys.NM),
// 					Z: float64(3 * phys.NM)},
// 				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			},
// 		},
// 		Camera: []phys.Camera{
// 			phys.OrthographicCamera{
// 				LookFrom: r3.Point{
// 					X: float64(-50 * phys.NM),
// 					Y: float64(-100 * phys.NM),
// 					Z: float64(100 * phys.NM)},
// 				LookAt:    r3.Point{X: 0, Y: 0, Z: 0},
// 				VUp:       r3.Vec{X: 0, Y: 0, Z: 1},
// 				FOVHeight: 1 * phys.NM,
// 				FOVWidth:  1 * phys.NM,
// 			},
// 		},
// 		Node: []phys.Node{
// 			{

// 				Name:     "statue",
// 				Material: phys.ShaderNormal{},
// 				// Material: phys.Lambertian{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}},
// 				Shape: phys.LoadOBJ("../../asset/nude.obj"),
// 			},
// 			{
// 				Name: "behind camera light",
// 				// Material: phys.NormalShader{},
// 				Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 				Shape: phys.Sphere{
// 					Center: r3.Point{
// 						X: float64(-100 * phys.NM),
// 						Y: float64(-200 * phys.NM),
// 						Z: float64(200 * phys.NM),
// 					},
// 					Radius: 100 * phys.NM,
// 				},
// 			},
// 			phys.PropSkybox(1*phys.M, phys.Emitter{Color: r3.Vec{X: 0.1, Y: 0.1, Z: 0.1}}),
// 		},
// 	}
// 	// scene.Add(phys.PropAxes(r3.Point{}, 0.01*phys.NM, 1*phys.NM)...)
// 	// scene.Add(phys.PropSkybox(1 * phys.M))
// 	fmt.Printf("bbox: %v\n", scene.Node[0].Shape.Bounds())

// 	// Render the scene and save it to ./output.png
// 	r, err := phys.Render(&scene)
// 	if err != nil {
// 		panic(err)
// 	}
// 	path := time.Now().Format("./out/out_20060102_150405.png")
// 	err = phys.SavePNG(path, r.Image)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Rendered scene to output.png")
// }
