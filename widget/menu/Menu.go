// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

import (
	"fmt"
	"github.com/richardwilkes/geom"
	"github.com/richardwilkes/ui"
	"github.com/richardwilkes/ui/border"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/draw"
	"github.com/richardwilkes/ui/event"
	"github.com/richardwilkes/ui/layout"
	"github.com/richardwilkes/ui/layout/flex"
	"github.com/richardwilkes/ui/menu"
	"github.com/richardwilkes/ui/widget"
	"github.com/richardwilkes/ui/widget/window"
)

type Menu struct {
	widget.Block
	item           *MenuItem
	wnd            ui.Window
	attachToBottom bool
}

func NewMenu(title string) *Menu {
	mnu := &Menu{item: NewMenuItem(title, 0, nil)}
	mnu.item.menu = mnu
	mnu.Describer = func() string {
		return fmt.Sprintf("Menu #%d (%s)", mnu.ID(), mnu.Title())
	}
	mnu.SetBorder(border.NewLine(color.Gray, geom.Insets{Top: 1, Left: 1, Bottom: 1, Right: 1}))
	mnu.item.EventHandlers().Add(event.SelectionType, mnu.open)
	flex.NewLayout(mnu).SetEqualColumns(true)
	return mnu
}

// Title returns the title of this menu.
func (mnu *Menu) Title() string {
	return mnu.item.Title()
}

// AddItem appends an item to the end of this menu.
func (mnu *Menu) AddItem(item menu.Item) {
	switch actual := item.(type) {
	case *MenuItem:
		mnu.AddChild(actual)
		actual.EventHandlers().Add(event.ClosingType, mnu.close)
		actual.SetLayoutData(flex.NewData().SetHorizontalGrab(true).SetHorizontalAlignment(draw.AlignFill))
	case *Separator:
		mnu.AddChild(actual)
		actual.SetLayoutData(flex.NewData().SetHorizontalGrab(true).SetHorizontalAlignment(draw.AlignFill))
	}
}

// AddMenu appends an item with a sub-menu to the end of this menu.
func (mnu *Menu) AddMenu(subMenu menu.Menu) {
	if actual, ok := subMenu.(*Menu); ok {
		mnu.AddItem(actual.item)
	}
}

// InsertItem inserts an item at the specified item index within this menu.
func (mnu *Menu) InsertItem(index int, item menu.Item) {
	switch actual := item.(type) {
	case *MenuItem:
		mnu.AddChildAtIndex(actual, index)
		actual.EventHandlers().Add(event.ClosingType, mnu.close)
		actual.SetLayoutData(flex.NewData().SetHorizontalGrab(true).SetHorizontalAlignment(draw.AlignFill))
	case *Separator:
		mnu.AddChildAtIndex(actual, index)
		actual.SetLayoutData(flex.NewData().SetHorizontalGrab(true).SetHorizontalAlignment(draw.AlignFill))
	}
}

// InsertMenu inserts an item with a sub-menu at the specified item index within this menu.
func (mnu *Menu) InsertMenu(index int, subMenu menu.Menu) {
	if actual, ok := subMenu.(*Menu); ok {
		mnu.InsertItem(index, actual.item)
	}
}

// Remove the item at the specified index from this menu.
func (mnu *Menu) Remove(index int) {
	mnu.RemoveChildAtIndex(index)
}

// Count of items in this menu.
func (mnu *Menu) Count() int {
	return len(mnu.Children())
}

// Item at the specified index, or nil.
func (mnu *Menu) Item(index int) menu.Item {
	switch actual := mnu.Children()[index].(type) {
	case *MenuItem:
		return actual
	case *Separator:
		return actual
	}
	panic("Invalid child")
}

func (mnu *Menu) adjustItems(evt event.Event) {
	var largest float64
	for _, child := range mnu.Children() {
		switch item := child.(type) {
		case *MenuItem:
			pos := item.calculateAcceleratorPosition()
			if largest < pos {
				largest = pos
			}
		}
	}
	for _, child := range mnu.Children() {
		switch item := child.(type) {
		case *MenuItem:
			item.pos = largest
		}
	}
}

func (mnu *Menu) open(evt event.Event) {
	bounds := mnu.item.Bounds()
	where := mnu.item.ToWindow(bounds.Point)
	where.Add(mnu.item.Window().ContentFrame().Point)
	if mnu.attachToBottom {
		where.Y += bounds.Height
	} else {
		where.X += bounds.Width
	}
	mnu.adjustItems(nil)
	_, pref, _ := mnu.Layout().Sizes(layout.NoHintSize)
	mnu.SetBounds(geom.Rect{Size: pref})
	mnu.Layout().Layout()
	wnd := window.NewWindowWithContentSize(where, pref, window.BorderlessWindowMask)
	wnd.RootWidget().AddChild(mnu)
	wnd.EventHandlers().Add(event.FocusLostType, mnu.close)
	wnd.ToFront()
	mnu.wnd = wnd
	mnu.item.menuOpen = true
	mnu.item.Repaint()
}

func (mnu *Menu) close(evt event.Event) {
	if mnu.wnd != nil {
		wnd := mnu.wnd
		mnu.wnd = nil
		wnd.Close()
		mnu.item.menuOpen = false
		mnu.item.Repaint()
	}
}

// Dispose releases any operating system resources associated with this menu. It will also
// call Dispose() on all menu items it contains.
func (mnu *Menu) Dispose() {
	for _, child := range mnu.Children() {
		switch item := child.(type) {
		case menu.Item:
			item.Dispose()
		}
	}
}
