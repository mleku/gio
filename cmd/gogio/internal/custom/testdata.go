// SPDX-License-Identifier: Unlicense OR MIT

//go:build linux
// +build linux

// This program demonstrates the use of a custom OpenGL ES context with
// app.Window.
package main

import (
	"image/color"
	"log"
	"os"

	"gio.mleku.dev/app"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
)

/*
#cgo linux pkg-config: egl wayland-egl
#cgo freebsd openbsd CFLAGS: -I/usr/local/include
#cgo openbsd CFLAGS: -I/usr/X11R6/include
#cgo freebsd openbsd LDFLAGS: -L/usr/local/lib
#cgo openbsd LDFLAGS: -L/usr/X11R6/lib
#cgo linux LDFLAGS: -lEGL -lwayland-egl
#cgo freebsd openbsd LDFLAGS: -lEGL

#include <EGL/egl.h>
#include <EGL/eglext.h>
#include <wayland-client.h>
#include <wayland-egl.h>

static EGLDisplay eglDisplay;
static EGLConfig eglConfig;
static EGLContext eglContext;
static struct wl_egl_window *eglWindow;
static EGLSurface eglSurface;

int initEGL(void *display) {
	eglDisplay = eglGetDisplay((EGLNativeDisplayType)display);
	if (eglDisplay == EGL_NO_DISPLAY) {
		return 0;
	}
	if (!eglInitialize(eglDisplay, NULL, NULL)) {
		return 0;
	}
	EGLint configAttribs[] = {
		EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
		EGL_BLUE_SIZE, 8,
		EGL_GREEN_SIZE, 8,
		EGL_RED_SIZE, 8,
		EGL_ALPHA_SIZE, 8,
		EGL_RENDERABLE_TYPE, EGL_OPENGL_ES2_BIT,
		EGL_NONE,
	};
	EGLint numConfigs;
	if (!eglChooseConfig(eglDisplay, configAttribs, &eglConfig, 1, &numConfigs)) {
		return 0;
	}
	EGLint contextAttribs[] = {
		EGL_CONTEXT_CLIENT_VERSION, 2,
		EGL_NONE,
	};
	eglContext = eglCreateContext(eglDisplay, eglConfig, EGL_NO_CONTEXT, contextAttribs);
	if (eglContext == EGL_NO_CONTEXT) {
		return 0;
	}
	return 1;
}

int createEGLSurface(void *surface, int width, int height) {
	eglWindow = wl_egl_window_create((struct wl_surface*)surface, width, height);
	if (eglWindow == EGL_NO_SURFACE) {
		return 0;
	}
	eglSurface = eglCreateWindowSurface(eglDisplay, eglConfig, (EGLNativeWindowType)eglWindow, NULL);
	if (eglSurface == EGL_NO_SURFACE) {
		return 0;
	}
	return 1;
}

void makeCurrent() {
	eglMakeCurrent(eglDisplay, eglSurface, eglSurface, eglContext);
}

void swapBuffers() {
	eglSwapBuffers(eglDisplay, eglSurface);
}

void destroyEGL() {
	if (eglSurface != EGL_NO_SURFACE) {
		eglDestroySurface(eglDisplay, eglSurface);
	}
	if (eglWindow) {
		wl_egl_window_destroy(eglWindow);
	}
	if (eglContext != EGL_NO_CONTEXT) {
		eglDestroyContext(eglDisplay, eglContext);
	}
	if (eglDisplay != EGL_NO_DISPLAY) {
		eglTerminate(eglDisplay);
	}
}
*/
import "C"

func main() {
	go func() {
		w := new(app.Window)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Paint black background
			paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 0, B: 0, A: 255})

			e.Frame(gtx.Ops)
		case pointer.Event:
			// Log mouse events
			log.Printf("Mouse event: %+v\n", e)
		}
	}
}
