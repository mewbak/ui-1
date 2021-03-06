package event

import (
	"bytes"
	"fmt"
)

// Resized is generated when a widget is resized.
type Resized struct {
	target   Target
	finished bool
}

// NewResized creates a new Resized event. 'target' is the widget that was resized.
func NewResized(target Target) *Resized {
	return &Resized{target: target}
}

// Type returns the event type ID.
func (e *Resized) Type() Type {
	return ResizedType
}

// Target the original target of the event.
func (e *Resized) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *Resized) Cascade() bool {
	return false
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *Resized) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *Resized) Finish() {
	e.finished = true
}

// String implements the fmt.Stringer interface.
func (e *Resized) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Resized[Target: %v", e.target))
	if e.finished {
		buffer.WriteString(", Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
