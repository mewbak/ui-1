// Copyright (c) 2016 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package ui

import (
	"github.com/richardwilkes/geom"
	"github.com/richardwilkes/ui/color"
	"github.com/richardwilkes/ui/draw"
	"github.com/richardwilkes/ui/event"
	"github.com/richardwilkes/ui/theme"
	"math"
	"time"
)

const (
	scrollBarNone scrollBarPart = iota
	scrollBarThumb
	scrollBarLineUp
	scrollBarLineDown
	scrollBarPageUp
	scrollBarPageDown
)

type scrollBarPart int

// Pager objects can provide line and page information for scrolling.
type Pager interface {
	// LineScrollAmount is called to determine how far to scroll in the given direction to reveal
	// a full 'line' of content. A positive value should be returned regardless of the direction,
	// although negative values will behave as if they were positive.
	LineScrollAmount(horizontal, towardsStart bool) float64
	// PageScrollAmount is called to determine how far to scroll in the given direction to reveal
	// a full 'page' of content. A positive value should be returned regardless of the direction,
	// although negative values will behave as if they were positive.
	PageScrollAmount(horizontal, towardsStart bool) float64
}

// Scrollable objects can respond to ScrollBars.
type Scrollable interface {
	Pager
	// ScrolledPosition is called to determine the current position of the Scrollable.
	ScrolledPosition(horizontal bool) float64
	// SetScrolledPosition is called to set the current position of the Scrollable.
	SetScrolledPosition(horizontal bool, position float64)
	// VisibleSize is called to determine the size of the visible portion of the Scrollable.
	VisibleSize(horizontal bool) float64
	// ContentSize is called to determine the total size of the Scrollable.
	ContentSize(horizontal bool) float64
}

// ScrollBar represents a widget for controlling scrolling.
type ScrollBar struct {
	Block
	Target     Scrollable       // The target of the scrollbar.
	Theme      *theme.ScrollBar // The theme the scrollbar will use to draw itself.
	pressed    scrollBarPart
	sequence   int
	thumbDown  float64
	horizontal bool
}

// NewScrollBar creates a new scrollbar.
func NewScrollBar(horizontal bool, target Scrollable) *ScrollBar {
	sb := &ScrollBar{}
	sb.Theme = theme.StdScrollBar
	sb.Target = target
	sb.horizontal = horizontal
	sb.SetSizer(sb)
	handlers := sb.EventHandlers()
	handlers.Add(event.PaintType, sb.paint)
	handlers.Add(event.MouseDownType, sb.mouseDown)
	handlers.Add(event.MouseDraggedType, sb.mouseDragged)
	handlers.Add(event.MouseUpType, sb.mouseUp)
	return sb
}

// Sizes implements the Sizer interface.
func (sb *ScrollBar) Sizes(hint geom.Size) (min, pref, max geom.Size) {
	if sb.horizontal {
		min.Width = sb.Theme.Size * 2
		min.Height = sb.Theme.Size
		pref.Width = sb.Theme.Size * 2
		pref.Height = sb.Theme.Size
		max.Width = DefaultMax
		max.Height = sb.Theme.Size
	} else {
		min.Width = sb.Theme.Size
		min.Height = sb.Theme.Size * 2
		pref.Width = sb.Theme.Size
		pref.Height = sb.Theme.Size * 2
		max.Width = sb.Theme.Size
		max.Height = DefaultMax
	}
	if border := sb.Border(); border != nil {
		insets := border.Insets()
		min.AddInsets(insets)
		pref.AddInsets(insets)
		max.AddInsets(insets)
	}
	return min, pref, max
}

func (sb *ScrollBar) paint(evt event.Event) {
	bounds := sb.LocalInsetBounds()
	if sb.horizontal {
		bounds.Height = sb.Theme.Size
	} else {
		bounds.Width = sb.Theme.Size
	}
	bgColor := sb.baseBackground(scrollBarNone)
	gc := evt.(*event.Paint).GC()
	gc.SetColor(bgColor)
	gc.FillRect(bounds)
	gc.SetColor(bgColor.AdjustBrightness(sb.Theme.OutlineAdjustment))
	gc.StrokeRect(bounds)
	sb.drawLineButton(gc, scrollBarLineDown)
	if sb.pressed == scrollBarPageUp || sb.pressed == scrollBarPageDown {
		bounds = sb.partRect(sb.pressed)
		if !bounds.IsEmpty() {
			if sb.horizontal {
				bounds.Y++
				bounds.Height -= 2
			} else {
				bounds.X++
				bounds.Width -= 2
			}
			gc.SetColor(sb.baseBackground(sb.pressed))
			gc.FillRect(bounds)
		}
	}
	sb.drawThumb(gc)
	sb.drawLineButton(gc, scrollBarLineUp)
}

