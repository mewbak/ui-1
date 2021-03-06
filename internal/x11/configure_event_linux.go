package x11

import (
	// #cgo pkg-config: x11
	// #include <X11/Xlib.h>
	"C"
	"unsafe"
)

type ConfigureEvent C.XConfigureEvent

func (evt *ConfigureEvent) Window() Window {
	return Window(evt.window)
}

func (evt *ConfigureEvent) ToEvent() *Event {
	return (*Event)(unsafe.Pointer(evt))
}
