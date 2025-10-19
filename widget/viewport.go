// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"gio.mleku.dev/gesture"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/unit"
)

// ScrollEvent represents a queued scroll event
type ScrollEvent struct {
	ScrollY float32
	Time    time.Time
}

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

	// Scroll event handling
	scrollQueue             []ScrollEvent
	scrollAnimating         bool
	scrollAnimationStart    time.Time
	scrollAnimationDuration time.Duration
	scrollAnimationStartPos image.Point
	scrollAnimationEndPos   image.Point

	// Gesture scroll handler
	scroller gesture.Scroll
}

// Viewport creates a new viewport widget with scrollbars.
func Viewport(contentWidth, contentHeight int) *ViewportStyle {
	horizontalPos := &Float{Value: 0.0}
	verticalPos := &Float{Value: 0.0}
	horizontalProportion := &Float{Value: 0.3}
	verticalProportion := &Float{Value: 0.4}

	return &ViewportStyle{
		ContentWidth:            contentWidth,
		ContentHeight:           contentHeight,
		ScrollbarWidth:          unit.Dp(4),
		ScrollbarColor:          color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80},
		TrackColor:              color.NRGBA{A: 0xFF, R: 0xC0, G: 0xC0, B: 0xC0},
		BorderColor:             color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0x80},
		BorderWidth:             unit.Dp(1),
		HorizontalPos:           horizontalPos,
		VerticalPos:             verticalPos,
		HorizontalProportion:    horizontalProportion,
		VerticalProportion:      verticalProportion,
		scrollQueue:             make([]ScrollEvent, 0),
		scrollAnimationDuration: 200 * time.Millisecond,
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

	// Handle scroll events
	v.handleScrollEvents(gtx)

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
						// Register gesture.Scroll handler FIRST to capture scroll events
						defer clip.Rect(image.Rectangle{
							Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
						}).Push(gtx.Ops).Pop()
						v.scroller.Add(gtx.Ops)

						// Clip to viewport size
						defer clip.Rect(image.Rectangle{
							Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
						}).Push(gtx.Ops).Pop()

						// Offset content based on scroll position
						trans := op.Offset(image.Point{X: -contentOffsetX, Y: -contentOffsetY}).Push(gtx.Ops)

						// Draw scrollable content
						content(gtx)

						trans.Pop()

						// Add overlay click handler LAST to have lowest precedence
						// This captures events that weren't consumed by content widgets
						defer clip.Rect(image.Rectangle{
							Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
						}).Push(gtx.Ops).Pop()
						event.Op(gtx.Ops, v)

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

