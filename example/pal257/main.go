// Copyright 2024 Scott Lawson. All rights reserved.
// Program main renders a simple orbiting scene and displays it in an OpenGL
// window. The program uses package [gl] for the window and publishes UI
// events onto an [instrument.Bus]. We log file drop, key up/down, and wheel
// events.
//
// Scrubbing behavior:
//   - Click (mouse down + up without dragging) toggles play/pause.
//   - While paused, click+drag adjusts the frame *relative* to pointer movement
//     (smooth, no jumps). When resuming play, playback continues from the
//     scrubbed frame (no jump back to the old animation position).
package main

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/gl"
	"github.com/scottlawsonbc/slam/code/photon/instrument"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/phys"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r3"
)

type config struct {
	windowTitle string
	windowX     int
	windowY     int

	wireframeDx           phys.Distance
	wireframeDy           phys.Distance
	wireframeDz           phys.Distance
	wireframeEdgeRadius   phys.Distance
	wireframeVertexRadius phys.Distance

	renderFPS int
	renderDx  int
	renderDy  int

	orbitLookAt    r3.Point
	orbitPeriod    time.Duration
	orbitRadius    phys.Distance
	orbitLookFromZ phys.Distance

	colorOrigin     phys.Spectrum
	colorVertex     phys.Spectrum
	colorWireframeX phys.Spectrum
	colorWireframeY phys.Spectrum
	colorWireframeZ phys.Spectrum

	fovDx         phys.Distance
	fovDy         phys.Distance
	fovCenter     r3.Point
	fovWireRadius phys.Distance

	propAxesRadius phys.Distance
	propAxesLength phys.Distance
}

type orbitCamera struct {
	lookAt     r3.Point
	lookFromZ  phys.Distance
	vup        r3.Vec
	intrinsics phys.CameraIntrinsics
	radius     phys.Distance
	frameSpan  int
}

func (o orbitCamera) at(i int) phys.Camera {
	theta := 2 * math.Pi * (float64(i) / float64(max(1, o.frameSpan)))
	lookFrom := r3.Point{
		X: float64(o.radius) * math.Cos(theta),
		Y: float64(o.radius) * math.Sin(theta),
		Z: float64(o.lookFromZ),
	}
	return phys.NewCalibratedCamera(
		o.intrinsics,
		phys.CameraExtrinsics{LookFrom: lookFrom, LookAt: o.lookAt, VUp: o.vup},
	)
}

type frameCache struct{ buf []*image.RGBA }

func newFrameCache(n int) *frameCache        { return &frameCache{buf: make([]*image.RGBA, n)} }
func (fc *frameCache) has(i int) bool        { return i >= 0 && i < len(fc.buf) && fc.buf[i] != nil }
func (fc *frameCache) get(i int) *image.RGBA { return fc.buf[i] }
func (fc *frameCache) put(i int, img *image.RGBA) {
	if i >= 0 && i < len(fc.buf) {
		fc.buf[i] = img
	}
}
func (fc *frameCache) full() bool {
	for _, f := range fc.buf {
		if f == nil {
			return false
		}
	}
	return true
}

// uiState holds interaction state for pause/play + scrubbing.
type uiState struct {
	mu sync.RWMutex

	paused    bool
	mouseDown bool
	dragging  bool

	downX float64
	downY float64

	// Relative-drag fields
	dragLastX float64 // last X processed (for incremental delta)
	dragAccum float64 // fractional frames accumulated (can be +/-)
	scrubIdx  int     // current frame to show when paused
	winW      int     // last known window width for scaling

	// latest frame index actually shown (render loop updates)
	displayedIdx int64

	// resume request (set when switching from pause->play)
	resumePending int32 // 0/1
	resumeToIdx   int32
}

func (s *uiState) setPaused(p bool) { s.mu.Lock(); s.paused = p; s.mu.Unlock() }

