// Package instrument defines a stable event protocol for scientific and
// industrial instrument drivers.
//
// Drivers communicate exclusively by publishing and subscribing to events on a
// shared bus. The bus decouples participants, keeps timing predictable, and
// allows drivers to evolve independently.
//
// # Mental model
//
// Each driver is analogous to a VI with a front panel: controls (settable
// state) and indicators (readable state). Devices appear and disappear as
// Arrived/Removed. UI rendering is requested with Paint and acknowledged with
// Painted. Local, low-latency pointer motion is reported as Hovered.
//
// # Design rules
//
//  1. Publish on the bus, not directly to peers.
//  2. React quickly. Under nominal conditions a Controlled or Indicated
//     response should arrive within 100 ms of its request.
//  3. Publish a response only when the requested condition holds. For control
//     changes, do not send Controlled until the device is actually at the
//     requested state, not merely scheduled to be applied.
//  4. When an operation fails or is refused, set Err to a clear,
//     human-readable message.
//  5. Keep events precise. Prefer one event per semantic step.
//
// # Canonical flows
//
// Control:
//
//	Sequencer publishes Control{Drv,Ctl,Val,Req}.
//	Target driver applies the change.
//	Target driver publishes Controlled{Req,Drv,Ctl,Val,Err?}.
//
// Acquisition (example):
//
//	Sequencer publishes a Control that initiates acquisition (e.g., a grab).
//	Device publishes Grabbed with the acquired data.
//
// UI paint:
//
//	A driver or component publishes Paint with an image or a ref.
//	The main thread draws and publishes Painted{Err?}.
//
// Pointer input:
//
//	A GL window publishes Hovered while the cursor is inside the window.
//	Handlers may use it for reactive visuals or joystick-like control.
//
// Sequencing:
//
//	A sequencer publishes Sequenced{Pos,Err?} to mark step completion.
//
// # Health checks
//
// At startup, drivers are checked to confirm they meet latency expectations.
// A driver that cannot satisfy the protocol deadlines prevents the system from
// starting. This front-loads failures and preserves determinism.
package instrument

import (
	"fmt"
	"time"

	"github.com/scottlawsonbc/slam/code/photon/camera"
)

/*
Event protocol (authoritative contract).

Drivers communicate only through Event values carried on the Bus. Each Event
has a Type and a Data payload with exactly one pointer field set, selected by
Type. Names and wire formats are stable.

Timing guidance:
  - Controlled and Indicated should arrive within 100 ms of the corresponding
    Control or Indicate request in nominal conditions.
  - Grabbed timing is device-dependent and may stream.

Correctness rule for control:
  - Controlled must not be published until the device is at the requested
    state. This keeps the sequencer simple and avoids out-of-band polling.

UI guidance:
  - Hovered is frequent and lightweight. Downstream consumers should avoid
    heavy work directly in response to each hover sample.

Error reporting:
  - When Err is non-empty in a response event, it contains a human-readable
    diagnostic suitable for logs and UIs.
*/

// EventType categorizes the kind of event carried on the bus.
//
// Each value below maps to a payload struct in Data. Only the field matching
// the Type should be non-nil in Event.Data; all others must be nil.
type EventType string

const (
	// Arrived announces a new instrument (driver) has arrived in the system.
	Arrived EventType = "Arrived"
	// Removed announces an instrument (driver) has been removed from the system.
	Removed EventType = "Removed"

	// Grabbed announces newly acquired data from a driver.
	Grabbed EventType = "Grabbed"

	// Paint requests main thread to paint the graphics window with an image.
	Paint EventType = "Paint"
	// Painted notifies that a display target refreshed with a frame.
	Painted EventType = "Painted"

	// Indicate requests to get a read-only indicator value from a driver.
	Indicate EventType = "Indicate"
	// Indicated replies to an Indicate request with the value or error.
	Indicated EventType = "Indicated"

	// Control requests to set a control on a driver.
	Control EventType = "Control"
	// Controlled replies to a Control request with the applied value.
	Controlled EventType = "Controlled"

	// Sequenced reports the completion of a step in a sequence.
	Sequenced EventType = "Sequenced"

	// Hovered reports pointer motion while the mouse cursor is over a GL window.
	// Hovered events are produced by UI surfaces (e.g., gl) rather than
	// hardware drivers. Consumers can treat Hovered as a low-latency, local
	// input stream to drive reactive visualizations or joystick-like controls.
	Hovered EventType = "Hovered"
)

