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
			app.Title("Selectable Primary Clipboard Demo"),
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
	th         *material.Theme
	selectable widget.Selectable
}

func run(w *app.Window) error {
	th := material.NewTheme()
	demo := &Demo{
		th: th,
	}

	// Set some initial text with lorem ipsum
	demo.selectable.SetText("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n\nSed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo. Nemo enim ipsam voluptatem quia voluptas sit aspernatur aut odit aut fugit, sed quia consequuntur magni dolores eos qui ratione voluptatem sequi nesciunt.\n\nAt vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium voluptatum deleniti atque corrupti quos dolores et quas molestias excepturi sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")

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
				title := material.H4(d.th, "Selectable Primary Clipboard Demo")
				return title.Layout(gtx)
			}),

			// Instructions
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				instructions := material.Body1(d.th, "Instructions:\n"+
					"1. Select text in the selectable widget below\n"+
					"2. Selected text will automatically copy to primary clipboard\n"+
					"3. Click middle mouse button to paste from primary clipboard\n"+
					"4. Use Ctrl+C for regular clipboard operations\n\n"+
					"Note: This is a read-only widget, so middle-click paste won't insert text,\n"+
					"but it will demonstrate the primary clipboard functionality.")
				return instructions.Layout(gtx)
			}),

			// Selectable widget
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				label := material.Label(d.th, d.th.TextSize, d.selectable.Text())
				label.State = &d.selectable
				return label.Layout(gtx)
			}),
		)
	})
}