func (s *uiState) setMouseDown(x, y float64) {
	s.mu.Lock()
	s.mouseDown, s.dragging = true, false
	s.downX, s.downY = x, y
	s.dragLastX = x
	s.dragAccum = 0
	s.mu.Unlock()
}

func (s *uiState) setMouseUp()  { s.mu.Lock(); s.mouseDown, s.dragging = false, false; s.mu.Unlock() }
func (s *uiState) setDragging() { s.mu.Lock(); s.dragging = true; s.mu.Unlock() }

func (s *uiState) setScrub(idx, w int) { s.mu.Lock(); s.scrubIdx, s.winW = idx, w; s.mu.Unlock() }

func (s *uiState) get() (paused, mouseDown, dragging bool, scrubIdx int) {
	s.mu.RLock()
	paused, mouseDown, dragging, scrubIdx = s.paused, s.mouseDown, s.dragging, s.scrubIdx
	s.mu.RUnlock()
	return
}

func (s *uiState) getDown() (x, y float64) {
	s.mu.RLock()
	x, y = s.downX, s.downY
	s.mu.RUnlock()
	return
}

func (s *uiState) storeDisplayed(idx int) { atomic.StoreInt64(&s.displayedIdx, int64(idx)) }
func (s *uiState) loadDisplayed() int     { return int(atomic.LoadInt64(&s.displayedIdx)) }

func (s *uiState) addRelativeDelta(deltaFrames float64, framesPerOrbit int) {
	s.mu.Lock()
	s.dragAccum += deltaFrames
	// apply only the integer portion; keep fractional remainder for smoothness
	if s.dragAccum >= 1 || s.dragAccum <= -1 {
		step := int(math.Trunc(s.dragAccum)) // towards zero
		s.dragAccum -= float64(step)
		s.scrubIdx = wrapIndex(s.scrubIdx+step, framesPerOrbit)
	}
	s.mu.Unlock()
}

func (s *uiState) requestResumeAt(idx int) {
	atomic.StoreInt32(&s.resumeToIdx, int32(idx))
	atomic.StoreInt32(&s.resumePending, 1)
}

func (s *uiState) takeResume() (idx int, ok bool) {
	if atomic.SwapInt32(&s.resumePending, 0) == 1 {
		return int(atomic.LoadInt32(&s.resumeToIdx)), true
	}
	return 0, false
}

