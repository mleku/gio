// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op/clip"
)

// InsetWidget represents a widget that provides margins and padding
type InsetWidget struct {
	*Widget
	// Margins (space outside the widget)
	MarginTop, MarginRight, MarginBottom, MarginLeft int
}

// NewInsetWidget creates a new inset widget
func NewInsetWidget() *InsetWidget {
	return &InsetWidget{
		Widget:    NewWidget(),
		MarginTop: 0, MarginRight: 0, MarginBottom: 0, MarginLeft: 0,
	}
}

// Margin sets all margins to the same value
func (i *InsetWidget) Margin(margin int) *InsetWidget {
	i.MarginTop = margin
	i.MarginRight = margin
	i.MarginBottom = margin
	i.MarginLeft = margin
	return i
}

// Margins sets all margins.
func (i *InsetWidget) Margins(top, right, bottom, left int) *InsetWidget {
	i.MarginTop = top
	i.MarginRight = right
	i.MarginBottom = bottom
	i.MarginLeft = left
	return i
}

// SetMarginVertical sets top and bottom margins
func (i *InsetWidget) SetMarginVertical(margin int) *InsetWidget {
	i.MarginTop = margin
	i.MarginBottom = margin
	return i
}

// SetMarginHorizontal sets left and right margins
func (i *InsetWidget) SetMarginHorizontal(margin int) *InsetWidget {
	i.MarginLeft = margin
	i.MarginRight = margin
	return i
}

// SetMarginTop sets the top margin
func (i *InsetWidget) SetMarginTop(margin int) *InsetWidget {
	i.MarginTop = margin
	return i
}

// SetMarginRight sets the right margin
func (i *InsetWidget) SetMarginRight(margin int) *InsetWidget {
	i.MarginRight = margin
	return i
}

// SetMarginBottom sets the bottom margin
func (i *InsetWidget) SetMarginBottom(margin int) *InsetWidget {
	i.MarginBottom = margin
	return i
}

// SetMarginLeft sets the left margin
func (i *InsetWidget) SetMarginLeft(margin int) *InsetWidget {
	i.MarginLeft = margin
	return i
}

// GetWidget returns the underlying widget (for Renderer interface)
func (i *InsetWidget) GetWidget() *Widget {
	return i.Widget
}

// RenderWidget renders the inset widget
func (i *InsetWidget) RenderWidget(gtx app.Context) {
	if !i.Visible {
		return
	}

	// Set up clipping for the entire widget area
	defer clip.Rect{Min: image.Point{X: i.X, Y: i.Y}, Max: image.Point{X: i.X + i.Width, Y: i.Y + i.Height}}.Push(gtx.Ops).Pop()

	// Call custom render function if provided
	if i.Render != nil {
		i.Render(gtx, i.Widget)
	}

	// Render children with padding applied
	for _, child := range i.Children {
		// Adjust child position to account for margins and padding
		child.X = i.X + i.MarginLeft
		child.Y = i.Y + i.MarginTop
		child.RenderWidget(gtx)
	}
}

// GetContentArea returns the area available for content (after margins and padding)
func (i *InsetWidget) GetContentArea() image.Rectangle {
	return image.Rect(
		i.X+i.MarginLeft,
		i.Y+i.MarginTop,
		i.X+i.Width-i.MarginRight,
		i.Y+i.Height-i.MarginBottom,
	)
}

// GetContentSize returns the size available for content
func (i *InsetWidget) GetContentSize() (width, height int) {
	contentArea := i.GetContentArea()
	return contentArea.Dx(), contentArea.Dy()
}

// Fluent methods for InsetWidget that delegate to the embedded Widget

// SetPosition sets the inset widget's position
func (i *InsetWidget) SetPosition(x, y int) *InsetWidget {
	i.Widget.SetPosition(x, y)
	return i
}

// SetSize sets the inset widget's size
func (i *InsetWidget) SetSize(width, height int) *InsetWidget {
	i.Widget.SetSize(width, height)
	return i
}

// SetVisible sets the inset widget's visibility
func (i *InsetWidget) SetVisible(visible bool) *InsetWidget {
	i.Widget.SetVisible(visible)
	return i
}

// AddChild adds a child widget
func (i *InsetWidget) AddChild(child *Widget) *InsetWidget {
	i.Widget.AddChild(child)
	return i
}