// Event carries a single event on the bus.
type Event struct {
	// Time is the timestamp when the event was created.
	// It is set by the publisher of the event, not by the bus.
	Time time.Time
	// From names the publisher that emitted the event.
	// From is expected to be a short, unique identifier
	// such as a driver name or UI surface name.
	From string
	// Type reports the kind of event. It is one of the EventType values.
	Type EventType
	// Data carries the payload for the event.
	// It should have exactly one non-nil pointer field, selected by Type.
	// This keeps the wire format stable and easy to marshal.
	Data Data
}

// String reports a compact textual form for logs.
func (e Event) String() string {
	return fmt.Sprintf(
		"Event{Time=%s, From=%s, Type=%s, Data=%s}",
		e.Time.Format(time.RFC3339Nano),
		e.From,
		string(e.Type),
		e.Data.String(),
	)
}

// Data groups payloads for all event kinds.
//
// Data should have exactly one pointer non-nil, selected by Event.Type. This
// keeps the wire format stable and easy to marshal. New event kinds should add
// a new pointer field rather than overloading an existing one.
type Data struct {
	Arrived *DataArrived
	Removed *DataRemoved

	Paint   *DataPaint
	Painted *DataPainted

	Indicate  *DataIndicate
	Indicated *DataIndicated

	Control    *DataControl
	Controlled *DataControlled

	Grabbed   *DataGrabbed
	Sequenced *DataSequenced

	Hovered *DataHovered
}

// String reports a compact textual form for logs.
func (d Data) String() string {
	switch {
	case d.Arrived != nil:
		return fmt.Sprintf("Data{Arrived=%s}", *d.Arrived)
	case d.Removed != nil:
		return fmt.Sprintf("Data{Removed=%s}", *d.Removed)
	case d.Paint != nil:
		return fmt.Sprintf("Data{Paint=%s}", *d.Paint)
	case d.Painted != nil:
		return fmt.Sprintf("Data{Painted=%s}", *d.Painted)
	case d.Indicate != nil:
		return fmt.Sprintf("Data{Indicate=%s}", *d.Indicate)
	case d.Indicated != nil:
		return fmt.Sprintf("Data{Indicated=%s}", *d.Indicated)
	case d.Control != nil:
		return fmt.Sprintf("Data{Control=%s}", *d.Control)
	case d.Controlled != nil:
		return fmt.Sprintf("Data{Controlled=%s}", *d.Controlled)
	case d.Grabbed != nil:
		return fmt.Sprintf("Data{Grabbed=%s}", *d.Grabbed)
	case d.Sequenced != nil:
		return fmt.Sprintf("Data{Sequenced=%s}", *d.Sequenced)
	case d.Hovered != nil:
		return fmt.Sprintf("Data{Hovered=%s}", *d.Hovered)
	default:
		return "Data{}"
	}
}

// DataArrived describes an Arrived event.
type DataArrived struct {
	// Drv names the driver that was added.
	Drv string
	// Pkg names the package that provided the driver.
	Pkg string
}

// String reports a compact textual form for logs.
func (d DataArrived) String() string {
	return fmt.Sprintf("DataArrived{Drv=%s, Pkg=%s}", d.Drv, d.Pkg)
}

// DataRemoved describes a Removed event.
type DataRemoved struct {
	// Drv names the driver that was removed.
	Drv string
	// Pkg names the package that provided the driver.
	Pkg string
}

