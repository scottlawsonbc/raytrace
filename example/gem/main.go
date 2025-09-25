// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

func main() {

}

// package main

// import (
// 	"fmt"

// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
// )

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func optimizePAL(old []phys.Node) (new []phys.Node) {
// 	on := []phys.Node{}
// 	offShapes := []phys.Shape{}
// 	for n, e := range old {
// 		if m, ok := e.Material.(phys.Emitter); ok {
// 			if m.Color == (r3.Vec{X: 0, Y: 0, Z: 0}) {
// 				offShapes = append(offShapes, e.Shape)
// 			} else {
// 				on = append(on, old[n])
// 			}
// 		}
// 	}
// 	offNode := phys.Node{
// 		Shape:    phys.NewBVH(offShapes, 2),
// 		Material: phys.Emitter{Texture: phys.TextureUniform{Color: r3.Vec{X: 0, Y: 0, Z: 0}}},
// 	}
// 	new = append(on, offNode)
// 	return
// }

// func importPAL(posz float64) (leds []phys.Node) {
// 	for _, led := range LEDS {
// 		leds = append(leds, phys.Node{
// 			Shape: phys.Sphere{
// 				Center: r3.Point{
// 					X: led.POSX * float64(phys.MM),
// 					Y: led.POSY * float64(phys.MM),
// 					Z: posz,
// 				},
// 				Radius: phys.MM * 1,
// 			},

// 			Material: phys.Emitter{Texture: phys.TextureUniform{Color: r3.Vec{X: 0.2, Y: 0, Z: 0}}},

// 		})
// 	}
// 	return
// }

// func main() {
// 	// colormap3 := map[int]r3.Vec{
// 	// 	1: {X: 228 / 255.0, Y: 26 / 255.0, Z: 28 / 255.0},
// 	// 	2: {X: 55 / 255.0, Y: 126 / 255.0, Z: 184 / 255.0},
// 	// 	3: {X: 77 / 255.0, Y: 175 / 255.0, Z: 74 / 255.0},
// 	// 	4: {X: 152 / 255.0, Y: 78 / 255.0, Z: 163 / 255.0},
// 	// 	5: {X: 255 / 255.0, Y: 127 / 255.0, Z: 0 / 255.0},
// 	// 	6: {X: 255 / 255.0, Y: 255 / 255.0, Z: 51 / 255.0},
// 	// 	7: {X: 166 / 255.0, Y: 86 / 255.0, Z: 40 / 255.0},
// 	// 	8: {X: 247 / 255.0, Y: 129 / 255.0, Z: 191 / 255.0},
// 	// }
// 	// leds := importPAL(0)
// 	// frames := []image.Image{}
// 	// for pcbn := 1; pcbn <= 257; pcbn += 1 {
// 	// 	fmt.Println("Rendering PCBN", pcbn)
// 	// 	for bitn, led := range LEDS {
// 	// 		if led.PCBN == pcbn {
// 	// 			leds[bitn].Material = phys.Light{
// 	// 				Color: colormap3[led.RING],
// 	// 			}
// 	// 		} else {
// 	// 		}
// 	// 	}
// 	// LookFrom:  //r3.Point{X: 100 + float64(float64(pcbn)*-3*float64(phys.MM)), Y: 300 + float64(float64(pcbn)*-2*float64(phys.MM)), Z: float64(500 * phys.MM)},

// 	var scene = phys.Scene{
// 		Camera: phys.OrthographicCamera{
// 			LookFrom:  r3.Point{X: 0, Y: 0, Z: float64(500 * phys.MM)},
// 			LookAt:    r3.Point{X: 0, Y: 0, Z: 0},
// 			VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
// 			FOVHeight: 125 * phys.MM,
// 			FOVWidth:  125 * phys.MM,
// 		},
// 		Node: []phys.Node{
// 			{
// 				Shape: phys.Cylinder{
// 					Origin:    r3.Point{X: 0, Y: 0, Z: 0},
// 					Direction: r3.Vec{X: 1, Y: 0, Z: 0},
// 					Radius:    0.25 * phys.MM,
// 					Height:    50 * phys.M,
// 				},
// 				Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 0, Z: 0}},
// 			},
// 			{
// 				Shape: phys.Cylinder{
// 					Origin:    r3.Point{X: 0, Y: 0, Z: 0},
// 					Direction: r3.Vec{X: 0, Y: 1, Z: 0},
// 					Radius:    0.25 * phys.MM,
// 					Height:    50 * phys.M,
// 				},
// 				Material: phys.Emitter{Color: r3.Vec{X: 0, Y: 1, Z: 0}},
// 			},
// 			{
// 				Shape: phys.Cylinder{
// 					Origin:    r3.Point{X: 0, Y: 0, Z: 0},
// 					Direction: r3.Vec{X: 0, Y: 0, Z: 1},
// 					Radius:    0.25 * phys.MM,
// 					Height:    50 * phys.M,
// 				},
// 				Material: phys.Emitter{Color: r3.Vec{X: 0, Y: 0, Z: 1}},
// 			},
// 			{
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 1 * phys.M},
// 				Material: phys.Emitter{Color: r3.Vec{X: 0.1, Y: 0.1, Z: 0.1}},
// 			},
// 			{
// 				Shape:    phys.Sphere{Center: r3.Point{X: float64(-60 * phys.MM), Y: float64(60 * phys.MM), Z: float64(50 * phys.MM)}, Radius: 40 * phys.MM},
// 				Material: phys.ShaderNormal{},
// 			},
// 			{
// 				Shape:    phys.Sphere{Center: r3.Point{X: float64(60 * phys.MM), Y: float64(60 * phys.MM), Z: float64(50 * phys.MM)}, Radius: 40 * phys.MM},
// 				Material: phys.ShaderNormal{},
// 			},
// 			{
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 10 * phys.MM},
// 				Material: phys.ShaderNormal{},
// 			},
// 		},
// 		RenderOptions: phys.RenderOptions{
// 			Seed:         0,
// 			RaysPerPixel: 4,
// 			RayDepth:     10,
// 			Dx:           256,
// 			Dy:           256,
// 		},
// 	}
// 	// scene.Add(optimizePAL(leds)...)
// 	phys.Render(&scene)
// 	// frames = append(frames, img)
// 	s := scene.JSON()
// 	fmt.Printf("\n%+v\n\n", s)
// 	// s_backwards, err := phys.NewSceneFromJSON(s)
// 	// check(err)
// 	// fmt.Printf("scene json: %+v\n", s_backwards.JSON())

// 	// fmt.Printf("scene json: %+v img.Dx: %+v img.Dy: %+v\n", s, img.Bounds().Dx(), img.Bounds().Dy())

// 	// 	if pcbn > 0 {
// 	// 		break
// 	// 	}
// 	// // }
// 	// g := phys.NewGIF(frames)
// 	// path := time.Now().Format("./out/out_20060102_150405.gif")
// 	// f, err := os.Create(path)
// 	// check(err)
// 	// err = gif.EncodeAll(f, g)
// 	// check(err)
// 	// f.Close()
// 	// fmt.Println("Saved to", path)
// }
