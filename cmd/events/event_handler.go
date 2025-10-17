package main

import (
	"fmt"

	"lol.mleku.dev/log"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

// EventHandler provides consistent event processing for pointer and scroll events
type EventHandler struct {
	// Track button state to show which button was released
	pressedButtons pointer.Buttons
	// Gesture scroll to capture scroll events
	scroll gesture.Scroll
	// Event logging function
	logEvent func(string)
}

// NewEventHandler creates a new event handler with the given event logging function
func NewEventHandler(logEvent func(string)) *EventHandler {
	return &EventHandler{
		logEvent: logEvent,
	}
}

// AddToOps registers the event handler for events in the given ops
func (eh *EventHandler) AddToOps(ops *op.Ops) {
	// Register for pointer events using the proper Gio event system
	event.Op(ops, eh)

	// Add scroll gesture to capture scroll events
	eh.scroll.Add(ops)
}

// ProcessEvents processes all events from the given context
func (eh *EventHandler) ProcessEvents(gtx layout.Context) {
	// Process scroll events first
	eh.processScrollEvents(gtx)

	// Process pointer events from the source
	eh.processPointerEvents(gtx)
}

// processScrollEvents processes scroll events using gesture.Scroll
func (eh *EventHandler) processScrollEvents(gtx layout.Context) {
	// Process scroll events using gesture.Scroll
	scrollDistance := eh.scroll.Update(gtx.Metric, gtx.Source, gtx.Now, gesture.Vertical,
		pointer.ScrollRange{Min: -1000, Max: 1000},
		pointer.ScrollRange{Min: -1000, Max: 1000})

	if scrollDistance != 0 {
		var direction string
		if scrollDistance > 0 {
			direction = "Down"
		} else {
			direction = "Up"
		}

		log.I.F("Found SCROLL gesture: %s, distance=%d", direction, scrollDistance)
		eh.logEvent(fmt.Sprintf("SCROLL: Direction=%s, Distance=%d", direction, scrollDistance))
	}
}

// processPointerEvents processes pointer events from the source
func (eh *EventHandler) processPointerEvents(gtx layout.Context) {
	log.I.F("processPointerEvents called")

	// Process pointer events from the source
	pointerCount := 0
	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: eh,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag | pointer.Move | pointer.Enter | pointer.Leave | pointer.Cancel,
		})
		if !ok {
			// Try to get any event without filter to see what's available
			ev2, ok2 := gtx.Source.Event(pointer.Filter{Target: eh})
			if ok2 {
				log.I.F("Found unfiltered event: %T", ev2)
			}
		}
		if !ok {
			break
		}
		pointerCount++
		if e, ok := ev.(pointer.Event); ok {
			// Track button state and determine which button was released
			var buttonInfo string

			switch e.Kind {
			case pointer.Press:
				eh.pressedButtons |= e.Buttons
				buttonInfo = e.Buttons.String()
			case pointer.Release:
				// For release, show which button was released (the difference between old and new state)
				releasedButton := eh.pressedButtons &^ e.Buttons
				if releasedButton != 0 {
					buttonInfo = fmt.Sprintf("Released: %s", releasedButton.String())
				} else {
					buttonInfo = "Released: Unknown"
				}
				eh.pressedButtons = e.Buttons
			default:
				buttonInfo = e.Buttons.String()
			}

			log.I.F("Found pointer event: %s at position (%.1f,%.1f), buttons=%s", e.Kind, e.Position.X, e.Position.Y, buttonInfo)

			// Create event description
			eventDesc := fmt.Sprintf("POINTER: Kind=%s, Source=%s, Position=(%.1f,%.1f), Scroll=(%.1f,%.1f), Buttons=%s, Modifiers=%v, Time=%v",
				e.Kind, e.Source, e.Position.X, e.Position.Y, e.Scroll.X, e.Scroll.Y, buttonInfo, e.Modifiers, e.Time)

			eh.logEvent(eventDesc)
		} else {
			log.I.F("Found non-pointer event: %T", ev)
		}
	}
	if pointerCount > 0 {
		log.I.F("Processed %d pointer events", pointerCount)
	}
}

// HandleEvent implements event.Handler to capture all events
func (eh *EventHandler) HandleEvent(ev event.Event) {
	log.I.F("HandleEvent called with: %T", ev)
	switch e := ev.(type) {
	case pointer.Event:
		// Track button state and determine which button was released
		var buttonInfo string

		switch e.Kind {
		case pointer.Press:
			eh.pressedButtons |= e.Buttons
			buttonInfo = e.Buttons.String()
		case pointer.Release:
			// For release, show which button was released (the difference between old and new state)
			releasedButton := eh.pressedButtons &^ e.Buttons
			if releasedButton != 0 {
				buttonInfo = fmt.Sprintf("Released: %s", releasedButton.String())
			} else {
				buttonInfo = "Released: Unknown"
			}
			eh.pressedButtons = e.Buttons
		default:
			buttonInfo = e.Buttons.String()
		}

		// Create event description
		eventDesc := fmt.Sprintf("POINTER: Kind=%s, Source=%s, Position=(%.1f,%.1f), Scroll=(%.1f,%.1f), Buttons=%s, Modifiers=%v, Time=%v",
			e.Kind, e.Source, e.Position.X, e.Position.Y, e.Scroll.X, e.Scroll.Y, buttonInfo, e.Modifiers, e.Time)

		eh.logEvent(eventDesc)
	default:
		log.I.F("Unknown event type: %T", ev)
	}
}
