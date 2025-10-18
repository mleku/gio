# Widget Position Tracking

This implementation adds position tracking to the Gio layout engine, allowing widgets to know their coordinates relative to the window.

## Changes Made

### 1. Enhanced `layout.Dimensions` struct
Added a `Position` field to track the top-left corner of each widget relative to the window:

```go
type Dimensions struct {
    Size     image.Point
    Baseline int
    Position image.Point  // NEW: Top-left corner relative to window
}
```

### 2. Enhanced `layout.Context` struct
Added position tracking to the layout context:

```go
type Context struct {
    // ... existing fields ...
    Position image.Point  // NEW: Cumulative position offset
}
```

### 3. New Context Methods
Added helper methods to work with positions:

- `WithPosition(pos image.Point) Context` - Set position
- `OffsetPosition(offset image.Point) Context` - Add offset to position
- `GetPosition() image.Point` - Get current position

### 4. Updated Layout Functions
All layout functions now propagate position information:

- `Inset.Layout()` - Accounts for inset offsets
- `Direction.Layout()` - Accounts for alignment offsets
- `Flex.Layout()` - Tracks position for flex children
- `Stack.Layout()` - Tracks position for stacked children
- `Spacer.Layout()` - Returns current position

## Usage Example

```go
func myWidget(gtx layout.Context) layout.Dimensions {
    // Get the current position
    pos := gtx.GetPosition()
    
    // Use the position for calculations
    log.Printf("Widget is at position: %v", pos)
    
    // Layout your widget content
    // ...
    
    return layout.Dimensions{
        Size:     image.Point{X: 100, Y: 50},
        Position: pos, // Position is automatically set by layout functions
    }
}

// Use with insets
inset := layout.Inset{Top: 10, Left: 20}
dims := inset.Layout(gtx, myWidget)
// dims.Position will be (20, 10) - the inset offset
```

## Benefits

1. **Widget Awareness**: Widgets can now know exactly where they are positioned
2. **Tooltip Positioning**: Easy to position tooltips relative to widgets
3. **Hit Testing**: Simplified hit testing for complex layouts
4. **Debugging**: Easier to debug layout issues by knowing exact positions
5. **Context Menus**: Better positioning of context menus relative to widgets

## Backward Compatibility

This change is fully backward compatible. Existing code will continue to work without modification, as the `Position` field defaults to `(0,0)` for widgets that don't explicitly set it.
