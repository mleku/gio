#!/bin/bash

echo "Testing Primary Clipboard Middle-Click Paste Demo"
echo "================================================"
echo ""
echo "Instructions for testing:"
echo "1. Run this demo: go run main.go"
echo "2. Type some text in the editor"
echo "3. Select the text with your mouse (it should copy to primary clipboard)"
echo "4. Click somewhere else in the editor"
echo "5. Middle-click to paste from primary clipboard"
echo ""
echo "Expected behavior:"
echo "- Text selection should automatically copy to primary clipboard"
echo "- Middle-click should paste the selected text at cursor position"
echo "- Debug output should show clipboard operations"
echo ""
echo "Note: This only works on X11 systems (Linux with X11, not Wayland)"
echo ""

# Check if we're on X11
if [ "$XDG_SESSION_TYPE" = "x11" ]; then
    echo "✅ Detected X11 session - primary clipboard should work"
elif [ "$XDG_SESSION_TYPE" = "wayland" ]; then
    echo "⚠️  Detected Wayland session - primary clipboard may not work"
else
    echo "❓ Unknown session type - primary clipboard may not work"
fi

echo ""
echo "Starting demo..."
go run main.go
