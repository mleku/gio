// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"
	"log"

	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
)

// ContextWidget represents a widget that can provide a context menu.
// When a widget implements this interface, it can specify what widget
// should be displayed when right-clicked.
type ContextWidget interface {
	// ContextMenu returns the widget to display as a context menu,
	// or nil if no context menu should be shown.
	ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget
}

// ContextManager manages the display of context widgets.
type ContextManager struct {
	// activeContext holds the currently displayed context widget
	activeContext *activeContext
	// viewportSize is the current viewport size for positioning calculations
	viewportSize image.Point
	// registeredWidgets stores widgets that can provide context menus
	registeredWidgets []registeredWidget
	// scrim clickable for dismissing context widgets
	scrimClickable *Clickable
}

type activeContext struct {
	widget      layout.Widget
	position    image.Point
	tag         event.Tag
	clickPos    image.Point
	dimensions  layout.Dimensions
	firstLayout bool // Track if this is the first layout call
	layoutCount int  // Track number of layout calls
}

type registeredWidget struct {
	widget   ContextWidget
	priority int
	bounds   image.Rectangle // Track widget bounds for hit detection
}

// NewContextManager creates a new context manager.
func NewContextManager() *ContextManager {
	return &ContextManager{
		scrimClickable: &Clickable{},
	}
}

// RegisterWidget registers a widget that can provide context menus.
func (cm *ContextManager) RegisterWidget(widget ContextWidget, priority int) {
	cm.registeredWidgets = append(cm.registeredWidgets, registeredWidget{
		widget:   widget,
		priority: priority,
		bounds:   image.Rectangle{}, // Will be updated during layout
	})
}

// UpdateWidgetBounds updates the bounds of a registered widget.
func (cm *ContextManager) UpdateWidgetBounds(widget ContextWidget, bounds image.Rectangle) {
	for i := range cm.registeredWidgets {
		if cm.registeredWidgets[i].widget == widget {
			cm.registeredWidgets[i].bounds = bounds
			break
		}
	}
}

// Update processes right-click events and manages context widget display.
func (cm *ContextManager) Update(gtx layout.Context) {
	log.Printf("*** CONTEXT MANAGER UPDATE CALLED ***")
	log.Printf("ContextManager.Update called, activeContext=%t", cm.activeContext != nil)

	// Always process right-click events first, regardless of active context state
	log.Printf("Processing right-click events, registeredWidgets=%d", len(cm.registeredWidgets))
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: cm,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		// Only handle right-clicks
		if e.Buttons != pointer.ButtonSecondary {
			continue
		}

		log.Printf("Right-click detected at position: %v", e.Position.Round())

		// Find the highest priority widget that contains the click
		var clickedWidget *registeredWidget
		highestPriority := -1
		clickPos := e.Position.Round()

		for _, reg := range cm.registeredWidgets {
			// Check if click is within widget bounds
			if clickPos.In(reg.bounds) && reg.priority > highestPriority {
				clickedWidget = &reg
				highestPriority = reg.priority
			}
		}

		// If no widget was hit, use the lowest priority widget (background)
		if clickedWidget == nil {
			lowestPriority := cm.registeredWidgets[0]
			for _, reg := range cm.registeredWidgets {
				if reg.priority < lowestPriority.priority {
					lowestPriority = reg
				}
			}
			clickedWidget = &lowestPriority
		}

		if clickedWidget.widget != nil {
			contextWidget := clickedWidget.widget.ContextMenu(gtx, clickPos)
			if contextWidget != nil {
				cm.ShowContextWidget(gtx, contextWidget, clickPos)
			}
		}
	}

	// If context menu is active, check for scrim clicks to dismiss it
	if cm.activeContext != nil {
		// Register scrim clickable for pointer events
		event.Op(gtx.Ops, cm.scrimClickable)

		// Check for clicks on scrim area
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: cm.scrimClickable,
				Kinds:  pointer.Press,
			})
			if !ok {
				break
			}
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}

			// Double-check that we still have an active context
			if cm.activeContext == nil {
				log.Printf("Scrim click received but no active context - ignoring")
				continue
			}

			clickPos := e.Position.Round()

			// Check if click is outside the context menu area
			contextMenuBounds := image.Rectangle{
				Min: cm.activeContext.position,
				Max: cm.activeContext.position.Add(cm.activeContext.dimensions.Size),
			}

			if !clickPos.In(contextMenuBounds) {
				log.Printf("Click outside context menu, dismissing")
				cm.dismissContextWidget()
				return
			}
		}
		return
	}

}

