// SPDX-License-Identifier: Unlicense OR MIT

package widget_test

import (
	"image"
	"testing"

	"gio.realy.lol/f32"
	"gio.realy.lol/io/input"
	"gio.realy.lol/io/pointer"
	"gio.realy.lol/io/semantic"
	"gio.realy.lol/layout"
	"gio.realy.lol/op"
	"gio.realy.lol/widget"
)

func TestBool(t *testing.T) {
	var (
		r input.Router
		b widget.Bool
	)
	gtx := layout.Context{
		Ops:    new(op.Ops),
		Source: r.Source(),
	}
	layout := func() {
		b.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			semantic.CheckBox.Add(gtx.Ops)
			semantic.DescriptionOp("description").Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(100, 100)}
		})
	}
	layout()
	r.Frame(gtx.Ops)
	r.Queue(
		pointer.Event{
			Source:   pointer.Touch,
			Kind:     pointer.Press,
			Position: f32.Pt(50, 50),
		},
		pointer.Event{
			Source:   pointer.Touch,
			Kind:     pointer.Release,
			Position: f32.Pt(50, 50),
		},
	)
	gtx.Reset()
	layout()
	r.Frame(gtx.Ops)
	tree := r.AppendSemantics(nil)
	n := tree[0].Children[0].Desc
	if n.Description != "description" {
		t.Errorf("unexpected semantic description: %s", n.Description)
	}
	if n.Class != semantic.CheckBox {
		t.Errorf("unexpected semantic class: %v", n.Class)
	}
	if !b.Value || !n.Selected {
		t.Error("click did not select")
	}
}
