package widget

import (
	"image"
	"image/color"

	"lol.mleku.dev/log"

	"gio.mleku.dev/layout"
	"gio.mleku.dev/op/paint"
)

// ContextWrapper wraps any widget and provides context menu functionality.
// It implements the ContextWidget interface and can contain any widget.
type ContextWrapper struct {
	// The widget to wrap
	Widget layout.Widget

	// Context menu function - called when right-clicked
	ContextMenuFunc func(gtx layout.Context, pos image.Point) layout.Widget

	// Priority for context menu handling (higher = more priority)
	Priority int

	// Internal state
	clickable Clickable
	bounds    image.Rectangle
}

// NewContextWrapper creates a new context wrapper that wraps the given widget.
func NewContextWrapper(widget layout.Widget, contextMenuFunc func(gtx layout.Context, pos image.Point) layout.Widget, priority int) *ContextWrapper {
	return &ContextWrapper{
		Widget:          widget,
		ContextMenuFunc: contextMenuFunc,
		Priority:        priority,
	}
}

// Layout lays out the wrapped widget and handles context menu events.
func (cw *ContextWrapper) Layout(gtx layout.Context) layout.Dimensions {
	// Layout the wrapped widget
	dims := cw.Widget(gtx)

	// Register for pointer events
	cw.clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return dims
	})

	return dims
}

// ContextMenu implements the ContextWidget interface.
func (cw *ContextWrapper) ContextMenu(gtx layout.Context, pos image.Point) layout.Widget {
	log.I.F("ContextWrapper.ContextMenu called at position %v", pos)

	// Check if the click is within our bounds
	if !pos.In(cw.bounds) {
		log.I.F("Click outside context wrapper bounds %v, ignoring", cw.bounds)
		return nil
	}

	log.I.F("Click within context wrapper bounds, showing context menu")

	// Call the context menu function
	if cw.ContextMenuFunc != nil {
		return cw.ContextMenuFunc(gtx, pos)
	}

	return nil
}

// GetPriority returns the priority of this context wrapper.
func (cw *ContextWrapper) GetPriority() int {
	return cw.Priority
}

// GetBounds returns the current bounds of the widget.
func (cw *ContextWrapper) GetBounds() image.Rectangle {
	return cw.bounds
}

// UpdateBounds updates the bounds of the widget (called by ContextManager).
func (cw *ContextWrapper) UpdateBounds(bounds image.Rectangle) {
	cw.bounds = bounds
	log.I.F("ContextWidget bounds updated to %v", bounds)
}

// WithContextMenu is a helper function to create a context wrapper with a simple context menu.
func WithContextMenu(widget layout.Widget, menuItems []string, onItemClick func(item string)) *ContextWrapper {
	return NewContextWrapper(widget, func(gtx layout.Context, pos image.Point) layout.Widget {
		return func(gtx layout.Context) layout.Dimensions {
			// Create a simple menu widget
			return layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Menu background
					paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Menu items - simple text for now
					var children []layout.FlexChild
					for range menuItems {
						children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// Simple text button
							btn := Clickable{}
							return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								// Draw text
								paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
								return layout.Dimensions{Size: image.Point{X: 100, Y: 30}}
							})
						}))
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
				}),
			)
		}
	}, 10) // Default priority
}
