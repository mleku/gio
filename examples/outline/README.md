# Outline Example

This example demonstrates an OutlineWidget inside nested Flex containers, automatically drawing at the constraint max boundary.

## Features Demonstrated

- **Nested Flex Layout**: OutlineWidget inside Flex inside Flexed container
- **Automatic Sizing**: OutlineWidget automatically sizes to constraint boundaries
- **Background Fill**: Light cyan background fills the outline widget
- **White Border**: 16px white border around the outline widget

## Layout Structure

The example creates a window with the following structure:

```
Window Root
└── Outer FlexWidget (Column direction)
    └── Inner FlexWidget (Column direction, Flexed)
        └── OutlineWidget (Flexed - fills available space)
```

## Code Example

```go
// Create an outline widget
outlineWidget := widget.NewOutlineWidget().
    SetThickness(16).
    SetCornerRadius(0).
    SetOutlineColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255})

// Set custom background rendering for the outline widget
outlineWidget.Render = func(gtx app.Context, w *widget.Widget) {
    // Fill background with light cyan
    paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 240, B: 240, A: 255})
}

// Create inner flex container
innerFlex := widget.NewFlexWidget().
    SetDirection(widget.FlexColumn)

// Add the outline widget to the inner flex
innerFlex.Flexed(outlineWidget)

// Create outer flex container
outerFlex := widget.NewFlexWidget().
    SetDirection(widget.FlexColumn)

// Add the inner flex to the outer flex as a flexed item
outerFlex.Flexed(innerFlex)

// Set the root widget to render the outer flex
windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
    // Update the outer flex size to match the window
    outerFlex.SetSize(w.Width, w.Height)
    
    // Render the outer flex (which will render children)
    outerFlex.RenderWidget(gtx)
}
```

## Key Features

- **Nested Flex Containers**: Demonstrates complex layout with nested flex widgets
- **Automatic Constraint Sizing**: OutlineWidget automatically sizes to available space
- **Flex Layout**: Uses column direction for vertical layout
- **Responsive Design**: Automatically adapts to window size changes
- **Clean Hierarchy**: Proper widget composition with flex containers

## Running the Example

```bash
go run ./examples/outline/main.go
```

## Visual Result

The example shows:
- Light cyan background filling the outline widget area
- White border with 16px thickness around the outline widget
- Outline widget automatically sized to fill the available space
- Nested flex layout working correctly

## Use Cases

- **Complex Layouts**: Understanding nested flex containers
- **Automatic Sizing**: Widgets that adapt to available space
- **Flex Layout**: Building responsive UI layouts
- **Widget Composition**: Combining multiple widget types
- **Constraint-Based Layout**: Understanding how widgets size themselves