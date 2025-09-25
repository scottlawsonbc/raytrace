// Copyright 2024 Scott Lawson scottlawsonbc@gmail.com. All rights reserved.
package phys

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"cmp"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type ray struct {
	radiance  Spectrum
	origin    r3.Point
	direction r3.Vec
	depth     int
	pixelX    int
	pixelY    int
	rand      *Rand
}

func (r ray) at(t Distance) r3.Point {
	p := r.origin.Add(r.direction.Muls(float64(t)))
	return p
}

// RenderStats collects runtime metrics for the rendering process.
type RenderStats struct {
	RaysExceededDepth uint64        // Total count of rays that exceeded max ray depth.
	RaysLeftScene     uint64        // Total count of rays that left the scene.
	TotalRays         uint64        // Total count of all rays generated.
	RenderTime        time.Duration // How long it took to render the scene.
	Dx                int           // Width of the rendered image.
	Dy                int           // Height of the rendered image.
}

func (stats RenderStats) String() string {
	return fmt.Sprintf("RenderStats{RaysExceededDepth=%d, RaysLeftScene=%d, TotalRays=%d, RenderTime=%s}",
		stats.RaysExceededDepth, stats.RaysLeftScene, stats.TotalRays, stats.RenderTime)
}

func (s RenderStats) PPrint() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}
	str := string(data)
	timePerPx := s.RenderTime / time.Duration(s.Dx*s.Dy)
	var maxDepthPercent float64
	var outScenePercent float64
	if s.TotalRays != 0 {
		maxDepthPercent = 100 * float64(s.RaysExceededDepth) / float64(s.TotalRays)
		outScenePercent = 100 * float64(s.RaysLeftScene) / float64(s.TotalRays)
	}
	str += "\n" + fmt.Sprintf("RenderTime: %s (%s per pixel)\n", s.RenderTime, timePerPx)
	str += fmt.Sprintf("TotalRays: %d\n", s.TotalRays)
	str += fmt.Sprintf("RaysExceedingDepth: %d (%.1f%%)\n", s.RaysExceededDepth, maxDepthPercent)
	str += fmt.Sprintf("RaysLeftScene: %d (%.1f%%)\n", s.RaysLeftScene, outScenePercent)
	str += fmt.Sprintf("Rendered %dx%d\n", s.Dx, s.Dy)
	return str
}

type RenderOptions struct {
	Seed         int64 // Random base seed.
	RaysPerPixel int   // Number of rays to generate for each pixel.
	MaxRayDepth  int   // Maximum number of collisions before terminating ray.
	Dx           int   // Width of the rendered image in pixels.
	Dy           int   // Height of the rendered image in pixels.
}

func (r RenderOptions) Validate() error {
	if r.Seed < 0 {
		return fmt.Errorf("bad Seed must be non-negative but got %d", r.Seed)
	}
	if r.RaysPerPixel <= 0 {
		return fmt.Errorf("bad RaysPerPixel must be positive but got %d", r.RaysPerPixel)
	}
	if r.MaxRayDepth <= 0 {
		return fmt.Errorf("bad MaxRayDepth must be positive but got %d", r.MaxRayDepth)
	}
	if r.Dx <= 0 {
		return fmt.Errorf("bad Dx must be positive but got %d", r.Dx)
	}
	if r.Dy <= 0 {
		return fmt.Errorf("bad Dy must be positive but got %d", r.Dy)
	}
	return nil
}

// RenderArtifact represents the output of a rendering process (a render artifact).
type RenderArtifact struct {
	Image *image.RGBA
	Stats RenderStats
}

type tile struct {
	x0, x1, y0, y1 int
}

func (t tile) String() string {
	return fmt.Sprintf("Tile{xStart=%d, xEnd=%d, yStart=%d, yEnd=%d}", t.x0, t.x1, t.y0, t.y1)
}

// min reports the smaller of a and b.
// It works for any ordered type: integers, floats, strings.
func min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// max reports the larger of a and b.
func max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// clamp reports a if a is in [min, max], min if a < min, and max if a > max.
func clamp[T cmp.Ordered](a, minVal, maxVal T) T {
	if a < minVal {
		return minVal
	}
	if a > maxVal {
		return maxVal
	}
	return a
}

