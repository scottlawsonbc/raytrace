// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.

package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
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
			Dx:           4096,
			Dy:           4096,
		},
		Light: []phys.Light{},
		Camera: []phys.Camera{
			phys.OrthographicCamera{
				LookFrom: r3.Point{
					X: float64(-100 * phys.M),
					Y: float64(-100 * phys.M),
					Z: float64(-200 * phys.M)},
				LookAt: r3.Point{
					X: float64(0 * phys.NM),
					Y: float64(0 * phys.NM),
					Z: float64(0 * phys.NM)},
				VUp:       r3.Vec{X: 0, Y: 0, Z: -1},
				FOVHeight: 80 * phys.NM,
				FOVWidth:  80 * phys.NM,
			},
		},
	}
	assetFS := os.DirFS("../../../../3d/scan/")
	// modelPath := "linear-stage-controller-marker-simplified.obj"
	modelPath := "module-psu-buck.obj"
	// modelPath := "lilygo-esp32-camera.obj"
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
	fmt.Printf("Render Stats: %v\n", renderArtifact.Stats.PPrint())
	pathPNG := time.Now().Format("./out/out_20060102_150405.png")
	err = phys.SavePNG(pathPNG, renderArtifact.Image)
	if err != nil {
		panic(err)
	}
	err = phys.SavePNG("./gltf.png", renderArtifact.Image)
	if err != nil {
		panic(err)
	}
	log.Printf("Saved to ./gltf.png\n")
}