func newScene(cfg config) phys.Scene {
	scene := phys.Scene{
		RenderOptions: phys.RenderOptions{
			Seed:         0,
			RaysPerPixel: 1,
			MaxRayDepth:  6,
			Dx:           cfg.renderDx,
			Dy:           cfg.renderDy,
		},
		Light:  []phys.Light{},
		Camera: []phys.Camera{},
		Node: []phys.Node{
			// Wireframe vertices (top)
			node("P1", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, -float64(cfg.wireframeDx)/2, -float64(cfg.wireframeDy)/2, float64(cfg.wireframeDz), cfg.wireframeVertexRadius),
			node("P2", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, +float64(cfg.wireframeDx)/2, -float64(cfg.wireframeDy)/2, float64(cfg.wireframeDz), cfg.wireframeVertexRadius),
			node("P3", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, +float64(cfg.wireframeDx)/2, +float64(cfg.wireframeDy)/2, float64(cfg.wireframeDz), cfg.wireframeVertexRadius),
			node("P4", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, -float64(cfg.wireframeDx)/2, +float64(cfg.wireframeDy)/2, float64(cfg.wireframeDz), cfg.wireframeVertexRadius),
			// Wireframe vertices (bottom)
			node("P5", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, -float64(cfg.wireframeDx)/2, -float64(cfg.wireframeDy)/2, 0, cfg.wireframeVertexRadius),
			node("P6", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, +float64(cfg.wireframeDx)/2, -float64(cfg.wireframeDy)/2, 0, cfg.wireframeVertexRadius),
			node("P7", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, +float64(cfg.wireframeDx)/2, +float64(cfg.wireframeDy)/2, 0, cfg.wireframeVertexRadius),
			node("P8", phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorVertex}}, -float64(cfg.wireframeDx)/2, +float64(cfg.wireframeDy)/2, 0, cfg.wireframeVertexRadius),

			// Axes at origin
			phys.PropAxes(r3.Point{}, cfg.propAxesRadius, cfg.propAxesLength, "")[0],
			phys.PropAxes(r3.Point{}, cfg.propAxesRadius, cfg.propAxesLength, "")[1],
			phys.PropAxes(r3.Point{}, cfg.propAxesRadius, cfg.propAxesLength, "")[2],

			// Edges of the box (12 cylinders)
			{Name: "AxisX1Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeX}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: float64(cfg.wireframeDz)}, Direction: r3.Vec{X: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDx}},
			{Name: "AxisX2Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeX}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: +float64(cfg.wireframeDy) / 2, Z: float64(cfg.wireframeDz)}, Direction: r3.Vec{X: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDx}},
			{Name: "AxisY1Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeY}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: float64(cfg.wireframeDz)}, Direction: r3.Vec{Y: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDy}},
			{Name: "AxisY2Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeY}},
				Shape: phys.Cylinder{Origin: r3.Point{X: +float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: float64(cfg.wireframeDz)}, Direction: r3.Vec{Y: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDy}},
			{Name: "AxisX1Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeX}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{X: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDx}},
			{Name: "AxisX2Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeX}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: +float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{X: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDx}},
			{Name: "AxisZ1Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeZ}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: +float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Z: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDz}},
			{Name: "AxisZ2Top", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeZ}},
				Shape: phys.Cylinder{Origin: r3.Point{X: +float64(cfg.wireframeDx) / 2, Y: +float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Z: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDz}},
			{Name: "AxisZ1Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeZ}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Z: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDz}},
			{Name: "AxisZ2Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeZ}},
				Shape: phys.Cylinder{Origin: r3.Point{X: +float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Z: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDz}},
			{Name: "AxisY1Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeY}},
				Shape: phys.Cylinder{Origin: r3.Point{X: -float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Y: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDy}},
			{Name: "AxisY2Bottom", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorWireframeY}},
				Shape: phys.Cylinder{Origin: r3.Point{X: +float64(cfg.wireframeDx) / 2, Y: -float64(cfg.wireframeDy) / 2, Z: 0}, Direction: r3.Vec{Y: 1}, Radius: cfg.wireframeEdgeRadius, Height: cfg.wireframeDy}},

			{Name: "FOVX1", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorOrigin}},
				Shape: phys.Cylinder{Origin: r3.Point{X: float64(cfg.fovCenter.X), Y: float64(cfg.fovCenter.Y), Z: float64(cfg.fovCenter.Z)}, Direction: r3.Vec{X: 1}, Radius: cfg.fovWireRadius, Height: cfg.fovDx}},
			{Name: "FOVX2", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorOrigin}},
				Shape: phys.Cylinder{Origin: r3.Point{X: float64(cfg.fovCenter.X), Y: float64(cfg.fovCenter.Y), Z: float64(cfg.fovCenter.Z)}, Direction: r3.Vec{X: -1}, Radius: cfg.fovWireRadius, Height: cfg.fovDx}},
			{Name: "FOVY1", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorOrigin}},
				Shape: phys.Cylinder{Origin: r3.Point{X: float64(cfg.fovCenter.X), Y: float64(cfg.fovCenter.Y), Z: float64(cfg.fovCenter.Z)}, Direction: r3.Vec{Y: 1}, Radius: cfg.fovWireRadius, Height: cfg.fovDy}},
			{Name: "FOVY2", Material: phys.Emitter{Texture: phys.TextureUniform{Color: cfg.colorOrigin}},
				Shape: phys.Cylinder{Origin: r3.Point{X: float64(cfg.fovCenter.X), Y: float64(cfg.fovCenter.Y), Z: float64(cfg.fovCenter.Z)}, Direction: r3.Vec{Y: -1}, Radius: cfg.fovWireRadius, Height: cfg.fovDy}},

			beam("Beam1", r3.Point{X: float64(26.15 * phys.MM), Y: float64(34.63 * phys.MM), Z: 0}, r3.Point{X: 0, Y: 0, Z: float64(50 * phys.MM)}),
			{
				Name:     "PAL257",
				Material: phys.Emitter{Texture: phys.MustNewTextureImage("./asset/pal257_top.png", "", "")},
				Shape: phys.Quad{
					Center: r3.Point{X: 0, Y: 0, Z: 0},
					Width:  100 * phys.MM,
					Height: 100 * phys.MM,
					Normal: r3.Vec{Z: 1},
				},
			},
		},
	}
	return scene
}