// String reports a compact textual form for logs.
func (d DataRemoved) String() string {
	return fmt.Sprintf("DataRemoved{Drv=%s, Pkg=%s}", d.Drv, d.Pkg)
}

// DataPaint describes a Paint request to draw an image or ref.
type DataPaint struct {
	// Drv names the driver requesting the paint.
	Drv string
	// Img is the image to paint.
	Img camera.Frame
	// Ref is the reference to the image to paint.
	Ref string
}

// String reports a compact textual form for logs.
func (d DataPaint) String() string {
	if d.Img.Image == nil {
		return fmt.Sprintf("DataPaint{Drv=%s, Img=empty, Ref=%q}", d.Drv, d.Ref)
	}
	b := d.Img.Image.Bounds()
	dx := b.Dx()
	dy := b.Dy()
	if dx == 0 || dy == 0 {
		return fmt.Sprintf("DataPaint{Drv=%s, Img=empty, Ref=%q}", d.Drv, d.Ref)
	}
	return fmt.Sprintf("DataPaint{Drv=%s, Img=%dx%d, Ref=%q}", d.Drv, dx, dy, d.Ref)
}

// DataPainted describes a Painted acknowledgement.
type DataPainted struct {
	// Drv names the driver that painted the image.
	Drv string
	// Ref is the reference to the image that was painted.
	Ref string
	// Err is the error message if the paint failed or was refused.
	Err string
}

// String reports a compact textual form for logs.
func (d DataPainted) String() string {
	if d.Err != "" {
		return fmt.Sprintf("DataPainted{Drv=%s, Ref=%q, Err=%s}", d.Drv, d.Ref, d.Err)
	}
	return fmt.Sprintf("DataPainted{Drv=%s, Ref=%q}", d.Drv, d.Ref)
}

// DataControl describes a Control request.
type DataControl struct {
	// Req identifier to match the request to the response.
	Req string
	// Drv names the driver to apply the control change.
	Drv string
	// Ctl names the control to set.
	Ctl string
	// Val sets the desired value for the control.
	Val string
}

// String reports a compact textual form for logs.
func (d DataControl) String() string {
	return fmt.Sprintf("DataControl{Req=%s, Drv=%s, Ctl=%s, Val=%s}", d.Req, d.Drv, d.Ctl, d.Val)
}

// DataControlled describes a Controlled response.
type DataControlled struct {
	// Req echoes the value passed in from the Control request.
	Req string
	// Drv names the driver to apply the control change.
	Drv string
	// Ctl names the control to set.
	Ctl string
	// Val sets the desired value for the control.
	Val string
	// Err is the error message if the request failed or was refused.
	Err string
}

// String reports a compact textual form for logs.
func (d DataControlled) String() string {
	if d.Err != "" {
		return fmt.Sprintf("DataControlled{Req=%s, Drv=%s, Ctl=%s, Val=%s, Err=%s}", d.Req, d.Drv, d.Ctl, d.Val, d.Err)
	}
	return fmt.Sprintf("DataControlled{Req=%s, Drv=%s, Ctl=%s, Val=%s}", d.Req, d.Drv, d.Ctl, d.Val)
}

// DataIndicate describes an Indicate request.
type DataIndicate struct {
	// Req identifier to match the request to the response.
	Req string
	// Drv names the driver that is requested to read the indicator.
	Drv string
	// Ind names the indicator to read.
	Ind string
}

// String reports a compact textual form for logs.
func (d DataIndicate) String() string {
	return fmt.Sprintf("DataIndicate{Req=%s, Drv=%s, Ind=%s}", d.Req, d.Drv, d.Ind)
}

// DataIndicated describes an Indicated response.
type DataIndicated struct {
	// Req echoes the value passed in from the Indicate request.
	Req string
	// Drv names the driver that read the indicator.
	Drv string
	// Ind is the indicator name that was read.
	Ind string
	// Val is the value that was read.
	Val string
	// Err is the error message if the request failed or was refused.
	Err string
}

