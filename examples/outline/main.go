// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates a border widget with a fill operation before the border operation.

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
	// Create and configure the outline widget with fluent interface
	outlineWidget := widget.NewOutlineWidget().
		Thickness(16).
		CornerRadius(0).
		OutlineColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}).
		Background(func(gtx app.Context, w *widget.Widget) {
			// Fill background with light cyan
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 240, B: 240, A: 255})
		})

	// Create and configure the layout hierarchy using fluent interface
	outerFlex := widget.Flex().
		Direction(widget.FlexColumn).
		Flexed(
			widget.Flex().
				Direction(widget.FlexColumn).
				Flexed(outlineWidget),
		)

	// Create window widget and set up root rendering
	windowWidget := widget.New(widget.DefaultConfig())
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Update the outer flex size to match the window
		outerFlex.SetSize(w.Width, w.Height)

		// Render the outer flex (which will render children)
		outerFlex.RenderWidget(gtx)
	}

	// Run the window
	return windowWidget.Run(w)
}