func main() {
	runtime.LockOSThread()
	cfg := config{
		renderDx:  1440 / 2,
		renderDy:  1080 / 2,
		renderFPS: 30,

		windowTitle: "pal257",
		windowX:     800,
		windowY:     100,

		wireframeDx:           100 * phys.MM,
		wireframeDy:           100 * phys.MM,
		wireframeDz:           50 * phys.MM,
		wireframeEdgeRadius:   0.5 * phys.MM,
		wireframeVertexRadius: 1 * phys.MM,

		orbitPeriod:    8 * time.Second,
		orbitLookAt:    r3.Point{},
		orbitLookFromZ: 100 * phys.MM,
		orbitRadius:    200 * phys.MM,

		colorOrigin:     phys.Spectrum{X: 0.5, Y: 0.5, Z: 0.5},
		colorVertex:     phys.Spectrum{X: 203.0, Y: 136.0, Z: 206.0}.Divs(255),
		colorWireframeX: phys.Spectrum{X: 255, Y: 0, Z: 157.0}.Divs(255),
		colorWireframeY: phys.Spectrum{X: 157, Y: 255, Z: 0}.Divs(255),
		colorWireframeZ: phys.Spectrum{X: 0, Y: 57.0, Z: 255}.Divs(255),

		fovDx:         5 * phys.MM,
		fovDy:         5 * phys.MM,
		fovCenter:     r3.Point{X: 0, Y: 0, Z: float64(50 * phys.MM)},
		fovWireRadius: 0.1 * phys.MM,

		propAxesRadius: 0.1 * phys.MM,
		propAxesLength: 50 * phys.MM,
	}
	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func run(cfg config) error {
	framesPerOrbit := int(math.Round(cfg.orbitPeriod.Seconds() * float64(cfg.renderFPS)))
	if framesPerOrbit < 1 {
		framesPerOrbit = 1
	}

	// Event bus + subscriber
	bus := instrument.NewBus(64)
	defer bus.Close()
	events := bus.SubscribeNamed("scene", 64)
	defer bus.Unsubscribe(events)

	var ui uiState
	const dragThreshold = 1.0 // pixels to consider a drag

	// Event handler goroutine
	go func() {
		for e := range events {
			switch e.Type {
			// ---------- Pointer ----------
			case instrument.PointerEnter:
				d := e.Data.PointerEnter
				log.Printf("ui.enter   from=%s x=%.1f y=%.1f w=%d h=%d", e.From, d.X, d.Y, d.W, d.H)

			case instrument.PointerLeave:
				d := e.Data.PointerLeave
				log.Printf("ui.leave   from=%s x=%.1f y=%.1f w=%d h=%d", e.From, d.X, d.Y, d.W, d.H)
				// Leaving does not change paused/position.

			case instrument.PointerDown:
				d := e.Data.PointerDown
				log.Printf("ui.down    from=%s x=%.1f y=%.1f w=%d h=%d btn=%d mods=%d clicks=%d",
					e.From, d.X, d.Y, d.W, d.H, d.Button, d.Mods, d.Clicks)
				ui.setMouseDown(d.X, d.Y)
				// If paused, initialize relative-drag baseline to current view
				if paused, _, _, _ := ui.get(); paused {
					if d.W > 0 {
						ui.setScrub(ui.loadDisplayed(), d.W)
					}
				}

			case instrument.PointerMove:
				d := e.Data.PointerMove
				downX, downY := ui.getDown()
				if dx, dy := d.X-downX, d.Y-downY; (dx*dx + dy*dy) >= (dragThreshold * dragThreshold) {
					ui.setDragging()
				}
				paused, mouseDown, _, _ := ui.get()
				if paused && mouseDown && d.W > 0 && framesPerOrbit > 0 {
					// Relative scrubbing: convert deltaX since last move into fractional frames
					framesPerPixel := float64(framesPerOrbit) / float64(d.W)
					ui.mu.Lock()
					dx := d.X - ui.dragLastX
					ui.dragLastX = d.X
					ui.mu.Unlock()
					ui.addRelativeDelta(dx*framesPerPixel, framesPerOrbit)
				}

			case instrument.PointerUp:
				d := e.Data.PointerUp
				log.Printf("ui.up      from=%s x=%.1f y=%.1f w=%d h=%d btn=%d mods=%d clicks=%d",
					e.From, d.X, d.Y, d.W, d.H, d.Button, d.Mods, d.Clicks)

				// Click (down+up without drag) toggles pause/play.
				paused, mouseDown, dragging, scrub := ui.get()
				if mouseDown && !dragging {
					if !paused {
						// Play -> Pause: freeze to currently displayed frame
						ui.setPaused(true)
						ui.setScrub(ui.loadDisplayed(), d.W)
						log.Printf("playback: paused @ %d", ui.loadDisplayed())
					} else {
						// Pause -> Play: resume from current scrub index (no jump)
						ui.requestResumeAt(scrub)
						ui.setPaused(false)
						log.Printf("playback: playing from %d", scrub)
					}
				}
				ui.setMouseUp()

			case instrument.PointerCancel:
				d := e.Data.PointerCancel
				log.Printf("ui.cancel  from=%s x=%.1f y=%.1f w=%d h=%d reason=%q", e.From, d.X, d.Y, d.W, d.H, d.Reason)
				ui.setMouseUp()

			// ---------- Keyboard ----------
			case instrument.KeyDown:
				d := e.Data.KeyDown
				log.Printf("key.down   from=%s key=%d scancode=%d mods=%d repeat=%v", e.From, d.Key, d.Scancode, d.Mods, d.Repeat)
			case instrument.KeyUp:
				d := e.Data.KeyUp
				log.Printf("key.up     from=%s key=%d scancode=%d mods=%d", e.From, d.Key, d.Scancode, d.Mods)

			// ---------- Wheel / Drop ----------
			case instrument.Wheel:
				d := e.Data.Wheel
				log.Printf("wheel      from=%s dx=%.2f dy=%.2f", e.From, d.OffX, d.OffY)
			case instrument.Dropped:
				d := e.Data.Dropped
				log.Printf("drop       from=%s files=%v", e.From, d.Names)
			}
		}
	}()

	// Create window
	win, err := gl.New(gl.Options{
		Name:   cfg.windowTitle,
		X:      cfg.windowX,
		Y:      cfg.windowY,
		Width:  cfg.renderDx,
		Height: cfg.renderDy,
		Bus:    bus,
	})
	if err != nil {
		return fmt.Errorf("window: %w", err)
	}
	defer win.Close()

	scene := newScene(cfg)
	cache := newFrameCache(framesPerOrbit)

	ctx := context.Background()
	frameDur := time.Second / time.Duration(cfg.renderFPS)
	ticker := time.NewTicker(frameDur)
	defer ticker.Stop()

	orbit := orbitCamera{
		intrinsics: phys.IntrinsicsFireflyDLGeneric6mm,
		lookAt:     cfg.orbitLookAt,
		vup:        r3.Vec{Z: -1},
		radius:     cfg.orbitRadius,
		lookFromZ:  cfg.orbitLookFromZ,
		frameSpan:  framesPerOrbit,
	}

	fmt.Printf("bbox: %v\n", scene.Node[0].Shape.Bounds())

	var firstSaved bool
	var fpsCount int
	lastFPS := time.Now()
	var frame int

	for {
		select {
		case <-ticker.C:
		default:
			<-ticker.C
		}

		if win.ShouldClose() {
			return nil
		}
		win.PollEvents()

		paused, _, _, scrub := ui.get()

		// If we just switched from pause->play, align animation start to scrub index.
		if !paused {
			if idx, ok := ui.takeResume(); ok {
				frame = clampInt(idx, 0, framesPerOrbit-1)
			}
		}

		// Compute animation index after possible resume alignment.
		idxAnim := frame % framesPerOrbit
		showIdx := idxAnim

		if paused && cache.full() {
			showIdx = clampInt(scrub, 0, framesPerOrbit-1)
		}

		var img image.Image
		if cache.has(showIdx) {
			img = cache.get(showIdx)
		} else {
			scene.Camera = []phys.Camera{orbit.at(showIdx)}
			res, err := phys.Render(ctx, &scene)
			if err != nil {
				return fmt.Errorf("render: %w", err)
			}
			img = res.Image
			cache.put(showIdx, cloneRGBA(img))
			if !firstSaved {
				path := time.Now().Format("./out/out_20060102_150405.png")
				if err := phys.SavePNG(path, res.Image); err != nil {
					return fmt.Errorf("save first frame: %w", err)
				}
				log.Printf("saved first frame --> %s\n", path)
				firstSaved = true
			}
		}

		win.Draw(img)
		ui.storeDisplayed(showIdx)

		// advance animation only when playing or still filling cache
		if !paused || !cache.full() {
			frame++
		}

		fpsCount++
		if time.Since(lastFPS) >= time.Second {
			mode := "render"
			if cache.full() {
				if paused {
					mode = "paused"
				} else {
					mode = "playback"
				}
			}
			fmt.Printf("fps≈%d  mode=%s  frame=%d/%d  showIdx=%d\n",
				fpsCount, mode, (idxAnim + 1), framesPerOrbit, showIdx)
			fpsCount = 0
			lastFPS = time.Now()
		}
	}
}

func beam(name string, from r3.Point, to r3.Point) phys.Node {
	dir := to.Sub(from).Unit()
	cylRadius := 0.1 * phys.MM
	cylHeight := phys.Distance(to.Sub(from).Length())
	node := phys.Node{
		Name: name,
		Shape: phys.Cylinder{
			Origin:    from,
			Direction: dir,
			Radius:    cylRadius,
			Height:    cylHeight,
		},
		Material: phys.Emitter{
			Texture: phys.TextureUniform{Color: phys.Spectrum{X: 0, Y: 1, Z: 1}},
		},
	}
	if err := node.Validate(); err != nil {
		panic(fmt.Sprintf("invalid beam %q: %v", name, err))
	}
	return node
}

func node(name string, mat phys.Material, x, y, z float64, radius phys.Distance) phys.Node {
	return phys.Node{
		Name:     name,
		Shape:    phys.Sphere{Center: r3.Point{X: x, Y: y, Z: z}, Radius: radius},
		Material: mat,
	}
}

func cloneRGBA(src image.Image) *image.RGBA {
	r := src.Bounds()
	dst := image.NewRGBA(r)
	draw.Draw(dst, r, src, r.Min, draw.Src)
	return dst
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func wrapIndex(i, n int) int {
	if n <= 0 {
		return 0
	}
	i %= n
	if i < 0 {
		i += n
	}
	return i
}
