// Package size defines an event for changing the size of the window.
package size

import (
	"github.com/scottlawsonbc/slam/code/photon/raytrace/r2"
)

type Event struct {
	Size r2.Vec
}
