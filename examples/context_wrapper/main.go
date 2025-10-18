package main

import (
	"image"
	"image/color"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
	"gio.mleku.dev/widget/material"
	"lol.mleku.dev/log"
)

func main() {
	log.I.F("Starting context wrapper example application")
	go func() {
		w := new(app.Window)
		log.I.F("Created new window")
		if err := run(w); err != nil {
			log.I.F("Error in run: %v", err)
			os.Exit(1)
		}
		log.I.F("Application exiting")
		os.Exit(0)
	}()
	log.I.F("Starting app.Main()")
	app.Main()
}

func run(w *app.Window) error {
	log.I.F("Starting run function")
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops
	log.I.F("Created theme and operations")

	// Create a context manager
	contextManager := widget.NewContextManager()
	log.I.F("Created context manager")

	// Create a simple button widget
	button := &widget.Clickable{}

	// Wrap the button with context menu functionality
	contextButton := widget.NewContextWrapper(
		func(gtx layout.Context) layout.Dimensions {
			// Layout the button
			return button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				// Draw button background
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x80, B: 0x00, A: 0xFF}) // Green

				// Draw button text
				label := material.Body1(th, "Right-click me!")
				label.Color = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF} // White text
				return label.Layout(gtx)
			})
		},
		func(gtx layout.Context, pos image.Point) layout.Widget {
			// Context menu function
			log.I.F("Context menu requested at position %v", pos)
			return func(gtx layout.Context) layout.Dimensions {
				// Create a simple context menu
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Menu background
						paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})
						return layout.Dimensions{Size: image.Point{X: 150, Y: 100}}
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Menu items
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								item := material.Body2(th, "Copy")
								item.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
								return item.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								item := material.Body2(th, "Paste")
								item.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
								return item.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								item := material.Body2(th, "Close")
								item.Color = color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
								return item.Layout(gtx)
							}),
						)
					}),
				)
			}
		},
		20, // Priority
	)
	log.I.F("Created context button")

	// Register the context button with the context manager
	contextManager.RegisterWidget(contextButton, 20)
	log.I.F("Registered context button")

	log.I.F("Starting main event loop")
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			log.I.F("Received destroy event: %v", e.Err)
			return e.Err
		case app.FrameEvent:
			log.I.F("Received frame event")
			gtx := app.NewContext(&ops, e)
			log.I.F("Created context with constraints: %v", gtx.Constraints.Max)

			// Register for pointer events over the entire window
			log.I.F("Setting up pointer event registration")
			r := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
			clip.Rect(r).Push(gtx.Ops)
			log.I.F("Pointer event rectangle: %v", r)

			// Add the context manager handler FIRST, before any widgets
			log.I.F("Adding context manager handler")
			contextManager.AddContextHandler(gtx.Ops)

			// Layout the context button
			log.I.F("Laying out context button")
			dims := layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				log.I.F("Layout.Center.Layout called")
				return layout.Inset{
					Top:    unit.Dp(20),
					Bottom: unit.Dp(20),
					Left:   unit.Dp(20),
					Right:  unit.Dp(20),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					log.I.F("Calling contextButton.Layout")
					return contextButton.Layout(gtx)
				})
			})
			log.I.F("Context button layout dimensions: %v", dims.Size)

			// Update widget bounds for hit detection
			centerX := gtx.Constraints.Max.X / 2
			centerY := gtx.Constraints.Max.Y / 2
			buttonX := centerX - 75 // Half of button width
			buttonY := centerY - 25 // Half of button height
			log.I.F("Context button calculated position: (%d,%d)", buttonX, buttonY)

			bounds := image.Rectangle{
				Min: image.Point{X: buttonX, Y: buttonY},
				Max: image.Point{X: buttonX + 150, Y: buttonY + 50},
			}
			log.I.F("Context button bounds: %v", bounds)
			contextManager.UpdateWidgetBounds(contextButton, bounds)

			// Update and layout the context manager AFTER widgets are laid out
			log.I.F("Updating context manager")
			log.I.F("About to call contextManager.Update")
			contextManager.Update(gtx)
			log.I.F("Finished calling contextManager.Update")
			log.I.F("Laying out context manager")
			contextManager.Layout(gtx)

			log.I.F("Calling e.Frame")
			e.Frame(gtx.Ops)
			log.I.F("Frame completed")
		}
	}
}
