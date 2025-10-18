package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
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

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	// Create overlay example
	overlayExample := NewOverlayExample(th, w)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Layout the overlay example
			overlayExample.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

// OverlayExample demonstrates the overlay widget functionality.
type OverlayExample struct {
	overlayStack *widget.OverlayStack
	theme        *material.Theme
	label        material.LabelStyle
	window       *app.Window
}

// NewOverlayExample creates a new overlay example.
func NewOverlayExample(theme *material.Theme, window *app.Window) *OverlayExample {
	overlayStack := widget.NewOverlayStack()
	overlayStack.SetInvalidateFunc(func() {
		window.Invalidate()
	})

	return &OverlayExample{
		overlayStack: overlayStack,
		theme:        theme,
		label:        material.Label(theme, unit.Sp(16), "Right-click to open overlay"),
		window:       window,
	}
}

// Layout renders the overlay example.
func (oe *OverlayExample) Layout(gtx layout.Context) layout.Dimensions {
	// Add pointer input for right-click detection
	area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, oe)
	area.Pop()

	// Handle right-click events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: oe,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if ev, ok := ev.(pointer.Event); ok {
			if ev.Kind == pointer.Press && ev.Buttons == pointer.ButtonSecondary {
				// Get mouse position
				mousePos := ev.Position
				windowDims := gtx.Constraints.Max

				// Calculate content size based on text (approximate)
				contentWidth := 200
				contentHeight := 100

				// Calculate which corner should be under the mouse
				centerX := float32(windowDims.X) / 2
				centerY := float32(windowDims.Y) / 2

				var contentX, contentY int

				// Determine position based on which quadrant the mouse is in
				if mousePos.X < centerX {
					// Mouse is in left half - position content to the right of mouse
					contentX = int(mousePos.X)
				} else {
					// Mouse is in right half - position content to the left of mouse
					contentX = int(mousePos.X) - contentWidth
				}

				if mousePos.Y < centerY {
					// Mouse is in top half - position content below mouse
					contentY = int(mousePos.Y)
				} else {
					// Mouse is in bottom half - position content above mouse
					contentY = int(mousePos.Y) - contentHeight
				}

				// Clamp to window bounds
				if contentX < 0 {
					contentX = 0
				} else if contentX+contentWidth > windowDims.X {
					contentX = windowDims.X - contentWidth
				}

				if contentY < 0 {
					contentY = 0
				} else if contentY+contentHeight > windowDims.Y {
					contentY = windowDims.Y - contentHeight
				}

				// Create content widget for overlay
				overlayContent := func(gtx layout.Context) layout.Dimensions {
					// Fill the entire content area with a bright background to make it visible
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 0, A: 255}, clip.Rect{Max: gtx.Constraints.Max}.Op())

					// Add a border to make it more visible
					borderRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Stroke{
						Path:  clip.Rect(borderRect).Path(),
						Width: 2,
					}.Op())

					// Center the label in the content area
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Label(oe.theme, unit.Sp(20), "OVERLAY CONTENT")
						label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Ensure text is black
						return label.Layout(gtx)
					})
				}

				// Show overlay using the stack
				overlayID := "main-overlay"
				oe.overlayStack.Push(overlayID, overlayContent, image.Point{X: contentX, Y: contentY}, image.Point{X: contentWidth, Y: contentHeight})
			}
		}
	}

	// Set up overlay click handler
	oe.overlayStack.SetClickHandler("main-overlay", func() {
		log.Println("lol.mleku.dev: Overlay closed by scrim click")
	})

	// Layout main content
	mainDims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return oe.label.Layout(gtx)
	})

	// Layout overlay stack on top
	oe.overlayStack.Layout(gtx)

	return mainDims
}
