// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"
	"time"

	"gio.mleku.dev/internal/f32color"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/key"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
)

// ScrollbarStyle defines the appearance of a scrollbar using the same approach as ProgressBar.
type ScrollbarStyle struct {
	Color      color.NRGBA
	Height     unit.Dp
	Radius     unit.Dp
	TrackColor color.NRGBA
	Proportion float32 // Value between 0 and 1 representing how much is visible
	Start      float32 // Value between 0 and 1 representing start position
	Vertical   bool    // Whether this is a vertical scrollbar

	// Interaction state
	StartWidget      *Float // Widget to update when start position changes
	ProportionWidget *Float // Widget to update when proportion changes

	// Drag state
	IsDragging bool    // Whether currently dragging
	DragOffset float32 // Offset from click position to thumb start

	// Animation state
	IsAnimating        bool      // Whether currently animating
	AnimationStart     float32   // Starting position for animation
	AnimationEnd       float32   // Target position for animation
	AnimationStartTime time.Time // When animation started

	// Long press state
	IsLongPressing     bool      // Whether currently in long press mode
	LongPressStartTime time.Time // When long press started
	LongPressTarget    float32   // Target position for long press scroll
}

// HorizontalScrollbar creates a horizontal scrollbar style with default values.
func HorizontalScrollbar(thickness unit.Dp, proportion, start float32, startWidget, proportionWidget *Float) ScrollbarStyle {
	return ScrollbarStyle{
		Proportion:       clamp1(proportion),
		Start:            clamp1(start),
		Height:           unit.Dp(4),                                      // Same as progressbar
		Radius:           unit.Dp(2),                                      // Same as progressbar
		Color:            color.NRGBA{A: 0xFF, R: 0x00, G: 0x00, B: 0x00}, // Black thumb
		TrackColor:       color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80}, // Gray track
		Vertical:         false,
		StartWidget:      startWidget,
		ProportionWidget: proportionWidget,
	}
}

// VerticalScrollbar creates a vertical scrollbar style with default values.
func VerticalScrollbar(thickness unit.Dp, proportion, start float32, startWidget, proportionWidget *Float) ScrollbarStyle {
	return ScrollbarStyle{
		Proportion:       clamp1(proportion),
		Start:            clamp1(start),
		Height:           unit.Dp(4),                                      // Same as progressbar
		Radius:           unit.Dp(2),                                      // Same as progressbar
		Color:            color.NRGBA{A: 0xFF, R: 0x00, G: 0x00, B: 0x00}, // Black thumb
		TrackColor:       color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80}, // Gray track
		Vertical:         true,
		StartWidget:      startWidget,
		ProportionWidget: proportionWidget,
	}
}

