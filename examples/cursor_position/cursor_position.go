// SPDX-License-Identifier: Unlicense OR MIT
package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/f32"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget/material"
)

func main() {
	// Create a new window.
	go func() {
		w := new(app.Window)
		w.Option(app.Size(800, 600))
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	// Initialize the mouse position.
	var mousePos f32.Point
	mousePresent := false
	// Create a material theme.
	th := material.NewTheme()
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			// Register for pointer move events over the entire window.
			r := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
			area := clip.Rect(r).Push(&ops)
			event.Op(&ops, &mousePos)
			area.Pop()
			for {
				ev, ok := gtx.Event(pointer.Filter{
					Target: &mousePos,
					Kinds:  pointer.Move | pointer.Enter | pointer.Leave,
				})
				if !ok {
					break
				}
				switch ev := ev.(type) {
				case pointer.Event:
					switch ev.Kind {
					case pointer.Enter:
						mousePresent = true
					case pointer.Leave:
						mousePresent = false
					}
					mousePos = ev.Position
				}
			}

			// Display the mouse coordinates.
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				coords := "Mouse is outside window"
				if mousePresent {
					coords = fmt.Sprintf("Mouse Position: (%.2f, %.2f)", mousePos.X, mousePos.Y)
				}
				lbl := material.Label(th, unit.Sp(24), coords)
				return lbl.Layout(gtx)
			})

			e.Frame(gtx.Ops)
		}
	}
}
