// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates a border widget with a fill operation before the border operation.

package main

import (
	"image"
	"image/color"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/f32"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/widget"
	"lol.mleku.dev/log"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.I.F("Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	// Create window widget
	windowWidget := widget.New(widget.DefaultConfig())

	// Set the root widget to render as a border widget with fill background
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// First, fill the entire area with a light cyan background
		paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 240, B: 240, A: 255}) // Light cyan background

		// Then, draw an 8Dp black border using Path operations to create a "hole"
		thickness := float32(8) // 8Dp wide border

		// Create outer rectangle for the border (inset from the window edges)
		outerRect := image.Rect(0, 0, w.Width-0, w.Height-0)

		// Create inner rectangle (spaced by border width)
		innerRect := image.Rect(
			outerRect.Min.X+int(thickness),
			outerRect.Min.Y+int(thickness),
			outerRect.Max.X-int(thickness),
			outerRect.Max.Y-int(thickness),
		)

		// Create a path that defines both outer and inner rectangles
		// The non-zero winding rule will create the border effect
		var path clip.Path
		path.Begin(gtx.Ops)

		// Draw outer rectangle (clockwise)
		path.MoveTo(f32.Pt(float32(outerRect.Min.X), float32(outerRect.Min.Y)))
		path.LineTo(f32.Pt(float32(outerRect.Max.X), float32(outerRect.Min.Y)))
		path.LineTo(f32.Pt(float32(outerRect.Max.X), float32(outerRect.Max.Y)))
		path.LineTo(f32.Pt(float32(outerRect.Min.X), float32(outerRect.Max.Y)))
		path.Close()

		// Draw inner rectangle (counter-clockwise to create hole)
		path.MoveTo(f32.Pt(float32(innerRect.Min.X), float32(innerRect.Min.Y)))
		path.LineTo(f32.Pt(float32(innerRect.Min.X), float32(innerRect.Max.Y)))
		path.LineTo(f32.Pt(float32(innerRect.Max.X), float32(innerRect.Max.Y)))
		path.LineTo(f32.Pt(float32(innerRect.Max.X), float32(innerRect.Min.Y)))
		path.Close()

		// Apply the path as a clip and fill
		defer clip.Outline{Path: path.End()}.Op().Push(gtx.Ops).Pop()
		paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255}) // Black border
	}

	// Run the window
	return windowWidget.Run(w)
}
