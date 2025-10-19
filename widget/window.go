// SPDX-License-Identifier: Unlicense OR MIT

package widget

import (
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
)

// Config holds configuration options for the Window widget
type Config struct {
	// BackgroundColor sets the window background color
	BackgroundColor color.NRGBA

	// OnFrame is called for each frame to render the UI
	OnFrame func(gtx app.Context)
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

// Window is a basic window widget that handles window management
type Window struct {
	config Config
	root   *Widget
}

// New creates a new Window widget with the given configuration
func New(config Config) *Window {
	return &Window{
		config: config,
		root:   NewWidget().SetSize(800, 600), // Default window size
	}
}

// Root returns the root widget for adding child widgets
func (w *Window) Root() *Widget {
	return w.root
}

// Run starts the window event loop
func (w *Window) Run(window *app.Window) error {
	var ops op.Ops

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Update root widget size to match window size
			w.root.SetSize(gtx.Size.X, gtx.Size.Y)

			// Paint background
			paint.Fill(gtx.Ops, w.config.BackgroundColor)

			// Render the root widget
			w.root.RenderWidget(gtx)

			// Call user-defined frame handler
			if w.config.OnFrame != nil {
				w.config.OnFrame(gtx)
			}

			e.Frame(gtx.Ops)
		}
	}
}
