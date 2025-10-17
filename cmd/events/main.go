package main

import (
	"fmt"
	"image"
	"os"
	"time"

	"lol.mleku.dev/log"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	log.I.Ln("Starting Event Logger application...")
	go func() {
		log.I.Ln("Creating window...")
		w := &app.Window{}
		w.Option(
			app.Title("Event Logger"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		log.I.Ln("Running window...")
		if err := run(w); err != nil {
			log.F.Ln(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

type App struct {
	th     *material.Theme
	events []string
	list   widget.List
	// Event handler for consistent event processing
	eventHandler *EventHandler
}

func run(w *app.Window) error {
	log.I.Ln("Initializing application...")
	th := material.NewTheme()
	myApp := &App{
		th: th,
		list: widget.List{
			List: layout.List{
				Axis: layout.Vertical,
			},
		},
	}
	// Initialize event handler after myApp is created
	myApp.eventHandler = NewEventHandler(myApp.addEvent)

	var ops op.Ops
	log.I.Ln("Entering event loop...")
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			log.I.Ln("Received DestroyEvent")
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			myApp.Layout(gtx)
			e.Frame(gtx.Ops)

			// Invalidate the window to trigger continuous frame events for event capture
			w.Invalidate()
		}
	}
}

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	log.I.F("Layout called, constraints: %v", gtx.Constraints.Max)

	// Create a scrollable area that covers the entire window
	area := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)

	// Register event handler for events
	a.eventHandler.AddToOps(gtx.Ops)

	area.Pop()

	// Process all events using the event handler
	a.eventHandler.ProcessEvents(gtx)

	// Layout the event log
	return a.layoutEventLog(gtx)
}

func (a *App) addEvent(event string) {
	timestamp := time.Now().Format("15:04:05.000")
	a.events = append(a.events, fmt.Sprintf("[%s] %s", timestamp, event))

	// Keep only last 100 events
	if len(a.events) > 100 {
		a.events = a.events[len(a.events)-100:]
	}

	// Log to console as well
	log.I.F("EVENT: %s", event)
}

func (a *App) layoutEventLog(gtx layout.Context) layout.Dimensions {
	// Background
	paint.Fill(gtx.Ops, a.th.Palette.Bg)

	// Title
	title := material.H5(a.th, "Pointer Event Logger - Focus on mouse events")
	title.Layout(gtx)

	// Instructions
	instructions := material.Body1(a.th, "Instructions:\n"+
		"• Move mouse around (should trigger POINTER Move events)\n"+
		"• Click with mouse (should trigger POINTER Press/Release events)\n"+
		"• Use mouse wheel to scroll (should trigger POINTER Scroll events)\n"+
		"• Watch console output for detailed pointer event logging")
	instructions.Layout(gtx)

	// Event list
	gtx.Constraints.Min = gtx.Constraints.Max
	gtx.Constraints.Max.Y -= 200 // Leave space for title and instructions

	return material.List(a.th, &a.list).Layout(gtx, len(a.events), func(gtx layout.Context, index int) layout.Dimensions {
		if index >= len(a.events) {
			return layout.Dimensions{}
		}

		eventText := material.Body2(a.th, a.events[index])
		eventText.Color = a.th.Palette.Fg
		return eventText.Layout(gtx)
	})
}
