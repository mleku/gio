// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"
	"image/color"
	"log"
	"time"

	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
)

// OverlayStack manages a stack of overlays that can be pushed and popped.
type OverlayStack struct {
	overlays       []*Overlay
	invalidateFunc func() // Function to call when invalidation is needed
}

// Overlay represents a single overlay in the stack.
type Overlay struct {
	id      string
	visible bool
	scrim   bool
	color   [4]float32 // RGBA color for scrim background

	// Position and size of the content widget
	contentPos  image.Point
	contentDims image.Point

	// Widget to display in the overlay
	content layout.Widget

	// Click handlers
	clickHandler func()

	// Z-index for ordering (higher values appear on top)
	zIndex int

	// Animation state
	opacity           float32 // Current opacity (0.0 to 1.0)
	targetOpacity     float32 // Target opacity
	animating         bool    // Whether currently animating
	animationDuration int64   // Animation duration in milliseconds
	startTime         int64   // Animation start time
}

// NewOverlayStack creates a new overlay stack.
func NewOverlayStack() *OverlayStack {
	return &OverlayStack{
		overlays: make([]*Overlay, 0),
	}
}

// SetInvalidateFunc sets the function to call when the window needs to be invalidated.
func (os *OverlayStack) SetInvalidateFunc(fn func()) {
	os.invalidateFunc = fn
}

// Push adds a new overlay to the stack.
func (os *OverlayStack) Push(id string, content layout.Widget, pos image.Point, dims image.Point) *Overlay {
	overlay := &Overlay{
		id:                id,
		visible:           true,
		scrim:             true,
		color:             [4]float32{0, 0, 0, 0.5}, // Semi-transparent black
		content:           content,
		contentPos:        pos,
		contentDims:       dims,
		zIndex:            len(os.overlays), // Higher z-index for newer overlays
		opacity:           0.0,              // Start invisible for fade in
		targetOpacity:     1.0,              // Fade to fully visible
		animating:         true,
		animationDuration: 200, // 200ms animation
		startTime:         time.Now().UnixMilli(),
	}

	// Remove any existing overlay with the same ID
	os.Remove(id)

	os.overlays = append(os.overlays, overlay)
	return overlay
}

// Pop removes the topmost overlay from the stack.
func (os *OverlayStack) Pop() {
	if len(os.overlays) == 0 {
		return
	}
	os.overlays = os.overlays[:len(os.overlays)-1]
}

// Remove removes an overlay by ID.
func (os *OverlayStack) Remove(id string) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			// Start fade out animation instead of immediate removal
			overlay.targetOpacity = 0.0
			overlay.animating = true
			overlay.startTime = time.Now().UnixMilli()
			overlay.animationDuration = 200 // 200ms fade out
			os.invalidate()
			return
		}
	}
}

// Show makes an overlay visible.
func (os *OverlayStack) Show(id string) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.visible = true
			os.invalidate()
			return
		}
	}
}

// Hide makes an overlay invisible.
func (os *OverlayStack) Hide(id string) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.visible = false
			os.invalidate()
			return
		}
	}
}

// SetScrim enables or disables the scrim background for a specific overlay.
func (os *OverlayStack) SetScrim(id string, enabled bool) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.scrim = enabled
			return
		}
	}
}

// SetScrimColor sets the color of the scrim background for a specific overlay.
func (os *OverlayStack) SetScrimColor(id string, r, g, b, a float32) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.color = [4]float32{r, g, b, a}
			return
		}
	}
}

// SetClickHandler sets the function to call when the scrim is clicked for a specific overlay.
func (os *OverlayStack) SetClickHandler(id string, handler func()) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.clickHandler = handler
			return
		}
	}
}

// UpdateContent updates the content widget for a specific overlay.
func (os *OverlayStack) UpdateContent(id string, content layout.Widget, pos image.Point, dims image.Point) {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			overlay.content = content
			overlay.contentPos = pos
			overlay.contentDims = dims
			return
		}
	}
}

// IsVisible returns whether an overlay is currently visible.
func (os *OverlayStack) IsVisible(id string) bool {
	for _, overlay := range os.overlays {
		if overlay.id == id {
			return overlay.visible
		}
	}
	return false
}

// Count returns the number of overlays in the stack.
func (os *OverlayStack) Count() int {
	return len(os.overlays)
}

// Clear removes all overlays from the stack.
func (os *OverlayStack) Clear() {
	os.overlays = os.overlays[:0]
}

// Layout renders all visible overlays in the stack.
func (os *OverlayStack) Layout(gtx layout.Context) layout.Dimensions {
	if len(os.overlays) == 0 {
		return layout.Dimensions{}
	}

	// Update animations
	os.updateAnimations()

	// Get the window dimensions
	windowDims := gtx.Constraints.Max

	// Render overlays in z-index order (lowest first)
	for _, overlay := range os.overlays {
		if !overlay.visible {
			continue
		}
		os.layoutOverlay(gtx, overlay, windowDims)
	}

	return layout.Dimensions{Size: windowDims}
}

