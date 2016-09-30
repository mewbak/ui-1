// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package menu

// Bar represents a set of menus.
type Bar interface {
	// AddMenu appends a menu to the end of this bar.
	AddMenu(menu Menu)
	// InsertMenu inserts a menu at the specified menu index within this bar.
	InsertMenu(index int, menu Menu)
	// Remove the menu at the specified index from this bar.
	Remove(index int)
	// Count of menus in this bar.
	Count() int
	// Menu at the specified index, or nil.
	Menu(index int) Menu
}