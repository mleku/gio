//go:build !android
// +build !android

package main

import "gio.mleku.dev/io/event"

func ProcessPlatformEvent(event event.Event) bool {
	return false
}
