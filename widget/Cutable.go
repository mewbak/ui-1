// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package widget

import (
	"github.com/richardwilkes/i18n"
	"github.com/richardwilkes/ui/event"
	"github.com/richardwilkes/ui/menu"
)

// Cutable defines the methods required of objects that can respond to the Cut menu item.
type Cutable interface {
	// CanCut returns true if Cut() can be called successfully.
	CanCut() bool
	// Cut the data from the object and copy it to the clipboard.
	Cut()
}

// AddCutItem adds the standard Cut menu item to the specified menu.
func AddCutItem(m *menu.Menu) *menu.Item {
	item := m.AddItem(i18n.Text("Cut"), "x")
	handlers := item.EventHandlers()
	handlers.Add(event.SelectionType, func(evt event.Event) {
		window := KeyWindow()
		if window != nil {
			focus := window.Focus()
			if c, ok := focus.(Cutable); ok {
				c.Cut()
			}
		}
	})
	handlers.Add(event.ValidateType, func(evt event.Event) {
		valid := false
		window := KeyWindow()
		if window != nil {
			focus := window.Focus()
			if c, ok := focus.(Cutable); ok {
				valid = c.CanCut()
			}
		}
		if !valid {
			evt.(*event.Validate).MarkInvalid()
		}
	})
	return item
}