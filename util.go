// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
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

func selectFirstRow(tree *gtk.TreeView) {
	selection, err := tree.GetSelection()
	if err == nil {
		path, err := gtk.TreePathNewFirst()
		if err == nil {
			selection.SelectPath(path)
		}
	}
}
