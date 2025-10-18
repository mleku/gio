// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
	"gio.mleku.dev/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

// layoutBorder adds a border outline around a widget
func layoutBorder(gtx layout.Context, w layout.Widget) layout.Dimensions {
	// Draw border outline
	borderColor := color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80}
	borderRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	paint.FillShape(gtx.Ops, borderColor, clip.Stroke{
		Path:  clip.Rect(borderRect).Path(),
		Width: 1,
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

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	// Slider widgets for controlling horizontal scrollbar values
	var horizontalProportionSlider widget.Float
	var horizontalStartSlider widget.Float

	// Slider widgets for controlling vertical scrollbar values
	var verticalProportionSlider widget.Float
	var verticalStartSlider widget.Float

	// Set initial values
	horizontalProportionSlider.Value = 0.25
	horizontalStartSlider.Value = 0.25
	verticalProportionSlider.Value = 0.25
	verticalStartSlider.Value = 0.25

	// Create persistent scrollbar instances using the widget component
	horizontalScrollbar := widget.HorizontalScrollbar(unit.Dp(4), horizontalProportionSlider.Value, horizontalStartSlider.Value, &horizontalStartSlider, &horizontalProportionSlider)
	verticalScrollbar := widget.VerticalScrollbar(unit.Dp(4), verticalProportionSlider.Value, verticalStartSlider.Value, &verticalStartSlider, &verticalProportionSlider)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Update scrollbar proportion from sliders
			horizontalScrollbar.Proportion = horizontalProportionSlider.Value
			verticalScrollbar.Proportion = verticalProportionSlider.Value

			// Handle long press timers for both scrollbars
			horizontalScrollbar.HandleLongPress(gtx)
			verticalScrollbar.HandleLongPress(gtx)

			// Invalidate if any long press is active to keep timer running
			if horizontalScrollbar.IsLongPressing || verticalScrollbar.IsLongPressing {
				gtx.Execute(op.InvalidateCmd{})
			}

			// Layout with scrollbars positioned at window edges using VFlex/HFlex structure
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				// Top row: Main content and vertical scrollbar
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{
						Axis: layout.Horizontal,
					}.Layout(gtx,
						// Main content area
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{
									Axis: layout.Horizontal,
								}.Layout(gtx,
									// Left side - Horizontal scrollbar controls
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.H3(th, "Horizontal Scrollbar").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body1(th, fmt.Sprintf("Proportion: %.2f, Start: %.2f", horizontalProportionSlider.Value, horizontalStartSlider.Value)).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body2(th, "Proportion Visible (0-1):").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints = layout.Exact(image.Point{X: 300, Y: gtx.Constraints.Max.Y})
												return material.Slider(th, &horizontalProportionSlider).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body2(th, "Start Position (0-1):").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints = layout.Exact(image.Point{X: 300, Y: gtx.Constraints.Max.Y})
												return material.Slider(th, &horizontalStartSlider).Layout(gtx)
											}),
										)
									}),
									// Right side - Vertical scrollbar controls
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										return layout.Flex{
											Axis: layout.Vertical,
										}.Layout(gtx,
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.H3(th, "Vertical Scrollbar").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body1(th, fmt.Sprintf("Proportion: %.2f, Start: %.2f", verticalProportionSlider.Value, verticalStartSlider.Value)).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body2(th, "Proportion Visible (0-1):").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints = layout.Exact(image.Point{X: 300, Y: gtx.Constraints.Max.Y})
												return material.Slider(th, &verticalProportionSlider).Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												return material.Body2(th, "Start Position (0-1):").Layout(gtx)
											}),
											layout.Rigid(func(gtx layout.Context) layout.Dimensions {
												gtx.Constraints = layout.Exact(image.Point{X: 300, Y: gtx.Constraints.Max.Y})
												return material.Slider(th, &verticalStartSlider).Layout(gtx)
											}),
										)
									}),
								)
							})
						}),
						// Vertical scrollbar box
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if verticalProportionSlider.Value >= 1.0 {
								return layout.Dimensions{}
							}
							gtx.Constraints = layout.Exact(image.Point{X: 20, Y: gtx.Constraints.Max.Y})
							return verticalScrollbar.Layout(gtx)
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
							if horizontalProportionSlider.Value >= 1.0 {
								return layout.Dimensions{}
							}
							gtx.Constraints = layout.Exact(image.Point{X: gtx.Constraints.Max.X, Y: 20})
							return horizontalScrollbar.Layout(gtx)
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

			e.Frame(gtx.Ops)
		}
	}
}