func tracePath(ctx context.Context, scene *Scene, r ray, stats *RenderStats) Spectrum {
	atomic.AddUint64(&stats.TotalRays, 1)
	if ctx.Err() != nil {
		return Spectrum{}
	}
	if r.origin.IsNaN() || r.origin.IsInf() || r.direction.IsNaN() || r.direction.IsInf() {
		log.Printf("invalid ray: %+v", r)
		return Spectrum{}
	}
	if r.depth > scene.RenderOptions.MaxRayDepth {
		atomic.AddUint64(&stats.RaysExceededDepth, 1)
		return Spectrum{}
	}
	var nearest surfaceInteraction
	var minDist = Distance(math.MaxFloat64)
	var hit bool
	for i := range scene.Node {
		node := &scene.Node[i]
		h, c := node.Shape.Collide(r, eps, minDist)
		if h && c.t < minDist {
			minDist = c.t
			nearest.collision = c
			nearest.incoming = r
			nearest.outgoing = r.direction.Muls(-1) // Direction towards camera.
			nearest.node = scene.Node[i]
			hit = true
		}
	}
	if !hit {
		atomic.AddUint64(&stats.RaysLeftScene, 1)
		return Spectrum{}
	}
	resolution := nearest.node.Material.Resolve(ctx, nearest)
	rgb := Spectrum{}
	rgb = rgb.Add(resolution.emission)
	rgb = rgb.Add(nearest.node.Material.ComputeDirectLighting(ctx, nearest, scene))
	for _, newRay := range resolution.scattered {
		outgoingColor := tracePath(ctx, scene, newRay, stats)
		rgb = rgb.Add(outgoingColor)
	}
	return rgb
}

// renderPixel renders a single pixel in the image. Any x, y outside the image bounds will be clamped.
func renderPixel(ctx context.Context, scene *Scene, camera Camera, rand *Rand, stats *RenderStats, x, y int, img *image.RGBA) {
	dx := scene.RenderOptions.Dx
	dy := scene.RenderOptions.Dy
	// Clamp pixel coordinates to image bounds.
	cx := clamp(x, 0, dx-1)
	cy := clamp(y, 0, dy-1)
	if x != cx || y != cy {
		log.Printf("clamped pixel coordinates: (x, y)=(%d, %d) to (%d, %d)", x, y, cx, cy)
	}
	imgy := dy - 1 - cy // Flip y-axis to match image coordinates.
	rgb := Spectrum{}
	for sample := 0; sample < scene.RenderOptions.RaysPerPixel; sample++ {
		if ctx.Err() != nil {
			return
		}
		var s, tSample float64
		if scene.RenderOptions.RaysPerPixel == 1 {
			// Sample center of the pixel.
			s = (float64(cx) + 0.5) / float64(dx)
			tSample = (float64(cy) + 0.5) / float64(dy)
		} else {
			// Sample randomly within the pixel.
			s = (float64(cx) + rand.Float64()) / float64(dx)
			tSample = (float64(cy) + rand.Float64()) / float64(dy)
		}
		// Cast the ray from the camera.
		ray := camera.Cast(s, tSample, rand)
		ray.pixelX = cx
		ray.pixelY = imgy
		color := tracePath(ctx, scene, ray, stats)
		rgb = rgb.Add(color)
	}
	rgb = rgb.Divs(float64(scene.RenderOptions.RaysPerPixel))
	img.Set(x, imgy, color.RGBA{
		R: uint8(math.Min(255, 255.99*rgb.X)),
		G: uint8(math.Min(255, 255.99*rgb.Y)),
		B: uint8(math.Min(255, 255.99*rgb.Z)),
		A: 255,
	})
}

func renderTile(ctx context.Context, scene *Scene, camera Camera, t tile, img *image.RGBA, stats *RenderStats) {
	for y := t.y0; y < t.y1; y++ {
		select {
		case <-ctx.Done():
			return
		default:
		}
		rand := NewRand(scene.RenderOptions.Seed + int64(y)*int64(scene.RenderOptions.Dx) + int64(t.x0))
		for x := t.x0; x < t.x1; x++ {
			renderPixel(ctx, scene, camera, rand, stats, x, y, img)
		}
	}
}

