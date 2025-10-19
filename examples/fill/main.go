// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates a simple FillWidget that fills the entire window.

package main

import (
	"image/color"
	"os"

	"gio.mleku.dev/app"
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
	// Create and configure the fill widget using fluent interface
	fillWidget := widget.Fill().
		Color(color.NRGBA{R: 100, G: 150, B: 200, A: 255})

	// Create window widget and set up root rendering
	windowWidget := widget.New(widget.DefaultConfig())
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Update the fill widget size to match the window
		fillWidget.Size(w.Width, w.Height)

		// Render the fill widget
		fillWidget.RenderWidget(gtx)
	}

	// Run the window
	return windowWidget.Run(w)
}
