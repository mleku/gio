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

// OverlayExample demonstrates the 2x2 grid layout.
type OverlayExample struct {
	overlayStack *widget.OverlayStack
	theme        *material.Theme
	window       *app.Window
	redTarget    *colorSquare
	yellowTarget *colorSquare
	greenTarget  *colorSquare
	blueTarget   *colorSquare
}

// colorSquare represents an individual colored square with its own event target
type colorSquare struct {
	color     color.NRGBA
	colorName string
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
		window:       window,
		redTarget:    &colorSquare{color: color.NRGBA{R: 255, G: 0, B: 0, A: 255}, colorName: "Red"},
		yellowTarget: &colorSquare{color: color.NRGBA{R: 255, G: 255, B: 0, A: 255}, colorName: "Yellow"},
		greenTarget:  &colorSquare{color: color.NRGBA{R: 0, G: 255, B: 0, A: 255}, colorName: "Green"},
		blueTarget:   &colorSquare{color: color.NRGBA{R: 0, G: 0, B: 255, A: 255}, colorName: "Blue"},
	}
}

// Layout renders the 2x2 grid with colored squares.
func (oe *OverlayExample) Layout(gtx layout.Context) layout.Dimensions {
	// Create a 2x2 grid layout
	gridDims := layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// First row
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Top-left square (red)
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return oe.drawSquareWithBorder(gtx, oe.redTarget)
				}),
				// Top-right square (yellow)
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return oe.drawSquareWithBorder(gtx, oe.yellowTarget)
				}),
			)
		}),
		// Second row
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis: layout.Horizontal,
			}.Layout(gtx,
				// Bottom-left square (green)
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return oe.drawSquareWithBorder(gtx, oe.greenTarget)
				}),
				// Bottom-right square (blue)
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return oe.drawSquareWithBorder(gtx, oe.blueTarget)
				}),
			)
		}),
	)

	// Layout overlay stack on top
	oe.overlayStack.Layout(gtx)

	return gridDims
}

// drawSquareWithBorder draws a square with 1px border and colored center square.
func (oe *OverlayExample) drawSquareWithBorder(gtx layout.Context, square *colorSquare) layout.Dimensions {
	// Add pointer input for right-click detection
	area := clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
	event.Op(gtx.Ops, square)
	area.Pop()

	// Handle right-click events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: square,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}
		if ev, ok := ev.(pointer.Event); ok {
			if ev.Kind == pointer.Press && ev.Buttons == pointer.ButtonSecondary {
				// Get mouse position in window coordinates
				mousePos := gtx.CursorPosition()
				windowDims := gtx.WindowSize()

				// Calculate content size based on text (approximate)
				contentWidth := 200
				contentHeight := 100

				// Calculate window center
				centerX := windowDims.X / 2
				centerY := windowDims.Y / 2

				var contentX, contentY int

				// Position overlay with mouse as one corner, opposite corner toward window center
				if mousePos.X < centerX {
					// Mouse is in left half - position content to the right of mouse
					contentX = mousePos.X
				} else {
					// Mouse is in right half - position content to the left of mouse
					contentX = mousePos.X - contentWidth
				}

				if mousePos.Y < centerY {
					// Mouse is in top half - position content below mouse
					contentY = mousePos.Y
				} else {
					// Mouse is in bottom half - position content above mouse
					contentY = mousePos.Y - contentHeight
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
					// Fill the entire content area with yellow background
					paint.FillShape(gtx.Ops, color.NRGBA{R: 255, G: 255, B: 0, A: 255}, clip.Rect{Max: gtx.Constraints.Max}.Op())

					// Add a black border
					borderRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
					paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Stroke{
						Path:  clip.Rect(borderRect).Path(),
						Width: 2,
					}.Op())

					// Center the color name in the content area
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						label := material.Label(oe.theme, unit.Sp(20), square.colorName)
						label.Color = color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Ensure text is black
						return label.Layout(gtx)
					})
				}

				// Show overlay using the stack
				overlayID := "color-overlay-" + square.colorName
				oe.overlayStack.Push(overlayID, overlayContent, image.Point{X: contentX, Y: contentY}, image.Point{X: contentWidth, Y: contentHeight})
			}
		}
	}

	// Set up overlay click handler
	oe.overlayStack.SetClickHandler("color-overlay-"+square.colorName, func() {
		log.Println("lol.mleku.dev: Overlay closed by scrim click")
	})

	// Draw the border (1px)
	borderRect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	paint.FillShape(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}, clip.Stroke{
		Path:  clip.Rect(borderRect).Path(),
		Width: 1,
	}.Op())

	// Draw the colored square in the center (100x100Dp)
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		size := gtx.Dp(100) // 100Dp
		rect := image.Rect(0, 0, size, size)
		paint.FillShape(gtx.Ops, square.color, clip.Rect(rect).Op())
		return layout.Dimensions{Size: image.Point{X: size, Y: size}}
	})
}
