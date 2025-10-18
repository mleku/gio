// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"

	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
)

// ViewportStyle defines the appearance and behavior of a scrollable viewport.
type ViewportStyle struct {
	// Content dimensions
	ContentWidth  int
	ContentHeight int

	// Scrollbar configuration
	ScrollbarWidth unit.Dp
	ScrollbarColor color.NRGBA
	TrackColor     color.NRGBA

	// Border configuration
	BorderColor color.NRGBA
	BorderWidth unit.Dp

	// Scroll state
	HorizontalPos        *Float
	VerticalPos          *Float
	HorizontalProportion *Float
	VerticalProportion   *Float

	// Scrollbar instances
	horizontalScrollbar ScrollbarStyle
	verticalScrollbar   ScrollbarStyle
}

// Viewport creates a new viewport widget with scrollbars.
func Viewport(contentWidth, contentHeight int) *ViewportStyle {
	horizontalPos := &Float{Value: 0.0}
	verticalPos := &Float{Value: 0.0}
	horizontalProportion := &Float{Value: 0.3}
	verticalProportion := &Float{Value: 0.4}

	return &ViewportStyle{
		ContentWidth:         contentWidth,
		ContentHeight:        contentHeight,
		ScrollbarWidth:       unit.Dp(4),
		ScrollbarColor:       color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80},
		TrackColor:           color.NRGBA{A: 0xFF, R: 0xC0, G: 0xC0, B: 0xC0},
		BorderColor:          color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80},
		BorderWidth:          unit.Dp(1),
		HorizontalPos:        horizontalPos,
		VerticalPos:          verticalPos,
		HorizontalProportion: horizontalProportion,
		VerticalProportion:   verticalProportion,
	}
}

// Layout renders the viewport widget.
func (v *ViewportStyle) Layout(gtx layout.Context, content layout.Widget) layout.Dimensions {
	// Initialize scrollbars if not already done
	if v.horizontalScrollbar.StartWidget == nil {
		v.horizontalScrollbar = HorizontalScrollbar(v.ScrollbarWidth, v.HorizontalProportion.Value, v.HorizontalPos.Value, v.HorizontalPos, v.HorizontalProportion)
	}
	if v.verticalScrollbar.StartWidget == nil {
		v.verticalScrollbar = VerticalScrollbar(v.ScrollbarWidth, v.VerticalProportion.Value, v.VerticalPos.Value, v.VerticalPos, v.VerticalProportion)
	}

	// Update scrollbar proportions based on viewport size
	viewportWidth := gtx.Constraints.Max.X - 20  // Account for vertical scrollbar
	viewportHeight := gtx.Constraints.Max.Y - 20 // Account for horizontal scrollbar

	if viewportWidth > 0 && v.ContentWidth > viewportWidth {
		v.HorizontalProportion.Value = float32(viewportWidth) / float32(v.ContentWidth)
	} else {
		v.HorizontalProportion.Value = 1.0
	}

	if viewportHeight > 0 && v.ContentHeight > viewportHeight {
		v.VerticalProportion.Value = float32(viewportHeight) / float32(v.ContentHeight)
	} else {
		v.VerticalProportion.Value = 1.0
	}

	// Update scrollbar proportions
	v.horizontalScrollbar.Proportion = v.HorizontalProportion.Value
	v.verticalScrollbar.Proportion = v.VerticalProportion.Value

	// Handle long press timers for both scrollbars
	v.horizontalScrollbar.HandleLongPress(gtx)
	v.verticalScrollbar.HandleLongPress(gtx)

	// Invalidate if any long press is active to keep timer running
	if v.horizontalScrollbar.IsLongPressing || v.verticalScrollbar.IsLongPressing {
		gtx.Execute(op.InvalidateCmd{})
	}

	// Calculate content offset based on scroll positions
	maxHScroll := max(0, v.ContentWidth-viewportWidth)
	maxVScroll := max(0, v.ContentHeight-viewportHeight)
	contentOffsetX := int(float32(maxHScroll) * v.HorizontalPos.Value)
	contentOffsetY := int(float32(maxVScroll) * v.VerticalPos.Value)

	// Layout with scrollbars positioned at window edges using VFlex/HFlex structure
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Top row: Main content and vertical scrollbar
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Main content area
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return v.layoutBorder(gtx, func(gtx layout.Context) layout.Dimensions {
						// Clip to viewport size
						defer clip.Rect(image.Rectangle{
							Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
						}).Push(gtx.Ops).Pop()

						// Offset content based on scroll position
						trans := op.Offset(image.Point{X: -contentOffsetX, Y: -contentOffsetY}).Push(gtx.Ops)

						// Draw scrollable content
						content(gtx)

						trans.Pop()

						return layout.Dimensions{Size: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
					})
				}),
				// Vertical scrollbar box
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if v.VerticalProportion.Value >= 1.0 {
						return layout.Dimensions{}
					}
					gtx.Constraints = layout.Exact(image.Point{X: 20, Y: gtx.Constraints.Max.Y})
					return v.verticalScrollbar.Layout(gtx)
				}),
			)
		}),
		// Bottom row: Horizontal scrollbar and corner
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Horizontal scrollbar box
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					if v.HorizontalProportion.Value >= 1.0 {
						return layout.Dimensions{}
					}
					gtx.Constraints = layout.Exact(image.Point{X: gtx.Constraints.Max.X, Y: 20})
					return v.horizontalScrollbar.Layout(gtx)
				}),
				// Corner box
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints = layout.Exact(image.Point{X: 20, Y: 20})
					// Empty corner - just shows the space
					return layout.Dimensions{Size: image.Point{X: 20, Y: 20}}
				}),
			)
		}),
	)
}

// layoutBorder adds a border outline around a widget
func (v *ViewportStyle) layoutBorder(gtx layout.Context, w layout.Widget) layout.Dimensions {
	// Draw border outline
	borderRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	paint.FillShape(gtx.Ops, v.BorderColor, clip.Stroke{
		Path:  clip.Rect(borderRect).Path(),
		Width: float32(v.BorderWidth),
	}.Op())

	// Layout the widget inside with padding
	gtx.Constraints.Max.X -= 2
	gtx.Constraints.Max.Y -= 2
	gtx.Constraints.Min.X = max(0, gtx.Constraints.Min.X-2)
	gtx.Constraints.Min.Y = max(0, gtx.Constraints.Min.Y-2)

	trans := op.Offset(image.Point{X: 1, Y: 1}).Push(gtx.Ops)
	dims := w(gtx)
	trans.Pop()

	return layout.Dimensions{
		Size: dims.Size.Add(image.Point{X: 2, Y: 2}),
	}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
