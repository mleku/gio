// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"log"

	"github.com/mleku/gio/app"
	"github.com/mleku/gio/layout"
	"github.com/mleku/gio/op"
	"github.com/mleku/gio/widget/material"
)

func main() {
	fmt.Println("Starting Gio WASM app...")
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

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			simple(th, gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func simple(th *material.Theme, gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.H3(th, "Simple WASM Test")
		label.Color = th.Palette.Fg
		return label.Layout(gtx)
	})
}