// startProgressBar displays the rendering progress.
func startProgressBar(ctx context.Context, totalTiles int, tilesCompleted *uint64) chan struct{} {
	progressDone := make(chan struct{})
	// go func() {
	// 	ticker := time.NewTicker(200 * time.Millisecond)
	// 	defer ticker.Stop()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		case <-progressDone:
	// 			fmt.Printf("\rRendering: 100%% complete\n")
	// 			return
	// 		case <-ticker.C:
	// 			completed := atomic.LoadUint64(tilesCompleted)
	// 			percent := float64(completed) / float64(totalTiles) * 100
	// 			fmt.Printf("\rRendering: %.2f%% complete", percent)
	// 		}
	// 	}
	// }()
	return progressDone
}

// fillRenderQueue populates the render queue with tiles.
func fillRenderQueue(ctx context.Context, dx, dy, tileSize int, renderQueue chan tile) {
	numTilesX := (dx + tileSize - 1) / tileSize
	numTilesY := (dy + tileSize - 1) / tileSize
	for ty := 0; ty < numTilesY; ty++ {
		for tx := 0; tx < numTilesX; tx++ {
			select {
			case <-ctx.Done():
				return
			case renderQueue <- tile{
				x0: tx * tileSize,
				x1: min((tx+1)*tileSize, dx),
				y0: ty * tileSize,
				y1: min((ty+1)*tileSize, dy),
			}:
			}
		}
	}
	close(renderQueue)
}

func renderScene(ctx context.Context, scene *Scene, camera Camera) (RenderArtifact, error) {
	t0 := time.Now()
	dx := scene.RenderOptions.Dx
	dy := scene.RenderOptions.Dy
	img := image.NewRGBA(image.Rect(0, 0, dx, dy))
	stats := RenderStats{}
	stats.Dx = dx
	stats.Dy = dy

	ctxScene, cancel := context.WithCancel(ctx)
	defer cancel()

	numWorkers := runtime.NumCPU()
	tileSize := 16
	numTilesX := (dx + tileSize - 1) / tileSize
	numTilesY := (dy + tileSize - 1) / tileSize
	totalTiles := numTilesX * numTilesY
	var tilesCompleted uint64

	renderQueue := make(chan tile, numWorkers)
	progressBar := startProgressBar(ctxScene, totalTiles, &tilesCompleted)

	// Start worker goroutines.
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			workerStats := RenderStats{}
			tilesCompleted := uint64(0)
			for t := range renderQueue {
				if ctxScene.Err() != nil {
					return
				}
				renderTile(ctxScene, scene, camera, t, img, &workerStats)
				tilesCompleted++
			}

			// Accumulate workerStats into main stats.
			atomic.AddUint64(&tilesCompleted, tilesCompleted)
			atomic.AddUint64(&stats.TotalRays, workerStats.TotalRays)
			atomic.AddUint64(&stats.RaysExceededDepth, workerStats.RaysExceededDepth)
			atomic.AddUint64(&stats.RaysLeftScene, workerStats.RaysLeftScene)
		}(i)
	}

	go fillRenderQueue(ctxScene, dx, dy, tileSize, renderQueue)

	// Wait for workers to finish or an error to occur.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		cancel()
		close(progressBar)
		return RenderArtifact{}, ctx.Err()
	case <-done:
		close(progressBar)
	}
	stats.RenderTime = time.Since(t0)
	return RenderArtifact{Image: img, Stats: stats}, nil
}

func Render(ctx context.Context, scene *Scene) (output RenderArtifact, err error) {
	err = scene.Validate()
	if err != nil {
		return RenderArtifact{}, fmt.Errorf("invalid scene: %v", err)
	}
	// Select the first camera in the scene.
	// We already know there is at least one camera in the scene.
	camera := scene.Camera[0]
	output, err = renderScene(ctx, scene, camera)
	if err != nil {
		return RenderArtifact{}, fmt.Errorf("failed to render scene: %v", err)
	}
	return output, nil
}