// String reports a compact textual form for logs.
func (d DataIndicated) String() string {
	if d.Err != "" {
		return fmt.Sprintf("DataIndicated{Req=%s, Drv=%s, Ind=%s, Err=%s}", d.Req, d.Drv, d.Ind, d.Err)
	}
	return fmt.Sprintf("DataIndicated{Req=%s, Drv=%s, Ind=%s, Val=%s}", d.Req, d.Drv, d.Ind, d.Val)
}

// DataGrabbed describes a Grabbed event with a camera frame.
type DataGrabbed struct {
	// Drv names the driver that grabbed the data.
	Drv string
	// Img holds the grabbed camera frame.
	Img camera.Frame
}

// String reports a compact textual form for logs.
func (d DataGrabbed) String() string {
	b := d.Img.Image.Bounds()
	dx := b.Dx()
	dy := b.Dy()
	if dx == 0 || dy == 0 {
		return fmt.Sprintf("DataGrabbed{Drv=%s, Img=empty}", d.Drv)
	}
	return fmt.Sprintf("DataGrabbed{Drv=%s, Img=%dx%d}", d.Drv, dx, dy)
}

// DataSequenced describes a Sequenced event.
type DataSequenced struct {
	// Drv names the driver that orchestrated the sequence step.
	Drv string
	// Pos is the position in the sequence.
	Pos int
	// Err is the error message if the sequence step failed or was refused.
	Err string
}

// String reports a compact textual form for logs.
func (d DataSequenced) String() string {
	if d.Err != "" {
		return fmt.Sprintf("DataSequenced{Drv=%s, Pos=%d, Err=%s}", d.Drv, d.Pos, d.Err)
	}
	return fmt.Sprintf("DataSequenced{Drv=%s, Pos=%d}", d.Drv, d.Pos)
}

// DataHovered describes a Hovered event emitted by a UI window.
//
// Coordinates use a window-local pixel space with the origin at the
// top-left corner of the client area. X increases to the right. Y increases
// downward. W and H describe the window's drawable size in pixels at the
// time of the event. Consumers can compute normalized coordinates as:
//
//	nx := X / float64(W)
//	ny := Y / float64(H)
//
// DataHovered is small and frequent. Consumers should avoid heavy work
// directly on Hovered and instead derive lightweight, low-latency reactions.
type DataHovered struct {
	// Drv names the logical owner of the window (typically a driver name).
	Drv string
	// X is the cursor X position in window pixels from the left edge.
	X float64
	// Y is the cursor Y position in window pixels from the top edge.
	Y float64
	// W is the window width in pixels at the time of the event.
	W int
	// H is the window height in pixels at the time of the event.
	H int
}

