// SPDX-License-Identifier: Unlicense OR MIT

package layout

import (
	"image"
	"time"

	"gio.mleku.dev/io/input"
	"gio.mleku.dev/io/system"
	"gio.mleku.dev/op"
	"gio.mleku.dev/unit"
)

// Context carries the state needed by almost all layouts and widgets.
// A zero value Context never returns events, map units to pixels
// with a scale of 1.0, and returns the zero time from Now.
type Context struct {
	// Constraints track the constraints for the active widget or
	// layout.
	Constraints Constraints

	Metric unit.Metric
	// Now is the animation time.
	Now time.Time

	// Locale provides information on the system's language preferences.
	// BUG(whereswaldon): this field is not currently populated automatically.
	// Interested users must look up and populate these values manually.
	Locale system.Locale

	// Values is a map of program global data associated with the context.
	// It is not for use by widgets.
	Values map[string]any

	// WindowDimensions is the dimensions of the window in Dp.
	WindowDimensions image.Point

	// MousePosition is the current mouse cursor position relative to the window
	// dimensions, or outside if the cursor is outside the window.
	MousePosition image.Point

	// WidgetPosition is the top-left position of the current widget as a
	// coordinate of the window.
	WidgetPosition image.Point

	input.Source
	*op.Ops
}

// Dp converts v to pixels.
func (c Context) Dp(v unit.Dp) int {
	return c.Metric.Dp(v)
}

// Sp converts v to pixels.
func (c Context) Sp(v unit.Sp) int {
	return c.Metric.Sp(v)
}

// Disabled returns a copy of this context that don't deliver any events.
func (c Context) Disabled() Context {
	c.Source = c.Source.Disabled()
	return c
}

// WindowSize returns the window dimensions in Dp.
func (c Context) WindowSize() image.Point {
	return c.WindowDimensions
}

// CursorPosition returns the current mouse cursor position relative to the window
// dimensions, or outside if the cursor is outside the window.
func (c Context) CursorPosition() image.Point {
	return c.MousePosition
}

// WidgetOffset returns the top-left position of the current widget as a
// coordinate of the window.
func (c Context) WidgetOffset() image.Point {
	return c.WidgetPosition
}

// WithOffset returns a copy of this context with the widget position
// offset by the given amount.
func (c Context) WithOffset(offset image.Point) Context {
	c.WidgetPosition = c.WidgetPosition.Add(offset)
	return c
}
