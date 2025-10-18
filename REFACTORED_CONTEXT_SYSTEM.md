# Refactored Context Menu and Popup System

This document describes the refactored context menu and popup system that now uses the new position tracking feature.

## Overview

The context menu and tooltip system has been refactored to use the new `layout.Dimensions.Position` field instead of manually tracking widget bounds. This provides more accurate positioning and simplifies the API.

## Key Changes

### 1. ContextManager Refactoring

**Before:**
```go
type registeredWidget struct {
    widget   ContextWidget
    priority int
    bounds   image.Rectangle // Manual bounds tracking
}

func (cm *ContextManager) UpdateWidgetBounds(widget ContextWidget, bounds image.Rectangle)
```

**After:**
```go
type registeredWidget struct {
    widget   ContextWidget
    priority int
    position image.Point    // From layout.Dimensions.Position
    size     image.Point    // From layout.Dimensions.Size
}

func (cm *ContextManager) UpdateWidgetPosition(widget ContextWidget, position, size image.Point)
```

### 2. ContextWrapper Refactoring

**Before:**
```go
type ContextWrapper struct {
    // ...
    bounds image.Rectangle // Manual bounds tracking
}

func (cw *ContextWrapper) UpdateBounds(bounds image.Rectangle)
```

**After:**
```go
type ContextWrapper struct {
    // ...
    position image.Point // From layout.Dimensions.Position
    size     image.Point // From layout.Dimensions.Size
}

func (cw *ContextWrapper) GetPosition() image.Point
func (cw *ContextWrapper) GetSize() image.Point
```

### 3. TooltipManager Refactoring

**Before:**
```go
type registeredTooltipWidget struct {
    widget TooltipWidget
    bounds image.Rectangle // Manual bounds tracking
}
```

**After:**
```go
type registeredTooltipWidget struct {
    widget   TooltipWidget
    position image.Point // From layout.Dimensions.Position
    size     image.Point // From layout.Dimensions.Size
}
```

## Benefits

1. **Automatic Position Tracking**: Widgets automatically report their position through `layout.Dimensions.Position`
2. **More Accurate Hit Testing**: Position information is always up-to-date with the actual widget layout
3. **Simplified API**: No need to manually track and update bounds
4. **Better Tooltip Positioning**: Tooltips can be positioned relative to the actual widget position
5. **Consistent Coordinate System**: All positioning uses the same coordinate system as the layout engine

## Usage Examples

### Basic Context Menu Setup

```go
// Create a context wrapper
contextWrapper := widget.NewContextWrapper(
    myWidget,
    func(gtx layout.Context, pos image.Point) layout.Widget {
        // Return context menu widget
        return myContextMenuWidget
    },
    10, // Priority
)

// Register with context manager
contextManager.RegisterWidget(contextWrapper, contextWrapper.GetPriority())

// In your layout loop:
dims := contextWrapper.Layout(gtx)
contextManager.UpdateWidgetPosition(contextWrapper, dims.Position, dims.Size)
contextManager.Update(gtx)
contextManager.Layout(gtx)
```

### Tooltip Setup

```go
// Create a tooltip widget
tooltipWidget := &MyTooltipWidget{}

// Register with tooltip manager
tooltipManager.RegisterWidget(tooltipWidget)

// In your layout loop:
dims := tooltipWidget.Layout(gtx)
tooltipManager.UpdateWidgetPosition(tooltipWidget, dims.Position, dims.Size)
tooltipManager.Update(gtx)
tooltipManager.Layout(gtx, shaper)
```

### Helper Functions

The `ContextWrapper` now provides helper functions for easier management:

```go
// Register with both managers
contextWrapper.RegisterWithManagers(contextManager, tooltipManager)

// Update both managers after layout
contextWrapper.UpdateManagers(contextManager, tooltipManager)
```

## Migration Guide

### For Existing Code

1. **Replace `UpdateWidgetBounds` calls**:
   ```go
   // Old
   contextManager.UpdateWidgetBounds(widget, bounds)
   
   // New
   contextManager.UpdateWidgetPosition(widget, dims.Position, dims.Size)
   ```

2. **Update widget registration**:
   ```go
   // Old
   contextWrapper.UpdateBounds(bounds)
   
   // New
   // Position and size are automatically updated in Layout()
   ```

3. **Use new helper methods**:
   ```go
   // Get position directly
   pos := contextWrapper.GetPosition()
   size := contextWrapper.GetSize()
   bounds := contextWrapper.GetBounds() // Computed from position + size
   ```

## Implementation Details

### Position Calculation

The position is calculated by the layout engine and propagated through the layout hierarchy:

1. **Root Layout**: Starts at `(0,0)`
2. **Inset Layout**: Adds inset offsets to position
3. **Direction Layout**: Adds alignment offsets to position
4. **Flex/Stack Layout**: Calculates positions for child widgets
5. **Widget**: Receives final position in `layout.Dimensions.Position`

### Hit Testing

Hit testing now uses the position and size directly:

```go
widgetBounds := image.Rectangle{
    Min: reg.position,
    Max: reg.position.Add(reg.size),
}
if clickPos.In(widgetBounds) {
    // Handle click
}
```

### Context Menu Positioning

Context menus are positioned relative to the clicked widget's position:

```go
// Calculate position relative to widget
pos := widgetPos.Add(image.Point{X: 0, Y: -10}) // Above the widget
```

### Tooltip Positioning

Tooltips are positioned relative to the hovered widget's position:

```go
// Position tooltip relative to the widget's position
pos := widgetPos.Add(image.Point{X: 0, Y: -10}) // Above the widget
```

## Testing

The refactored system maintains backward compatibility while providing more accurate positioning. All existing tests pass, and new tests verify the position tracking functionality.

## Future Enhancements

1. **Automatic Manager Updates**: Could automatically update managers when widgets are laid out
2. **Position-based Animations**: Use position information for smooth animations
3. **Collision Detection**: Use position information for advanced collision detection
4. **Layout Debugging**: Visual debugging tools using position information
