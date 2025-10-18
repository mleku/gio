package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
	"gio.mleku.dev/widget/material"
)

func main() {
	go func() {
		w := new(app.Window)
		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	// Create a context manager
	contextManager := widget.NewContextManager()

	// Create two example widgets with context menus
	button1 := &ExampleButton{
		Clickable: widget.Clickable{},
		Label:     "red",
		Color:     color.NRGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // Pure red
	}

	button2 := &ExampleButton{
		Clickable: widget.Clickable{},
		Label:     "blue",
		Color:     color.NRGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // Pure blue
	}

	// Create main background widget
	mainBackground := &MainBackgroundWidget{}

	// Register widgets with the context manager (higher priority = higher number)
	contextManager.RegisterWidget(button1, 20)       // Higher priority for red button
	contextManager.RegisterWidget(button2, 10)       // Lower priority for blue button
	contextManager.RegisterWidget(mainBackground, 1) // Lowest priority for background

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Register for pointer events over the entire window
			r := image.Rectangle{Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
			clip.Rect(r).Push(gtx.Ops)

			// Add the context manager handler FIRST, before any widgets
			contextManager.AddContextHandler(gtx.Ops)

			// Layout the main UI - two buttons side by side
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Horizontal,
					Spacing: layout.SpaceEvenly,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Give each button a fixed width with margin
						buttonGtx := gtx
						buttonGtx.Constraints = layout.Constraints{
							Min: image.Point{X: 150, Y: 80},
							Max: image.Point{X: 150, Y: 80},
						}
						dims := layout.Inset{
							Top:    unit.Dp(20),
							Bottom: unit.Dp(20),
							Left:   unit.Dp(20),
							Right:  unit.Dp(20),
						}.Layout(buttonGtx, func(gtx layout.Context) layout.Dimensions {
							return button1.Layout(gtx, th)
						})

						// Update widget bounds for hit detection
						// Calculate the actual position of the button
						centerX := gtx.Constraints.Max.X / 2
						centerY := gtx.Constraints.Max.Y / 2
						buttonX := centerX - 150 - 20 // Left side of screen, accounting for margins
						buttonY := centerY - 40 - 20  // Center vertically, accounting for margins

						bounds := image.Rectangle{
							Min: image.Point{X: buttonX, Y: buttonY},
							Max: image.Point{X: buttonX + 150, Y: buttonY + 80},
						}
						contextManager.UpdateWidgetBounds(button1, bounds)

						return dims
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						// Give each button a fixed width with margin
						buttonGtx := gtx
						buttonGtx.Constraints = layout.Constraints{
							Min: image.Point{X: 150, Y: 80},
							Max: image.Point{X: 150, Y: 80},
						}
						dims := layout.Inset{
							Top:    unit.Dp(20),
							Bottom: unit.Dp(20),
							Left:   unit.Dp(20),
							Right:  unit.Dp(20),
						}.Layout(buttonGtx, func(gtx layout.Context) layout.Dimensions {
							return button2.Layout(gtx, th)
						})

						// Update widget bounds for hit detection
						// Calculate the actual position of the button
						centerX := gtx.Constraints.Max.X / 2
						centerY := gtx.Constraints.Max.Y / 2
						buttonX := centerX + 20      // Right side of screen, accounting for margins
						buttonY := centerY - 40 - 20 // Center vertically, accounting for margins

						bounds := image.Rectangle{
							Min: image.Point{X: buttonX, Y: buttonY},
							Max: image.Point{X: buttonX + 150, Y: buttonY + 80},
						}
						contextManager.UpdateWidgetBounds(button2, bounds)

						return dims
					}),
				)
			})

			// Layout the main background widget (invisible, covers entire area)
			mainBackground.Layout(gtx, th)

			// Update main background bounds (covers entire viewport)
			mainBounds := image.Rectangle{
				Min: image.Point{},
				Max: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y},
			}
			contextManager.UpdateWidgetBounds(mainBackground, mainBounds)

			// Update and layout the context manager AFTER widgets are laid out
			contextManager.Update(gtx)
			contextManager.Layout(gtx)

			e.Frame(gtx.Ops)
		}
	}
}

