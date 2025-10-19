// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"gio.mleku.dev/app"
)

// FlexDirection defines the direction of flex layout
type FlexDirection int

const (
	FlexRow FlexDirection = iota
	FlexColumn
)

// FlexItem represents a single item in a flex container
type FlexItem struct {
	Renderer Renderer // Changed from Widget to Renderer interface
	Weight   float32  // Weight for flexed items (default 1.0)
	Rigid    bool     // Whether this item has a fixed size
}

// FlexWidget represents a flex container widget
type FlexWidget struct {
	*Widget
	Direction FlexDirection
	Items     []FlexItem
}

// NewFlexWidget creates a new flex widget
func NewFlexWidget() *FlexWidget {
	return &FlexWidget{
		Widget:    NewWidget(),
		Direction: FlexRow,
		Items:     make([]FlexItem, 0),
	}
}

// SetDirection sets the flex direction
func (f *FlexWidget) SetDirection(direction FlexDirection) *FlexWidget {
	f.Direction = direction
	return f
}

// AddRigid adds a rigid (fixed-size) item to the flex container
func (f *FlexWidget) Rigid(renderer Renderer) *FlexWidget {
	f.Items = append(f.Items, FlexItem{
		Renderer: renderer,
		Rigid:    true,
		Weight:   0,
	})
	return f
}

// AddFlexed adds a flexed item with the given weight
func (f *FlexWidget) Flexed(renderer Renderer) *FlexWidget {
	f.Items = append(f.Items, FlexItem{
		Renderer: renderer,
		Rigid:    false,
		Weight:   1.0,
	})
	return f
}

// Weight sets the weight of a Flexed, sets the value for the most recently created Flexed.
func (f *FlexWidget) Weight(weight float32) *FlexWidget {
	f.Items[len(f.Items)-1].Weight = weight
	return f
}

// Layout calculates and sets positions and sizes for all flex items
func (f *FlexWidget) Layout() {
	if len(f.Items) == 0 {
		return
	}

	if f.Direction == FlexRow {
		f.layoutRow()
	} else {
		f.layoutColumn()
	}
}

// layoutRow arranges items horizontally
func (f *FlexWidget) layoutRow() {
	availableWidth := f.Width
	totalWeight := float32(0)
	rigidWidth := 0

	// Calculate total weight and rigid width
	for _, item := range f.Items {
		if item.Rigid {
			rigidWidth += item.Renderer.GetWidget().Width
		} else {
			totalWeight += item.Weight
		}
	}

	availableFlexWidth := availableWidth - rigidWidth
	if availableFlexWidth < 0 {
		availableFlexWidth = 0
	}

	// Position items
	currentX := f.X
	for _, item := range f.Items {
		item.Renderer.GetWidget().X = currentX
		item.Renderer.GetWidget().Y = f.Y

		if item.Rigid {
			// Use widget's natural width
			currentX += item.Renderer.GetWidget().Width
		} else {
			// Calculate flexed width
			if totalWeight > 0 {
				flexWidth := int(float32(availableFlexWidth) * (item.Weight / totalWeight))
				item.Renderer.GetWidget().Width = flexWidth
			}
			currentX += item.Renderer.GetWidget().Width
		}

		// Set height to match container
		item.Renderer.GetWidget().Height = f.Height
	}
}

// layoutColumn arranges items vertically
func (f *FlexWidget) layoutColumn() {
	availableHeight := f.Height
	totalWeight := float32(0)
	rigidHeight := 0

	// Calculate total weight and rigid height
	for _, item := range f.Items {
		if item.Rigid {
			rigidHeight += item.Renderer.GetWidget().Height
		} else {
			totalWeight += item.Weight
		}
	}

	availableFlexHeight := availableHeight - rigidHeight
	if availableFlexHeight < 0 {
		availableFlexHeight = 0
	}

	// Position items
	currentY := f.Y
	for _, item := range f.Items {
		item.Renderer.GetWidget().X = f.X
		item.Renderer.GetWidget().Y = currentY

		if item.Rigid {
			// Use widget's natural height
			currentY += item.Renderer.GetWidget().Height
		} else {
			// Calculate flexed height
			if totalWeight > 0 {
				flexHeight := int(float32(availableFlexHeight) * (item.Weight / totalWeight))
				item.Renderer.GetWidget().Height = flexHeight
			}
			currentY += item.Renderer.GetWidget().Height
		}

		// Set width to match container
		item.Renderer.GetWidget().Width = f.Width
	}
}

// RenderFlexWidget renders the flex widget and its items
func (f *FlexWidget) RenderFlexWidget(gtx app.Context) {
	// Layout items first
	f.Layout()

	// Render the flex container itself
	f.RenderWidget(gtx)

	// Render all flex items
	for _, item := range f.Items {
		if item.Renderer != nil {
			item.Renderer.RenderWidget(gtx)
		}
	}
}

// Fluent methods for FlexWidget that delegate to the embedded Widget

// SetPosition sets the flex widget's position
func (f *FlexWidget) SetPosition(x, y int) *FlexWidget {
	f.Widget.SetPosition(x, y)
	return f
}

// SetSize sets the flex widget's size
func (f *FlexWidget) SetSize(width, height int) *FlexWidget {
	f.Widget.SetSize(width, height)
	return f
}

// SetVisible sets the flex widget's visibility
func (f *FlexWidget) SetVisible(visible bool) *FlexWidget {
	f.Widget.SetVisible(visible)
	return f
}
