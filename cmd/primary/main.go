package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(
			app.Title("Primary Clipboard Demo"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type Demo struct {
	th     *material.Theme
	editor widget.Editor
}

func run(w *app.Window) error {
	th := material.NewTheme()
	demo := &Demo{
		th: th,
		editor: widget.Editor{
			SingleLine: false,
		},
	}

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			demo.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (d *Demo) Layout(gtx layout.Context) layout.Dimensions {
	return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			// Header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				title := material.H4(d.th, "Primary Clipboard Demo")
				return title.Layout(gtx)
			}),

			// Instructions
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				instructions := material.Body1(d.th, "Instructions:\n"+
					"1. Type some text in the editor below\n"+
					"2. Select text with mouse - it will automatically copy to primary clipboard\n"+
					"3. Click middle mouse button to paste from primary clipboard\n"+
					"4. Use Ctrl+C/V for regular clipboard operations\n\n"+
					"Note: Primary clipboard only works on X11 systems")
				return instructions.Layout(gtx)
			}),

			// Editor
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				editor := material.Editor(d.th, &d.editor, "Type some text here and then select it with your mouse...")
				editor.Editor.SingleLine = false
				return editor.Layout(gtx)
			}),
		)
	})
}
