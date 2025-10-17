# Primary Clipboard Demo - Complete Implementation

This directory contains a complete demo of the primary clipboard functionality implemented in gio.

## Files Created

- **`main.go`** - Main demo application showcasing primary clipboard functionality
- **`example.go`** - Code example showing how to use primary clipboard in your own apps
- **`go.mod`** - Go module configuration
- **`build.sh`** - Build script for easy compilation
- **`README.md`** - Documentation and usage instructions

## Features Demonstrated

### 1. Automatic Selection Copy
When users select text in the editor widget, it's automatically copied to the primary clipboard (X11 PRIMARY selection). This happens transparently without any user intervention.

### 2. Middle-Click Paste
Middle mouse button click pastes content from the primary clipboard into the editor at the cursor position.

### 3. Regular Clipboard Support
Traditional Ctrl+C/V operations continue to work for regular clipboard operations, providing backward compatibility.

### 4. Cross-Platform Support
- **X11**: Full primary clipboard support using X11 PRIMARY selection
- **Other platforms**: Graceful fallback with stub implementations

## How to Run

```bash
cd cmd/primary
go mod tidy
go run .
```

Or use the build script:
```bash
cd cmd/primary
chmod +x build.sh
./build.sh
./primary-demo
```

## Usage Instructions

1. **Type text**: Enter some text in the editor
2. **Select text**: Use mouse to select text - it automatically copies to primary clipboard
3. **Middle-click paste**: Click middle mouse button to paste from primary clipboard
4. **Regular clipboard**: Use Ctrl+C/V for traditional clipboard operations

## Technical Implementation

The demo showcases the following gio enhancements:

- **Enhanced clipboard package**: Added `WritePrimaryCmd` and `ReadPrimaryCmd`
- **Updated input router**: Handles primary clipboard operations
- **Platform-specific implementations**: X11 integration with PRIMARY selection
- **Enhanced editor widget**: Automatic selection copying and middle-click paste
- **Event handling**: Proper integration with gio's event system

## Code Example

```go
// Write to primary clipboard
gtx.Execute(clipboard.WritePrimaryCmd{Text: "Hello primary clipboard!"})

// Read from primary clipboard
gtx.Execute(clipboard.ReadPrimaryCmd{Tag: &myTag})
```

This implementation provides seamless primary clipboard functionality following the same patterns as the gel library, making it easy for developers to integrate primary clipboard support into their gio applications.
