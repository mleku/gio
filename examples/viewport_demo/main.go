// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
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

	// Content dimensions (larger than viewport to require scrolling)
	contentWidth := 800
	contentHeight := 600

	// Create viewport widget
	viewport := widget.Viewport(contentWidth, contentHeight)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Layout the viewport with scrollable content
			viewport.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return drawScrollableContent(gtx, th, contentWidth, contentHeight)
			})

			e.Frame(gtx.Ops)
		}
	}
}

// drawScrollableContent draws the actual content that can be scrolled
func drawScrollableContent(gtx layout.Context, th *material.Theme, contentWidth, contentHeight int) layout.Dimensions {
	// Draw a grid pattern to show scrolling
	gridSize := 50
	gridColor := color.NRGBA{A: 0xFF, R: 0xE0, G: 0xE0, B: 0xE0}

	// Draw vertical grid lines
	for x := 0; x <= contentWidth; x += gridSize {
		paint.FillShape(gtx.Ops, gridColor, clip.Stroke{
			Path: clip.Rect(image.Rectangle{
				Min: image.Point{X: x, Y: 0},
				Max: image.Point{X: x + 1, Y: contentHeight},
			}).Path(),
			Width: 1,
		}.Op())
	}

	// Draw horizontal grid lines
	for y := 0; y <= contentHeight; y += gridSize {
		paint.FillShape(gtx.Ops, gridColor, clip.Stroke{
			Path: clip.Rect(image.Rectangle{
				Min: image.Point{X: 0, Y: y},
				Max: image.Point{X: contentWidth, Y: y + 1},
			}).Path(),
			Width: 1,
		}.Op())
	}

	// Draw some text content
	textColor := color.NRGBA{A: 0xFF, R: 0x00, G: 0x00, B: 0x00}
	paint.ColorOp{Color: textColor}.Add(gtx.Ops)

	// Draw title
	titleText := "Scrollable Viewport Example"
	titleDims := material.H3(th, titleText).Layout(gtx)

	// Draw some sample content
	sampleText := "This is a scrollable viewport with horizontal and vertical scrollbars.\n" +
		"You can scroll by dragging the scrollbar thumbs or clicking on the track.\n" +
		"The content is larger than the viewport, so scrolling is required to see everything."

	// Position sample text below title
	trans := op.Offset(image.Point{X: 0, Y: titleDims.Size.Y + 20}).Push(gtx.Ops)
	material.Body1(th, sampleText).Layout(gtx)
	trans.Pop()

	// Draw some colored rectangles to make scrolling more interesting
	rects := []struct {
		x, y, w, h int
		color      color.NRGBA
	}{
		{100, 150, 100, 80, color.NRGBA{A: 0xFF, R: 0xFF, G: 0x80, B: 0x80}},
		{300, 200, 120, 60, color.NRGBA{A: 0xFF, R: 0x80, G: 0xFF, B: 0x80}},
		{500, 100, 80, 120, color.NRGBA{A: 0xFF, R: 0x80, G: 0x80, B: 0xFF}},
		{200, 400, 150, 100, color.NRGBA{A: 0xFF, R: 0xFF, G: 0xFF, B: 0x80}},
		{600, 350, 100, 80, color.NRGBA{A: 0xFF, R: 0xFF, G: 0x80, B: 0xFF}},
	}

	for _, rect := range rects {
		paint.FillShape(gtx.Ops, rect.color, clip.Rect(image.Rectangle{
			Min: image.Point{X: rect.x, Y: rect.y},
			Max: image.Point{X: rect.x + rect.w, Y: rect.y + rect.h},
		}).Op())
	}

	return layout.Dimensions{Size: image.Point{X: contentWidth, Y: contentHeight}}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
