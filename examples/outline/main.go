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
	// Create window widget
	windowWidget := widget.New(widget.DefaultConfig())

	// Create an outline widget
	outlineWidget := widget.NewOutlineWidget().
		SetThickness(16).
		SetCornerRadius(0).
		SetOutlineColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	// Set custom background rendering for the outline widget
	outlineWidget.Render = func(gtx app.Context, w *widget.Widget) {
		// Fill background with light cyan
		paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 240, B: 240, A: 255})
	}

	// Create inner flex container
	innerFlex := widget.NewFlexWidget().
		SetDirection(widget.FlexColumn)

	// Add the outline widget to the inner flex
	innerFlex.Flexed(outlineWidget)

	// Create outer flex container
	outerFlex := widget.NewFlexWidget().
		SetDirection(widget.FlexColumn)

	// Add the inner flex to the outer flex as a flexed item
	outerFlex.Flexed(innerFlex)

	// Set the root widget to render the outer flex
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Update the outer flex size to match the window
		outerFlex.SetSize(w.Width, w.Height)

		// Render the outer flex (which will render children)
		outerFlex.RenderWidget(gtx)
	}

	// Run the window
	return windowWidget.Run(w)
}