// ExampleButton is a button widget that implements ContextWidget
type ExampleButton struct {
	widget.Clickable
	Label string
	Color color.NRGBA
}

// ContextMenu implements the ContextWidget interface
func (b *ExampleButton) ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget {
	th := material.NewTheme()
	return func(gtx layout.Context) layout.Dimensions {
		// Create context menu items - button label and close button
		menuItems := []string{b.Label, "Close"}
		clickables := make([]*widget.Clickable, len(menuItems))
		for i := range clickables {
			clickables[i] = &widget.Clickable{}
		}

		// Create close button
		closeButton := &widget.Clickable{}

		// Calculate menu size based on items
		itemHeight := 25
		padding := 10
		menuWidth := 150
		menuHeight := len(menuItems)*itemHeight + padding*2

		// Handle close button events BEFORE clipping
		closeButtonSize := 20 // Square button, 2 text heights
		closeButtonOffset := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Add clickable area for close button first
		closeButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: closeButtonSize, Y: closeButtonSize}}
		})

		// Check for close button click
		if closeButton.Clicked(gtx) {
			log.Printf("Close button clicked for Button %s", b.Label)
			// Note: In a real implementation, you'd dismiss the context menu here
		}

		closeButtonOffset.Pop()

		// Clip to the exact menu size (including close button area)
		clipStack := clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Push(gtx.Ops)

		// Draw the menu frame background
		paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})

		// Draw the menu frame border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button visuals AFTER clipping
		closeButtonOffset2 := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Draw close button background
		if closeButton.Hovered() {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF})
		} else {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xF8, G: 0xF8, B: 0xF8, A: 0xFF})
		}

		// Draw close button border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: closeButtonSize, Y: closeButtonSize}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button "X"
		closeLabel := material.Body2(th, "×")
		closeLabel.Color = color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}
		layout.Center.Layout(gtx, closeLabel.Layout)

		closeButtonOffset2.Pop()

		// Draw menu items
		itemOffset := op.Offset(image.Point{X: padding, Y: padding}).Push(gtx.Ops)
		for i, item := range menuItems {
			itemClickable := clickables[i]

			// Add clickable area first - this enables hover detection
			itemClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Point{X: menuWidth - padding*2, Y: itemHeight}}
			})

			// Check for clicks on this item
			if itemClickable.Clicked(gtx) {
				log.Printf("Button %s clicked: %s", b.Label, item)
			}

			// Draw item text first
			label := material.Body2(th, item)
			label.Color = color.NRGBA{A: 0xFF}
			layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, label.Layout)
			})

			// Draw hover effect overlay after text (10% opacity dark overlay)
			if itemClickable.Hovered() {
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1A}) // 10% opacity black
			}

			// Move to next item
			op.Offset(image.Point{Y: itemHeight}).Add(gtx.Ops)
		}
		itemOffset.Pop()

		// Pop the clip
		clipStack.Pop()

		return layout.Dimensions{Size: image.Point{X: menuWidth, Y: menuHeight}}
	}
}

// Layout renders the button
func (b *ExampleButton) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Draw the button background with proper clipping
	clipStack := clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops)
	paint.Fill(gtx.Ops, b.Color)
	clipStack.Pop()

	// Draw the button label
	label := material.Body1(th, b.Label)
	label.Color = color.NRGBA{A: 0xFF}

	// Center the label and get dimensions
	return layout.Center.Layout(gtx, label.Layout)
}

// handleRightClick handles right-click events for this specific button
func (b *ExampleButton) handleRightClick(gtx layout.Context, contextManager *widget.ContextManager) {
	log.Printf("Button %s: Checking for right-click events", b.Label)

	// Process right-click events for this button
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: b,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		log.Printf("Button %s: Received pointer event, buttons=%v", b.Label, e.Buttons)

		// Only handle right-clicks
		if e.Buttons != pointer.ButtonSecondary {
			continue
		}

		log.Printf("Button %s: Right-click detected, showing context menu", b.Label)

		// Show context menu for this specific button
		contextWidget := b.ContextMenu(gtx, e.Position.Round())
		if contextWidget != nil {
			contextManager.ShowContextWidget(gtx, contextWidget, e.Position.Round())
		}
		break // Only handle one event per frame
	}
}

