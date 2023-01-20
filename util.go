// Copyright © 2022-23 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func versionInfo() string {
	return fmt.Sprintf("%s %s • Go %s • Gtk %d.%d.%d", appName,
		strings.TrimSpace(Version),
		strings.TrimPrefix(runtime.Version(), "go"), gtk.GetMajorVersion(),
		gtk.GetMinorVersion(), gtk.GetMicroVersion())
}

func getPixbuf(name string, size int) *gdk.Pixbuf {
	raw, err := Images.ReadFile(name)
	if err == nil {
		img, err := gdk.PixbufNewFromBytesOnly(raw)
		if err == nil {
			if size > 0 {
				img, err = img.ScaleSimple(size, size, gdk.INTERP_NEAREST)
			}
			if err == nil {
				return img
			}
		}
	}
	return nil
}

func getImage(name string) *gtk.Image {
	if pixbuf := getPixbuf(name, iconSize); pixbuf != nil {
		if img, err := gtk.ImageNewFromPixbuf(pixbuf); err == nil {
			return img
		}
	}
	return nil
}
