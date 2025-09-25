// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/obj"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

func walk(fsys fs.FS, msg string) {
	fmt.Printf("Walking %v msg=%s\n", fsys, msg)
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			st, err := fs.Stat(fsys, p)
			if err != nil {
				return err
			}
			r, err := fsys.Open(p)
			if err != nil {
				return err
			}
			defer r.Close()

			// Read prefix.
			var buf [md5.Size]byte
			n, _ := io.ReadFull(r, buf[:])

			// Hash remainder.
			h := md5.New()
			_, err = io.Copy(h, r)
			if err != nil {
				return err
			}
			s := h.Sum(nil)
			log.Printf("| %s %d %x %x\n", p, st.Size(), buf[:n], s)
		}
		return nil
	})
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func animate(n int, nmax int) phys.Camera {
	theta := 2 * math.Pi * float64(n) / float64(nmax)
	phi := math.Pi/3 + math.Pi/4*math.Sin(2*math.Pi*float64(n)/float64(nmax))
	r := float64(1000 * phys.M)
	x := r * math.Sin(phi) * math.Cos(theta)
	z := r * math.Sin(phi) * math.Sin(theta)
	y := -r * math.Cos(phi)
	return phys.OrthographicCamera{
		LookFrom: r3.Point{
			X: float64(x),
			Y: float64(y),
			Z: float64(z)},
		LookAt: r3.Point{
			X: float64(0 * phys.NM),
			Y: float64(0 * phys.NM),
			Z: float64(0 * phys.NM)},
		VUp:       r3.Vec{X: 0, Y: -1, Z: 0},
		FOVHeight: 250 * phys.NM,
		FOVWidth:  250 * phys.NM,
	}
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		fmt.Println("Writing CPU profile to", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,
			MaxRayDepth:  10,
			Dx:           1024,
			Dy:           1024,
		},
		Light:  []phys.Light{},
		Camera: []phys.Camera{},
		Node:   []phys.Node{},
	}

	assetFS := os.DirFS("../../../../3d/scan/")
	modelPath := "bottle.obj"
	modelFS, err := fs.Sub(assetFS, modelPath)
	walk(modelFS, "main.modelFS")
	if err != nil {
		log.Fatalf("Failed to create model FS: %v", err)
	}
	parsedObj, err := obj.ParseFS(modelFS, modelPath)
	if err != nil {
		log.Fatalf("Failed to parse OBJ: %v", err)
	}
	nodes, err := phys.ConvertObjectToNodes(parsedObj, modelFS)
	if err != nil {
		log.Fatalf("Failed to convert OBJ to Nodes: %v", err)
	}
	fmt.Printf("got %d nodes\n", len(nodes))
	if len(nodes) == 0 {
		log.Fatalf("No nodes found in OBJ file")
	}
	for i, node := range nodes {
		fmt.Printf("node %d bounds %v\n", i, node.Shape.Bounds())
		scene.Add(node)
	}

	// Remove ./gif directory and create a new one.
	err = os.RemoveAll("./out")
	if err != nil {
		panic(err)
	}
	err = os.Mkdir("./out", 0755)
	if err != nil {
		panic(err)
	}

	artifacts := []image.Image{}
	nmax := 96
	for n := 0; n < nmax; n++ {
		camera := animate(n, nmax)
		scene.Camera = []phys.Camera{camera}
		artifact, err := phys.Render(context.Background(), &scene)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Render Stats: %v\n", artifact.Stats.PPrint())
		artifacts = append(artifacts, artifact.Image)
		// Save the number indexed frame to a file.
		pathFrame := fmt.Sprintf("./gif/frame_%d.png", n)
		err = phys.SavePNG(pathFrame, artifact.Image)
		if err != nil {
			panic(err)
		}
		log.Printf("Saved to %s\n", pathFrame)
	}

	pathGIF := time.Now().Format("./out/out_20060102_150405.gif")
	g := phys.NewGIF(artifacts)
	err = phys.SaveGIF(pathGIF, g)
	if err != nil {
		panic(err)
	}
	log.Printf("Saved to %s\n", pathGIF)

	pathPNG := time.Now().Format("./out/out_20060102_150405.png")
	err = phys.SavePNG(pathPNG, artifacts[0])
	if err != nil {
		panic(err)
	}
	err = phys.SavePNG("./bottle.png", artifacts[0])
	if err != nil {
		panic(err)
	}

	log.Printf("Saved to ./bottle.png\n")
}
