// Package paint defines an event for the app being ready to paint.
package paint

// Event indicates that the app is ready to paint the next frame of the GUI.
//
// A frame is completed by calling the App's Publish method.
type Event struct {
	// Reason is the reason the paint event was sent. For debugging purposes.
	Reason string

	// External is true for paint events sent by the screen driver.
	//
	// An external event may be sent at any time in response to an
	// operating system event, for example the window opened, was
	// resized, or the screen memory was lost.
	//
	// Programs actively drawing to the screen as fast as vsync allows
	// should ignore external paint events to avoid a backlog of paint
	// events building up.
	External bool
}
