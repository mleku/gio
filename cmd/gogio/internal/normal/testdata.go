// SPDX-License-Identifier: Unlicense OR MIT

// A simple app used for gogio's end-to-end tests.
package main

import (
	"fmt"
	"image/color"
	"log"

	"gio.mleku.dev/app"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Paint black background
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			e.Frame(gtx.Ops)
		case pointer.Event:
			// Log mouse events
			fmt.Printf("lol.mleku.dev: Mouse event: %+v\n", e)
		}
	}
}