// Layout renders the active context widget if any.
func (cm *ContextManager) Layout(gtx layout.Context) layout.Dimensions {
	cm.viewportSize = gtx.Constraints.Max

	if cm.activeContext == nil {
		return layout.Dimensions{}
	}

	// Calculate dimensions on first layout call if not already calculated
	if cm.activeContext.dimensions.Size.X == 0 && cm.activeContext.dimensions.Size.Y == 0 {
		cm.activeContext.layoutCount++

		macro := op.Record(gtx.Ops)
		tempGtx := gtx
		tempGtx.Constraints = layout.Constraints{
			Min: image.Point{},
			Max: image.Point{X: 300, Y: 200}, // Reasonable max size for context menus
		}
		dims := cm.activeContext.widget(tempGtx)
		macro.Stop()

		// If dimensions are zero, dismiss the context menu
		if dims.Size.X == 0 && dims.Size.Y == 0 {
			if cm.activeContext.layoutCount <= 3 {
				// Don't dismiss for first few layout calls with zero dimensions
			} else {
				cm.dismissContextWidget()
				return layout.Dimensions{Size: cm.viewportSize}
			}
		}

		cm.activeContext.dimensions = dims
	}

	// Calculate the optimal position for the context widget first
	pos := cm.calculateOptimalPosition(gtx, cm.activeContext.clickPos, cm.activeContext.dimensions)
	widgetSize := cm.activeContext.dimensions.Size

	// Handle scrim - clickable area covering viewport but excluding context menu
	// Layout the scrim clickable to cover the entire viewport
	cm.scrimClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: cm.viewportSize}
	})

	// Check if click is outside the context menu area
	// This will be handled in the Update method by checking click position

	// Draw scrim areas (semi-transparent overlay)
	// Top scrim
	topScrim := image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: cm.viewportSize.X, Y: pos.Y},
	}
	if topScrim.Max.Y > 0 {
		clipStack := clip.Rect(topScrim).Push(gtx.Ops)
		paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33}) // 20% opacity black
		clipStack.Pop()
	}

	// Bottom scrim
	bottomScrim := image.Rectangle{
		Min: image.Point{X: 0, Y: pos.Y + widgetSize.Y},
		Max: cm.viewportSize,
	}
	if bottomScrim.Min.Y < cm.viewportSize.Y {
		clipStack := clip.Rect(bottomScrim).Push(gtx.Ops)
		paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33}) // 20% opacity black
		clipStack.Pop()
	}

	// Left scrim
	leftScrim := image.Rectangle{
		Min: image.Point{X: 0, Y: pos.Y},
		Max: image.Point{X: pos.X, Y: pos.Y + widgetSize.Y},
	}
	if leftScrim.Max.X > 0 {
		clipStack := clip.Rect(leftScrim).Push(gtx.Ops)
		paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33}) // 20% opacity black
		clipStack.Pop()
	}

	// Right scrim
	rightScrim := image.Rectangle{
		Min: image.Point{X: pos.X + widgetSize.X, Y: pos.Y},
		Max: image.Point{X: cm.viewportSize.X, Y: pos.Y + widgetSize.Y},
	}
	if rightScrim.Min.X < cm.viewportSize.X {
		clipStack := clip.Rect(rightScrim).Push(gtx.Ops)
		paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x33}) // 20% opacity black
		clipStack.Pop()
	}

	// Record the context widget
	macro := op.Record(gtx.Ops)
	contextGtx := gtx
	contextGtx.Constraints = layout.Exact(cm.activeContext.dimensions.Size)
	dims := cm.activeContext.widget(contextGtx)
	call := macro.Stop()

	// If the widget returns zero dimensions, dismiss it
	if dims.Size.X == 0 && dims.Size.Y == 0 {
		cm.dismissContextWidget()
		return layout.Dimensions{Size: cm.viewportSize}
	}

	// Position and render the context widget
	trans := op.Offset(pos).Push(gtx.Ops)
	defer trans.Pop()

	// Add clipping to ensure the context widget doesn't go outside viewport
	clipRect := image.Rectangle{
		Min: image.Point{},
		Max: cm.viewportSize,
	}
	defer clip.Rect(clipRect).Push(gtx.Ops).Pop()

	call.Add(gtx.Ops)

	return layout.Dimensions{Size: cm.viewportSize}
}

