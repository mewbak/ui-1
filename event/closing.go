package event

import (
	"bytes"
	"fmt"
)

// Closing is generated when a window is asked to close.
type Closing struct {
	target   Target
	finished bool
	aborted  bool
}

// NewClosing creates a new Closing event. 'target' is the window that is being asked to close.
func NewClosing(target Target) *Closing {
	return &Closing{target: target}
}

// Type returns the event type ID.
func (e *Closing) Type() Type {
	return ClosingType
}

// Target the original target of the event.
func (e *Closing) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *Closing) Cascade() bool {
	return false
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *Closing) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *Closing) Finish() {
	e.finished = true
}

// Aborted returns true if closing should not proceed.
func (e *Closing) Aborted() bool {
	return e.aborted
}

// Abort marks this event as being aborted and done.
func (e *Closing) Abort() {
	e.aborted = true
	e.finished = true
}

// String implements the fmt.Stringer interface.
func (e *Closing) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Closing[")
	if e.aborted {
		buffer.WriteString("Aborted, ")
	}
	buffer.WriteString(fmt.Sprintf("Target: %v", e.target))
	if e.finished {
		buffer.WriteString(", Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
