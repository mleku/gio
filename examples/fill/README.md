# Fill Widget Example

This example demonstrates a simple FillWidget that fills the entire window with a single color.

## Features Demonstrated

- **Simple Fill**: A single FillWidget that fills the entire window
- **Basic Usage**: Minimal example showing the core FillWidget functionality
- **Window Integration**: How to add a FillWidget to the window root

## Layout Structure

The example creates a window with the following structure:

```
Window Root
└── Fill Widget (Light blue background, fills entire window)
```

## Code Example

```go
// Create window widget
windowWidget := widget.New(widget.DefaultConfig())

// Create a single fill widget that fills the entire window
fillWidget := widget.NewFillWidget().
    Color(color.NRGBA{R: 100, G: 150, B: 200, A: 255}) // Light blue background

// Add the fill widget to the window root
windowWidget.Root().AddChild(fillWidget.Widget)

// Run the window
return windowWidget.Run(w)
```

## Key Features

- **Simplicity**: Minimal example showing just the FillWidget
- **Full Window**: Demonstrates filling the entire window area
- **Basic API**: Shows the core `Color()` method usage
- **Window Integration**: How to integrate with the window system

## Running the Example

```bash
go run ./examples/fill/main.go
```

## Visual Result

The example shows:
- A light blue background filling the entire window
- Clean, simple demonstration of the FillWidget
- No complex layouts or nested widgets

## Use Cases

- **Background Colors**: Simple colored backgrounds for applications
- **Theme Support**: Basic theme color application
- **Testing**: Minimal test case for FillWidget functionality
- **Learning**: Simple example for understanding FillWidget basics