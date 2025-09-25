// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package main

func main() {

}

// package main

// import (
// 	"image"
// 	"image/color/palette"
// 	"image/draw"
// 	"image/gif"
// 	"image/png"
// 	"log"
// 	"math"
// 	"os"

// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
// )

// func bead(center r3.Point, diameter float64) phys.Node {
// 	return phys.Node{
// 		Shape:    phys.Sphere{Center: center, Radius: diameter / 2},
// 		Material: phys.Lambertian{Texture: r3.Vec{X: 0.5, Y: 0.2, Z: 0.5}},
// 	}
// }

// var exampleSnapshot = phys.Scene{
// 	Detector: phys.Detector{
// 		Name:   "TestDetector",
// 		Width:  250,
// 		Height: 250,
// 		// PixelSize: 9.2 * phys.Micrometer,
// 		PixelSize: 2,
// 	},
// 	// Lens: phys.Lens{
// 	// 	Name: "Nikon 20x",
// 	// 	// WorkingDistance:   1.0, // mm. // TODO: this is not used yet.
// 	// 	NumericalAperture: 0.75,
// 	// 	Magnification:     20,
// 	// 	FieldNumber:       25, // TODO: this is not used yet.
// 	// },
// 	LookAt:   r3.Point{X: 0, Y: 0, Z: 0}.Add(r3.Vec{X: 0.5, Y: -0.5, Z: 0}),
// 	LookFrom: r3.Point{X: 0, Y: 20, Z: 100}.Add(r3.Vec{X: 0.5, Y: -0.5, Z: 0}),

// 	Illuminant: []phys.Emitter{
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 20000000000}, Color: r3.Vec{X: 1, Y: 1, Z: 1}},

// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 		// // X axis.
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 1000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 2000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 3000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		// // Y axis.
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 1000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 2000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 3000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		// // Z axis.
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 1000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 2000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 		// {Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 3000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 	},
// 	Node: []phys.Node{
// 		{Shape: phys.LoadOBJ("../asset/array-coarse.obj"), Material: phys.Lambertian{Texture: r3.Vec{X: 0.5, Y: 0.5, Z: 0.5}}},
// 		// bead(r3.Point{X: 0, Y: 0, Z: 0}, 4000),
// 		// bead(r3.Point{X: 0, Y: 4, Z: 0}, 40),
// 		// bead(r3.Point{X: 4, Y: 0, Z: 0}, 40),
// 		// bead(r3.Point{X: 4, Y: 4, Z: 0}, 40),
// 	},
// }

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func save(path string, img image.Image) {
// 	f, err := os.Create(path)
// 	check(err)
// 	err = (&png.Encoder{CompressionLevel: png.NoCompression}).Encode(f, img)
// 	f.Close()
// 	check(err)
// }

// func makeGif(imgs []image.Image) *gif.GIF {
// 	palettedImages := []*image.Paletted{}
// 	delays := []int{}
// 	for _, img := range imgs {
// 		pImg := image.NewPaletted(img.Bounds(), palette.Plan9)
// 		draw.Draw(pImg, pImg.Rect, img, img.Bounds().Min, draw.Over)
// 		palettedImages = append(palettedImages, pImg)
// 		delays = append(delays, 0)
// 	}
// 	return &gif.GIF{
// 		Image: palettedImages,
// 		Delay: delays,
// 	}
// }

// func main() {
// 	opts := phys.RenderOptions{
// 		Seed:         0,
// 		RaysPerPixel: 1,
// 		MaxRayDepth:     100,
// 	}
// 	log.Printf("bounds %+v\n", exampleSnapshot.Node[0].Shape.Bounds())

// 	imgs := []image.Image{}
// 	r := 500000.0
// 	for theta := 0.0; theta < 2*math.Pi; theta += 2 * math.Pi / 8 {
// 		exampleSnapshot.LookFrom.X = r * math.Cos(theta)
// 		exampleSnapshot.LookFrom.Y = r * math.Sin(theta)
// 		exampleSnapshot.LookFrom.Z = r * 2
// 		img := phys.Render(&exampleSnapshot, opts)
// 		imgs = append(imgs, img)
// 	}
// 	log.Println("generating gif")
// 	g := makeGif(imgs)
// 	f, err := os.Create("out.gif")
// 	check(err)
// 	defer f.Close()
// 	err = gif.EncodeAll(f, g)
// 	check(err)
// 	println("done")
// }
