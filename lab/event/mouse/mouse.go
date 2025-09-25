// Package mouse provide a mouse event type.
package mouse

import (
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
	"github.com/scottlawsonbc/slam/code/photon/raytrace/lab/event/key"
)

// Event is a mouse event.
type Event struct {
	// Point is the position of the mouse event in pixels.
	Point r2.Point

	// Button is the mouse button being pressed or released. Its value may be
	// zero, for a mouse move or drag without any button change.
	Button Button

	// Modifiers is a bitmask representing a set of modifier keys:
	// key.ModShift, key.ModAlt, etc.
	Modifiers key.Modifiers

	// Direction is the direction of the mouse event: DirPress, DirRelease,
	// or DirNone (for mouse moves or drags).
	Direction key.Direction
}

// Button is a mouse button.
type Button int32

const (
	ButtonLeft     Button = 0
	ButtonMiddle   Button = 1
	ButtonRight    Button = 2
	ButtonBackward Button = 3
	ButtonForward  Button = 4
)
