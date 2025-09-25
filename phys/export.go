// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package phys

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
)

func MustLoadPNG(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		panic(err)
	}
	return img
}

func MustLoadJPEG(path string) image.Image {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		panic(err)
	}
	return img
}

func SavePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	err = (&png.Encoder{CompressionLevel: png.NoCompression}).Encode(f, img)
	f.Close()
	return err
}

func NewGIF(imgs []image.Image) *gif.GIF {
	palettedImages := []*image.Paletted{}
	delays := []int{}
	for _, img := range imgs {
		pImg := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(pImg, pImg.Rect, img, img.Bounds().Min, draw.Over)
		palettedImages = append(palettedImages, pImg)
		delays = append(delays, 0)
	}
	return &gif.GIF{
		Image: palettedImages,
		Delay: delays,
	}
}

func SaveGIF(path string, g *gif.GIF) error {
	fmt.Printf("Saving GIF to %s with n=%d\n", path, len(g.Image))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}

func Montage(imgs []image.Image) image.Image {
	if len(imgs) == 0 {
		return nil
	}
	if len(imgs) == 1 {
		return imgs[0]
	}

	// Montage multiple images into a single image.
	dx := imgs[0].Bounds().Dx()
	dy := imgs[0].Bounds().Dy()
	for _, img := range imgs {
		if img.Bounds().Dx() != dx || img.Bounds().Dy() != dy {
			panic(fmt.Sprintf("image sizes do not match: %v vs %v", img.Bounds(), imgs[0].Bounds()))
		}
	}
	montage := image.NewRGBA(image.Rect(0, 0, dx*len(imgs), dy))
	for i, img := range imgs {
		for y := 0; y < dy; y++ {
			for x := 0; x < dx; x++ {
				montage.Set(x+i*dx, y, img.At(x, y))
			}
		}
	}
	return montage
}
