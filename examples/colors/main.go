// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates the complete Material Design color system with labeled color swatches.

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
	// Create a color scheme
	colors := widget.NewColorsWithMode(widget.ThemeModeLight)
	palette := colors.Palette()

	// Create color swatches for each category using Fill widgets
	primarySwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.Primary50),
		widget.Fill().Color(palette.Primary100),
		widget.Fill().Color(palette.Primary200),
		widget.Fill().Color(palette.Primary300),
		widget.Fill().Color(palette.Primary400),
		widget.Fill().Color(palette.Primary500),
		widget.Fill().Color(palette.Primary600),
		widget.Fill().Color(palette.Primary700),
		widget.Fill().Color(palette.Primary800),
		widget.Fill().Color(palette.Primary900),
		widget.Fill().Color(palette.Primary950),
	}

	secondarySwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.Secondary50),
		widget.Fill().Color(palette.Secondary100),
		widget.Fill().Color(palette.Secondary200),
		widget.Fill().Color(palette.Secondary300),
		widget.Fill().Color(palette.Secondary400),
		widget.Fill().Color(palette.Secondary500),
		widget.Fill().Color(palette.Secondary600),
		widget.Fill().Color(palette.Secondary700),
		widget.Fill().Color(palette.Secondary800),
		widget.Fill().Color(palette.Secondary900),
		widget.Fill().Color(palette.Secondary950),
	}

	tertiarySwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.Tertiary50),
		widget.Fill().Color(palette.Tertiary100),
		widget.Fill().Color(palette.Tertiary200),
		widget.Fill().Color(palette.Tertiary300),
		widget.Fill().Color(palette.Tertiary400),
		widget.Fill().Color(palette.Tertiary500),
		widget.Fill().Color(palette.Tertiary600),
		widget.Fill().Color(palette.Tertiary700),
		widget.Fill().Color(palette.Tertiary800),
		widget.Fill().Color(palette.Tertiary900),
		widget.Fill().Color(palette.Tertiary950),
	}

	errorSwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.Error50),
		widget.Fill().Color(palette.Error100),
		widget.Fill().Color(palette.Error200),
		widget.Fill().Color(palette.Error300),
		widget.Fill().Color(palette.Error400),
		widget.Fill().Color(palette.Error500),
		widget.Fill().Color(palette.Error600),
		widget.Fill().Color(palette.Error700),
		widget.Fill().Color(palette.Error800),
		widget.Fill().Color(palette.Error900),
		widget.Fill().Color(palette.Error950),
	}

	neutralSwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.Neutral50),
		widget.Fill().Color(palette.Neutral100),
		widget.Fill().Color(palette.Neutral200),
		widget.Fill().Color(palette.Neutral300),
		widget.Fill().Color(palette.Neutral400),
		widget.Fill().Color(palette.Neutral500),
		widget.Fill().Color(palette.Neutral600),
		widget.Fill().Color(palette.Neutral700),
		widget.Fill().Color(palette.Neutral800),
		widget.Fill().Color(palette.Neutral900),
		widget.Fill().Color(palette.Neutral950),
	}

	neutralVariantSwatches := []*widget.FillWidget{
		widget.Fill().Color(palette.NeutralVariant50),
		widget.Fill().Color(palette.NeutralVariant100),
		widget.Fill().Color(palette.NeutralVariant200),
		widget.Fill().Color(palette.NeutralVariant300),
		widget.Fill().Color(palette.NeutralVariant400),
		widget.Fill().Color(palette.NeutralVariant500),
		widget.Fill().Color(palette.NeutralVariant600),
		widget.Fill().Color(palette.NeutralVariant700),
		widget.Fill().Color(palette.NeutralVariant800),
		widget.Fill().Color(palette.NeutralVariant900),
		widget.Fill().Color(palette.NeutralVariant950),
	}

	// Create row containers for each color category using Flex widgets with Flexed items
	primaryRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range primarySwatches {
		primaryRow.Flexed(swatch)
	}

	secondaryRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range secondarySwatches {
		secondaryRow.Flexed(swatch)
	}

	tertiaryRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range tertiarySwatches {
		tertiaryRow.Flexed(swatch)
	}

	errorRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range errorSwatches {
		errorRow.Flexed(swatch)
	}

	neutralRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range neutralSwatches {
		neutralRow.Flexed(swatch)
	}

	neutralVariantRow := widget.Flex().Direction(widget.FlexRow)
	for _, swatch := range neutralVariantSwatches {
		neutralVariantRow.Flexed(swatch)
	}

	// Create main column container
	mainColumn := widget.Flex().
		Direction(widget.FlexColumn).
		Flexed(primaryRow).
		Flexed(secondaryRow).
		Flexed(tertiaryRow).
		Flexed(errorRow).
		Flexed(neutralRow).
		Flexed(neutralVariantRow)

	// Create window widget and set up root rendering
	windowWidget := widget.New(widget.DefaultConfig())
	windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
		// Fill background with the theme's background color
		paint.Fill(gtx.Ops, colors.Background())

		// Update the main column size to match the window
		mainColumn.SetSize(w.Width, w.Height)

		// Render the main column (which will render all rows and swatches)
		mainColumn.RenderWidget(gtx)
	}

	// Run the window
	return windowWidget.Run(w)
}
