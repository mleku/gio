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
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Fill background
			paint.Fill(gtx.Ops, th.Background())

			// Simple layout with padding
			inset := layout.Inset{Top: 50, Bottom: 50, Left: 50, Right: 50}
			inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						title := material.H2(th, "Material Design Colors")
						title.Color = th.OnSurface()
						return title.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layoutColorTest(gtx, th)
					}),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}

func layoutColorTest(gtx layout.Context, th *material.Theme) layout.Dimensions {
	colors := []struct {
		name  string
		color color.NRGBA
	}{
		{"Primary", th.Primary()},
		{"Secondary", th.Secondary()},
		{"Error", th.Error()},
		{"Surface", th.Surface()},
	}

	return layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceEnd,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, th, "Primary", th.Primary())
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, th, "Secondary", th.Secondary())
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, th, "Error", th.Error())
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layoutColorSquare(gtx, th, "Surface", th.Surface())
		}),
	)
}

func layoutColorSquare(gtx layout.Context, th *material.Theme, name string, color color.NRGBA) layout.Dimensions {
	const size = 80

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			// Draw color square directly
			gtx.Constraints.Min = image.Point{X: size, Y: size}
			paint.Fill(gtx.Ops, color)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(th, name)
			label.Color = th.OnSurface()
			label.Alignment = text.Middle
			return label.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			hex := material.Caption(th, colorToHex(color))
			hex.Color = th.OnSurfaceVariant()
			hex.Alignment = text.Middle
			return hex.Layout(gtx)
		}),
	)
}

func colorToHex(c color.NRGBA) string {
	return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B)
}
