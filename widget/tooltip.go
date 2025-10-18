package widget

import (
	"image"
	"image/color"
	"time"

	"gio.mleku.dev/font"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
)

// TooltipWidget represents a widget that can provide a tooltip.
type TooltipWidget interface {
	// Tooltip returns the tooltip text to display, or empty string if no tooltip should be shown.
	Tooltip() string
}

// TooltipManager manages the display of tooltips.
type TooltipManager struct {
	// activeTooltip holds the currently displayed tooltip
	activeTooltip *activeTooltip
	// viewportSize is the current viewport size for positioning calculations
	viewportSize image.Point
	// registeredWidgets stores widgets that can provide tooltips
	registeredWidgets []registeredTooltipWidget
	// hoverTimeout is the delay before showing tooltip
	hoverTimeout time.Duration
	// hoverStartTime tracks when hover started
	hoverStartTime time.Time
	// isHovering tracks if we're currently hovering
	isHovering bool
	// hoveredWidget tracks which widget is being hovered
	hoveredWidget TooltipWidget
	// mousePosition tracks current mouse position
	mousePosition image.Point
	// invalidateFunc is called to trigger repaints
	invalidateFunc func()
}

type activeTooltip struct {
	text       string
	position   image.Point
	tag        event.Tag
	dimensions layout.Dimensions
	theme      *text.Shaper
}

type registeredTooltipWidget struct {
	widget TooltipWidget
	bounds image.Rectangle // Track widget bounds for hit detection
}

// NewTooltipManager creates a new tooltip manager.
func NewTooltipManager(invalidateFunc func()) *TooltipManager {
	return &TooltipManager{
		hoverTimeout:   500 * time.Millisecond, // 500ms delay before showing tooltip
		invalidateFunc: invalidateFunc,
	}
}

// RegisterWidget registers a widget that can provide tooltips.
func (tm *TooltipManager) RegisterWidget(widget TooltipWidget) {
	tm.registeredWidgets = append(tm.registeredWidgets, registeredTooltipWidget{
		widget: widget,
		bounds: image.Rectangle{}, // Will be updated during layout
	})
}

// UpdateWidgetBounds updates the bounds of a registered widget.
func (tm *TooltipManager) UpdateWidgetBounds(widget TooltipWidget, bounds image.Rectangle) {
	for i := range tm.registeredWidgets {
		if tm.registeredWidgets[i].widget == widget {
			tm.registeredWidgets[i].bounds = bounds
			break
		}
	}
}

// Update processes hover events and manages tooltip display.
func (tm *TooltipManager) Update(gtx layout.Context) {

	// Process hover events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: tm,
			Kinds:  pointer.Move | pointer.Enter | pointer.Leave,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		switch e.Kind {
		case pointer.Enter, pointer.Move:
			// Update mouse position
			tm.mousePosition = e.Position.Round()

			// Find the widget under the cursor
			var hoveredWidget *registeredTooltipWidget
			pos := e.Position.Round()

			for _, reg := range tm.registeredWidgets {
				if pos.In(reg.bounds) {
					hoveredWidget = &reg
					break
				}
			}

			if hoveredWidget != nil {
				// Check if we're hovering a different widget
				if tm.hoveredWidget != hoveredWidget.widget {
					// If we were hovering a different widget, dismiss the current tooltip
					if tm.isHovering && tm.activeTooltip != nil {
						tm.dismissTooltip()
					}
					tm.hoveredWidget = hoveredWidget.widget
					tm.hoverStartTime = time.Now()
					tm.isHovering = true
				}
			} else {
				// Not hovering any widget - dismiss tooltip if we were hovering
				if tm.isHovering {
					tm.dismissTooltip()
				}
			}

		case pointer.Leave:
			// Mouse left the area
			if tm.isHovering {
				tm.dismissTooltip()
			}
		}
	}

	// Check if we should show tooltip after hover timeout
	if tm.isHovering && tm.hoveredWidget != nil && tm.activeTooltip == nil {
		elapsed := time.Since(tm.hoverStartTime)
		if elapsed >= tm.hoverTimeout {
			tooltipText := tm.hoveredWidget.Tooltip()
			if tooltipText != "" {
				tm.showTooltip(gtx, tooltipText)
			}
		} else {
			// Keep invalidating to check the timeout on the next frame
			if tm.invalidateFunc != nil {
				tm.invalidateFunc()
			}
		}
	}
}

