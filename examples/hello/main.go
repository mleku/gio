// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"log"

	"github.com/mleku/gio/app"
	"github.com/mleku/gio/font/gofont"
	"github.com/mleku/gio/layout"
	"github.com/mleku/gio/op"
	"github.com/mleku/gio/text"
	"github.com/mleku/gio/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	th := material.NewTheme()

	// Configure the theme with fonts for WASM
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Regular()))

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			hello(th, gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func hello(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H1(th, "Hello, Gio!")
				title.Color = th.Palette.Fg
				return title.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: 20}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				body := material.Body1(th, "Welcome to Gio running in WebAssembly!")
				body.Color = th.Palette.Fg
				return body.Layout(gtx)
			}),
		)
	})
}