// Layout renders the scrollbar using the exact same approach as ProgressBar.
func (s *ScrollbarStyle) Layout(gtx layout.Context) layout.Dimensions {
	// Don't render if proportion is 1 (all visible)
	if s.Proportion >= 1.0 {
		return layout.Dimensions{}
	}

	// Handle animation
	if s.IsAnimating {
		// Calculate animation progress based on time
		animationDuration := 200 * time.Millisecond
		elapsed := gtx.Now.Sub(s.AnimationStartTime)
		progress := float32(elapsed) / float32(animationDuration)

		if progress >= 1.0 {
			// Animation complete
			s.IsAnimating = false
			s.StartWidget.Value = s.AnimationEnd
			s.Start = s.AnimationEnd // Update the scrollbar's Start field too
		} else {
			// Interpolate between start and end positions
			currentPos := s.AnimationStart + (s.AnimationEnd-s.AnimationStart)*easeOutCubic(progress)
			s.StartWidget.Value = currentPos
			s.Start = currentPos // Update the scrollbar's Start field too

			// Invalidate to trigger repaint for next frame
			gtx.Execute(op.InvalidateCmd{})
		}
	} else {
		// Not animating or long pressing - update Start from slider value
		s.Start = s.StartWidget.Value
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

	var trackOffset image.Point
	if s.Vertical {
		// Track is 4dp wide, centered in 20dp area
		trackWidth := gtx.Dp(s.Height)        // 4dp
		trackAreaWidth := gtx.Dp(unit.Dp(20)) // 20dp
		trackOffset = image.Point{X: (trackAreaWidth - trackWidth) / 2, Y: 0}

		trackSize = image.Point{X: trackWidth, Y: gtx.Constraints.Max.Y}
		thumbHeight := int(float32(trackSize.Y) * s.Proportion)
		thumbSize = image.Point{X: trackSize.X, Y: thumbHeight}
		thumbPos = image.Point{X: trackOffset.X, Y: int(float32(trackSize.Y-thumbHeight) * s.Start)}
	} else {
		// Track is 4dp tall, centered in 20dp area
		trackHeight := gtx.Dp(s.Height)        // 4dp
		trackAreaHeight := gtx.Dp(unit.Dp(20)) // 20dp
		trackOffset = image.Point{X: 0, Y: (trackAreaHeight - trackHeight) / 2}

		trackSize = image.Point{X: gtx.Constraints.Max.X, Y: trackHeight}
		thumbWidth := int(float32(trackSize.X) * s.Proportion)
		thumbSize = image.Point{X: thumbWidth, Y: trackSize.Y}
		thumbPos = image.Point{X: int(float32(trackSize.X-thumbWidth) * s.Start), Y: trackOffset.Y}
	}

	// Use the same alignment as ProgressBar (West for horizontal, North for vertical)
	alignment := layout.W
	if s.Vertical {
		alignment = layout.N
	}

	return layout.Stack{Alignment: alignment}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			trans := op.Offset(trackOffset).Push(gtx.Ops)
			dims := shader(trackSize.X, trackSize.Y, s.TrackColor)
			trans.Pop()
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			fillColor := s.Color
			if !gtx.Enabled() {
				fillColor = f32color.Disabled(fillColor)
			}

			// Handle mouse events with sophisticated interaction
			if s.StartWidget != nil {
				// Use the full 20dp track area as hit area
				var widgetRect image.Rectangle
				if s.Vertical {
					trackAreaWidth := gtx.Dp(unit.Dp(20))
					widgetRect = image.Rectangle{
						Min: image.Point{X: 0, Y: 0},
						Max: image.Point{X: trackAreaWidth, Y: trackSize.Y},
					}
				} else {
					trackAreaHeight := gtx.Dp(unit.Dp(20))
					widgetRect = image.Rectangle{
						Min: image.Point{X: 0, Y: 0},
						Max: image.Point{X: trackSize.X, Y: trackAreaHeight},
					}
				}
				defer clip.Rect(widgetRect).Push(gtx.Ops).Pop()
				event.Op(gtx.Ops, s)

				for {
					evt, ok := gtx.Source.Event(pointer.Filter{
						Target: s,
						Kinds:  pointer.Press | pointer.Drag | pointer.Release,
					})
					if !ok {
						break
					}
					if ev, ok := evt.(pointer.Event); ok {
						switch ev.Kind {
						case pointer.Press:
							// Determine if click is on thumb or track
							clickPos := ev.Position
							var clickCoord float32
							var thumbStart, thumbEnd float32

							if s.Vertical {
								clickCoord = clickPos.Y
								thumbStart = float32(thumbPos.Y)
								thumbEnd = float32(thumbPos.Y + thumbSize.Y)
							} else {
								clickCoord = clickPos.X
								thumbStart = float32(thumbPos.X)
								thumbEnd = float32(thumbPos.X + thumbSize.X)
							}

							if clickCoord >= thumbStart && clickCoord <= thumbEnd {
								// Clicked on thumb - start dragging with offset
								s.IsDragging = true
								s.DragOffset = clickCoord - thumbStart
								gtx.Execute(key.FocusCmd{Tag: s})
							} else {
								// Clicked on track - immediate one-thumb scroll + long press timer
								var immediateTarget float32
								var longPressTarget float32
								var thumbWidth float32
								var maxPos float32
								var currentPos float32

								if s.Vertical {
									currentPos = float32(thumbPos.Y)
									thumbWidth = float32(thumbSize.Y)
									maxPos = float32(trackSize.Y - thumbSize.Y)
								} else {
									currentPos = float32(thumbPos.X)
									thumbWidth = float32(thumbSize.X)
									maxPos = float32(trackSize.X - thumbSize.X)
								}

								if maxPos <= 0 {
									// No space to move
									continue
								}

								// Calculate immediate movement (one thumb width)
								if clickCoord < currentPos {
									// Clicked before thumb - move left/up by thumb width
									newPos := currentPos - thumbWidth
									if newPos < 0 {
										newPos = 0 // Move to start
									}
									immediateTarget = clamp1(newPos / maxPos)
									longPressTarget = 0.0 // Long press goes to start
								} else {
									// Clicked after thumb - move right/down by thumb width
									newPos := currentPos + thumbWidth
									if newPos > maxPos {
										newPos = maxPos // Move to end
									}
									immediateTarget = clamp1(newPos / maxPos)
									longPressTarget = 1.0 // Long press goes to end
								}

								// Start immediate animation
								s.IsAnimating = true
								s.AnimationStart = s.StartWidget.Value
								s.AnimationEnd = immediateTarget
								s.AnimationStartTime = gtx.Now

								// Start long press timer for end position
								s.IsLongPressing = true
								s.LongPressStartTime = gtx.Now
								s.LongPressTarget = longPressTarget

								// Invalidate to start timer
								gtx.Execute(op.InvalidateCmd{})
							}

						case pointer.Drag:
							if s.IsDragging {
								// Update position based on drag with offset
								var newStart float32
								if s.Vertical {
									// Calculate new position maintaining drag offset
									dragY := ev.Position.Y - s.DragOffset
									maxPos := float32(trackSize.Y - thumbSize.Y)
									if maxPos > 0 {
										newStart = clamp1(dragY / maxPos)
									}
								} else {
									// Calculate new position maintaining drag offset
									dragX := ev.Position.X - s.DragOffset
									maxPos := float32(trackSize.X - thumbSize.X)
									if maxPos > 0 {
										newStart = clamp1(dragX / maxPos)
									}
								}
								s.StartWidget.Value = newStart
							}

						case pointer.Release:
							// Stop dragging
							s.IsDragging = false
							s.DragOffset = 0
							// Cancel long press if it was active
							s.IsLongPressing = false
						}
					}
				}
			}

			trans := op.Offset(thumbPos).Push(gtx.Ops)
			dims := shader(thumbSize.X, thumbSize.Y, fillColor)
			trans.Pop()
			return dims
		}),
	)
}

// HandleLongPress handles long press timers for scrollbars
func (s *ScrollbarStyle) HandleLongPress(gtx layout.Context) {
	if s.IsLongPressing && !s.IsAnimating {
		longPressDuration := 750 * time.Millisecond
		elapsed := gtx.Now.Sub(s.LongPressStartTime)
		if elapsed >= longPressDuration {
			// Long press threshold reached - start animation to target
			s.IsLongPressing = false
			s.IsAnimating = true
			s.AnimationStart = s.StartWidget.Value
			s.AnimationEnd = s.LongPressTarget
			s.AnimationStartTime = gtx.Now
		}
	}
}

// easeOutCubic provides smooth easing for animations
func easeOutCubic(t float32) float32 {
	if t >= 1.0 {
		return 1.0
	}
	return 1.0 - pow(1.0-t, 3)
}

// pow calculates x^y for float32
func pow(x, y float32) float32 {
	// Simple implementation for small integer powers
	if y == 3 {
		return x * x * x
	}
	return x // fallback
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
