// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux || js
// +build linux js

package app

import "os"

func dataDir() (string, error) {
	return os.UserConfigDir()
}
