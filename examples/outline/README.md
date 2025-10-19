# Outline Widget Example

This example demonstrates a simple OutlineWidget that draws a 2px square corner outline box around the entire window.

## Features Demonstrated

- **Simple Outline**: A 2px black outline drawn around the entire window
- **Square Corners**: No rounded corners, just straight edges
- **Window Integration**: How to draw an outline around the window root
- **Edge Drawing**: Drawing individual edges (top, bottom, left, right)

## Layout Structure

The example creates a window with the following structure:

```
Window Root
└── Custom Render Function (Draws 2px black outline around entire window)
```

## Code Example

```go
// Set the root widget to render as an outline widget
windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
    // Draw a 2px black outline around the entire window using Gio's stroke operations
    thickness := float32(2)
    
    // Create rectangle for the outline
    r := image.Rect(0, 0, w.Width, w.Height)
    
    // Draw the outline using Gio's stroke operation
    paint.FillShape(gtx.Ops,
        color.NRGBA{R: 0, G: 0, B: 0, A: 255}, // Black outline
        clip.Stroke{
            Path:  clip.RRect{Rect: r, NW: 0, NE: 0, SW: 0, SE: 0}.Path(gtx.Ops),
            Width: thickness,
        }.Op(),
    )
}
```

## Key Features

- **Gio Stroke Operations**: Uses Gio's native `clip.Stroke` and `clip.RRect` for proper outline drawing
- **Square Corners**: No rounded corners, just straight rectangular edges
- **Native Performance**: Leverages Gio's optimized stroke rendering
- **Window Integration**: Uses the root widget's render function for full window coverage
- **Thickness Control**: Easy to change outline thickness by modifying the `thickness` variable
- **Proper API Usage**: Follows Gio's recommended patterns for drawing outlines

## Running the Example

```bash
go run ./examples/outline/main.go
```

## Visual Result

The example shows:
- A 2px black outline around the entire window
- Square corners (no rounding)
- Clean demonstration of outline drawing
- No background fill, just the outline

## Use Cases

- **Window Borders**: Adding custom borders around application windows
- **Debugging**: Visualizing widget boundaries
- **UI Framing**: Creating frames around content areas
- **Testing**: Simple test case for outline functionality
- **Learning**: Understanding how to draw rectangular outlines
