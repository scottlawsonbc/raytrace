//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"math"
	"sync"
	"syscall/js"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/obj"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type worker struct {
	scene        phys.Scene
	camera       phys.OrthographicCamera
	cameraTheta  float64
	cameraPhi    float64
	cameraRadius float64
	renderMutex  sync.Mutex
	renderDirty  bool
	isRendering  bool
}

func loadNodes(fsys fs.FS, objPath string) ([]phys.Node, error) {
	parsedObj, err := obj.ParseFS(fsys, objPath)
	if err != nil {
		return nil, err
	}
	nodes, err := phys.ConvertObjectToNodes(parsedObj, fsys)
	if err != nil {
		return nil, err
	}
	fmt.Printf("got %d nodes\n", len(nodes))
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found in OBJ file")
	}
	for i, n := range nodes {
		fmt.Printf("Node %d: %v\n", i, n.Shape.Bounds())
	}
	return nodes, nil
}

func (w *worker) init() {
	var fsys fs.FS
	var objPath string
	fsys = NewHTTPFS("/assets/camera-small.obj") // Adjust the base URL as needed
	objPath = "camera-small.obj"
	nodes, err := loadNodes(fsys, objPath)
	if err != nil {
		panic(err)
	}

	// // Initialize the scene and camera
	// nodes := []phys.Node{
	// 	{
	// 		Name: "triangle",
	// 		Shape: phys.TriangleUV{
	// 			P0:     r3.Point{X: 0, Y: 0, Z: 0},
	// 			P1:     r3.Point{X: 1, Y: 0, Z: 0},
	// 			P2:     r3.Point{X: 0, Y: 1, Z: 0},
	// 			UV0:    r2.Point{X: 0, Y: 0},
	// 			UV1:    r2.Point{X: 1, Y: 0},
	// 			UV2:    r2.Point{X: 0, Y: 1},
	// 			Normal: r3.Vec{X: 0, Y: 0, Z: 1},
	// 		},
	// 		Material: phys.DebugUV{},
	// 	},
	// }

	scene := phys.Scene{
		Node:   nodes,
		Camera: []phys.Camera{}, // Filled in automatically.
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,
			MaxRayDepth:  10,
			Dx:           256,
			Dy:           256,
		},
	}
	bounds := scene.Bounds()
	db := bounds.Max.Sub(bounds.Min).Get(bounds.LongestAxis())

	cam := phys.OrthographicCamera{
		LookFrom:  r3.Point{X: 0.5, Y: 0.5, Z: 20.0},
		LookAt:    r3.Point{X: 0.5, Y: 0.5, Z: 0},
		VUp:       r3.Vec{X: 0, Y: 1, Z: 0},
		FOVHeight: phys.Distance(db * 1.5),
		FOVWidth:  phys.Distance(db * 1.5),
	}

	w.camera = cam
	w.scene = scene
	w.computeSphericalCoordinates()
	w.isRendering = false
	w.renderDirty = false
}

func (w *worker) setRenderDirty(dirty bool) {
	w.renderMutex.Lock()
	w.renderDirty = dirty
	w.renderMutex.Unlock()
}

func (w *worker) scheduleRender() {
	// Use setTimeout to schedule the render after the current event loop
	js.Global().Call("setTimeout", js.FuncOf(func(js.Value, []js.Value) interface{} {
		w.render()
		return nil
	}), 0)
}

func (w *worker) onMessage(this js.Value, args []js.Value) interface{} {
	message := args[0].Get("data")
	messageType := message.Get("type").String()

	// Lock only when modifying shared state
	w.renderMutex.Lock()
	// Process the message and update the camera
	switch messageType {
	case "rotateCamera":
		dx := message.Get("dx").Float()
		dy := message.Get("dy").Float()
		w.rotateCamera(dx, dy)
	case "zoomCamera":
		delta := message.Get("delta").Float()
		w.zoomCamera(delta)
	case "translateCamera":
		dx := message.Get("dx").Float()
		dy := message.Get("dy").Float()
		w.translateCamera(dx, dy)
	default:
		log.Println("Unknown message type:", messageType)
		w.renderMutex.Unlock()
		return nil
	}

	if w.isRendering {
		// If rendering is in progress, mark as dirty
		w.renderDirty = true
		w.renderMutex.Unlock()
	} else {
		// Start a new render asynchronously
		w.isRendering = true
		w.renderMutex.Unlock()
		w.scheduleRender()
	}

	return nil
}

// computeSphericalCoordinates calculates the spherical coordinates (theta, phi, radius)
// based on the current camera position relative to its target.
func (w *worker) computeSphericalCoordinates() {
	dx := w.camera.LookFrom.X - w.camera.LookAt.X
	dy := w.camera.LookFrom.Y - w.camera.LookAt.Y
	dz := w.camera.LookFrom.Z - w.camera.LookAt.Z
	w.cameraRadius = math.Sqrt(dx*dx + dy*dy + dz*dz)
	w.cameraTheta = math.Atan2(dz, dx)           // azimuthal angle
	w.cameraPhi = math.Acos(dy / w.cameraRadius) // polar angle
}

