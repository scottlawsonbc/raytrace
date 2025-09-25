// Program owl renders an animated view of an OBJ model and saves GIF/PNG artifacts.
// It demonstrates loading and converting an OBJ to scene nodes, rendering frames,
// and writing the results to disk. The program logs basic file metadata while
// walking the asset directory.
//
// Concurrency guarantees:
// The program performs all work on the main goroutine. It does not spawn worker
// goroutines or share mutable state across goroutines.
//
// Zero value semantics:
// The zero value of the package-level variables is not used directly; flags must
// be parsed before use.
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

// walk walks the provided filesystem root and logs md5 information for files.
// walk reports any traversal error returned by the underlying filesystem and
// avoids dereferencing a nil DirEntry by handling err first in the callback.
func walk(fsys fs.FS, msg string) {
	if fsys == nil {
		fmt.Printf("walk got nil fsys msg=%s\n", msg)
		return
	}
	fmt.Printf("Walking %v msg=%s\n", fsys, msg)
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			// Propagate the error to stop at the failing subtree.
			return err
		}
		if d == nil {
			// Defensive guard; WalkDir may pass nil d only when err != nil,
			// but we guard anyway in case of unusual FS behavior.
			return nil
		}
		if d.IsDir() {
			return nil
		}

		st, err := fs.Stat(fsys, p)
		if err != nil {
			return err
		}
		r, err := fsys.Open(p)
		if err != nil {
			return err
		}
		defer r.Close()

		var prefix [md5.Size]byte
		n, _ := io.ReadFull(r, prefix[:])

		h := md5.New()
		if _, err := io.Copy(h, r); err != nil {
			return err
		}
		sum := h.Sum(nil)
		log.Printf("| %s %d %x %x\n", p, st.Size(), prefix[:n], sum)
		return nil
	})
	if err != nil {
		log.Printf("walk(%s) encountered error: %v", msg, err)
	}
}

// cpuprofile is the optional path to write a CPU profile file.
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// animate constructs a camera on a circular path around the origin for frame n.
// animate reports an orthographic camera whose position sweeps 360 degrees
// across nmax frames. The function has no side effects.
func animate(n int, nmax int) phys.Camera {
	theta := 2 * math.Pi * float64(n) / float64(nmax)
	phi := math.Pi / 2
	r := float64(1000 * phys.M)
	x := r * math.Sin(phi) * math.Cos(theta)
	z := r * math.Sin(phi) * math.Sin(theta)
	y := r * math.Cos(phi)
	return phys.OrthographicCamera{
		LookFrom: r3.Point{
			X: float64(x),
			Y: float64(y),
			Z: float64(z)},
		LookAt: r3.Point{
			X: float64(0 * phys.NM),
			Y: float64(0 * phys.NM),
			Z: float64(0 * phys.NM)},
		VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
		FOVHeight: 180 * phys.NM,
		FOVWidth:  180 * phys.NM,
	}
}

// // loadNodes loads an OBJ from the filesystem path and converts it to scene nodes.
// // loadNodes reports the parsed node list or an error. The function does not
// // mutate global state. The zero value for the returned slice is nil when an
// // error occurs.
// func loadNodes(path string) (nodes []phys.Node, err error) {
// 	dir := filepath.Dir(path)
// 	name := filepath.Base(path)

// 	assetFS := os.DirFS(dir)
// 	walk(assetFS, "loadNodes.assetFS")

// 	parsedObj, err := obj.ParseFS(assetFS, name)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nodes, err = phys.ConvertObjectToNodes(parsedObj, assetFS)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(nodes) == 0 {
// 		return nil, fmt.Errorf("no nodes found in OBJ file")
// 	}
// 	fmt.Printf("got %d nodes\n", len(nodes))
// 	return nodes, nil
// }

// main parses flags, loads the model, renders an animation, and saves outputs.
// main exits with a non-zero status on error. It reports whether artifacts were
// written via log messages.
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

	// Configure model location.
	// assetDir is a directory; modelName is a file inside that directory.
	assetDir := "../../../../../3d/scan"
	modelName := "owl.obj"

	assetFS := os.DirFS(assetDir)
	walk(assetFS, "main.assetFS")

	parsedObj, err := obj.ParseFS(assetFS, modelName)
	if err != nil {
		log.Fatalf("Failed to parse OBJ: %v", err)
	}
	nodes, err := phys.ConvertObjectToNodes(parsedObj, assetFS)
	if err != nil {
		log.Fatalf("Failed to convert OBJ to Nodes: %v", err)
	}
	if len(nodes) == 0 {
		log.Fatalf("No nodes found in OBJ file")
	}

	for i, node := range nodes {
		fmt.Printf("node %d bounds %v\n", i, node.Shape.Bounds())
		scene.Add(node)
	}

	artifacts := []image.Image{}
	nmax := 96
	for n := 0; n < nmax; n++ {
		camera := animate(n, nmax)
		scene.Camera = []phys.Camera{camera}
		artifact, err := phys.Render(context.Background(), &scene)
		if err != nil {
			log.Fatalf("Render failed on frame %d: %v", n, err)
		}
		fmt.Printf("Render Stats: %v\n", artifact.Stats.PPrint())
		artifacts = append(artifacts, artifact.Image)
	}

	pathGIF := time.Now().Format("./out/out_20060102_150405.gif")
	g := phys.NewGIF(artifacts)
	if err := phys.SaveGIF(pathGIF, g); err != nil {
		log.Fatalf("Failed to save GIF: %v", err)
	}
	log.Printf("Saved to %s\n", pathGIF)

	pathPNG := time.Now().Format("./out/out_20060102_150405.png")
	if err := phys.SavePNG(pathPNG, artifacts[0]); err != nil {
		log.Fatalf("Failed to save PNG: %v", err)
	}
	if err := phys.SavePNG("./owl.png", artifacts[0]); err != nil {
		log.Fatalf("Failed to save PNG ./owl.png: %v", err)
	}

	log.Printf("Saved to ./owl.png\n")
}
