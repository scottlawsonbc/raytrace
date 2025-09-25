package wheel

import "github.com/scottlawsonbc/slam/code/photon/raytrace/r2"

// Event represents a mouse wheel event.
type Event struct {
	// Delta is the amount the wheel was scrolled.
	// X is the horizontal scroll amount, and Y is the vertical scroll amount.
	Delta r2.Vec
}

func (e Event) String() string {
	return "wheel.Event{" + e.Delta.String() + "}"
}
