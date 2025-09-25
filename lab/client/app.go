//go:build js && wasm

package main

import (
	"fmt"
	"log"
	"syscall/js"

	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/client/canvas"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/key"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/lifecycle"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/mouse"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/paint"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/size"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/wheel"
)

type app struct {
	// filters holds a list of event filter functions applied to incoming events.
	filters []func(any) any

	// eventOut is the channel through which processed events are emitted to the main loop.
	eventOut chan any

	// eventIn is the channel where raw events are received from various sources.
	eventIn chan any

	// glctx is the canvas rendering context used for drawing images.
	glctx canvas.Context

	// Camera and scene configuration.
	isLeftButtonDown   bool
	isMiddleButtonDown bool
	lastMouse          r2.Point

	worker js.Value
}

func (a *app) initWorker() {
	if !js.Global().Get("Worker").Truthy() {
		log.Println("Web Workers are not supported in this environment.")
		return
	}
	var worker js.Value
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("failed to create worker: %v", r)
			}
		}()
		worker = js.Global().Get("Worker").New("worker.js")
	}()
	if err != nil {
		log.Println(err)
		return
	}
	a.worker = worker
	worker.Set("onmessage", js.FuncOf(a.onWorkerMessage))
	worker.Set("onerror", js.FuncOf(a.onWorkerError))
	logf("green", "Worker initialized")
}

func (a *app) onWorkerError(this js.Value, args []js.Value) interface{} {
	errorMessage := args[0].Get("message").String()
	log.Printf("Worker error: %s", errorMessage)
	return nil
}

func (a *app) onWorkerMessage(this js.Value, args []js.Value) interface{} {
	message := args[0].Get("data")
	// Check if the message is a string (e.g., "Worker started").
	if message.Type() == js.TypeString {
		logf("green", "recv from worker `%s`", message.String())
		a.Send(paint.Event{Reason: "worker start", External: true})
		return nil
	}

	width := message.Get("width").Int()
	height := message.Get("height").Int()
	jsPixelData := message.Get("pixelData")
	logf("green", "received image from worker %dx%d", width, height)

	imgData := js.Global().Get("ImageData").New(jsPixelData, width, height)
	a.glctx.Val.Call("putImageData", imgData, 0, 0)
	a.setRenderStatus(false)
	return nil
}

// setRenderStatus updates the visibility of the render status indicator in the DOM.
func (a *app) setRenderStatus(visible bool) {
	doc := js.Global().Get("document")
	statusElem := doc.Call("getElementById", "render-status")
	if visible {
		statusElem.Get("style").Set("visibility", "visible")
	} else {
		statusElem.Get("style").Set("visibility", "hidden")
	}
}

// RegisterFilter adds a new event filter to the application's filter chain.
// Filters are applied in the order they are registered to incoming events.
func (a *app) RegisterFilter(f func(any) any) {
	a.filters = append(a.filters, f)
}

// Filter sequentially applies all registered filters to an event.
// If any filter returns nil, the event is discarded.
func (a *app) Filter(e any) any {
	for _, f := range a.filters {
		if e = f(e); e == nil {
			return nil
		}
	}
	return e
}

// Events returns a read-only channel from which processed events can be received.
func (a *app) Events() <-chan any {
	return a.eventOut
}

// Send enqueues an event to the eventIn channel for processing.
// It utilizes the pump to ensure non-blocking behavior.
func (a *app) Send(e any) {
	a.eventIn <- e
}

func (a *app) sendRotateCameraMessage(dx, dy float64) {
	message := js.Global().Get("Object").New()
	message.Set("type", "rotateCamera")
	message.Set("dx", dx)
	message.Set("dy", dy)
	a.worker.Call("postMessage", message)
}

func (a *app) sendTranslateCameraMessage(dx, dy float64) {
	message := js.Global().Get("Object").New()
	message.Set("type", "translateCamera")
	message.Set("dx", dx)
	message.Set("dy", dy)
	a.worker.Call("postMessage", message)
}

