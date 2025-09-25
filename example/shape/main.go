// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
)

// This example scene demonstrates using different primitive shapes.

// func node(col, row int, name string, mat phys.Material) phys.Node {
// 	dx := float64(10 * phys.MM)
// 	dy := float64(10 * phys.MM)
// 	cx := float64(col)*dx - float64(45*phys.MM)
// 	cy := float64(row)*dy - float64(25*phys.MM)
// 	radius := 4 * phys.MM
// 	return phys.Node{
// 		Name:     name,
// 		Shape:    phys.Sphere{Center: r3.Point{X: cx, Y: cy, Z: 0}, Radius: radius},
// 		Material: mat,
// 	}
// }

var MM = phys.MM

func mm(v float64) float64 {
	return v * float64(phys.MM)
}

func main() {
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 5,
			MaxRayDepth:  10,
			Dx:           1024,
			Dy:           1024,
		},
		Light: []phys.Light{
			phys.PointLight{
				Position: r3.Point{
					X: float64(500 * phys.MM),
					Y: float64(500 * phys.MM),
					Z: float64(500 * phys.MM)},
				RadiantIntensity: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3},
			},
		},
		Camera: []phys.Camera{
			phys.OrthographicCamera{
				LookFrom: r3.Point{
					X: float64(100 * phys.MM),
					Y: float64(100 * phys.MM),
					Z: float64(100 * phys.MM)},
				LookAt:    r3.Point{X: 0, Y: 0, Z: 0},
				VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
				FOVHeight: 200 * phys.MM,
				FOVWidth:  200 * phys.MM,
			},
		},
		Node: []phys.Node{

			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, phys.MM/2, phys.MM*100, "")[0],
			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, phys.MM/2, phys.MM*100, "")[1],
			phys.PropAxes(r3.Point{X: 0, Y: 0, Z: 0}, phys.MM/2, phys.MM*100, "")[2],

			// {
			// 	Name: "diffuse plane",
			// 	Shape: phys.Plane{
			// 		Center: r3.Point{X: float64(30 * phys.MM), Y: float64(30 * phys.MM), Z: 0},
			// 		Normal: r3.Vec{X: 0, Y: 0, Z: 1},
			// 		Width:  20 * phys.MM,
			// 		Height: 20 * phys.MM,
			// 	},
			// 	// Material: phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}},
			// 	Material: phys.Emitter{Texture: phys.TextureImage{FilePath: "./happy.png", Image: phys.MustLoadPNG("./happy.png")}},
			// },
			// {
			// 	Name: "diffuse ground plane",
			// 	Shape: phys.Plane{
			// 		Center: r3.Point{X: mm(50), Y: 0, Z: mm(50)},
			// 		Normal: r3.Vec{X: 0, Y: 1, Z: 0},
			// 		Width:  40 * phys.MM,
			// 		Height: 40 * phys.MM,
			// 	},
			// 	Material: phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.2, Y: 0.2, Z: 0.5}}},
			// },
			// {
			// 	Name: "decal plane",
			// 	Shape: phys.Plane{
			// 		Center: r3.Point{X: mm(0), Y: mm(20), Z: mm(20)},
			// 		Normal: r3.Vec{X: 1, Y: 0, Z: 0},
			// 		Width:  40 * phys.MM,
			// 		Height: 40 * phys.MM,
			// 	},
			// 	Material: phys.Emitter{Texture: phys.TextureImage{FilePath: "./happy.png", Image: phys.MustLoadPNG("./happy.png")}},
			// },

			// {
			// 	Name: "dielectric",
			// 	Shape: phys.Sphere{
			// 		Center: r3.Point{X: mm(0), Y: mm(0), Z: mm(0)},
			// 		Radius: 35 * phys.MM,
			// 	},
			// 	Material: phys.Dielectric{RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0, Roughness: 0.00},
			// },
			// {
			// 	Name: "sky globe 1",
			// 	Shape: phys.Sphere{
			// 		Center: r3.Point{X: mm(-50), Y: mm(0), Z: mm(0)},
			// 		Radius: 30 * phys.MM,
			// 	},
			// 	Material: phys.Emitter{Texture: phys.TextureImage{FilePath: "./star-1.png", Interp: "bilinear", Image: phys.MustLoadPNG("./star-1.png")}},
			// },
			// {
			// 	Name: "sky globe 2",
			// 	Shape: phys.Sphere{
			// 		Center: r3.Point{X: mm(50), Y: mm(0), Z: mm(0)},
			// 		Radius: 10 * phys.MM,
			// 	},
			// 	Material: phys.Emitter{Texture: phys.TextureImage{FilePath: "./star-1.png", Image: phys.MustLoadPNG("./star-1.png")}},
			// },
			// {
			// 	Name: "cylinder",
			// 	Shape: phys.Cylinder{
			// 		Origin:    r3.Point{X: mm(-60), Y: mm(-40), Z: mm(0)},
			// 		Radius:    10 * phys.MM,
			// 		Direction: r3.Vec{X: 1, Y: 0, Z: 0},
			// 		Height:    50 * phys.MM,
			// 	},
			// 	Material: phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}},
			// },
			// // phys.PropSkybox(1200*phys.MM, phys.Emitter{Texture: phys.TextureImage{FilePath: "./abstract-2.png", Image: phys.MustLoadPNG("./abstract-2.png")}}),
			// // phys.PropSkybox(40*phys.MM, phys.Emitter{Texture: phys.TextureImage{FilePath: "./star-1.png", Image: phys.MustLoadPNG("./star-1.png")}}),
			// {
			// 	Name: "plane",
			// 	Shape: phys.Plane{
			// 		Center: r3.Point{X: mm(70), Y: 0, Z: 0},
			// 		Normal: r3.Vec{X: 0, Y: 1, Z: 0},
			// 		Width:  10 * phys.MM,
			// 		Height: 20 * phys.MM,
			// 	},
			// 	Material: phys.Emitter{Texture: phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}},
			// },

			// {
			// 	Name: "plane with texture",
			// 	Shape: phys.Plane{
			// 		Center: r3.Point{X: mm(0), Y: 0, Z: mm(-100)},
			// 		Normal: r3.Vec{X: 0, Y: 1, Z: 0},
			// 		Width:  10 * phys.MM,
			// 		Height: 20 * phys.MM,
			// 	},
			// 	Material: phys.Emitter{Texture: phys.TextureImage{FilePath: "./happy.png", Image: phys.MustLoadPNG("./happy.png")}},
			// },

			{
				Name: "rotated2 plane with texture",
				Shape: phys.TransformedShape{
					Shape: phys.Quad{
						Center: r3.Point{},
						Normal: r3.Vec{X: 0, Y: 1, Z: 0},
						Width:  100 * phys.MM,
						Height: 100 * phys.MM,
					},
					Transform: phys.Transform{
						Translation: r3.Vec{X: mm(0), Y: mm(50), Z: mm(50)},
						Rotation:    r3.RotationMatrixZ(math.Pi / 2), // Rotate 45 degrees around Z-axis
						Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
					},
				},
				Material: phys.Emitter{
					Texture: phys.TextureImage{
						FilePath: "./asset/wood.png",
						Image:    phys.MustLoadPNG("./asset/wood.png"),
						Interp:   "bilinear",
						WrapMode: "repeat",
					},
				},
			},
			// Concrete031_2K-PNG_Color
			{
				Name: "rotated plane with texture",
				Shape: phys.TransformedShape{
					Shape: phys.Quad{
						Center: r3.Point{X: mm(0), Y: mm(0), Z: mm(0)},
						Normal: r3.Vec{X: 0, Y: 1, Z: 0},
						Width:  100 * phys.MM,
						Height: 100 * phys.MM,
					},
					Transform: phys.Transform{
						Translation: r3.Vec{X: mm(50), Y: mm(50), Z: mm(0)},
						Rotation:    r3.RotationMatrixY(math.Pi / 2).Mul(r3.RotationMatrixZ(math.Pi / 2)), // Rotate 45 degrees around Z-axis
						Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
					},
				},
				Material: phys.Emitter{
					Texture: phys.TextureImage{
						FilePath: "./asset/stars.png",
						Image:    phys.MustLoadPNG("./asset/stars.png"),
						Interp:   "bilinear",
						WrapMode: "repeat",
					},
				},
			},

			{
				Name: "rotated plane with texture 3",
				Shape: phys.TransformedShape{
					Shape: phys.Quad{
						Center: r3.Point{X: mm(0), Y: mm(0), Z: mm(0)},
						Normal: r3.Vec{X: 0, Y: 1, Z: 0},
						Width:  100 * phys.MM,
						Height: 100 * phys.MM,
					},
					Transform: phys.Transform{
						Translation: r3.Vec{X: mm(50), Y: mm(0), Z: mm(50)},
						Rotation:    r3.IdentityMat3x3(),
						Scale:       r3.Vec{X: 1, Y: 1, Z: 1},
					},
				},
				Material: phys.Emitter{
					Texture: phys.TextureImage{
						FilePath: "./asset/rocks.png",
						Image:    phys.MustLoadPNG("./asset/rocks.png"),
						Interp:   "bilinear",
						WrapMode: "repeat",
					},
				},
			},

			// return phys.Node{
			// 	Name:     name,
			// 	Shape:    phys.Sphere{Center: r3.Point{X: cx, Y: cy, Z: 0}, Radius: radius},
			// 	Material: mat,
			// }

			// // Lambertian materials.
			// node(0, 0, "checker red white", phys.Lambertian{Albedo: phys.TextureCheckerboard{
			// 	Odd:       phys.TextureUniform{Color: r3.Vec{X: 0.2, Y: 0.5, Z: 0.6}},
			// 	Even:      phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.0, Z: 0.0}},
			// 	Frequency: 50,
			// }}),
			// node(1, 0, "checker complementary colors", phys.Lambertian{Albedo: phys.TextureCheckerboard{
			// 	Odd:       phys.TextureUniform{Color: r3.Vec{X: 0, Y: 0, Z: 123}.Divs(255)},
			// 	Even:      phys.TextureUniform{Color: r3.Vec{X: 242, Y: 0, Z: 0}.Divs(255)},
			// 	Frequency: 70,
			// }}),
			// node(2, 0, "image texture", phys.Lambertian{
			// 	Albedo: phys.TextureImage{
			// 		Image:    phys.MustLoadPNG("./texture.png"),
			// 		FilePath: "./texture.png",
			// 		Interp:   "bilinear",
			// 		WrapMode: "repeat",
			// 	},
			// }),
			// node(3, 0, "face texture", phys.Lambertian{
			// 	Albedo: phys.TextureImage{
			// 		Image:    phys.MustLoadPNG("./faces.png"),
			// 		FilePath: "./faces.png",
			// 		Interp:   "bilinear",
			// 		WrapMode: "repeat",
			// 	},
			// }),
			// node(0, 1, "Lambertian min gray", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.05, Y: 0.05, Z: 0.05}}}),
			// node(0, 2, "Lambertian min red", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.05, Y: 0.0, Z: 0.0}}}),
			// node(0, 3, "Lambertian min green", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.05, Z: 0.0}}}),
			// node(0, 4, "Lambertian min blue", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.0, Z: 0.05}}}),
			// node(1, 1, "Lambertian med gray", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}}),
			// node(1, 2, "Lambertian med red", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.5, Y: 0.0, Z: 0.0}}}),
			// node(1, 3, "Lambertian med green", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.5, Z: 0.0}}}),
			// node(1, 4, "Lambertian med blue", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.0, Z: 0.5}}}),
			// node(2, 1, "Lambertian max gray", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.9, Y: 0.9, Z: 0.9}}}),
			// node(2, 2, "Lambertian max red", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.9, Y: 0.0, Z: 0.0}}}),
			// node(2, 3, "Lambertian max green", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.9, Z: 0.0}}}),
			// node(2, 4, "Lambertian max blue", phys.Lambertian{Albedo: phys.TextureUniform{Color: r3.Vec{X: 0.0, Y: 0.0, Z: 0.9}}}),

			// // Dielectric.
			// node(3, 1, "dielectric clear n=1.2", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			// node(3, 2, "dielectric clear n=1.5", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			// node(3, 3, "dielectric clear n=1.8", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			// node(3, 4, "dielectric clear n=2.0", phys.Dielectric{Roughness: 0.0, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			// node(4, 1, "dielectric good n=1.2", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			// node(4, 2, "dielectric good n=1.5", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			// node(4, 3, "dielectric good n=1.8", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			// node(4, 4, "dielectric good n=2.0", phys.Dielectric{Roughness: 0.05, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			// node(5, 1, "dielectric frosty n=1.2", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.2, RefractiveIndexExterior: 1.0}),
			// node(5, 2, "dielectric frosty n=1.5", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.5, RefractiveIndexExterior: 1.0}),
			// node(5, 3, "dielectric frosty n=1.8", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 1.8, RefractiveIndexExterior: 1.0}),
			// node(5, 4, "dielectric frosty n=2.0", phys.Dielectric{Roughness: 0.2, RefractiveIndexInterior: 2, RefractiveIndexExterior: 1.0}),

			// // // Metal materials.
			// node(6, 1, "metal fine gray", phys.Metal{Albedo: r3.Vec{X: 0.9, Y: 0.9, Z: 0.9}, Fuzz: 0.025}),
			// node(6, 2, "metal smooth gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.05}),
			// node(6, 3, "metal medium gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.1}),
			// node(6, 4, "metal rough gray", phys.Metal{Albedo: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}, Fuzz: 0.15}),

			// node(7, 1, "metal fine gray", phys.Metal{Albedo: r3.Vec{X: 0.3, Y: 0.3, Z: 0.3}, Fuzz: 0.025}),
			// node(7, 2, "metal smooth red", phys.Metal{Albedo: r3.Vec{X: 0.3, Y: 0.0, Z: 0.0}, Fuzz: 0.05}),
			// node(7, 3, "metal medium green", phys.Metal{Albedo: r3.Vec{X: 0.0, Y: 0.3, Z: 0.0}, Fuzz: 0.1}),
			// node(7, 4, "metal rough blue", phys.Metal{Albedo: r3.Vec{X: 0.0, Y: 0.0, Z: 0.3}, Fuzz: 0.15}),

			// // Shaders for debug and visualization.
			// node(9, 1, "ShaderNormal", phys.ShaderNormal{}),
			// node(9, 2, "ShaderUV", phys.ShaderUV{}),

			// phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4)[0],
			// phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4)[1],
			// phys.PropAxes(r3.Point{X: float64(-45 * phys.MM), Y: float64(-45 * phys.MM)}, phys.MM*0.15, phys.MM*4)[2],
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
	err = phys.SavePNG("./shape.png", r.Image)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Saved to %s\n", path)
}
