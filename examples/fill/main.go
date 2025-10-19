// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates a simple FillWidget that fills the entire window.

package main

import (
	"image/color"
	"os"

	"gio.mleku.dev/app"
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

	// Set the root widget to render as a fill widget
	windowWidget.Root().Render = func(gtx app.Context, widget *widget.Widget) {
		// Fill the entire root widget with blue color
		paint.Fill(gtx.Ops, color.NRGBA{R: 100, G: 150, B: 200, A: 255})
	}

	// Run the window
	return windowWidget.Run(w)
}
