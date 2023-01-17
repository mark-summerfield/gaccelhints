// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mark-summerfield/gong"
)

func main() {
	log.SetFlags(log.Lmsgprefix)
	config := NewConfig()
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	gong.CheckError("Failed to create application:", err)
	app.Connect(sigActivate, func() {
		appWindow := newAppWindow(app, config)
		appWindow.window.Move(config.getPosition())
		appWindow.window.SetDefaultSize(config.width, config.height)
		appWindow.onTextChanged()
		appWindow.window.ShowAll()
	})
	os.Exit(app.Run(os.Args))
}