func (a *app) sendZoomCameraMessage(delta float64) {
	message := js.Global().Get("Object").New()
	message.Set("type", "zoomCamera")
	message.Set("delta", delta)
	a.worker.Call("postMessage", message)
}

func (a *app) main() {
	for e := range a.Events() {
		switch e := a.Filter(e).(type) {
		case lifecycle.Event:
			// reserved
		case size.Event:
			// reserved
		case mouse.Event:
			switch e.Direction {
			case key.DirPress:
				if e.Button == mouse.ButtonLeft {
					a.isLeftButtonDown = true
				} else if e.Button == mouse.ButtonMiddle {
					a.isMiddleButtonDown = true
				}
				a.lastMouse = e.Point
			case key.DirRelease:
				if e.Button == mouse.ButtonLeft {
					a.isLeftButtonDown = false
				} else if e.Button == mouse.ButtonMiddle {
					a.isMiddleButtonDown = false
				}
			case key.DirNone:
				if a.isLeftButtonDown || a.isMiddleButtonDown {
					dx := e.Point.X - a.lastMouse.X
					dy := e.Point.Y - a.lastMouse.Y
					a.lastMouse = e.Point
					if a.isLeftButtonDown {
						a.sendRotateCameraMessage(dx, dy)
					} else if a.isMiddleButtonDown {
						a.sendTranslateCameraMessage(dx, dy)
					}
				}
			}
		case wheel.Event:
			a.sendZoomCameraMessage(e.Delta.Y)
		case key.Event:
			if e.Direction == key.DirPress {
				switch e.Code {
				case "ArrowLeft":
					a.sendTranslateCameraMessage(-0.1, 0)
				case "ArrowRight":
					a.sendTranslateCameraMessage(0.1, 0)
				case "ArrowUp":
					a.sendTranslateCameraMessage(0, -0.1)
				case "ArrowDown":
					a.sendTranslateCameraMessage(0, 0.1)
				}
			}
		default:
		}
	}
}

// stopPumping is a sentinel value used to signal the pump to stop forwarding events
// and close the destination channel gracefully.
type stopPumping struct{}

// pump creates an intermediary channel (src) that buffers events before forwarding them
// to the destination channel (dst). This mechanism ensures that sending to src is non-blocking,
// even if dst is temporarily unable to receive events.
//
// The pump function handles an internal buffer that dynamically resizes to accommodate
// bursts of incoming events. Sending a stopPumping value to src will initiate a graceful
// shutdown, ensuring all queued events are forwarded before closing dst.
//
// Usage:
//
//	src := pump(dst)
//	src <- event1
//	src <- event2
//	src <- stopPumping{}
func pump(dst chan interface{}) (src chan interface{}) {
	src = make(chan interface{})
	go func() {
		// initialSize is the initial size of the circular buffer. It must be a
		// power of 2.
		const initialSize = 16
		i, j, buf, mask := 0, 0, make([]interface{}, initialSize), initialSize-1
		srcActive := true
		for {
			maybeDst := dst
			if i == j {
				maybeDst = nil
			}
			if maybeDst == nil && !srcActive {
				break // Pump is stopped and empty.
			}
			select {
			case maybeDst <- buf[i&mask]:
				buf[i&mask] = nil
				i++
			case e := <-src:
				if _, ok := e.(stopPumping); ok {
					srcActive = false
					continue
				}
				if !srcActive {
					continue
				}
				// Allocate a bigger buffer if necessary.
				if i+len(buf) == j {
					b := make([]interface{}, 2*len(buf))
					n := copy(b, buf[j&mask:])
					copy(b[n:], buf[:j&mask])
					i, j = 0, len(buf)
					buf, mask = b, len(b)-1
				}
				buf[j&mask] = e
				j++
			}
		}
		close(dst)
		// Block forever to prevent goroutine from exiting.
		for range src {
		}
	}()
	return src
}
