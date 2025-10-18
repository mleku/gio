package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/widget/material"
	"gio.mleku.dev/x/colorpicker"
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

type (
	C = layout.Context
	D = layout.Dimensions
)

var white = color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

func loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	background := white
	current := color.NRGBA{R: 255, G: 128, B: 75, A: 255}
	picker := colorpicker.State{}
	picker.SetColor(current)
	muxState := colorpicker.NewMuxState(
		[]colorpicker.MuxOption{
			{
				Label: "current",
				Value: &current,
			},
			{
				Label: "background",
				Value: &th.Palette.Bg,
			},
			{
				Label: "foreground",
				Value: &th.Palette.Fg,
			},
		}...)
	background = *muxState.Color()
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if muxState.Update(gtx) {
				background = *muxState.Color()
			}
			if picker.Update(gtx) {
				current = picker.Color()
				background = *muxState.Color()
			}
			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return colorpicker.PickerStyle{
						Label:         "Current",
						Theme:         th,
						State:         &picker,
						MonospaceFace: "Go Mono",
					}.Layout(gtx)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return colorpicker.Mux(th, &muxState, "Display Right:").Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							size := gtx.Constraints.Max
							paint.FillShape(gtx.Ops, background, clip.Rect(image.Rectangle{Max: size}).Op())
							return D{Size: size}
						}),
					)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
