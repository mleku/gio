// SPDX-License-Identifier: Unlicense OR MIT

package input

import (
	"io"
	"slices"

	"gio.mleku.dev/io/clipboard"
	"gio.mleku.dev/io/event"
)

// clipboardState contains the state for clipboard event routing.
type clipboardState struct {
	receivers []event.Tag
}

// primaryClipboardState contains the state for primary clipboard event routing.
type primaryClipboardState struct {
	receivers []event.Tag
}

type clipboardQueue struct {
	// request avoid read clipboard every frame while waiting.
	requested bool
	mime      string
	text      []byte
	// primary clipboard state
	primaryState     primaryClipboardState
	primaryText      string
	primaryRequested bool
}

// WriteClipboard returns the most recent data to be copied
// to the clipboard, if any.
func (q *clipboardQueue) WriteClipboard() (mime string, content []byte, ok bool) {
	if q.text == nil {
		return "", nil, false
	}
	content = q.text
	q.text = nil
	return q.mime, content, true
}

// WritePrimaryClipboard returns the most recent text to be copied
// to the primary clipboard, if any.
func (q *clipboardQueue) WritePrimaryClipboard() (text string, ok bool) {
	if q.primaryText == "" {
		return "", false
	}
	text = q.primaryText
	q.primaryText = ""
	return text, true
}

// ClipboardRequested reports if any new handler is waiting
// to read the clipboard.
func (q *clipboardQueue) ClipboardRequested(state clipboardState) bool {
	req := len(state.receivers) > 0 && q.requested
	q.requested = false
	return req
}

// PrimaryClipboardRequested reports if any new handler is waiting
// to read the primary clipboard.
func (q *clipboardQueue) PrimaryClipboardRequested(state primaryClipboardState) bool {
	req := len(state.receivers) > 0 && q.primaryRequested
	q.primaryRequested = false
	return req
}

func (q *clipboardQueue) Push(state clipboardState, e event.Event) (clipboardState, []taggedEvent) {
	var evts []taggedEvent
	for _, r := range state.receivers {
		evts = append(evts, taggedEvent{tag: r, event: e})
	}
	state.receivers = nil
	return state, evts
}

func (q *clipboardQueue) PushPrimary(state primaryClipboardState, e event.Event) (primaryClipboardState, []taggedEvent) {
	var evts []taggedEvent
	for _, r := range state.receivers {
		evts = append(evts, taggedEvent{tag: r, event: e})
	}
	state.receivers = nil
	return state, evts
}

func (q *clipboardQueue) ProcessWriteClipboard(req clipboard.WriteCmd) {
	defer req.Data.Close()
	content, err := io.ReadAll(req.Data)
	if err != nil {
		return
	}
	q.mime = req.Type
	q.text = content
}

func (q *clipboardQueue) ProcessWritePrimaryClipboard(req clipboard.WritePrimaryCmd) {
	q.primaryText = req.Text
}

func (q *clipboardQueue) ProcessReadClipboard(state clipboardState, tag event.Tag) clipboardState {
	if slices.Contains(state.receivers, tag) {
		return state
	}
	n := len(state.receivers)
	state.receivers = append(state.receivers[:n:n], tag)
	q.requested = true
	return state
}

func (q *clipboardQueue) ProcessReadPrimaryClipboard(state primaryClipboardState, tag event.Tag) primaryClipboardState {
	if slices.Contains(state.receivers, tag) {
		return state
	}
	n := len(state.receivers)
	state.receivers = append(state.receivers[:n:n], tag)
	q.primaryRequested = true
	return state
}
