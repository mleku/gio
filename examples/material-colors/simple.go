// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
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

			// Fill background with theme background
			paint.Fill(gtx.Ops, th.Background())

			// Add a simple title
			l := material.H1(th, "Material Design Colors")
			l.Color = th.OnSurface()
			l.Alignment = text.Middle
			l.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}