// Layout renders the active tooltip if any.
func (tm *TooltipManager) Layout(gtx layout.Context, shaper *text.Shaper) layout.Dimensions {
	tm.viewportSize = gtx.Constraints.Max

	if tm.activeTooltip == nil {
		return layout.Dimensions{}
	}

	// Calculate tooltip dimensions if not already calculated
	if tm.activeTooltip.dimensions.Size.X == 0 && tm.activeTooltip.dimensions.Size.Y == 0 {
		tm.activeTooltip.theme = shaper

		// Create tooltip widget
		tooltipWidget := tm.createTooltipWidget(tm.activeTooltip.text)

		// Calculate dimensions
		macro := op.Record(gtx.Ops)
		tempGtx := gtx
		tempGtx.Constraints = layout.Constraints{
			Min: image.Point{},
			Max: image.Point{X: 300, Y: 100}, // Reasonable max size for tooltips
		}
		dims := tooltipWidget(tempGtx)
		macro.Stop()

		tm.activeTooltip.dimensions = dims
	}

	// Calculate the optimal position for the tooltip
	pos := tm.calculateOptimalPosition(gtx, tm.activeTooltip.position, tm.activeTooltip.dimensions)

	// Create and render the tooltip widget
	tooltipWidget := tm.createTooltipWidget(tm.activeTooltip.text)

	// Record the tooltip widget
	macro := op.Record(gtx.Ops)
	contextGtx := gtx
	contextGtx.Constraints = layout.Exact(tm.activeTooltip.dimensions.Size)
	tooltipWidget(contextGtx)
	call := macro.Stop()

	// Position and render the tooltip widget
	trans := op.Offset(pos).Push(gtx.Ops)
	defer trans.Pop()

	// Add clipping to ensure the tooltip doesn't go outside viewport
	clipRect := image.Rectangle{
		Min: image.Point{},
		Max: tm.viewportSize,
	}
	defer clip.Rect(clipRect).Push(gtx.Ops).Pop()

	call.Add(gtx.Ops)

	return layout.Dimensions{Size: tm.viewportSize}
}

// createTooltipWidget creates a tooltip widget with the given text.
func (tm *TooltipManager) createTooltipWidget(tooltipText string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		// Create tooltip with 0.5 text height inset around it
		inset := unit.Dp(8) // 0.5 text height (assuming 16dp text height)

		return layout.Inset{
			Top:    inset,
			Bottom: inset,
			Left:   inset,
			Right:  inset,
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// Draw tooltip background - yellow
			clipStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}) // Yellow background
			clipStack.Pop()

			// Draw tooltip text using proper text rendering
			if tm.activeTooltip != nil && tm.activeTooltip.theme != nil {
				// Create text color macro
				textColorMacro := op.Record(gtx.Ops)
				paint.ColorOp{Color: color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}}.Add(gtx.Ops) // Black text
				textColor := textColorMacro.Stop()

				// Create label widget for text rendering
				label := Label{
					Alignment: text.Start,
					MaxLines:  1,
				}

				// Layout the text
				textSize := unit.Sp(14) // Small text size for tooltips
				font := font.Font{}     // Use default font

				label.Layout(gtx, tm.activeTooltip.theme, font, textSize, tooltipText, textColor)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			}

			// Fallback if no theme available
			return layout.Dimensions{Size: gtx.Constraints.Max}
		})
	}
}

// showTooltip displays a tooltip with the specified text.
func (tm *TooltipManager) showTooltip(gtx layout.Context, text string) {
	// Position tooltip near mouse cursor
	pos := tm.mousePosition.Add(image.Point{X: 10, Y: -10}) // Offset from cursor

	tm.activeTooltip = &activeTooltip{
		text:     text,
		position: pos,
		tag:      &tm.activeTooltip,
	}

	// Invalidate to trigger repaint
	if tm.invalidateFunc != nil {
		tm.invalidateFunc()
	}
}

// dismissTooltip hides the currently displayed tooltip.
func (tm *TooltipManager) dismissTooltip() {
	tm.activeTooltip = nil
	tm.isHovering = false
	tm.hoveredWidget = nil

	// Invalidate to trigger repaint
	if tm.invalidateFunc != nil {
		tm.invalidateFunc()
	}
}

// calculateOptimalPosition calculates the best position for a tooltip
// to avoid going outside the viewport.
func (tm *TooltipManager) calculateOptimalPosition(gtx layout.Context, preferredPos image.Point, dims layout.Dimensions) image.Point {
	widgetSize := dims.Size
	viewportSize := tm.viewportSize

	// Start with preferred position
	pos := preferredPos

	// Ensure the tooltip doesn't go outside the viewport
	if pos.X < 0 {
		pos.X = 0
	}
	if pos.Y < 0 {
		pos.Y = 0
	}
	if pos.X+widgetSize.X > viewportSize.X {
		pos.X = viewportSize.X - widgetSize.X
	}
	if pos.Y+widgetSize.Y > viewportSize.Y {
		pos.Y = viewportSize.Y - widgetSize.Y
	}

	return pos
}

// AddTooltipHandler adds a tooltip handler to the operations list.
func (tm *TooltipManager) AddTooltipHandler(ops *op.Ops) {
	event.Op(ops, tm)
}
