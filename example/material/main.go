// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

// This example scene demonstrates an array of different materials.

func node(col, row int, name string, mat phys.Material) phys.Node {
	dx := float64(10 * phys.MM)
	dy := float64(10 * phys.MM)
	cx := float64(col)*dx - float64(45*phys.MM)
	cy := float64(row)*dy - float64(25*phys.MM)
	radius := 4 * phys.MM
	return phys.Node{
		Name:     name,
		Shape:    phys.Sphere{Center: r3.Point{X: cx, Y: cy, Z: 0}, Radius: radius},
		Material: mat,
	}
}

func main() {
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 64,
			MaxRayDepth:  5,
			Dx:           1024,
			Dy:           1024,
		},
		Light: []phys.Light{
			phys.PointLight{
				Position: r3.Point{
					X: float64(1000 * phys.MM),
					Y: float64(1000 * phys.MM),
					Z: float64(1000 * phys.MM)},
				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
			},
		},
		Camera: []phys.Camera{
			phys.OrthographicCamera{
				LookFrom: r3.Point{
					X: float64(40 * phys.MM),
					Y: float64(50 * phys.MM),
					Z: float64(400 * phys.MM)},
				LookAt:    r3.Point{X: 0, Y: 0, Z: 0},
				VUp:       r3.Vec{X: 1, Y: 0, Z: 0},
				FOVHeight: 110 * phys.MM,
				FOVWidth:  110 * phys.MM,
			},
		},
		Node: []phys.Node{

			// Lambertian materials.
			node(0, 0, "checker red white", phys.Lambertian{Texture: phys.TextureCheckerboard{
				Odd:       phys.TextureUniform{Color: phys.Spectrum{X: 0.2, Y: 0.5, Z: 0.6}},
				Even:      phys.TextureUniform{Color: phys.Spectrum{X: 0.5, Y: 0.0, Z: 0.0}},
				Frequency: 5,
			}}),
			node(1, 0, "checker complementary colors", phys.Lambertian{Texture: phys.TextureCheckerboard{
				Odd:       phys.TextureUniform{Color: phys.Spectrum{X: 0, Y: 0, Z: 123}.Divs(255)},
				Even:      phys.TextureUniform{Color: phys.Spectrum{X: 242, Y: 0, Z: 0}.Divs(255)},
				Frequency: 5,
			}}),

			node(2, 0, "image texture", phys.Lambertian{Texture: phys.MustNewTextureImage("./texture.png", "bilinear", "repeat")}),
			node(3, 0, "face texture", phys.Lambertian{Texture: phys.MustNewTextureImage("./faces.png", "bilinear", "repeat")}),

			node(0, 1, "Lambertian min gray", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.05, Y: 0.05, Z: 0.05}}}),
			node(0, 2, "Lambertian min red", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.05, Y: 0.0, Z: 0.0}}}),
			node(0, 3, "Lambertian min green", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.05, Z: 0.0}}}),
			node(0, 4, "Lambertian min blue", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.0, Z: 0.05}}}),
			node(1, 1, "Lambertian med gray", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.5, Y: 0.5, Z: 0.5}}}),
			node(1, 2, "Lambertian med red", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.5, Y: 0.0, Z: 0.0}}}),
			node(1, 3, "Lambertian med green", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.5, Z: 0.0}}}),
			node(1, 4, "Lambertian med blue", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.0, Z: 0.5}}}),
			node(2, 1, "Lambertian max gray", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.9, Y: 0.9, Z: 0.9}}}),
			node(2, 2, "Lambertian max red", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.9, Y: 0.0, Z: 0.0}}}),
			node(2, 3, "Lambertian max green", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.9, Z: 0.0}}}),
			node(2, 4, "Lambertian max blue", phys.Lambertian{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.0, Y: 0.0, Z: 0.9}}}),

			// Dielectric.
			node(3, 1, "dielectric clear n=1.2", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			node(3, 2, "dielectric clear n=1.5", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			node(3, 3, "dielectric clear n=1.8", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			node(3, 4, "dielectric clear n=2.0", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			node(4, 1, "dielectric good n=1.2", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			node(4, 2, "dielectric good n=1.5", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			node(4, 3, "dielectric good n=1.8", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			node(4, 4, "dielectric good n=2.0", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			node(5, 1, "dielectric frosty n=1.2", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			node(5, 2, "dielectric frosty n=1.5", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			node(5, 3, "dielectric frosty n=1.8", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			node(5, 4, "dielectric frosty n=2.0", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			// // Metal materials.
			node(6, 1, "shiny metal fine gray", phys.Metal{Albedo: r3.Vec{X: 0.9, Y: 0.9, Z: 0.9}, Fuzz: 0.025}),
			node(6, 2, "shiny metal smooth gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.05}),
			node(6, 3, "shiny metal medium gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.1}),
			node(6, 4, "shiny metal rough gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.15}),

			node(7, 1, "metal fine gray", phys.Metal{Albedo: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3}, Fuzz: 0.025}),
			node(7, 2, "metal smooth red", phys.Metal{Albedo: r3.Vec{X: 0.3, Y: 0.0, Z: 0.0}, Fuzz: 0.05}),
			node(7, 3, "metal medium green", phys.Metal{Albedo: r3.Vec{X: 0.0, Y: 0.3, Z: 0.0}, Fuzz: 0.1}),
			node(7, 4, "metal rough blue", phys.Metal{Albedo: r3.Vec{X: 0.0, Y: 0.0, Z: 0.3}, Fuzz: 0.15}),

			// Shaders for debug and visualization.
			node(9, 1, "ShaderNormal", phys.DebugNormal{}),
			node(9, 2, "ShaderUV", phys.DebugUV{}),

			phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4, "")[0],
			phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4, "")[1],
			phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4, "")[2],
		},
	}
	// scene.Add(phys.PropAxes(r3.Point{}, 0.01*phys.NM, 1*phys.NM)...)
	// scene.Add(phys.PropSkybox(1 * phys.M))
	fmt.Printf("bbox: %v\n", scene.Node[0].Shape.Bounds())

	// Render the scene and save it to ./output.png
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
	err = phys.SavePNG("./material.png", r.Image)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Saved to %s\n", path)
	// Render the stats path length image.
	// path = time.Now().Format("./out/out_20060102_150405_path_length.png")
	// // r.Stats.RaysPerPixel is a [][]uint32, we want this as an image.
	// // Convert the path length to a grayscale image.
	// statsImage := image.NewGray(image.Rect(0, 0, scene.RenderOptions.Dx, scene.RenderOptions.Dy))
	// for x := 0; x < scene.RenderOptions.Dx; x++ {
	// 	for y := 0; y < scene.RenderOptions.Dy; y++ {
	// 		rays := r.Stats.RaysPerPixel[x][y]
	// 		if rays > 255 {
	// 			rays = 255
	// 		}
	// 		statsImage.SetGray(x, y, color.Gray{Y: uint8(rays)})
	// 	}
	// }
	// err = phys.SavePNG(path, statsImage)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Saved to %s\n", path)
}
