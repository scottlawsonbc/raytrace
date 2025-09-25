// Package key defines an event for physical keyboard keys.
package key

import (
	"fmt"
	"strings"
)

// Event is a key event.
type Event struct {
	// Rune is the meaning of the key event as determined by the
	// operating system. The mapping is determined by system-dependent
	// current layout, modifiers, lock-states, etc.
	//
	// If non-negative, it is a Unicode codepoint: pressing the 'a' key
	// generates different Runes 'a' or 'A' (but the same Code) depending on
	// the state of the shift key.
	//
	// If -1, the key does not generate a Unicode codepoint. To distinguish
	// them, look at Code.
	Rune rune

	// Code is the identity of the physical key relative to a notional
	// "standard" keyboard, independent of current layout, modifiers,
	// lock-states, etc
	Code string

	// Modifiers is a bitmask representing a set of modifier keys: ModShift,
	// ModAlt, etc.
	Modifiers Modifiers

	// Direction is the direction of the key event: DirPress, DirRelease,
	// or DirNone (for key repeats).
	Direction Direction
}

func (e Event) String() string {
	if e.Rune >= 0 {
		return fmt.Sprintf("key.Event{%q (%v), %v, %v}", e.Rune, e.Code, e.Modifiers, e.Direction)
	}
	return fmt.Sprintf("key.Event{(%v), %v, %v}", e.Code, e.Modifiers, e.Direction)
}

// Direction is the direction of the key event.
type Direction uint8

const (
	DirNone    Direction = 0
	DirPress   Direction = 1
	DirRelease Direction = 2
)

// Modifiers is a bitmask representing a set of modifier keys.
type Modifiers uint32

const (
	ModShift   Modifiers = 1 << 0
	ModControl Modifiers = 1 << 1
	ModAlt     Modifiers = 1 << 2
	ModMeta    Modifiers = 1 << 3 // called "Command" on OS X.
)

var mods = [...]struct {
	m Modifiers
	s string
}{
	{ModShift, "Shift"},
	{ModControl, "Control"},
	{ModAlt, "Alt"},
	{ModMeta, "Meta"},
}

func (m Modifiers) String() string {
	var match []string
	for _, mod := range mods {
		if mod.m&m != 0 {
			match = append(match, mod.s)
		}
	}
	return "key.Modifiers(" + strings.Join(match, "|") + ")"
}

func (d Direction) String() string {
	switch d {
	case DirNone:
		return "None"
	case DirPress:
		return "Press"
	case DirRelease:
		return "Release"
	default:
		return fmt.Sprintf("key.Direction(%d)", d)
	}
}
