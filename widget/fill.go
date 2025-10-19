// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
)

// FillWidget represents a widget that fills its boundaries with a color
type FillWidget struct {
	*Widget
	FillColor color.NRGBA
}

// NewFillWidget creates a new fill widget
func NewFillWidget() *FillWidget {
	return &FillWidget{
		Widget:    NewWidget(),
		FillColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White fill
	}
}

// Color sets the fill color
func (f *FillWidget) Color(color color.NRGBA) *FillWidget {
	f.FillColor = color
	return f
}

// GetWidget returns the underlying widget (for Renderer interface)
func (f *FillWidget) GetWidget() *Widget {
	return f.Widget
}

// RenderWidget renders the fill widget
func (f *FillWidget) RenderWidget(gtx app.Context) {
	if !f.Visible {
		return
	}

	// Set up clipping for this widget
	defer clip.Rect{Min: image.Point{X: f.X, Y: f.Y}, Max: image.Point{X: f.X + f.Width, Y: f.Y + f.Height}}.Push(gtx.Ops).Pop()

	// Fill the widget with the specified color
	if f.FillColor.A > 0 {
		paint.Fill(gtx.Ops, f.FillColor)
	}

	// Call custom render function if provided
	if f.Render != nil {
		f.Render(gtx, f.Widget)
	}

	// Render children
	for _, child := range f.Children {
		child.RenderWidget(gtx)
	}
}

// Fluent methods for FillWidget that delegate to the embedded Widget

// SetPosition sets the fill widget's position
func (f *FillWidget) SetPosition(x, y int) *FillWidget {
	f.Widget.SetPosition(x, y)
	return f
}

// SetSize sets the fill widget's size
func (f *FillWidget) SetSize(width, height int) *FillWidget {
	f.Widget.SetSize(width, height)
	return f
}

// SetVisible sets the fill widget's visibility
func (f *FillWidget) SetVisible(visible bool) *FillWidget {
	f.Widget.SetVisible(visible)
	return f
}

// AddChild adds a child widget
func (f *FillWidget) AddChild(child *Widget) *FillWidget {
	f.Widget.AddChild(child)
	return f
}
