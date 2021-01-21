// SPDX-License-Identifier: Unlicense OR MIT

package layout

import (
	"image"
	"testing"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/router"
	"gioui.org/op"
)

func TestListPosition(t *testing.T) {
	_s := func(e ...event.Event) []event.Event { return e }
	r := new(router.Router)
	gtx := Context{
		Ops: new(op.Ops),
		Constraints: Constraints{
			Max: image.Pt(20, 10),
		},
		Queue: r,
	}
	el := func(gtx Context, idx int) Dimensions {
		return Dimensions{Size: image.Pt(10, 10)}
	}
	for _, tc := range []struct {
		label  string
		num    int
		scroll []event.Event
		first  int
		count  int
	}{
		{label: "no item"},
		{label: "1 visible 0 hidden", num: 1, count: 1},
		{label: "2 visible 0 hidden", num: 2, count: 2},
		{label: "2 visible 1 hidden", num: 3, count: 2},
		{label: "3 visible 0 hidden small scroll", num: 3, count: 3,
			scroll: _s(
				pointer.Event{
					Source:   pointer.Mouse,
					Buttons:  pointer.ButtonLeft,
					Type:     pointer.Press,
					Position: f32.Pt(0, 0),
				},
				pointer.Event{
					Source: pointer.Mouse,
					Type:   pointer.Scroll,
					Scroll: f32.Pt(5, 0),
				},
				pointer.Event{
					Source:   pointer.Mouse,
					Buttons:  pointer.ButtonLeft,
					Type:     pointer.Release,
					Position: f32.Pt(5, 0),
				},
			)},
		{label: "2 visible 1 hidden large scroll", num: 3, count: 2, first: 1,
			scroll: _s(
				pointer.Event{
					Source:   pointer.Mouse,
					Buttons:  pointer.ButtonLeft,
					Type:     pointer.Press,
					Position: f32.Pt(0, 0),
				},
				pointer.Event{
					Source: pointer.Mouse,
					Type:   pointer.Scroll,
					Scroll: f32.Pt(10, 0),
				},
				pointer.Event{
					Source:   pointer.Mouse,
					Buttons:  pointer.ButtonLeft,
					Type:     pointer.Release,
					Position: f32.Pt(15, 0),
				},
			)},
	} {
		t.Run(tc.label, func(t *testing.T) {
			gtx.Ops.Reset()

			var list List
			// Initialize the list.
			list.Layout(gtx, tc.num, el)
			// Generate the scroll events.
			r.Frame(gtx.Ops)
			r.Add(tc.scroll...)
			// Let the list process the events.
			list.Layout(gtx, tc.num, el)

			pos := list.Position
			if got, want := pos.First, tc.first; got != want {
				t.Errorf("List: invalid first position: got %v; want %v", got, want)
			}
			if got, want := pos.Count, tc.count; got != want {
				t.Errorf("List: invalid number of visible children: got %v; want %v", got, want)
			}
		})
	}
}