// ExampleContainer is a container widget that also implements ContextWidget
type ExampleContainer struct {
	Color     color.NRGBA
	Buttons   []*ExampleButton
	Clickable widget.Clickable
}

// ContextMenu implements the ContextWidget interface
func (c *ExampleContainer) ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget {
	th := material.NewTheme()
	return func(gtx layout.Context) layout.Dimensions {
		// Container context menu - just one close button
		menuItems := []string{"Close"}
		clickables := make([]*widget.Clickable, len(menuItems))
		for i := range clickables {
			clickables[i] = &widget.Clickable{}
		}

		// Create close button
		closeButton := &widget.Clickable{}

		// Calculate menu size based on items
		itemHeight := 25
		padding := 10
		menuWidth := 150
		menuHeight := len(menuItems)*itemHeight + padding*2

		// Handle close button events BEFORE clipping
		closeButtonSize := 20 // Square button, 2 text heights
		closeButtonOffset := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Add clickable area for close button first
		closeButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: closeButtonSize, Y: closeButtonSize}}
		})

		// Check for close button click
		if closeButton.Clicked(gtx) {
			log.Printf("Close button clicked for Container")
			// Note: In a real implementation, you'd dismiss the context menu here
		}

		closeButtonOffset.Pop()

		// Clip to the exact menu size (including close button area)
		clipStack := clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Push(gtx.Ops)

		// Draw the menu frame background
		paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})

		// Draw the menu frame border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button visuals AFTER clipping
		closeButtonOffset2 := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Draw close button background
		if closeButton.Hovered() {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF})
		} else {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xF8, G: 0xF8, B: 0xF8, A: 0xFF})
		}

		// Draw close button border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: closeButtonSize, Y: closeButtonSize}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button "X"
		closeLabel := material.Body2(th, "×")
		closeLabel.Color = color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}
		layout.Center.Layout(gtx, closeLabel.Layout)

		closeButtonOffset2.Pop()

		// Draw menu items
		itemOffset := op.Offset(image.Point{X: padding, Y: padding}).Push(gtx.Ops)
		for i, item := range menuItems {
			itemClickable := clickables[i]

			// Add clickable area first - this enables hover detection
			itemClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Point{X: menuWidth - padding*2, Y: itemHeight}}
			})

			// Check for clicks on this item
			if itemClickable.Clicked(gtx) {
				log.Printf("Container clicked: %s", item)
			}

			// Draw item text first
			label := material.Body2(th, item)
			label.Color = color.NRGBA{A: 0xFF}
			layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, label.Layout)
			})

			// Draw hover effect overlay after text (10% opacity dark overlay)
			if itemClickable.Hovered() {
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1A}) // 10% opacity black
			}

			// Move to next item
			op.Offset(image.Point{Y: itemHeight}).Add(gtx.Ops)
		}
		itemOffset.Pop()

		// Pop the clip
		clipStack.Pop()

		return layout.Dimensions{Size: image.Point{X: menuWidth, Y: menuHeight}}
	}
}

// Layout renders the container
func (c *ExampleContainer) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Draw container background
	paint.Fill(gtx.Ops, c.Color)

	// Layout buttons in a flex
	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceBetween,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.Buttons[0].Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.Buttons[1].Layout(gtx, th)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return c.Buttons[2].Layout(gtx, th)
		}),
	)
}

// MainBackgroundWidget is a widget that provides context menu for the main background
type MainBackgroundWidget struct {
	widget.Clickable
}

