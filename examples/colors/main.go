// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates the color scheme from fromage integrated into the gio widget package.

package main

import (
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

	// Create a color scheme
	colors := widget.NewColorsWithMode(widget.ThemeModeLight)

	// Set the root widget to render with the color scheme
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Fill background with the theme's background color
		paint.Fill(gtx.Ops, colors.Background())
	}

	// Run the window
	return windowWidget.Run(w)
}
