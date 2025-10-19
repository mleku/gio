# Border Widget Example

This example demonstrates a border widget with a fill operation before the border operation, showing how the render queue paints.

## Features Demonstrated

- **Fill Background**: Light gray background fills the entire window
- **Border Overlay**: 8Dp black border with 50% transparency drawn on the inside of the background
- **Render Queue**: Shows the order of operations in the render queue
- **Layered Rendering**: Background fill followed by border stroke

## Layout Structure

The example creates a window with the following structure:

```
Window Root
└── Custom Render Function (Fill background, then draw border)
```

## Code Example

```go
// Set the root widget to render as a border widget with fill background
windowWidget.Root().Render = func(gtx app.Context, w *widget.Widget) {
    // First, fill the entire area with a light gray background
    paint.Fill(gtx.Ops, color.NRGBA{R: 240, G: 240, B: 240, A: 255}) // Light gray background
    
    // Then, draw an 8Dp black border with 50% transparency on the inside of the filled area
    thickness := float32(8) // 8Dp wide border
    inset := int(thickness) // Inset the border from the edges
    
    // Create rectangle for the border (inset from the window edges)
    r := image.Rect(inset, inset, w.Width-inset, w.Height-inset)
    
    // Draw the border using Gio's stroke operation
    paint.FillShape(gtx.Ops,
        color.NRGBA{R: 0, G: 0, B: 0, A: 128}, // Black border with 50% transparency (128/255)
        clip.Stroke{
            Path:  clip.RRect{Rect: r, NW: 0, NE: 0, SW: 0, SE: 0}.Path(gtx.Ops),
            Width: thickness,
        }.Op(),
    )
}
```

## Key Features

- **Render Queue Order**: Demonstrates the order of operations in the render queue
- **Background Fill**: Light gray background fills the entire window area
- **Border Overlay**: Black border drawn on the inside of the background area
- **Gio Stroke Operations**: Uses Gio's native `clip.Stroke` and `clip.RRect`
- **Visual Contrast**: Background and border colors provide clear visual separation
- **Layered Composition**: Shows how multiple operations compose together

## Running the Example

```bash
go run ./examples/border/main.go
```

## Visual Result

The example shows:
- Light cyan background filling the entire window
- 8Dp black border with 50% transparency inset from the window edges
- Clear visual demonstration of the render queue order
- Background fill operation followed by border stroke operation

## Use Cases

- **Render Queue Testing**: Verifying that the render queue is working correctly
- **Layered UI**: Understanding how multiple operations compose
- **Background + Border**: Common pattern for UI elements
- **Visual Debugging**: Seeing the order of rendering operations
- **Learning**: Understanding Gio's rendering pipeline