// String reports a compact textual form for logs.
func (d DataHovered) String() string {
	return fmt.Sprintf("DataHovered{Drv=%s, X=%.2f, Y=%.2f, W=%d, H=%d}", d.Drv, d.X, d.Y, d.W, d.H)
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"math"
// 	"syscall/js"
// 	"time"
// 	"unicode/utf8"

// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/key"
// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/mouse"
// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/size"
// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/wheel"
// 	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
// )

// func newMouseEvent(e js.Value) mouse.Event {
// 	event := mouse.Event{
// 		Point: r2.Point{
// 			X: e.Get("clientX").Float(),
// 			Y: e.Get("clientY").Float(),
// 		},
// 		Button: mouse.Button(e.Get("button").Int()),
// 	}
// 	switch e.Get("type").String() {
// 	case "mousedown":
// 		event.Direction = key.DirPress
// 	case "mouseup":
// 		event.Direction = key.DirRelease
// 	case "mousemove":
// 		event.Direction = key.DirNone
// 	default:
// 		panic(fmt.Errorf("unsupported event type: %q", e.Get("type").String()))
// 	}
// 	if e.Get("altKey").Bool() {
// 		event.Modifiers |= key.ModAlt
// 	}
// 	if e.Get("ctrlKey").Bool() {
// 		event.Modifiers |= key.ModControl
// 	}
// 	if e.Get("shiftKey").Bool() {
// 		event.Modifiers |= key.ModShift
// 	}
// 	if e.Get("metaKey").Bool() {
// 		event.Modifiers |= key.ModMeta
// 	}
// 	return event
// }

// func newKeyEvent(e js.Value) key.Event {
// 	event := key.Event{Code: e.Get("code").String()}
// 	switch {
// 	case e.Get("repeat").Bool():
// 		event.Direction = key.DirNone
// 	case e.Get("type").String() == "keydown":
// 		event.Direction = key.DirPress
// 	case e.Get("type").String() == "keyup":
// 		event.Direction = key.DirRelease
// 	default:
// 		panic(fmt.Errorf("unsupported event type: %q", e.Get("type").String()))
// 	}
// 	if e.Get("altKey").Bool() {
// 		event.Modifiers |= key.ModAlt
// 	}
// 	if e.Get("ctrlKey").Bool() {
// 		event.Modifiers |= key.ModControl
// 	}
// 	if e.Get("shiftKey").Bool() {
// 		event.Modifiers |= key.ModShift
// 	}
// 	if e.Get("metaKey").Bool() {
// 		event.Modifiers |= key.ModMeta
// 	}
// 	k := e.Get("key").String()
// 	if utf8.RuneCountInString(k) == 1 {
// 		event.Rune = []rune(k)[0]
// 	} else {
// 		event.Rune = -1
// 	}
// 	return event
// }

// func newWheelEvent(e js.Value) wheel.Event {
// 	var deltaX, deltaY float64
// 	// We just map any positive or negative value to 1 or -1 or 0.
// 	if e.Get("deltaX").Float() > 0 {
// 		deltaX = 1
// 	} else if e.Get("deltaX").Float() < 0 {
// 		deltaX = -1
// 	}
// 	if e.Get("deltaY").Float() > 0 {
// 		deltaY = 1
// 	} else if e.Get("deltaY").Float() < 0 {
// 		deltaY = -1
// 	}
// 	event := wheel.Event{Delta: r2.Vec{X: deltaX, Y: deltaY}}
// 	return event
// }

// func redirectKeyboardEvent(to chan any) js.Func {
// 	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		to <- newKeyEvent(args[0])
// 		return nil
// 	})
// }

// func redirectMouseEvent(to chan any) js.Func {
// 	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		to <- newMouseEvent(args[0])
// 		return nil
// 	})
// }

// func redirectWheelEvent(to chan any) js.Func {
// 	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		to <- newWheelEvent(args[0])
// 		return nil
// 	})
// }

// // autoResize continuously resizes the canvas to fit the window.
// func autoResize(events chan<- any) {
// 	doc := js.Global().Get("document")
// 	canvas := doc.Call("getElementById", "canvas")
// 	prevWidth := doc.Get("body").Get("clientWidth").Float()
// 	prevHeight := doc.Get("body").Get("clientHeight").Float()
// 	canvas.Set("width", prevWidth)
// 	canvas.Set("height", prevHeight)
// 	ticker := time.NewTicker(10 * time.Millisecond)
// 	for range ticker.C {
// 		width := doc.Get("body").Get("clientWidth").Float()
// 		height := doc.Get("body").Get("clientHeight").Float()
// 		if width != prevWidth {
// 			canvas.Set("width", width)
// 			prevWidth = width
// 			events <- size.Event{Size: r2.Vec{X: math.Min(width, height), Y: math.Min(width, height)}}
// 		}
// 		if height != prevHeight {
// 			canvas.Set("height", height)
// 			prevHeight = height
// 			events <- size.Event{Size: r2.Vec{X: math.Min(width, height), Y: math.Min(width, height)}}
// 		}
// 	}
// }

// // logEventFilter logs the event in a human-readable format and passes it on.
// func logEventFilter(e any) any {
// 	// Print the type name and the event.
// 	log.Printf("%T%+v\n", e, e)
// 	return e
// }
