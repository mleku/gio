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
	"gio.mleku.dev/internal/f32color"
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

// ScrollbarStyle defines the appearance of a scrollbar using the same approach as ProgressBar.
type ScrollbarStyle struct {
	Color      color.NRGBA
	Height     unit.Dp
	Radius     unit.Dp
	TrackColor color.NRGBA
	Proportion float32 // Value between 0 and 1 representing how much is visible
	Start      float32 // Value between 0 and 1 representing start position
	Vertical   bool    // Whether this is a vertical scrollbar
}

// HorizontalScrollbar creates a horizontal scrollbar style with default values.
func HorizontalScrollbar(thickness unit.Dp, proportion, start float32) ScrollbarStyle {
	return ScrollbarStyle{
		Proportion: clamp1(proportion),
		Start:      clamp1(start),
		Height:     unit.Dp(4), // Same as progressbar
		Radius:     unit.Dp(2), // Same as progressbar
		Color:      color.NRGBA{A: 0x80, R: 0x80, G: 0x80, B: 0x80},
		TrackColor: color.NRGBA{A: 0x40, R: 0x80, G: 0x80, B: 0x80},
		Vertical:   false,
	}
}

// VerticalScrollbar creates a vertical scrollbar style with default values.
func VerticalScrollbar(thickness unit.Dp, proportion, start float32) ScrollbarStyle {
	return ScrollbarStyle{
		Proportion: clamp1(proportion),
		Start:      clamp1(start),
		Height:     unit.Dp(4), // Same as progressbar
		Radius:     unit.Dp(2), // Same as progressbar
		Color:      color.NRGBA{A: 0x80, R: 0x80, G: 0x80, B: 0x80},
		TrackColor: color.NRGBA{A: 0x40, R: 0x80, G: 0x80, B: 0x80},
		Vertical:   true,
	}
}

// Layout renders the scrollbar using the exact same approach as ProgressBar.
func (s ScrollbarStyle) Layout(gtx layout.Context) layout.Dimensions {
	// Don't render if proportion is 1 (all visible)
	if s.Proportion >= 1.0 {
		return layout.Dimensions{}
	}

	shader := func(width, height int, color color.NRGBA) layout.Dimensions {
		d := image.Point{X: width, Y: height}
		rr := gtx.Dp(s.Radius)

		defer clip.UniformRRect(image.Rectangle{Max: d}, rr).Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		return layout.Dimensions{Size: d}
	}

	var trackSize image.Point
	var thumbSize image.Point
	var thumbPos image.Point

	if s.Vertical {
		trackSize = image.Point{X: gtx.Dp(s.Height), Y: gtx.Constraints.Max.Y}
		thumbHeight := int(float32(trackSize.Y) * s.Proportion)
		thumbSize = image.Point{X: trackSize.X, Y: thumbHeight}
		thumbPos = image.Point{X: 0, Y: int(float32(trackSize.Y-thumbHeight) * s.Start)}
	} else {
		trackSize = image.Point{X: gtx.Constraints.Max.X, Y: gtx.Dp(s.Height)}
		thumbWidth := int(float32(trackSize.X) * s.Proportion)
		thumbSize = image.Point{X: thumbWidth, Y: trackSize.Y}
		thumbPos = image.Point{X: int(float32(trackSize.X-thumbWidth) * s.Start), Y: 0}
	}

	// Use the same alignment as ProgressBar (West for horizontal, North for vertical)
	alignment := layout.W
	if s.Vertical {
		alignment = layout.N
	}

	return layout.Stack{Alignment: alignment}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return shader(trackSize.X, trackSize.Y, s.TrackColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			fillColor := s.Color
			if !gtx.Enabled() {
				fillColor = f32color.Disabled(fillColor)
			}
			trans := op.Offset(thumbPos).Push(gtx.Ops)
			dims := shader(thumbSize.X, thumbSize.Y, fillColor)
			trans.Pop()
			return dims
		}),
	)
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

// clamp1 limits v to range [0..1].
func clamp1(v float32) float32 {
	if v >= 1 {
		return 1
	} else if v <= 0 {
		return 0
	} else {
		return v
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

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Create scrollbars with values from sliders
			horizontalScrollbar := HorizontalScrollbar(unit.Dp(20), horizontalProportionSlider.Value, horizontalStartSlider.Value)
			verticalScrollbar := VerticalScrollbar(unit.Dp(20), verticalProportionSlider.Value, verticalStartSlider.Value)

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
							// Add outline around main content box
							return layoutBorder(gtx, func(gtx layout.Context) layout.Dimensions {
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
							})
						}),
						// Vertical scrollbar box
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if verticalProportionSlider.Value >= 1.0 {
								return layout.Dimensions{}
							}
							// Add outline around vertical scrollbar box
							return layoutBorder(gtx, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints = layout.Exact(image.Point{X: 20, Y: gtx.Constraints.Max.Y})
								return verticalScrollbar.Layout(gtx)
							})
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
							// Add outline around horizontal scrollbar box
							return layoutBorder(gtx, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints = layout.Exact(image.Point{X: gtx.Constraints.Max.X, Y: 20})
								return horizontalScrollbar.Layout(gtx)
							})
						}),
						// Corner box
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							// Add outline around corner box
							return layoutBorder(gtx, func(gtx layout.Context) layout.Dimensions {
								gtx.Constraints = layout.Exact(image.Point{X: 20, Y: 20})
								// Empty corner - just shows the outline
								return layout.Dimensions{Size: image.Point{X: 20, Y: 20}}
							})
						}),
					)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
