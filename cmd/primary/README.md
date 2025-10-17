# Primary Clipboard Demo

This demo showcases the primary clipboard functionality implemented in gio, which provides:

- **Automatic selection copy**: When you select text in the editor or selectable widget, it's automatically copied to the primary clipboard (X11 PRIMARY selection)
- **Middle-click paste**: Middle mouse button click pastes content from the primary clipboard
- **Regular clipboard support**: Ctrl+C/V still work for regular clipboard operations

## Features Demonstrated

1. **Text Selection**: Select text with your mouse and it will automatically be copied to the primary clipboard
2. **Middle-Click Paste**: Click the middle mouse button to paste from the primary clipboard
3. **Multi-line Editor**: The demo includes a multi-line text editor for testing
4. **Selectable Widget**: Also includes a read-only selectable widget demo
5. **Cross-platform**: Works on X11 systems, gracefully handles other platforms

## How to Run

1. Navigate to the gio root directory
2. Run the editor demo:
   ```bash
   cd cmd/primary
   go mod tidy
   go run main.go
   ```
3. Run the selectable widget demo:
   ```bash
   go run selectable_demo.go
   ```

## Platform Support

- **X11**: Full primary clipboard support using X11 PRIMARY selection
- **Wayland**: Primary clipboard not supported (stub implementation)
- **Windows/macOS/iOS/Android/Web**: Primary clipboard not supported (stub implementation)

## Usage Instructions

### Editor Demo (main.go)
1. Type some text in the editor
2. Select text with your mouse - it will automatically copy to primary clipboard
3. Click middle mouse button to paste from primary clipboard
4. Use Ctrl+C/V for regular clipboard operations

### Selectable Demo (selectable_demo.go)
1. Select text in the selectable widget - it will automatically copy to primary clipboard
2. Click middle mouse button to demonstrate primary clipboard functionality
3. Note: This is a read-only widget, so middle-click won't insert text, but shows the clipboard working
4. Use Ctrl+C for regular clipboard operations

## Technical Details

The implementation includes:

- Enhanced clipboard package with `WritePrimaryCmd` and `ReadPrimaryCmd`
- Updated input router to handle primary clipboard operations
- Platform-specific implementations for all supported operating systems
- Enhanced editor widget with automatic selection copying and middle-click paste
- Enhanced selectable widget with automatic selection copying and middle-click detection
- X11 integration using X11 PRIMARY selection mechanism
- **Fixed X11 SelectionNotify handler** to properly handle both regular and primary clipboard events
- **Improved middle-click handling** based on p9c implementation:
  - Middle-click handled in `pointer.Event` case (not gestures)
  - Cursor positioned at click location before paste
  - Selection cleared before paste
  - Direct clipboard command execution
- **Enhanced mouse button support**:
  - Gesture system now accepts all mouse buttons (not just primary)
  - ClickEvent includes Button field for proper button identification
  - Both gesture and pointer event handling for comprehensive coverage
- **Fixed middle-click gesture conflict**:
  - Middle-click gestures are filtered out to prevent being treated as left-clicks
  - Only left-clicks are handled by gesture system
  - Middle-clicks are handled by pointer.Event system for paste functionality
- **Fixed middle-click event consumption**:
  - Modified gesture system to not consume middle-click events
  - Middle-click events now reach pointer.Event handlers properly
  - Ensures middle-click detection works correctly
- **Fixed transfer event filtering**:
  - Editor now properly filters transfer events by MIME type
  - Both "application/text" and "text/plain" events are handled
  - Ensures clipboard data events are properly received
- **Fixed primary clipboard data delivery**:
  - Router now delivers transfer.DataEvent to both regular and primary clipboard receivers
  - Primary clipboard data events are properly routed to requesting widgets
  - Ensures middle-click paste receives the clipboard data
- **Cross-platform support**: Works on X11 systems, gracefully handles other platforms

This demo is based on the gel library's clipboard implementation and improved using p9c's proven middle-click approach, providing seamless primary clipboard functionality for X11 environments.
