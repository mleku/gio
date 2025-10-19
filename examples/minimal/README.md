# Window-Level Mouse Enter/Exit Detection Example

This example demonstrates how to detect when the mouse enters or leaves the entire window in Gio. This is useful for enabling/disabling hover effects when the mouse is outside the window.

## Key Features

- **Window-Level Event Detection**: Uses `pointer.Enter` and `pointer.Leave` events to detect mouse enter/exit
- **Clear Logging**: Prominently logs window-level mouse events with clear indicators
- **Hover Effect Control**: Shows when hover effects should be enabled or disabled

## How It Works

1. **Event Handler Setup**: Creates a window-level event handler that covers the entire window:
   ```go
   area := clip.Rect{Max: gtx.Size}.Push(gtx.Ops)
   event.Op(gtx.Ops, w)
   area.Pop()
   ```

2. **Event Filtering**: Listens specifically for enter/leave events:
   ```go
   ev, ok := gtx.Source.Event(pointer.Filter{
       Target: w,
       Kinds:  pointer.Enter | pointer.Leave,
   })
   ```

3. **Event Processing**: Handles enter/leave events separately from other pointer events

## Expected Output

When you run this example, you'll first see instructions:

```
lol.mleku.dev: ===== INSTRUCTIONS =====
lol.mleku.dev: Window size: 800x600 pixels
lol.mleku.dev: To see Enter/Leave events, move your mouse:
lol.mleku.dev: 1. COMPLETELY OUTSIDE the window (to see Leave event)
lol.mleku.dev: 2. BACK INSIDE the window (to see Enter event)
lol.mleku.dev: Just moving within the window won't trigger these events.
lol.mleku.dev: =========================
```

**Important**: You must move your mouse **completely outside the window bounds** to see Leave events, and then back inside to see Enter events. Just moving within the window will only show Move events.

When you move your mouse completely outside and back inside, you'll see:

```
lol.mleku.dev: ===== MOUSE LEFT WINDOW ===== Position=(200.0,300.0), Source=Mouse @ 14:30:28.456
lol.mleku.dev: Hover effects should be DISABLED

lol.mleku.dev: ===== MOUSE ENTERED WINDOW ===== Position=(100.0,50.0), Source=Mouse @ 14:30:25.123
lol.mleku.dev: Hover effects should be ENABLED
```

## Usage in Your Application

To use this pattern in your own application:

1. Set up a window-level event handler covering the entire window
2. Listen for `pointer.Enter` and `pointer.Leave` events
3. Enable/disable hover effects based on these events
4. Handle other pointer events (moves, clicks) separately

This approach works across all platforms supported by Gio and provides reliable window-level mouse tracking.