func (sb *ScrollBar) mouseDown(evt event.Event) {
	sb.sequence++
	where := sb.FromWindow(evt.(*event.MouseDown).Where())
	part := sb.over(where)
	if sb.partEnabled(part) {
		sb.pressed = part
		switch part {
		case scrollBarThumb:
			if sb.horizontal {
				sb.thumbDown = where.X - sb.partRect(part).X
			} else {
				sb.thumbDown = where.Y - sb.partRect(part).Y
			}
		case scrollBarLineUp, scrollBarLineDown, scrollBarPageUp, scrollBarPageDown:
			sb.scheduleRepeat(part, sb.Theme.InitialRepeatDelay)
		}
		sb.Repaint()
	}
}

func (sb *ScrollBar) mouseDragged(evt event.Event) {
	if sb.pressed == scrollBarThumb {
		var pos float64
		rect := sb.partRect(scrollBarLineUp)
		where := sb.FromWindow(evt.(*event.MouseDragged).Where())
		if sb.horizontal {
			pos = where.X - (sb.thumbDown + rect.X + rect.Width - 1)
		} else {
			pos = where.Y - (sb.thumbDown + rect.Y + rect.Height - 1)
		}
		sb.SetScrolledPosition(pos / sb.thumbScale())
	}
}

func (sb *ScrollBar) mouseUp(evt event.Event) {
	sb.pressed = scrollBarNone
	sb.Repaint()
}

func (sb *ScrollBar) scheduleRepeat(part scrollBarPart, delay time.Duration) {
	window := sb.Window()
	if window.Valid() {
		current := sb.sequence
		switch part {
		case scrollBarLineUp:
			sb.SetScrolledPosition(sb.Target.ScrolledPosition(sb.horizontal) - math.Abs(sb.Target.LineScrollAmount(sb.horizontal, true)))
		case scrollBarLineDown:
			sb.SetScrolledPosition(sb.Target.ScrolledPosition(sb.horizontal) + math.Abs(sb.Target.LineScrollAmount(sb.horizontal, false)))
		case scrollBarPageUp:
			sb.SetScrolledPosition(sb.Target.ScrolledPosition(sb.horizontal) - math.Abs(sb.Target.PageScrollAmount(sb.horizontal, true)))
		case scrollBarPageDown:
			sb.SetScrolledPosition(sb.Target.ScrolledPosition(sb.horizontal) + math.Abs(sb.Target.PageScrollAmount(sb.horizontal, false)))
		default:
			return
		}
		window.InvokeAfter(func() {
			if current == sb.sequence && sb.pressed == part {
				sb.scheduleRepeat(part, sb.Theme.RepeatDelay)
			}
		}, delay)
	}
}

func (sb *ScrollBar) over(where geom.Point) scrollBarPart {
	for i := scrollBarThumb; i <= scrollBarPageDown; i++ {
		rect := sb.partRect(i)
		if rect.Contains(where) {
			return i
		}
	}
	return scrollBarNone
}

func (sb *ScrollBar) thumbScale() float64 {
	var scale float64 = 1
	content := sb.Target.ContentSize(sb.horizontal)
	visible := sb.Target.VisibleSize(sb.horizontal)
	if content-visible > 0 {
		var size float64
		min := sb.Theme.Size * 0.75
		bounds := sb.LocalInsetBounds()
		if sb.horizontal {
			size = bounds.Width
		} else {
			size = bounds.Height
		}
		size -= sb.Theme.Size*2 + 2
		if size > 0 {
			scale = size / content
			visible *= scale
			if visible < min {
				scale = (size + visible - min) / content
			}
		}
	}
	return scale
}

