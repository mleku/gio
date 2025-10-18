package widget

import (
	"image"
	"testing"

	"gio.mleku.dev/layout"
)

func TestContextManager(t *testing.T) {
	// Create a context manager
	cm := NewContextManager()

	// Test that it starts with no active context
	if cm.activeContext != nil {
		t.Error("Expected no active context initially")
	}

	// Test viewport size tracking
	cm.viewportSize = image.Point{X: 800, Y: 600}
	if cm.viewportSize.X != 800 || cm.viewportSize.Y != 600 {
		t.Error("Viewport size not set correctly")
	}
}

func TestRegisterWidget(t *testing.T) {
	// Create a context manager
	cm := NewContextManager()

	// Create a mock context widget
	mockWidget := &mockContextWidget{}

	// Register the widget
	cm.RegisterWidget(mockWidget, 5)

	if len(cm.registeredWidgets) != 1 {
		t.Error("Expected exactly one registered widget")
	}

	if cm.registeredWidgets[0].widget != mockWidget {
		t.Error("Expected widget to match")
	}
	if cm.registeredWidgets[0].priority != 5 {
		t.Error("Expected priority to match")
	}
}

func TestCalculateOptimalPosition(t *testing.T) {
	cm := &ContextManager{
		viewportSize: image.Point{X: 800, Y: 600},
	}

	// Test click in top-left quadrant
	clickPos := image.Point{X: 100, Y: 100}
	dims := layout.Dimensions{Size: image.Point{X: 200, Y: 150}}

	pos := cm.calculateOptimalPosition(layout.Context{}, clickPos, dims)

	// Should position widget to the right and below the click (towards center)
	expectedX := clickPos.X
	expectedY := clickPos.Y

	if pos.X != expectedX || pos.Y != expectedY {
		t.Errorf("Expected position (%d, %d), got (%d, %d)", expectedX, expectedY, pos.X, pos.Y)
	}

	// Test click in bottom-right quadrant
	clickPos = image.Point{X: 700, Y: 500}
	pos = cm.calculateOptimalPosition(layout.Context{}, clickPos, dims)

	// Should position widget to the left and above the click (towards center)
	expectedX = clickPos.X - dims.Size.X
	expectedY = clickPos.Y - dims.Size.Y

	if pos.X != expectedX || pos.Y != expectedY {
		t.Errorf("Expected position (%d, %d), got (%d, %d)", expectedX, expectedY, pos.X, pos.Y)
	}
}

// Mock context widget for testing
type mockContextWidget struct{}

func (m *mockContextWidget) ContextMenu(gtx layout.Context, clickPos image.Point) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Point{X: 100, Y: 50}}
	}
}
