// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"log"
	"os"

	"image"
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/widget"
	"gio.mleku.dev/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	var theme *material.Theme
	var themeMode material.ThemeMode = material.ThemeModeLight

	// Initialize theme
	theme = material.NewTheme()
	theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Fill background
			paint.Fill(gtx.Ops, theme.Background())

			render(gtx, theme, &themeMode)
			e.Frame(gtx.Ops)
		}
	}
}

func render(gtx layout.Context, theme *material.Theme, themeMode *material.ThemeMode) {
	// Update theme mode
	theme.SetThemeMode(*themeMode)

	// Add padding around the content
	inset := layout.Inset{Top: 20, Bottom: 20, Left: 20, Right: 20}
	inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Create layout with proper constraints
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// Header with theme toggle
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layoutHeader(gtx, theme, themeMode)
			}),
			// Color showcase
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layoutColorShowcase(gtx, theme)
			}),
		)
	})
}

func layoutHeader(gtx layout.Context, theme *material.Theme, themeMode *material.ThemeMode) layout.Dimensions {
	var toggleButton widget.Clickable

	return layout.Flex{
		Axis:      layout.Horizontal,
		Alignment: layout.Middle,
		Spacing:   layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H4(theme, "Material Design Color System")
			title.Color = theme.OnSurface()
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if toggleButton.Clicked(gtx) {
				if *themeMode == material.ThemeModeLight {
					*themeMode = material.ThemeModeDark
				} else {
					*themeMode = material.ThemeModeLight
				}
			}

			var buttonText string
			if *themeMode == material.ThemeModeLight {
				buttonText = "Switch to Dark"
			} else {
				buttonText = "Switch to Light"
			}

			button := material.Button(theme, &toggleButton, buttonText)
			return button.Layout(gtx)
		}),
	)
}

func layoutColorShowcase(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		// Color Roles Section
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorRoles(gtx, theme)
		}),
		// Palette Colors Section
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutPaletteColors(gtx, theme)
		}),
	)
}

func layoutColorRoles(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	title := material.H5(theme, "Color Roles")
	title.Color = theme.OnSurface()

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorGrid(gtx, theme, getColorRoles(theme))
		}),
	)
}

func layoutPaletteColors(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	title := material.H5(theme, "Color Palette")
	title.Color = theme.OnSurface()

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutPaletteGrid(gtx, theme)
		}),
	)
}

type ColorItem struct {
	Name  string
	Color color.NRGBA
}

func getColorRoles(theme *material.Theme) []ColorItem {
	return []ColorItem{
		{"Primary", theme.Primary()},
		{"OnPrimary", theme.OnPrimary()},
		{"PrimaryContainer", theme.PrimaryContainer()},
		{"OnPrimaryContainer", theme.OnPrimaryContainer()},
		{"Secondary", theme.Secondary()},
		{"OnSecondary", theme.OnSecondary()},
		{"SecondaryContainer", theme.SecondaryContainer()},
		{"OnSecondaryContainer", theme.OnSecondaryContainer()},
		{"Tertiary", theme.Tertiary()},
		{"OnTertiary", theme.OnTertiary()},
		{"TertiaryContainer", theme.TertiaryContainer()},
		{"OnTertiaryContainer", theme.OnTertiaryContainer()},
		{"Error", theme.Error()},
		{"OnError", theme.OnError()},
		{"ErrorContainer", theme.ErrorContainer()},
		{"OnErrorContainer", theme.OnErrorContainer()},
		{"Background", theme.Background()},
		{"OnBackground", theme.OnBackground()},
		{"Surface", theme.Surface()},
		{"OnSurface", theme.OnSurface()},
		{"SurfaceVariant", theme.SurfaceVariant()},
		{"OnSurfaceVariant", theme.OnSurfaceVariant()},
		{"Outline", theme.Outline()},
		{"OutlineVariant", theme.OutlineVariant()},
		{"Shadow", theme.Shadow()},
		{"Scrim", theme.Scrim()},
		{"InverseSurface", theme.InverseSurface()},
		{"InverseOnSurface", theme.InverseOnSurface()},
		{"InversePrimary", theme.InversePrimary()},
		{"SurfaceTint", theme.SurfaceTint()},
	}
}

func layoutColorGrid(gtx layout.Context, theme *material.Theme, colors []ColorItem) layout.Dimensions {
	const columns = 3
	const colorSize = 120
	const spacing = 8

	var children []layout.FlexChild

	for i := 0; i < len(colors); i += columns {
		rowColors := colors[i:]
		if len(rowColors) > columns {
			rowColors = rowColors[:columns]
		}

		rowColorsCopy := rowColors // capture for closure
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorRow(gtx, theme, rowColorsCopy, colorSize, spacing)
		}))
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func layoutColorRow(gtx layout.Context, theme *material.Theme, colors []ColorItem, size, spacing int) layout.Dimensions {
	var children []layout.FlexChild

	for _, colorItem := range colors {
		colorItem := colorItem // capture for closure
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, theme, colorItem, size)
		}))
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx, children...)
}

func layoutColorSquare(gtx layout.Context, theme *material.Theme, colorItem ColorItem, size int) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Color square
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorBox(gtx, colorItem.Color, size)
		}),
		// Color name and hex
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorInfo(gtx, theme, colorItem)
		}),
	)
}

func layoutColorBox(gtx layout.Context, color color.NRGBA, size int) layout.Dimensions {
	gtx.Constraints.Min = image.Point{X: size, Y: size}
	paint.Fill(gtx.Ops, color)
	return layout.Dimensions{Size: gtx.Constraints.Min}
}

