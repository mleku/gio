// SPDX-License-Identifier: Unlicense OR MIT

package clipboard

import (
	"io"

	"gioui.org/io/event"
)

// WriteCmd copies Text to the clipboard.
type WriteCmd struct {
	Type string
	Data io.ReadCloser
}

// ReadCmd requests the text of the clipboard, delivered to
// the handler through an [io/transfer.DataEvent].
type ReadCmd struct {
	Tag event.Tag
}

// WritePrimaryCmd copies Text to the primary clipboard (X11 PRIMARY selection).
type WritePrimaryCmd struct {
	Text string
}

// ReadPrimaryCmd requests the text of the primary clipboard, delivered to
// the handler through an [io/transfer.DataEvent].
type ReadPrimaryCmd struct {
	Tag event.Tag
}

func (WriteCmd) ImplementsCommand()        {}
func (ReadCmd) ImplementsCommand()         {}
func (WritePrimaryCmd) ImplementsCommand() {}
func (ReadPrimaryCmd) ImplementsCommand()  {}
