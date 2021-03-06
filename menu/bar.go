package menu

import (
	"github.com/richardwilkes/ui/event"
)

// Special menu ids.
const (
	ServicesMenu SpecialMenuType = iota
	WindowMenu
	HelpMenu
)

// SpecialMenuType identifies which special menu is being referred to.
type SpecialMenuType int

// Bar represents a set of menus.
type Bar interface {
	// AppendMenu appends a menu at the end of this bar.
	AppendMenu(menu Menu)
	// InsertMenu inserts a menu at the specified item index within this bar. Pass in a negative
	// index to append to the end.
	InsertMenu(menu Menu, index int)
	// Remove the menu at the specified index from this bar.
	Remove(index int)
	// Count of menus in this bar.
	Count() int
	// Menu at the specified index, or nil.
	Menu(index int) Menu
	// SpecialMenu returns the specified special menu, or nil if it has not been setup.
	SpecialMenu(which SpecialMenuType) Menu
	// SetupSpecialMenu sets up the specified special menu, which must have already been installed
	// into the menu bar.
	SetupSpecialMenu(which SpecialMenuType, menu Menu)
	// ProcessKeyDown is called to process KeyDown events prior to anything else receiving them.
	ProcessKeyDown(evt *event.KeyDown)
}

var (
	// AppBar returns the menu bar for the given window id. On some platforms, the menu bar is a
	// global entity and the same value will be returned for all window ids.
	AppBar func(id uint64) Bar
	// Global returns true if the menu bar is global.
	Global func() bool
)
