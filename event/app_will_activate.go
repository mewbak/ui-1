package event

import (
	"bytes"
)

// AppWillActivate is generated immediately prior to the application being brought to the
// foreground.
type AppWillActivate struct {
	target   Target
	finished bool
}

// SendAppWillActivate sends a new AppWillActivate event.
func SendAppWillActivate() {
	Dispatch(&AppWillActivate{target: GlobalTarget()})
}

// Type returns the event type ID.
func (e *AppWillActivate) Type() Type {
	return AppWillActivateType
}

// Target the original target of the event.
func (e *AppWillActivate) Target() Target {
	return e.target
}

// Cascade returns true if this event should be passed to its target's parent if not marked done.
func (e *AppWillActivate) Cascade() bool {
	return false
}

// Finished returns true if this event has been handled and should no longer be processed.
func (e *AppWillActivate) Finished() bool {
	return e.finished
}

// Finish marks this event as handled and no longer eligible for processing.
func (e *AppWillActivate) Finish() {
	e.finished = true
}

// String implements the fmt.Stringer interface.
func (e *AppWillActivate) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("AppWillActivate[")
	if e.finished {
		buffer.WriteString("Finished")
	}
	buffer.WriteString("]")
	return buffer.String()
}
