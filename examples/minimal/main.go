// SPDX-License-Identifier: Unlicense OR MIT

// This example demonstrates window-level mouse enter/exit event detection.
// It logs when the mouse enters or leaves the entire window, which is useful
// for enabling/disabling hover effects when the mouse is outside the window.

package main

import (
	"image/color"
	"os"
	"time"

	"gio.mleku.dev/app"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/key"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/io/transfer"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"lol.mleku.dev/log"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.I.F("Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	var hasShownInitialEnter bool
	var frameCount int
	var enterLeaveCount int
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			frameCount++
			gtx := app.NewContext(&ops, e)

			// Paint black background
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			// Subscribe to window-level events
			area := clip.Rect{Max: gtx.Size}.Push(gtx.Ops)
			event.Op(gtx.Ops, w)
			area.Pop()

			// Register event filters for other pointer events
			otherEventsFilter := pointer.Filter{
				Target: w,
				Kinds:  pointer.Cancel | pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Scroll,
			}

			// Focus the window to receive key events (no text input registration)
			gtx.Source.Execute(key.FocusCmd{Tag: w})

			// Track modifier state from pointer events
			var lastModifiers key.Modifiers

			// Force an initial mouse enter event if this is the first frame
			if !hasShownInitialEnter {
				hasShownInitialEnter = true
				// Simulate a mouse enter event at the center of the window
				centerX := float32(gtx.Size.X) / 2
				centerY := float32(gtx.Size.Y) / 2
				log.I.F("MOUSE ENTER WINDOW (initial): Position=(%.1f,%.1f), Source=Mouse @ %s\n",
					centerX, centerY, time.Now().Format("15:04:05.000"))
				log.I.F("===== WINDOW EVENT TEST =====\n")
				log.I.F("Window size: %dx%d pixels\n", gtx.Size.X, gtx.Size.Y)
				log.I.F("Testing window-level Enter/Leave events\n")
				log.I.F("To test Enter/Leave events:\n")
				log.I.F("1. Move mouse COMPLETELY OUTSIDE the window (should see Leave event)\n")
				log.I.F("2. Move mouse BACK INSIDE the window (should see Enter event)\n")
				log.I.F("Moving within the window will only show Move events\n")
				log.I.F("================================================\n")
			}

			// Handle window-level mouse enter/exit events first (most important for hover effects)
			// These events are now handled in the main event loop above

			// Handle other pointer events (moves, clicks, etc.)
			for {
				ev, ok := gtx.Source.Event(otherEventsFilter)
				if !ok {
					break
				}
				if ev, ok := ev.(pointer.Event); ok {
					// Check for modifier changes
					if ev.Modifiers != lastModifiers {
						// Modifier state changed - detect which ones
						if lastModifiers&key.ModShift != 0 && ev.Modifiers&key.ModShift == 0 {
							log.I.F("MODIFIER UP: Shift @ %s\n", time.Now().Format("15:04:05.000"))
						}
						if lastModifiers&key.ModCtrl != 0 && ev.Modifiers&key.ModCtrl == 0 {
							log.I.F("MODIFIER UP: Ctrl @ %s\n", time.Now().Format("15:04:05.000"))
						}
						if lastModifiers&key.ModAlt != 0 && ev.Modifiers&key.ModAlt == 0 {
							log.I.F("MODIFIER UP: Alt @ %s\n", time.Now().Format("15:04:05.000"))
						}
						if lastModifiers&key.ModSuper != 0 && ev.Modifiers&key.ModSuper == 0 {
							log.I.F("MODIFIER UP: Super @ %s\n", time.Now().Format("15:04:05.000"))
						}
						lastModifiers = ev.Modifiers
					}

					// Show window bounds for debugging
					if ev.Kind == pointer.Move && ev.Position.Y < 5 {
						log.I.F("MOUSE NEAR TOP EDGE: Position=(%.1f,%.1f), Window=(%dx%d) @ %s\n",
							ev.Position.X, ev.Position.Y, gtx.Size.X, gtx.Size.Y, time.Now().Format("15:04:05.000"))
					}
					log.I.F("POINTER event: Kind=%s, Buttons=%s, Position=(%.1f,%.1f), Source=%s, Scroll=(%.1f,%.1f), Modifiers=%s\n",
						ev.Kind, ev.Buttons, ev.Position.X, ev.Position.Y, ev.Source, ev.Scroll.X, ev.Scroll.Y, ev.Modifiers)
				}
			}

			// Handle raw key events - now with fixed filtering that allows modifiers

			// Handle focused key events (should now work with modifiers)
			for {
				ev, ok := gtx.Source.Event(key.Filter{Focus: w})
				if !ok {
					break
				}
				if ev, ok := ev.(key.Event); ok {
					timestamp := time.Now().Format("15:04:05.000")
					stateStr := "DOWN"
					if ev.State == key.Release {
						stateStr = "UP"
					}
					log.I.F("KEY %s: %s (KeyCode=%d, Timestamp=%d, Modifiers=%s) @ %s\n",
						stateStr, ev.Name, ev.KeyCode, ev.Timestamp, ev.Modifiers, timestamp)
				}
			}

			// Alternative approach: Handle unfocused key events (catch-all)
			for {
				ev, ok := gtx.Source.Event(key.Filter{})
				if !ok {
					break
				}
				if ev, ok := ev.(key.Event); ok {
					timestamp := time.Now().Format("15:04:05.000")
					stateStr := "DOWN"
					if ev.State == key.Release {
						stateStr = "UP"
					}
					log.I.F("KEY (unfocused) %s: %s (KeyCode=%d, Timestamp=%d, Modifiers=%s) @ %s\n",
						stateStr, ev.Name, ev.KeyCode, ev.Timestamp, ev.Modifiers, timestamp)
				}
			}

			// Handle transfer events (drag & drop and clipboard) using eliasnaur/gio pattern
			transferFilters := []event.Filter{
				transfer.SourceFilter{Target: w, Type: "text/plain"},
				transfer.TargetFilter{Target: w, Type: "text/plain"},
			}
			for {
				ev, ok := gtx.Source.Event(transferFilters...)
				if !ok {
					break
				}
				switch ev := ev.(type) {
				case transfer.InitiateEvent:
					log.I.F("TRANSFER INITIATE event\n")
				case transfer.RequestEvent:
					log.I.F("TRANSFER REQUEST event: Type=%s\n", ev.Type)
				case transfer.CancelEvent:
					log.I.F("TRANSFER CANCEL event\n")
				case transfer.DataEvent:
					log.I.F("TRANSFER DATA event: Type=%s\n", ev.Type)
				}
			}

			e.Frame(gtx.Ops)
		case app.WindowMouseEvent:
			enterLeaveCount++
			switch e.Kind {
			case app.WindowMouseEnter:
				log.I.F("===== MOUSE ENTERED WINDOW ===== Position=(%.1f,%.1f) @ %s\n",
					e.Position.X, e.Position.Y, time.Now().Format("15:04:05.000"))
				log.I.F("✅ Mouse is INSIDE the window - Hover effects ENABLED\n")
			case app.WindowMouseLeave:
				log.I.F("===== MOUSE LEFT WINDOW ===== Position=(%.1f,%.1f) @ %s\n",
					e.Position.X, e.Position.Y, time.Now().Format("15:04:05.000"))
				log.I.F("❌ Mouse is OUTSIDE the window - Hover effects DISABLED\n")
			}
		}
	}
}
