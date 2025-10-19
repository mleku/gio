// SPDX-License-Identifier: Unlicense OR MIT

/*
Package app provides a platform-independent interface to operating system
functionality for running graphical user interfaces.

See https://github.com/mleku/gio for instructions to set up and run Gio programs.

# Window Management

A Window is run by calling its Event method in a loop. The first time a
method on Window is called, a new GUI window is created and shown. On JS/WASM
platforms, Window merely connects with a previously created GUI window.

The most important event is [FrameEvent] that prompts an update of the window
contents.

For example:

	w := new(app.Window)
	for {
		e := w.Event()
		if e, ok := e.(app.FrameEvent); ok {
			ops.Reset()
			// Add operations to ops.
			...
			// Completely replace the window contents and state.
			e.Frame(ops)
		}
	}

A program must keep receiving events from the event channel until
[DestroyEvent] is received.

# Main

The Main function must be called from a program's main function, to hand over
control of the main thread to operating systems that need it.

Because Main is also blocking on Linux/X11, the event loop of a Window must run in a goroutine.

For example, to display a blank but otherwise functional window:

	package main

	import "github.com/mleku/gio/app"

	func main() {
		go func() {
			w := app.NewWindow()
			for {
				w.Event()
			}
		}()
		app.Main()
	}

# Permissions

The packages under github.com/mleku/gio/app/permission should be imported
by a Gio program or by one of its dependencies to indicate that specific
operating-system permissions are required.  Please see documentation for
package github.com/mleku/gio/app/permission for more information.

Note: Permission packages are primarily for mobile platforms and may not be
applicable to Linux/X11 or JS/WASM platforms.
*/
package app
