// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
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
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Paint black background
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			// Subscribe to pointer events for the entire window
			area := clip.Rect{Max: gtx.Size}.Push(gtx.Ops)
			event.Op(gtx.Ops, w)
			area.Pop()

			// Handle pointer events using the source directly
			for {
				ev, ok := gtx.Source.Event(pointer.Filter{
					Target: w,
					Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Enter | pointer.Leave,
				})
				if !ok {
					break
				}
				if ev, ok := ev.(pointer.Event); ok {
					fmt.Printf("lol.mleku.dev: Mouse event: Kind=%s, Buttons=%s, Position=(%.1f,%.1f), Source=%s\n",
						ev.Kind, ev.Buttons, ev.Position.X, ev.Position.Y, ev.Source)
				}
			}

			e.Frame(gtx.Ops)
		}
	}
}
