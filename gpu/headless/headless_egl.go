// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux
// +build linux

package headless

import (
	"github.com/mleku/gio/internal/egl"
)

func init() {
	newContextPrimary = func() (context, error) {
		return egl.NewContext(egl.EGL_DEFAULT_DISPLAY)
	}
}
