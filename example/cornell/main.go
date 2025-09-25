// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

func main() {

}

// func main() {
// 	scene := phys.Scene{
// 		RenderOptions: phys.RenderOptions{
// 			Seed:         0,
// 			RaysPerPixel: 500,
// 			MaxRayDepth:     20,
// 			Dx:           256,
// 			Dy:           256,
// 		},
// 		Light: []phys.Light{
// 			phys.PointLight{
// 				Position: r3.Point{
// 					X: float64(200 * phys.MM),
// 					Y: float64(200 * phys.MM),
// 					Z: float64(200 * phys.MM)},
// 				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			},
// 			// phys.PointLight{
// 			// 	Position: r3.Point{
// 			// 		X: float64(0 * phys.NM),
// 			// 		Y: float64(-4 * phys.NM),
// 			// 		Z: float64(8 * phys.NM)},
// 			// 	RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			// },
// 			// phys.PointLight{
// 			// 	Position: r3.Point{
// 			// 		X: float64(4 * phys.NM),
// 			// 		Y: float64(-4 * phys.NM),
// 			// 		Z: float64(8 * phys.NM)},
// 			// 	RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
// 			// },
// 			// phys.PointLight{
// 			// 	Position: r3.Point{
// 			// 		X: float64(4 * phys.NM),
// 			// 		Y: float64(-4 * phys.NM),
// 			// 		Z: float64(-8 * phys.NM)},
// 			// 	RadiantIntensity: r3.Vec{X: 0.2, Y: 0.2, Z: 0.2},
// 			// },

// 			// phys.PointLight{
// 			// 	Position: r3.Point{
// 			// 		X: float64(-4 * phys.NM),
// 			// 		Y: float64(0 * phys.NM),
// 			// 		Z: float64(4 * phys.NM)},
// 			// 	RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.8},
// 			// },
// 		},
// 		Camera: []phys.Camera{
// 			phys.OrthographicCamera{
// 				LookFrom: r3.Point{
// 					X: float64(0 * phys.MM),
// 					Y: float64(0 * phys.MM),
// 					Z: float64(100 * phys.MM)},
// 				LookAt:    r3.Point{X: 0, Y: 0, Z: 0},
// 				VUp:       r3.Vec{X: 1, Y: 0, Z: 0},
// 				FOVHeight: 200 * phys.MM,
// 				FOVWidth:  200 * phys.MM,
// 			},
// 		},
// 		Node: []phys.Node{
// 			{
// 				Name:     "diffuse lambertian",
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: float64(-60 * phys.MM), Z: 0}, Radius: 20 * phys.MM},
// 				Material: phys.Lambertian{Texture: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}},
// 			},
// 			{
// 				Name:     "diffuse cosine lambertian",
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: float64(60 * phys.MM), Z: 0}, Radius: 20 * phys.MM},
// 				Material: phys.Lambertian{Texture: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}},
// 				// Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 			},
// 			{
// 				Name:     "frosty dielectric",
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 20 * phys.MM},
// 				Material: phys.Dielectric{Roughness: 0.1, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0},
// 				// Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 			},
// 			{
// 				Name:     "frosty metal",
// 				Shape:    phys.Sphere{Center: r3.Point{X: float64(-60 * phys.MM), Y: 0, Z: 0}, Radius: 20 * phys.MM},
// 				Material: phys.Metal{Fuzz: 0.1, Albedo: r3.Vec{X: 0.8, Y: 0.8, Z: 0.8}},
// 				// Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 			},
// 			{
// 				Name:     "indirect backlight 1",
// 				Shape:    phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: float64(200 * phys.MM)}, Radius: 60 * phys.MM},
// 				Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 				// Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 			},
// 			{
// 				Name:     "indirect backlight",
// 				Shape:    phys.Sphere{Center: r3.Point{X: float64(200 * phys.MM), Y: 0, Z: float64(200 * phys.MM)}, Radius: 60 * phys.MM},
// 				Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 				// Material: phys.Emitter{Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 			},
// 			// phys.PropAxes(r3.Point{}, phys.MM*2, phys.MM*50)[0],
// 			// phys.PropAxes(r3.Point{}, phys.MM*2, phys.MM*50)[1],
// 			// phys.PropAxes(r3.Point{}, phys.MM*2, phys.MM*50)[2],
// 			phys.PropSkybox(1*phys.M, phys.Emitter{Color: r3.Vec{X: 0.0, Y: 0.0, Z: 0.0}}),
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
// 	// Save another copy with the same filename so that for debugging the
// 	// image can be opened in one pane and automatically reloads when rendered.
// 	err = phys.SavePNG("./cornell.png", r.Image)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Saved to %s\n", path)
// 	// Render the stats path length image.
// 	path = time.Now().Format("./out/out_20060102_150405_path_length.png")
// 	// r.Stats.RaysPerPixel is a [][]uint32, we want this as an image.
// 	// Convert the path length to a grayscale image.
// 	statsImage := image.NewGray(image.Rect(0, 0, scene.RenderOptions.Dx, scene.RenderOptions.Dy))
// 	for x := 0; x < scene.RenderOptions.Dx; x++ {
// 		for y := 0; y < scene.RenderOptions.Dy; y++ {
// 			rays := r.Stats.RaysPerPixel[x][y]
// 			if rays > 255 {
// 				rays = 255
// 			}
// 			statsImage.SetGray(x, y, color.Gray{Y: uint8(rays)})
// 		}
// 	}
// 	err = phys.SavePNG(path, statsImage)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Saved to %s\n", path)
// }
