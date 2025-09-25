// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
//
// This example shows how to render a 3D scanned model and its color texture
// from a Wavefront .obj file and associated .mtl file.

package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
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

			// Read prefix
			var buf [md5.Size]byte
			n, _ := io.ReadFull(r, buf[:])

			// Hash remainder
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

	// light := 0.2
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,
			MaxRayDepth:  10,
			Dx:           512,
			Dy:           512,
		},
		Light: []phys.Light{
			// phys.PointLight{
			// 	Position: r3.Point{
			// 		X: float64(5000 * phys.NM),
			// 		Y: float64(5000 * phys.NM),
			// 		Z: float64(5000 * phys.NM)},
			// 	RadiantIntensity: r3.Vec{X: light, Y: light, Z: light},
			// },
		},
		// Camera: []phys.Camera{
		// 	phys.OrthographicCamera{
		// 		LookFrom: r3.Point{
		// 			X: float64(0 * phys.NM),
		// 			Y: float64(0 * phys.NM),
		// 			Z: float64(100 * phys.M)},
		// 		LookAt: r3.Point{
		// 			X: float64(0 * phys.NM),
		// 			Y: float64(0 * phys.NM),
		// 			Z: float64(0 * phys.NM)},
		// 		VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
		// 		FOVHeight: 120 * phys.NM,
		// 		FOVWidth:  120 * phys.NM,
		// 	},
		// },
		// Node: []phys.Node{
		// 	phys.PropSkySphere(1000*phys.M, phys.Emitter{Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0.5, Y: 0.5, Z: 0.5}}}),
		// },
	}
	cameraSequence := []phys.Camera{}
	// Orbit 360 around the origin about z axis in 32 steps
	const nstep = 4
	for i := 0; i < nstep; i++ {
		theta := float64(i) * 2 * math.Pi / nstep
		camera := phys.OrthographicCamera{
			LookFrom: r3.Point{
				X: float64(100 * float64(phys.M) * math.Cos(theta)),
				Z: float64(100*phys.M) - 80*math.Sin(theta)*float64(phys.M),
				Y: float64(100 * float64(phys.M) * math.Sin(theta))},
			LookAt: r3.Point{
				X: float64(0 * phys.M),
				Y: float64(0 * phys.M),
				Z: float64(0 * phys.M)},
			VUp:       r3.Vec{X: 0, Y: 0, Z: 1},
			FOVHeight: 250 * phys.NM,
			FOVWidth:  250 * phys.NM,
		}
		cameraSequence = append(cameraSequence, camera)
	}
	scene.Camera = cameraSequence

	assetFS := os.DirFS("../../../../../3d/scan/")
	// modelPath := "ttgo-t-camera-simplified.obj"
	modelPath := "linear-stage-controller-marker-simplified.obj"

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
	renderArtifact, err := phys.Render(context.Background(), &scene)
	if err != nil {
		panic(err)
	}
	// pathGIF := time.Now().Format("./out/out_20060102_150405.gif")
	pathPNG := time.Now().Format("./out/out_20060102_150405.png")
	// g := phys.NewGIF(renderArtifact.Color)
	// f, err := os.Create(pathGIF)
	// if err != nil {
	// 	panic(err)
	// }
	// err = gif.EncodeAll(f, g)
	// if err != nil {
	// 	panic(err)
	// }
	// f.Close()
	// fmt.Println("Saved to", pathGIF)
	err = phys.SavePNG(pathPNG, renderArtifact.Image)
	if err != nil {
		panic(err)
	}
	// Save another copy with the same filename so that for debugging the
	// image can be opened in one pane and automatically reloads when rendered.
	err = phys.SavePNG("./scan.png", renderArtifact.Image)
	if err != nil {
		panic(err)
	}
	// traceSizeMB := len(renderArtifact.Stats.Trace) / 1024 / 1024
	// log.Printf("Saved to %s (%d MB)\n", pathGIF, traceSizeMB)
	log.Printf("Saved to ./scan.png\n")
}