func (sb *ScrollBar) partRect(part scrollBarPart) geom.Rect {
	var result geom.Rect
	switch part {
	case scrollBarThumb:
		if sb.Target != nil {
			content := sb.Target.ContentSize(sb.horizontal)
			visible := sb.Target.VisibleSize(sb.horizontal)
			if content-visible > 0 {
				pos := sb.Target.ScrolledPosition(sb.horizontal)
				full := sb.LocalInsetBounds()
				if sb.horizontal {
					full.X += sb.Theme.Size - 1
					full.Width -= sb.Theme.Size*2 - 2
					full.Height = sb.Theme.Size
					if full.Width > 0 {
						scale := full.Width / content
						visible *= scale
						min := sb.Theme.Size * 0.75
						if visible < min {
							scale = (full.Width + visible - min) / content
							visible = min
						}
						pos *= scale
						full.X += pos
						full.Width = visible + 1
						result = full
					}
				} else {
					full.Y += sb.Theme.Size - 1
					full.Height -= sb.Theme.Size*2 - 2
					full.Width = sb.Theme.Size
					if full.Height > 0 {
						scale := full.Height / content
						visible *= scale
						min := sb.Theme.Size * 0.75
						if visible < min {
							scale = (full.Height + visible - min) / content
							visible = min
						}
						pos *= scale
						full.Y += pos
						full.Height = visible + 1
						result = full
					}
				}
			}
		}
	case scrollBarLineUp:
		result = sb.LocalInsetBounds()
		result.Width = sb.Theme.Size
		result.Height = sb.Theme.Size
	case scrollBarLineDown:
		result = sb.LocalInsetBounds()
		if sb.horizontal {
			result.X += result.Width - sb.Theme.Size
		} else {
			result.Y += result.Height - sb.Theme.Size
		}
		result.Width = sb.Theme.Size
		result.Height = sb.Theme.Size
	case scrollBarPageUp:
		result = sb.partRect(scrollBarLineUp)
		thumb := sb.partRect(scrollBarThumb)
		if sb.horizontal {
			result.X += result.Width
			result.Width = thumb.X - result.X
		} else {
			result.Y += result.Height
			result.Height = thumb.Y - result.Y
		}
	case scrollBarPageDown:
		result = sb.partRect(scrollBarLineDown)
		thumb := sb.partRect(scrollBarThumb)
		if sb.horizontal {
			x := thumb.X + thumb.Width
			result.Width = result.X - x
			result.X = x
		} else {
			y := thumb.Y + thumb.Height
			result.Height = result.Y - y
			result.Y = y
		}
	}
	return result
}

func (sb *ScrollBar) drawThumb(g *draw.Graphics) {
	bounds := sb.partRect(scrollBarThumb)
	if !bounds.IsEmpty() {
		bgColor := sb.baseBackground(scrollBarThumb)
		g.Rect(bounds)
		var paint draw.Paint
		if sb.horizontal {
			paint = draw.NewLinearGradientPaint(sb.Theme.Gradient(bgColor), bounds.X, bounds.Y, bounds.X, bounds.Y+bounds.Height)
		} else {
			paint = draw.NewLinearGradientPaint(sb.Theme.Gradient(bgColor), bounds.X, bounds.Y, bounds.X+bounds.Width, bounds.Y)
		}
		g.SetPaint(paint)
		g.FillPath()
		paint.Dispose()
		g.SetColor(bgColor.AdjustBrightness(sb.Theme.OutlineAdjustment))
		g.StrokeRect(bounds)
		g.SetColor(sb.markColor(scrollBarThumb))
		var v0, v1, v2 float64
		if sb.horizontal {
			v0 = math.Floor(bounds.X + bounds.Width/2)
			d := math.Ceil(bounds.Height * 0.2)
			v1 = bounds.Y + d
			v2 = bounds.Y + bounds.Height - (d + 1)
		} else {
			v0 = math.Floor(bounds.Y + bounds.Height/2)
			d := math.Ceil(bounds.Width * 0.2)
			v1 = bounds.X + d
			v2 = bounds.X + bounds.Width - (d + 1)
		}
		for i := -1; i < 2; i++ {
			if sb.horizontal {
				x := v0 + float64(i*2)
				g.StrokeLine(x, v1, x, v2)
			} else {
				y := v0 + float64(i*2)
				g.StrokeLine(v1, y, v2, y)
			}
		}
	}
}

