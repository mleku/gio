// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
)

// OutlineWidget represents a widget that draws an outline around its bounds
type OutlineWidget struct {
	*Widget
	OutlineColor color.NRGBA
	Thickness    int
	CornerRadius int
}

// NewOutlineWidget creates a new outline widget
func NewOutlineWidget() *OutlineWidget {
	return &OutlineWidget{
		Widget:       NewWidget(),
		OutlineColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255}, // Black outline
		Thickness:    1,
		CornerRadius: 0,
	}
}

// SetOutlineColor sets the outline color
func (o *OutlineWidget) SetOutlineColor(color color.NRGBA) *OutlineWidget {
	o.OutlineColor = color
	return o
}

// SetThickness sets the outline thickness
func (o *OutlineWidget) SetThickness(thickness int) *OutlineWidget {
	o.Thickness = thickness
	return o
}

// SetCornerRadius sets the corner radius for rounded corners
func (o *OutlineWidget) SetCornerRadius(radius int) *OutlineWidget {
	o.CornerRadius = radius
	return o
}

// GetWidget returns the underlying widget (for Renderer interface)
func (o *OutlineWidget) GetWidget() *Widget {
	return o.Widget
}

// RenderWidget renders the outline widget (overrides Widget.RenderWidget)
func (o *OutlineWidget) RenderWidget(gtx app.Context) {
	if !o.Visible {
		return
	}

	// Set up clipping for this widget
	defer clip.Rect{Min: image.Point{X: o.X, Y: o.Y}, Max: image.Point{X: o.X + o.Width, Y: o.Y + o.Height}}.Push(gtx.Ops).Pop()

	// Call custom render function if provided (for background)
	if o.Render != nil {
		// Draw background with inset to leave room for outline
		inset := o.Thickness
		if inset > 0 {
			defer clip.Rect{
				Min: image.Point{X: o.X + inset, Y: o.Y + inset},
				Max: image.Point{X: o.X + o.Width - inset, Y: o.Y + o.Height - inset},
			}.Push(gtx.Ops).Pop()
		}
		o.Render(gtx, o.Widget)
	}

	// Draw the outline
	o.drawOutline(gtx)

	// Render children
	for _, child := range o.Children {
		child.RenderWidget(gtx)
	}
}

// drawOutline draws the outline around the widget bounds
func (o *OutlineWidget) drawOutline(gtx app.Context) {
	if o.Thickness <= 0 {
		return
	}

	// Create clipping rectangle for the outline
	bounds := o.Bounds()

	if o.CornerRadius > 0 {
		// Draw rounded rectangle outline
		o.drawRoundedOutline(gtx, bounds)
	} else {
		// Draw rectangular outline
		o.drawRectangularOutline(gtx, bounds)
	}
}

// drawRectangularOutline draws a rectangular outline using Gio's stroke operations
func (o *OutlineWidget) drawRectangularOutline(gtx app.Context, bounds image.Rectangle) {
	thickness := float32(o.Thickness)

	// Create a rectangle for the outline
	r := bounds

	// Draw the outline using Gio's stroke operation
	paint.FillShape(gtx.Ops,
		o.OutlineColor,
		clip.Stroke{
			Path:  clip.RRect{Rect: r, NW: 0, NE: 0, SW: 0, SE: 0}.Path(gtx.Ops),
			Width: thickness,
		}.Op(),
	)
}

// drawRoundedOutline draws a rounded rectangle outline using Gio's stroke operations
func (o *OutlineWidget) drawRoundedOutline(gtx app.Context, bounds image.Rectangle) {
	thickness := float32(o.Thickness)
	radius := float32(o.CornerRadius)

	// Create a rounded rectangle for the outline
	r := bounds

	// Draw the rounded outline using Gio's stroke operation
	paint.FillShape(gtx.Ops,
		o.OutlineColor,
		clip.Stroke{
			Path:  clip.RRect{Rect: r, NW: int(radius), NE: int(radius), SW: int(radius), SE: int(radius)}.Path(gtx.Ops),
			Width: thickness,
		}.Op(),
	)
}

// Fluent methods for OutlineWidget that delegate to the embedded Widget

// SetPosition sets the outline widget's position
func (o *OutlineWidget) SetPosition(x, y int) *OutlineWidget {
	o.Widget.SetPosition(x, y)
	return o
}

// SetSize sets the outline widget's size
func (o *OutlineWidget) SetSize(width, height int) *OutlineWidget {
	o.Widget.SetSize(width, height)
	return o
}

// SetVisible sets the outline widget's visibility
func (o *OutlineWidget) SetVisible(visible bool) *OutlineWidget {
	o.Widget.SetVisible(visible)
	return o
}

// AddChild adds a child widget
func (o *OutlineWidget) AddChild(child *Widget) *OutlineWidget {
	o.Widget.AddChild(child)
	return o
}
