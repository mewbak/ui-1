// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

#include <Cocoa/Cocoa.h>
#include "_cgo_export.h"
#include "Cursor.h"

uiCursor platformSystemCursor(int id) {
	switch (id) {
		case platformArrowID:
			return [NSCursor arrowCursor];
		case platformTextID:
			return [NSCursor IBeamCursor];
		case platformVerticalTextID:
			return [NSCursor IBeamCursorForVerticalLayout];
		case platformCrossHairID:
			return [NSCursor crosshairCursor];
		case platformClosedHandID:
			return [NSCursor closedHandCursor];
		case platformOpenHandID:
			return [NSCursor openHandCursor];
		case platformPointingHandID:
			return [NSCursor pointingHandCursor];
		case platformResizeLeftID:
			return [NSCursor resizeLeftCursor];
		case platformResizeRightID:
			return [NSCursor resizeRightCursor];
		case platformResizeLeftRightID:
			return [NSCursor resizeLeftRightCursor];
		case platformResizeUpID:
			return [NSCursor resizeUpCursor];
		case platformResizeDownID:
			return [NSCursor resizeDownCursor];
		case platformResizeUpDownID:
			return [NSCursor resizeUpDownCursor];
		case platformDisappearingItemID:
			return [NSCursor disappearingItemCursor];
		case platformNotAllowedID:
			return [NSCursor operationNotAllowedCursor];
		case platformDragLinkID:
			return [NSCursor dragLinkCursor];
		case platformDragCopyID:
			return [NSCursor dragCopyCursor];
		case platformContextMenuID:
			return [NSCursor contextualMenuCursor];
		default:
			return NULL;
	}
}

uiCursor platformNewCursor(void *img, float hotX, float hotY) {
	return [[[NSCursor alloc] initWithImage:img hotSpot:NSMakePoint(hotX,hotY)] retain];
}

void platformDisposeCursor(uiCursor cursor) {
	[((NSCursor *)cursor) release];
}
