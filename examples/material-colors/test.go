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

			// Simple test layout
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Test title
					title := material.H4(theme, "Material Design Colors Test")
					title.Color = theme.OnSurface()
					return title.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// Test a few colors
					return layoutTestColors(gtx, theme)
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func layoutTestColors(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	colors := []struct {
		name  string
		color color.NRGBA
	}{
		{"Primary", theme.Primary()},
		{"Secondary", theme.Secondary()},
		{"Error", theme.Error()},
		{"Surface", theme.Surface()},
		{"Background", theme.Background()},
	}

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(theme, "Test Colors")
			title.Color = theme.OnSurface()
			return title.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorRow(gtx, theme, colors)
		}),
	)
}

func layoutColorRow(gtx layout.Context, theme *material.Theme, colors []struct {
	name  string
	color color.NRGBA
}) layout.Dimensions {
	var children []layout.FlexChild

	for _, colorItem := range colors {
		colorItem := colorItem // capture for closure
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, theme, colorItem.name, colorItem.color)
		}))
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx, children...)
}

func layoutColorSquare(gtx layout.Context, theme *material.Theme, name string, color color.NRGBA) layout.Dimensions {
	const size = 100

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorBox(gtx, color, size)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorInfo(gtx, theme, name, color)
		}),
	)
}

func layoutColorBox(gtx layout.Context, color color.NRGBA, size int) layout.Dimensions {
	gtx.Constraints.Min = image.Point{X: size, Y: size}

	return layout.Background{}.Layout(gtx,
		func(gtx layout.Context) layout.Dimensions {
			paint.Fill(gtx.Ops, color)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: gtx.Constraints.Min}
		},
	)
}

func layoutColorInfo(gtx layout.Context, theme *material.Theme, name string, color color.NRGBA) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(theme, name)
			label.Color = theme.OnSurface()
			label.Alignment = text.Middle
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hex := material.Caption(theme, colorToHex(color))
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
