// SPDX-License-Identifier: Unlicense OR MIT

package widget_test

import (
	"image"
	"testing"

	"github.com/mleku/gio/f32"
	"github.com/mleku/gio/io/input"
	"github.com/mleku/gio/io/pointer"
	"github.com/mleku/gio/io/semantic"
	"github.com/mleku/gio/layout"
	"github.com/mleku/gio/op"
	"github.com/mleku/gio/widget"
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
