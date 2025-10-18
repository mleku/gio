# Context Widget System

This example demonstrates the context widget system for Gio, which allows any widget to provide a right-click context menu that appears at the cursor position.

## Features

- **Right-click Detection**: Automatically detects right-clicks anywhere in the UI
- **Smart Positioning**: Context widgets are positioned to avoid viewport edges, always opening towards the center
- **Widget Inheritance**: Widgets can inherit context menus from their parent containers
- **Priority System**: Higher priority widgets override lower priority ones
- **Easy Integration**: Simple interface that any widget can implement
- **Framed Menus**: Context menus have proper frames with borders
- **Hover Effects**: Menu items highlight on hover
- **Click Handling**: Left-clicking menu items prints the action and closes the popup

## How It Works

### 1. ContextWidget Interface

Any widget can implement the `ContextWidget` interface to provide a context menu:

```go
type ContextWidget interface {
    ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget
}
```

### 2. ContextManager

The `ContextManager` handles all right-click events and manages context widget display:

```go
contextManager := widget.NewContextManager()

// In your main loop:
contextManager.AddContextHandler(gtx.Ops)
contextManager.Update(gtx)
contextManager.Layout(gtx)
```

### 3. Widget Registration

Widgets register themselves with the context manager:

```go
contextManager.RegisterWidget(myWidget, priority)
```

### 4. Context Menu Implementation

Context menus are regular Gio widgets with frames, hover effects, and click handling:

```go
func (w *MyWidget) ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget {
    th := material.NewTheme()
    return func(gtx layout.Context) layout.Dimensions {
        // Create menu items
        menuItems := []string{"Copy", "Paste", "Delete"}
        clickables := make([]*widget.Clickable, len(menuItems))
        for i := range clickables {
            clickables[i] = &widget.Clickable{}
        }

        // Draw the menu frame
        layout.Flex{Axis: layout.Vertical}.Layout(gtx,
            layout.Rigid(func(gtx layout.Context) layout.Dimensions {
                // Draw frame background
                paint.Fill(gtx.Ops, color.NRGBA{R: 0xF0, G: 0xF0, B: 0xF0, A: 0xFF})
                
                // Draw frame border
                paint.FillShape(gtx.Ops, color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF},
                    clip.Stroke{
                        Path:  clip.Rect(image.Rectangle{Max: image.Point{X: 150, Y: 100}}).Path(),
                        Width: 1,
                    }.Op())
                
                return layout.Dimensions{Size: image.Point{X: 150, Y: 100}}
            }),
        )

        // Draw menu items with hover effects
        itemOffset := op.Offset(image.Point{X: 10, Y: 10}).Push(gtx.Ops)
        for i, item := range menuItems {
            itemClickable := clickables[i]

            // Check for clicks
            if itemClickable.Clicked(gtx) {
                log.Printf("Clicked: %s", item)
            }

            // Draw hover effect
            if itemClickable.Hovered() {
                paint.Fill(gtx.Ops, color.NRGBA{R: 0xE0, G: 0xE0, B: 0xFF, A: 0xFF})
            }

            // Draw item text
            label := material.Body2(th, item)
            label.Color = color.NRGBA{A: 0xFF}
            layout.W.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.Inset{Top: unit.Dp(5), Bottom: unit.Dp(5), Left: unit.Dp(10), Right: unit.Dp(10)}.Layout(gtx, label.Layout)
            })

            // Add clickable area
            itemClickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
                return layout.Dimensions{Size: image.Point{X: 130, Y: 25}}
            })

            // Move to next item
            op.Offset(image.Point{Y: 25}).Add(gtx.Ops)
        }
        itemOffset.Pop()

        return layout.Dimensions{Size: image.Point{X: 150, Y: 100}}
    }
}
```

## Priority System

- Higher priority widgets (larger numbers) override lower priority ones
- Child widgets typically have higher priority than their containers
- This allows buttons to show their own context menu instead of inheriting from their container

## Smart Positioning

The context widget system automatically positions context menus to:

1. Open towards the center of the viewport
2. Avoid going outside the viewport bounds
3. Choose the optimal corner based on click position relative to viewport center

## Running the Example

```bash
cd examples/popup
go run popup.go
```

Right-click on any of the colored buttons or the gray container to see different context menus appear at your cursor position. The menus have:

- **Frames**: Proper borders and backgrounds
- **Hover Effects**: Items highlight when you hover over them
- **Click Handling**: Left-click any item to see it logged and the menu closes

## Integration Notes

- The context widget system integrates seamlessly with existing Gio widgets
- No modifications to core Gio code are required
- Context widgets are regular Gio widgets, so they can contain any UI elements
- The system handles event propagation and widget dismissal automatically
- Each context menu is a complete widget with its own styling and behavior
