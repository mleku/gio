// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates a simple OutlineWidget that draws a 2px square corner outline box.

package main

import (
	"image"
	"image/color"
	"os"

	"gio.mleku.dev/app"
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

	// Set the root widget to render as an outline widget
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Draw a 2px black outline around the entire window using Gio's stroke operations
		thickness := float32(2)

		// Create rectangle for the outline
		r := image.Rect(0, 0, w.Width, w.Height)

		// Draw the outline using Gio's stroke operation
		paint.FillShape(gtx.Ops,
			color.NRGBA{R: 0, G: 0, B: 0, A: 255}, // Black outline
			clip.Stroke{
				Path:  clip.RRect{Rect: r, NW: 0, NE: 0, SW: 0, SE: 0}.Path(gtx.Ops),
				Width: thickness,
			}.Op(),
		)
	}

	// Run the window
	return windowWidget.Run(w)
}