func (sb *ScrollBar) drawLineButton(g *draw.Graphics, linePart scrollBarPart) {
	bounds := sb.partRect(linePart)
	g.Save()
	g.Rect(bounds)
	bgColor := sb.baseBackground(linePart)
	paint := draw.NewLinearGradientPaint(sb.Theme.Gradient(bgColor), bounds.X, bounds.Y, bounds.X, bounds.Y+bounds.Height)
	g.SetPaint(paint)
	g.FillPath()
	paint.Dispose()
	g.SetColor(bgColor.AdjustBrightness(sb.Theme.OutlineAdjustment))
	g.StrokeRect(bounds)
	bounds.InsetUniform(1)
	if sb.horizontal {
		triHeight := (bounds.Width * 0.75)
		triWidth := triHeight / 2
		g.BeginPath()
		left := bounds.X + (bounds.Width-triWidth)/2
		right := left + triWidth
		top := bounds.Y + (bounds.Height-triHeight)/2
		bottom := top + triHeight
		if linePart == scrollBarLineUp {
			left, right = right, left
		}
		g.MoveTo(left, top)
		g.LineTo(left, bottom)
		g.LineTo(right, top+(bottom-top)/2)
	} else {
		triWidth := (bounds.Height * 0.75)
		triHeight := triWidth / 2
		g.BeginPath()
		left := bounds.X + (bounds.Width-triWidth)/2
		right := left + triWidth
		top := bounds.Y + (bounds.Height-triHeight)/2
		bottom := top + triHeight
		if linePart == scrollBarLineUp {
			top, bottom = bottom, top
		}
		g.MoveTo(left, top)
		g.LineTo(right, top)
		g.LineTo(left+(right-left)/2, bottom)
	}
	g.ClosePath()
	g.SetColor(sb.markColor(linePart))
	g.FillPath()
	g.Restore()
}

func (sb *ScrollBar) baseBackground(part scrollBarPart) color.Color {
	switch {
	case !sb.Enabled():
		return sb.Theme.Background.AdjustBrightness(sb.Theme.DisabledAdjustment)
	case part != scrollBarNone && sb.pressed == part:
		return sb.Theme.BackgroundWhenPressed
	case sb.Focused():
		return sb.Theme.Background.Blend(color.KeyboardFocus, 0.5)
	default:
		return sb.Theme.Background
	}
}

func (sb *ScrollBar) markColor(part scrollBarPart) color.Color {
	if sb.partEnabled(part) {
		if sb.baseBackground(part).Luminance() > 0.65 {
			return sb.Theme.MarkWhenLight
		}
		return sb.Theme.MarkWhenDark
	}
	return sb.Theme.MarkWhenDisabled
}

func (sb *ScrollBar) partEnabled(part scrollBarPart) bool {
	if sb.Enabled() && sb.Target != nil {
		switch part {
		case scrollBarLineUp, scrollBarPageUp:
			return sb.Target.ScrolledPosition(sb.horizontal) > 0
		case scrollBarLineDown, scrollBarPageDown:
			return sb.Target.ScrolledPosition(sb.horizontal) < sb.Target.ContentSize(sb.horizontal)-sb.Target.VisibleSize(sb.horizontal)
		case scrollBarThumb:
			pos := sb.Target.ScrolledPosition(sb.horizontal)
			return pos > 0 || pos < sb.Target.ContentSize(sb.horizontal)-sb.Target.VisibleSize(sb.horizontal)
		default:
		}
	}
	return false
}

// SetScrolledPosition attempts to set the current scrolled position of this ScrollBar to the
// specified value. The value will be clipped to the available range. If no target has been set,
// then nothing will happen.
func (sb *ScrollBar) SetScrolledPosition(position float64) {
	if sb.Target != nil {
		position = math.Max(math.Min(position, sb.Target.ContentSize(sb.horizontal)-sb.Target.VisibleSize(sb.horizontal)), 0)
		if sb.Target.ScrolledPosition(sb.horizontal) != position {
			sb.Target.SetScrolledPosition(sb.horizontal, position)
			sb.Repaint()
		}
	}
}
