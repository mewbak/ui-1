package event

import (
	"bytes"
	"fmt"

	"github.com/richardwilkes/toolbox/xmath/geom"
	"github.com/richardwilkes/ui/keys"
)

// MouseWheel is generated when the mouse wheel is used over a widget.
type MouseWheel struct {
	target    Target
	delta     geom.Point
	where     geom.Point
	modifiers keys.Modifiers
	finished  bool
}

// NewMouseWheel creates a new MouseWheel event. 'target' is the widget that mouse is over. 'delta'
// is the amount the wheel was moved on each axis. 'where' is the location in the window where the
// mouse is. 'modifiers' are the keyboard modifiers keys that were down.
func NewMouseWheel(target Target, delta geom.Point, where geom.Point, modifiers keys.Modifiers) *MouseWheel {
	return &MouseWheel{target: target, delta: delta, where: where, modifiers: modifiers}
}

// Type returns the event type ID.
func (e *MouseWheel) Type() Type {
	return MouseWheelType
}

// Target the original target of the event.
func (e *MouseWheel) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *MouseWheel) Cascade() bool {
	return true
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *MouseWheel) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *MouseWheel) Finish() {
	e.finished = true
}

// Delta returns the amount the wheel was moved on each axis.
func (e *MouseWheel) Delta() geom.Point {
	return e.delta
}

// Where returns the location in the window the mouse is.
func (e *MouseWheel) Where() geom.Point {
	return e.where
}

// Modifiers returns the key modifiers that were down.
func (e *MouseWheel) Modifiers() keys.Modifiers {
	return e.modifiers
}

// String implements the fmt.Stringer interface.
func (e *MouseWheel) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("MouseWheel[Delta: [%v], Where: [%v], Target: %v", e.delta, e.where, e.target))
	modifiers := e.modifiers.String()
	if modifiers != "" {
		buffer.WriteString(", ")
		buffer.WriteString(modifiers)
	}
	if e.finished {
		buffer.WriteString(", Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
