package main

import (
	"flag"
	"log"
	"os"

	"gio.mleku.dev/app"
	page "gio.mleku.dev/examples/component/pages"
	"gio.mleku.dev/examples/component/pages/about"
	"gio.mleku.dev/examples/component/pages/appbar"
	"gio.mleku.dev/examples/component/pages/discloser"
	"gio.mleku.dev/examples/component/pages/menu"
	"gio.mleku.dev/examples/component/pages/navdrawer"
	"gio.mleku.dev/examples/component/pages/textfield"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/text"
	"gio.mleku.dev/widget/material"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	flag.Parse()
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

	router := page.NewRouter()
	router.Register(0, appbar.New(&router))
	router.Register(1, navdrawer.New(&router))
	router.Register(2, textfield.New(&router))
	router.Register(3, menu.New(&router))
	router.Register(4, discloser.New(&router))
	router.Register(5, about.New(&router))

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			router.Layout(gtx, th)
			e.Frame(gtx.Ops)
		}
	}
}
