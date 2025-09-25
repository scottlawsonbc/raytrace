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
// 	"log"
// 	"math"
// 	"os"

// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
// )

// var exampleSnapshot = phys.Scene{
// 	Detector: phys.Detector{
// 		Name:      "CMV50000Bin2x2",
// 		Width:     3960,
// 		Height:    3002,
// 		PixelSize: 9.2 * phys.Micrometer,
// 	},
// 	LookAt:   r3.Point{X: 0, Y: 0, Z: 0},
// 	LookFrom: r3.Point{X: -1, Y: -1, Z: -1},

// 	Illuminant: []phys.Emitter{
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 1, Z: 1}},
// 		// X axis.
// 		{Shape: phys.Sphere{Center: r3.Point{X: 1000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 2000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 3000, Y: 0, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 1, Y: 0.2, Z: 0.2}},
// 		// Y axis.
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 1000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 2000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 3000, Z: 0}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 1, Z: 0.2}},
// 		// Z axis.
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 1000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 2000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 		{Shape: phys.Sphere{Center: r3.Point{X: 0, Y: 0, Z: 3000}, Radius: 20 * 5}, Color: r3.Vec{X: 0.2, Y: 0.2, Z: 1}},
// 	},
// }

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
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
// 	imgs := []image.Image{}
// 	r := 500000.0
// 	for theta := 0.0; theta < 2*math.Pi; theta += 2 * math.Pi / 32 {
// 		exampleSnapshot.LookFrom.X = r * math.Cos(theta)
// 		exampleSnapshot.LookFrom.Y = r * math.Sin(theta)
// 		exampleSnapshot.LookFrom.Z = r
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