func (w *worker) updateCameraPosition() {
	x := w.cameraRadius * math.Sin(w.cameraPhi) * math.Cos(w.cameraTheta)
	y := w.cameraRadius * math.Cos(w.cameraPhi)
	z := w.cameraRadius * math.Sin(w.cameraPhi) * math.Sin(w.cameraTheta)
	w.camera.LookFrom = r3.Point{
		X: w.camera.LookAt.X + x,
		Y: w.camera.LookAt.Y + y,
		Z: w.camera.LookAt.Z + z,
	}
}

func (w *worker) rotateCamera(dx, dy float64) {
	const sensitivity = 0.005
	w.cameraTheta += dx * sensitivity
	w.cameraPhi -= dy * sensitivity
	// Clamp phi to avoid gimbal lock.
	w.cameraPhi = math.Max(0.01, math.Min(math.Pi-0.01, w.cameraPhi))
	w.updateCameraPosition()
}

func (w *worker) zoomCamera(delta float64) {
	prevFOVHeight := w.camera.FOVHeight
	prevFOVWidth := w.camera.FOVWidth
	zoomFactor := math.Exp(delta * 0.1) // Exponential zoom for smoother scaling
	w.camera.FOVHeight = phys.Distance(math.Min(math.Max(0.1, float64(w.camera.FOVHeight)*zoomFactor), 100))
	w.camera.FOVWidth = phys.Distance(math.Min(math.Max(0.1, float64(w.camera.FOVWidth)*zoomFactor), 100))
	log.Printf("zoomed camera: %f -> %f, %f -> %f\n", prevFOVHeight, w.camera.FOVHeight, prevFOVWidth, w.camera.FOVWidth)
}

func (w *worker) translateCamera(dx, dy float64) {
	sensitivity := w.getSensitivity()
	// Calculate right and up vectors.
	right := r3.Vec{X: -math.Sin(w.cameraTheta), Y: 0, Z: math.Cos(w.cameraTheta)}
	up := r3.Vec{X: 0, Y: 1, Z: 0}
	// Compute translation vector.
	delta := right.Muls(dx * sensitivity).Add(up.Muls(dy * sensitivity))
	// Update camera positions.
	w.camera.LookFrom = w.camera.LookFrom.Add(delta)
	w.camera.LookAt = w.camera.LookAt.Add(delta)
}

func (w *worker) getSensitivity() float64 {
	return float64(w.camera.FOVHeight) * 0.001 // Adjust the multiplier as needed
}

func (w *worker) render() {
	for {
		// Prepare the scene data.
		w.renderMutex.Lock()
		sceneCopy := w.scene // Make a copy of the scene to work with
		sceneCopy.Camera = []phys.Camera{w.camera}
		w.renderMutex.Unlock()

		// Render the scene.
		artifact, err := phys.Render(context.Background(), &sceneCopy)
		if err != nil {
			logf("red", "phys.Render error: %v", err)
			w.renderMutex.Lock()
			w.isRendering = false
			w.renderMutex.Unlock()
			return
		}

		// Convert the image to a format that can be sent back to the main thread.
		bounds := artifact.Image.Bounds()
		width, height := bounds.Dx(), bounds.Dy()
		pixelData := artifact.Image.Pix
		jsPixelData := js.Global().Get("Uint8ClampedArray").New(len(pixelData))
		n := js.CopyBytesToJS(jsPixelData, pixelData)
		if n != len(pixelData) {
			logf("red", "copying pixel data failed: %d != %d", n, len(pixelData))
			w.renderMutex.Lock()
			w.isRendering = false
			w.renderMutex.Unlock()
			return
		}

		// Send result to main thread
		result := js.Global().Get("Object").New()
		result.Set("width", width)
		result.Set("height", height)
		result.Set("pixelData", jsPixelData)
		js.Global().Call("postMessage", result)

		w.renderMutex.Lock()
		if w.renderDirty {
			// Reset the dirty flag and continue the loop to render again
			w.renderDirty = false
			w.renderMutex.Unlock()
			// Allow other events to be processed before the next render
			w.scheduleRender()
			return // Exit current render to schedule next render
		} else {
			w.isRendering = false
			w.renderMutex.Unlock()
			break
		}
	}
}

func main() {
	var workerApp worker
	workerApp.init()
	js.Global().Call("postMessage", "Worker started")
	js.Global().Set("onmessage", js.FuncOf(workerApp.onMessage))
	js.Global().Set("onerror", js.FuncOf(onError))
	select {} // Keep the worker running.
}

func onError(this js.Value, args []js.Value) interface{} {
	errorMessage := args[0].Get("message").String()
	logf("red", "Worker internal error: %s", errorMessage)
	return nil
}

// logf logs a formatted message to the browser console with the specified color.
func logf(color string, format string, args ...interface{}) {
	console := js.Global().Get("console")
	if console.IsUndefined() {
		return
	}
	message := fmt.Sprintf(format, args...)
	styledMessage := fmt.Sprintf("%%c%s", message)
	css := fmt.Sprintf("color: %s;", color)
	console.Call("log", styledMessage, css)
}
