// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op/clip"
)

// Renderer interface for widgets that can render themselves
type Renderer interface {
	RenderWidget(gtx app.Context)
	GetWidget() *Widget // Method to get the underlying widget for layout
}

// Widget represents a basic widget with position, size, and children
type Widget struct {
	// Position and size
	X, Y, Width, Height int

	// Children widgets
	Children []*Widget

	// Render function for this widget
	Render func(gtx app.Context, w *Widget)

	// Whether this widget is visible
	Visible bool
}

// NewWidget creates a new widget with default values
func NewWidget() *Widget {
	return &Widget{
		Width:   100,
		Height:  100,
		Visible: true,
	}
}

// SetPosition sets the widget's position
func (w *Widget) SetPosition(x, y int) *Widget {
	w.X = x
	w.Y = y
	return w
}

// SetSize sets the widget's size
func (w *Widget) SetSize(width, height int) *Widget {
	w.Width = width
	w.Height = height
	return w
}

// SetVisible sets the widget's visibility
func (w *Widget) SetVisible(visible bool) *Widget {
	w.Visible = visible
	return w
}

// AddChild adds a child widget
func (w *Widget) AddChild(child *Widget) *Widget {
	w.Children = append(w.Children, child)
	return w
}

// Bounds returns the widget's bounds as an image.Rectangle
func (w *Widget) Bounds() image.Rectangle {
	return image.Rect(w.X, w.Y, w.X+w.Width, w.Y+w.Height)
}

// GetWidget returns the widget itself (for Renderer interface)
func (w *Widget) GetWidget() *Widget {
	return w
}

// RenderWidget renders the widget and its children
func (w *Widget) RenderWidget(gtx app.Context) {
	if !w.Visible {
		return
	}

	// Set up clipping for this widget
	defer clip.Rect{Min: image.Point{X: w.X, Y: w.Y}, Max: image.Point{X: w.X + w.Width, Y: w.Y + w.Height}}.Push(gtx.Ops).Pop()

	// Call custom render function if provided
	if w.Render != nil {
		w.Render(gtx, w)
	}

	// Render children
	for _, child := range w.Children {
		child.RenderWidget(gtx)
	}
}