// ContextMenu implements the ContextWidget interface
func (m *MainBackgroundWidget) ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget {
	th := material.NewTheme()
	return func(gtx layout.Context) layout.Dimensions {
		// Create context menu items - main label and close button
		menuItems := []string{"main", "Close"}
		clickables := make([]*widget.Clickable, len(menuItems))
		for i := range clickables {
			clickables[i] = &widget.Clickable{}
		}

		// Create close button
		closeButton := &widget.Clickable{}

		// Calculate menu size based on items
		itemHeight := 25
		padding := 10
		menuWidth := 150
		menuHeight := len(menuItems)*itemHeight + padding*2

		// Handle close button events BEFORE clipping
		closeButtonSize := 20 // Square button, 2 text heights
		closeButtonOffset := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Add clickable area for close button first
		closeButton.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Point{X: closeButtonSize, Y: closeButtonSize}}
		})

		// Check for close button click
		if closeButton.Clicked(gtx) {
			log.Printf("Close button clicked for main background")
			// Note: In a real implementation, you'd dismiss the context menu here
		}

		closeButtonOffset.Pop()

		// Clip to the exact menu size (including close button area)
		clipStack := clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Push(gtx.Ops)

		// Draw the menu frame background
		paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})

		// Draw the menu frame border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: menuWidth, Y: menuHeight}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button visuals AFTER clipping
		closeButtonOffset2 := op.Offset(image.Point{X: menuWidth - closeButtonSize - 2, Y: 2}).Push(gtx.Ops)

		// Draw close button background
		if closeButton.Hovered() {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF})
		} else {
			paint.Fill(gtx.Ops, color.NRGBA{R: 0xF8, G: 0xF8, B: 0xF8, A: 0xFF})
		}

		// Draw close button border
		paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
			clip.Stroke{
				Path:  clip.Rect(image.Rectangle{Max: image.Point{X: closeButtonSize, Y: closeButtonSize}}).Path(),
				Width: 1,
			}.Op())

		// Draw close button "X"
		closeLabel := material.Body2(th, "×")
		closeLabel.Color = color.NRGBA{R: 0x60, G: 0x60, B: 0x60, A: 0xFF}
		layout.Center.Layout(gtx, closeLabel.Layout)

		closeButtonOffset2.Pop()

		// Draw menu items
		itemOffset := op.Offset(image.Point{X: padding, Y: padding}).Push(gtx.Ops)
		for i, item := range menuItems {
			itemClickable := clickables[i]

			// Add clickable area first - this enables hover detection
			itemClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Dimensions{Size: image.Point{X: menuWidth - padding*2, Y: itemHeight}}
			})

			// Check for clicks on this item
			if itemClickable.Clicked(gtx) {
				log.Printf("Main background clicked: %s", item)
			}

			// Draw item text first
			label := material.Body2(th, item)
			label.Color = color.NRGBA{A: 0xFF}
			layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, label.Layout)
			})

			// Draw hover effect overlay after text (10% opacity dark overlay)
			if itemClickable.Hovered() {
				paint.Fill(gtx.Ops, color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x1A}) // 10% opacity black
			}

			// Move to next item
			op.Offset(image.Point{Y: itemHeight}).Add(gtx.Ops)
		}
		itemOffset.Pop()

		// Pop the clip
		clipStack.Pop()

		return layout.Dimensions{Size: image.Point{X: menuWidth, Y: menuHeight}}
	}
}

// Layout renders the main background (invisible, just for context menu)
func (m *MainBackgroundWidget) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// This widget doesn't render anything visible, it just provides context menu functionality
	// It covers the entire available space
	return layout.Dimensions{Size: gtx.Constraints.Max}
}

// handleRightClick handles right-click events for the main background
func (m *MainBackgroundWidget) handleRightClick(gtx layout.Context, contextManager *widget.ContextManager) {
	log.Printf("Main background: Checking for right-click events")

	// Process right-click events for the main background
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: m,
			Kinds:  pointer.Press,
		})
		if !ok {
			break
		}

		e, ok := ev.(pointer.Event)
		if !ok {
			continue
		}

		log.Printf("Main background: Received pointer event, buttons=%v", e.Buttons)

		// Only handle right-clicks
		if e.Buttons != pointer.ButtonSecondary {
			continue
		}

		log.Printf("Main background: Right-click detected, showing context menu")

		// Show context menu for main background
		contextWidget := m.ContextMenu(gtx, e.Position.Round())
		if contextWidget != nil {
			contextManager.ShowContextWidget(gtx, contextWidget, e.Position.Round())
		}
		break // Only handle one event per frame
	}
}