func layoutColorInfo(gtx layout.Context, theme *material.Theme, colorItem ColorItem) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			name := material.Body2(theme, colorItem.Name)
			name.Color = theme.OnSurface()
			name.Alignment = text.Middle
			return name.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hex := material.Caption(theme, colorToHex(colorItem.Color))
			hex.Color = theme.OnSurfaceVariant()
			hex.Alignment = text.Middle
			return hex.Layout(gtx)
		}),
	)
}

func layoutPaletteGrid(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	palette := theme.Palette()

	var children []layout.FlexChild

	// Primary palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Primary", []ColorItem{
			{"Primary50", palette.Primary50},
			{"Primary100", palette.Primary100},
			{"Primary200", palette.Primary200},
			{"Primary300", palette.Primary300},
			{"Primary400", palette.Primary400},
			{"Primary500", palette.Primary500},
			{"Primary600", palette.Primary600},
			{"Primary700", palette.Primary700},
			{"Primary800", palette.Primary800},
			{"Primary900", palette.Primary900},
			{"Primary950", palette.Primary950},
		})
	}))

	// Secondary palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Secondary", []ColorItem{
			{"Secondary50", palette.Secondary50},
			{"Secondary100", palette.Secondary100},
			{"Secondary200", palette.Secondary200},
			{"Secondary300", palette.Secondary300},
			{"Secondary400", palette.Secondary400},
			{"Secondary500", palette.Secondary500},
			{"Secondary600", palette.Secondary600},
			{"Secondary700", palette.Secondary700},
			{"Secondary800", palette.Secondary800},
			{"Secondary900", palette.Secondary900},
			{"Secondary950", palette.Secondary950},
		})
	}))

	// Tertiary palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Tertiary", []ColorItem{
			{"Tertiary50", palette.Tertiary50},
			{"Tertiary100", palette.Tertiary100},
			{"Tertiary200", palette.Tertiary200},
			{"Tertiary300", palette.Tertiary300},
			{"Tertiary400", palette.Tertiary400},
			{"Tertiary500", palette.Tertiary500},
			{"Tertiary600", palette.Tertiary600},
			{"Tertiary700", palette.Tertiary700},
			{"Tertiary800", palette.Tertiary800},
			{"Tertiary900", palette.Tertiary900},
			{"Tertiary950", palette.Tertiary950},
		})
	}))

	// Error palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Error", []ColorItem{
			{"Error50", palette.Error50},
			{"Error100", palette.Error100},
			{"Error200", palette.Error200},
			{"Error300", palette.Error300},
			{"Error400", palette.Error400},
			{"Error500", palette.Error500},
			{"Error600", palette.Error600},
			{"Error700", palette.Error700},
			{"Error800", palette.Error800},
			{"Error900", palette.Error900},
			{"Error950", palette.Error950},
		})
	}))

	// Neutral palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Neutral", []ColorItem{
			{"Neutral50", palette.Neutral50},
			{"Neutral100", palette.Neutral100},
			{"Neutral200", palette.Neutral200},
			{"Neutral300", palette.Neutral300},
			{"Neutral400", palette.Neutral400},
			{"Neutral500", palette.Neutral500},
			{"Neutral600", palette.Neutral600},
			{"Neutral700", palette.Neutral700},
			{"Neutral800", palette.Neutral800},
			{"Neutral900", palette.Neutral900},
			{"Neutral950", palette.Neutral950},
		})
	}))

	// Neutral Variant palette
	children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layoutPaletteSection(gtx, theme, "Neutral Variant", []ColorItem{
			{"NeutralVariant50", palette.NeutralVariant50},
			{"NeutralVariant100", palette.NeutralVariant100},
			{"NeutralVariant200", palette.NeutralVariant200},
			{"NeutralVariant300", palette.NeutralVariant300},
			{"NeutralVariant400", palette.NeutralVariant400},
			{"NeutralVariant500", palette.NeutralVariant500},
			{"NeutralVariant600", palette.NeutralVariant600},
			{"NeutralVariant700", palette.NeutralVariant700},
			{"NeutralVariant800", palette.NeutralVariant800},
			{"NeutralVariant900", palette.NeutralVariant900},
			{"NeutralVariant950", palette.NeutralVariant950},
		})
	}))

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
}

func layoutPaletteSection(gtx layout.Context, theme *material.Theme, sectionName string, colors []ColorItem) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(theme, sectionName)
			title.Color = theme.OnSurface()
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutPaletteRow(gtx, theme, colors)
		}),
	)
}

func layoutPaletteRow(gtx layout.Context, theme *material.Theme, colors []ColorItem) layout.Dimensions {
	const colorSize = 80
	const spacing = 4

	var children []layout.FlexChild

	for _, colorItem := range colors {
		colorItem := colorItem // capture for closure
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutPaletteColor(gtx, theme, colorItem, colorSize)
		}))
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx, children...)
}

func layoutPaletteColor(gtx layout.Context, theme *material.Theme, colorItem ColorItem, size int) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		// Color square
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorBox(gtx, colorItem.Color, size)
		}),
		// Color name and hex
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutPaletteInfo(gtx, theme, colorItem)
		}),
	)
}

func layoutPaletteInfo(gtx layout.Context, theme *material.Theme, colorItem ColorItem) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			name := material.Caption(theme, colorItem.Name)
			name.Color = theme.OnSurface()
			name.Alignment = text.Middle
			return name.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hex := material.Caption(theme, colorToHex(colorItem.Color))
			hex.Color = theme.OnSurfaceVariant()
			hex.Alignment = text.Middle
			return hex.Layout(gtx)
		}),
	)
}

// colorToHex converts a color.NRGBA to hex string format
func colorToHex(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