// ShowContextWidget displays a context widget at the specified position.
func (cm *ContextManager) ShowContextWidget(gtx layout.Context, widget layout.Widget, clickPos image.Point) {
	// Don't calculate dimensions immediately - defer until layout phase
	// This prevents close button clicks from being processed during dimension calculation
	cm.activeContext = &activeContext{
		widget:      widget,
		clickPos:    clickPos,
		dimensions:  layout.Dimensions{}, // Will be calculated during layout
		tag:         &cm.activeContext,
		firstLayout: true, // Mark as first layout call
		layoutCount: 0,    // Initialize layout count
	}
}

// dismissContextWidget hides the currently displayed context widget.
func (cm *ContextManager) dismissContextWidget() {
	cm.activeContext = nil
	// Reset scrim clickable to ensure it doesn't interfere with new events
	cm.scrimClickable = &Clickable{}
}

// isPointInContextWidget checks if a point is within the active context widget.
func (cm *ContextManager) isPointInContextWidget(pos image.Point) bool {
	if cm.activeContext == nil {
		return false
	}

	// Calculate the bounds of the context widget
	bounds := image.Rectangle{
		Min: cm.activeContext.position,
		Max: cm.activeContext.position.Add(cm.activeContext.dimensions.Size),
	}

	return pos.In(bounds)
}

// calculateOptimalPosition calculates the best position for a context widget
// to avoid going outside the viewport.
func (cm *ContextManager) calculateOptimalPosition(gtx layout.Context, clickPos image.Point, dims layout.Dimensions) image.Point {
	widgetSize := dims.Size
	viewportSize := cm.viewportSize

	// Calculate the center of the viewport
	center := image.Point{
		X: viewportSize.X / 2,
		Y: viewportSize.Y / 2,
	}

	// Determine which corner to use based on click position relative to center
	var corner image.Point

	if clickPos.X < center.X {
		// Click is on the left side, use left edge
		corner.X = clickPos.X
	} else {
		// Click is on the right side, use right edge
		corner.X = clickPos.X - widgetSize.X
	}

	if clickPos.Y < center.Y {
		// Click is on the top side, use top edge
		corner.Y = clickPos.Y
	} else {
		// Click is on the bottom side, use bottom edge
		corner.Y = clickPos.Y - widgetSize.Y
	}

	// Ensure the widget doesn't go outside the viewport
	if corner.X < 0 {
		corner.X = 0
	}
	if corner.Y < 0 {
		corner.Y = 0
	}
	if corner.X+widgetSize.X > viewportSize.X {
		corner.X = viewportSize.X - widgetSize.X
	}
	if corner.Y+widgetSize.Y > viewportSize.Y {
		corner.Y = viewportSize.Y - widgetSize.Y
	}

	// Update the active context position
	if cm.activeContext != nil {
		cm.activeContext.position = corner
	}

	return corner
}

// AddContextHandler adds a context handler to the operations list.
func (cm *ContextManager) AddContextHandler(ops *op.Ops) {
	event.Op(ops, cm)
}