// HandleEvent implements event.Handler for the viewport
func (v *ViewportStyle) HandleEvent(ev event.Event) {
	// This method is called by the event system when events are registered
	// The actual scroll event processing is done in handleScrollEvents
	if pointerEvent, ok := ev.(pointer.Event); ok {
		if pointerEvent.Kind == pointer.Scroll {
			fmt.Printf("Viewport HandleEvent: Received scroll event - Y=%.1f, Position=(%.1f,%.1f)\n",
				pointerEvent.Scroll.Y, pointerEvent.Position.X, pointerEvent.Position.Y)
		}
	}
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

// handleScrollEvents processes scroll events and manages smooth scrolling animation
func (v *ViewportStyle) handleScrollEvents(gtx layout.Context) {
	// Calculate scroll ranges based on viewport dimensions
	viewportWidth := gtx.Constraints.Max.X - 20
	viewportHeight := gtx.Constraints.Max.Y - 20

	maxHScroll := max(0, v.ContentWidth-viewportWidth)
	maxVScroll := max(0, v.ContentHeight-viewportHeight)

	var scrollX, scrollY pointer.ScrollRange
	if maxHScroll > 0 {
		currentScrollX := int(float32(maxHScroll) * v.HorizontalPos.Value)
		scrollX.Min = -currentScrollX
		scrollX.Max = maxHScroll - currentScrollX
	}
	if maxVScroll > 0 {
		currentScrollY := int(float32(maxVScroll) * v.VerticalPos.Value)
		scrollY.Min = -currentScrollY
		scrollY.Max = maxVScroll - currentScrollY
	}

	// Use gesture.Scroll to handle both wheel and drag events
	sdist := v.scroller.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical, scrollX, scrollY)

	if sdist != 0 {
		fmt.Printf("Viewport: Received gesture.Scroll event - distance=%d\n", sdist)
		// Convert scroll distance to scroll amount
		scrollAmount := float32(sdist)
		v.processScrollEvent(gtx, scrollAmount)
	}

	// Also try to capture direct scroll events as a fallback
	for {
		event, ok := gtx.Source.Event(pointer.Filter{
			Target: v,
			Kinds:  pointer.Scroll,
		})
		if !ok {
			break
		}
		if pointerEvent, ok := event.(pointer.Event); ok {
			if pointerEvent.Kind == pointer.Scroll {
				fmt.Printf("Viewport: Received direct scroll pointer event - Y=%.1f, Position=(%.1f,%.1f)\n",
					pointerEvent.Scroll.Y, pointerEvent.Position.X, pointerEvent.Position.Y)
				v.processScrollEvent(gtx, pointerEvent.Scroll.Y)
			}
		}
	}

	// Process scroll animation
	if v.scrollAnimating {
		elapsed := gtx.Now.Sub(v.scrollAnimationStart)
		progress := float32(elapsed) / float32(v.scrollAnimationDuration)

		if progress >= 1.0 {
			// Animation complete
			v.scrollAnimating = false
			v.VerticalPos.Value = float32(v.scrollAnimationEndPos.Y) / float32(v.ContentHeight)
			v.HorizontalPos.Value = float32(v.scrollAnimationEndPos.X) / float32(v.ContentWidth)
		} else {
			// Interpolate between start and end positions
			currentY := v.scrollAnimationStartPos.Y + int(float32(v.scrollAnimationEndPos.Y-v.scrollAnimationStartPos.Y)*easeOutCubic(progress))
			currentX := v.scrollAnimationStartPos.X + int(float32(v.scrollAnimationEndPos.X-v.scrollAnimationStartPos.X)*easeOutCubic(progress))

			v.VerticalPos.Value = float32(currentY) / float32(v.ContentHeight)
			v.HorizontalPos.Value = float32(currentX) / float32(v.ContentWidth)

			// Invalidate to trigger repaint for next frame
			gtx.Execute(op.InvalidateCmd{})
		}
	}

	// Invalidate if we are animating
	if v.scrollAnimating {
		gtx.Execute(op.InvalidateCmd{})
	}
}

// processScrollEvent handles a scroll event with smooth animation
func (v *ViewportStyle) processScrollEvent(gtx layout.Context, scrollY float32) {
	fmt.Printf("Viewport: Processing scroll event - scrollY=%.1f\n", scrollY)

	// Calculate viewport dimensions
	viewportWidth := gtx.Constraints.Max.X - 20
	viewportHeight := gtx.Constraints.Max.Y - 20

	maxHScroll := max(0, v.ContentWidth-viewportWidth)
	maxVScroll := max(0, v.ContentHeight-viewportHeight)

	// Calculate scroll amount (10% of viewport for mouse wheel, direct for drag)
	var scrollAmount float32
	if absFloat32(scrollY) < 10 { // Mouse wheel events are typically small values
		scrollAmount = float32(viewportHeight) * 0.1
		if scrollY < 0 {
			scrollAmount = -scrollAmount // Scroll up
		}
		fmt.Printf("Viewport: Mouse wheel scroll - calculated amount=%.1f\n", scrollAmount)
	} else { // Drag events are larger values
		scrollAmount = scrollY
		fmt.Printf("Viewport: Drag scroll - using direct amount=%.1f\n", scrollAmount)
	}

	// Calculate current scroll position in pixels
	currentScrollX := int(float32(maxHScroll) * v.HorizontalPos.Value)
	currentScrollY := int(float32(maxVScroll) * v.VerticalPos.Value)

	// Calculate target scroll position
	targetScrollY := currentScrollY + int(scrollAmount)
	targetScrollX := currentScrollX

	// Clamp to valid ranges
	if targetScrollY < 0 {
		targetScrollY = 0
	}
	if targetScrollY > maxVScroll {
		targetScrollY = maxVScroll
	}
	if targetScrollX < 0 {
		targetScrollX = 0
	}
	if targetScrollX > maxHScroll {
		targetScrollX = maxHScroll
	}

	// Start smooth animation
	v.scrollAnimating = true
	v.scrollAnimationStart = gtx.Now
	v.scrollAnimationStartPos = image.Point{X: currentScrollX, Y: currentScrollY}
	v.scrollAnimationEndPos = image.Point{X: targetScrollX, Y: targetScrollY}

	// Invalidate to start animation
	gtx.Execute(op.InvalidateCmd{})
}

// absFloat32 returns the absolute value of a float32
func absFloat32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