// updateAnimations updates all overlay animations and removes completed fade-outs.
func (os *OverlayStack) updateAnimations() {
	currentTime := time.Now().UnixMilli()
	needsInvalidation := false

	// Process overlays in reverse order to safely remove items
	for i := len(os.overlays) - 1; i >= 0; i-- {
		overlay := os.overlays[i]

		if !overlay.animating {
			continue
		}

		elapsed := currentTime - overlay.startTime
		progress := float32(elapsed) / float32(overlay.animationDuration)

		if progress >= 1.0 {
			// Animation complete
			overlay.opacity = overlay.targetOpacity
			overlay.animating = false

			// If fade out is complete, remove the overlay
			if overlay.targetOpacity == 0.0 {
				os.overlays = append(os.overlays[:i], os.overlays[i+1:]...)
			}
			needsInvalidation = true
		} else {
			// Interpolate opacity
			overlay.opacity = overlay.opacity + (overlay.targetOpacity-overlay.opacity)*progress
			needsInvalidation = true
		}
	}

	// Invalidate if any animations are running
	if needsInvalidation {
		os.invalidate()
	}
}

// layoutOverlay renders a single overlay.
func (os *OverlayStack) layoutOverlay(gtx layout.Context, overlay *Overlay, windowDims image.Point) {
	// Draw scrim background if enabled
	if overlay.scrim {
		scrimColor := color.NRGBA{
			R: uint8(overlay.color[0] * 255),
			G: uint8(overlay.color[1] * 255),
			B: uint8(overlay.color[2] * 255),
			A: uint8(overlay.color[3] * overlay.opacity * 255),
		}

		// Draw scrim above content
		if overlay.contentPos.Y > 0 {
			topScrim := image.Rect(0, 0, windowDims.X, overlay.contentPos.Y)
			area := clip.Rect{Min: topScrim.Min, Max: topScrim.Max}.Push(gtx.Ops)
			paint.FillShape(gtx.Ops, scrimColor, clip.Rect{Min: topScrim.Min, Max: topScrim.Max}.Op())
			event.Op(gtx.Ops, overlay)
			area.Pop()
		}

		// Draw scrim below content
		if overlay.contentPos.Y+overlay.contentDims.Y < windowDims.Y {
			bottomScrim := image.Rect(0, overlay.contentPos.Y+overlay.contentDims.Y, windowDims.X, windowDims.Y)
			area := clip.Rect{Min: bottomScrim.Min, Max: bottomScrim.Max}.Push(gtx.Ops)
			paint.FillShape(gtx.Ops, scrimColor, clip.Rect{Min: bottomScrim.Min, Max: bottomScrim.Max}.Op())
			event.Op(gtx.Ops, overlay)
			area.Pop()
		}

		// Draw scrim to the left of content
		if overlay.contentPos.X > 0 {
			leftScrim := image.Rect(0, overlay.contentPos.Y, overlay.contentPos.X, overlay.contentPos.Y+overlay.contentDims.Y)
			area := clip.Rect{Min: leftScrim.Min, Max: leftScrim.Max}.Push(gtx.Ops)
			paint.FillShape(gtx.Ops, scrimColor, clip.Rect{Min: leftScrim.Min, Max: leftScrim.Max}.Op())
			event.Op(gtx.Ops, overlay)
			area.Pop()
		}

		// Draw scrim to the right of content
		if overlay.contentPos.X+overlay.contentDims.X < windowDims.X {
			rightScrim := image.Rect(overlay.contentPos.X+overlay.contentDims.X, overlay.contentPos.Y, windowDims.X, overlay.contentPos.Y+overlay.contentDims.Y)
			area := clip.Rect{Min: rightScrim.Min, Max: rightScrim.Max}.Push(gtx.Ops)
			paint.FillShape(gtx.Ops, scrimColor, clip.Rect{Min: rightScrim.Min, Max: rightScrim.Max}.Op())
			event.Op(gtx.Ops, overlay)
			area.Pop()
		}
	}

	// Save the current transformation matrix and translate to content position
	transformStack := op.Offset(image.Point{X: overlay.contentPos.X, Y: overlay.contentPos.Y}).Push(gtx.Ops)
	defer transformStack.Pop()

	// Clip to content area for events and rendering
	contentArea := clip.Rect{Min: image.Point{}, Max: overlay.contentDims}.Push(gtx.Ops)

	// Add event handler for content area to prevent scrim clicks
	event.Op(gtx.Ops, &overlay.contentPos) // Use contentPos as unique tag

	// Apply opacity to content
	opacityStack := paint.PushOpacity(gtx.Ops, overlay.opacity)
	defer opacityStack.Pop()

	// Create new context for content with proper constraints
	contentGtx := gtx
	contentGtx.Constraints = layout.Exact(overlay.contentDims)

	// Render the content widget
	overlay.content(contentGtx)

	// Pop the clip area
	contentArea.Pop()

	// Handle pointer events for this overlay
	// Process all available events
	for {
		// Check for content area clicks first
		ev, ok := gtx.Event(pointer.Filter{
			Target: &overlay.contentPos,
			Kinds:  pointer.Press,
		})
		if ok {
			if ev, ok := ev.(pointer.Event); ok && ev.Kind == pointer.Press {
				log.Println("lol.mleku.dev: Content clicked - overlay stays open")
			}
			continue // Process next event
		}

		// Check for scrim clicks
		if overlay.scrim {
			ev, ok := gtx.Event(pointer.Filter{
				Target: overlay,
				Kinds:  pointer.Press,
			})
			if ok {
				if ev, ok := ev.(pointer.Event); ok && ev.Kind == pointer.Press {
					log.Println("lol.mleku.dev: Scrim clicked - closing overlay")
					if overlay.clickHandler != nil {
						overlay.clickHandler()
					}
					os.Remove(overlay.id)
					return // Exit immediately
				}
				continue // Process next event
			}
		}

		// No more events to process
		break
	}
}

// invalidate triggers a window redraw if the invalidate function is set.
func (os *OverlayStack) invalidate() {
	if os.invalidateFunc != nil {
		os.invalidateFunc()
	}
}
