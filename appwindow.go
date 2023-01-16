// Copyright © 2022 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mark-summerfield/gong"
)

type AppWindow struct {
	config      *Config
	window      *gtk.ApplicationWindow
	container   *gtk.Widget
	statusLabel *gtk.Label
}

func newAppWindow(app *gtk.Application, config *Config) *AppWindow {
	appWindow := &AppWindow{}
	appWindow.config = config
	appWindow.makeWidgets(app)
	//appWindow.makeLayout()
	//appWindow.populateConfigWidgets()
	appWindow.makeConnections()
	appWindow.window.SetTitle(appName)
	raw, err := Images.ReadFile(icon)
	if err == nil {
		img, err := gdk.PixbufNewFromBytesOnly(raw)
		if err == nil {
			appWindow.window.SetIcon(img)
		}
	}
	appWindow.window.SetBorderWidth(stdMargin)
	appWindow.window.Add(appWindow.container)
	return appWindow
}

func (me *AppWindow) makeWidgets(app *gtk.Application) {
	var err error
	me.window, err = gtk.ApplicationWindowNew(app)
	gong.CheckError("Failed to create window:", err)
	me.statusLabel, err = gtk.LabelNew("Ready")
	gong.CheckError("Failed to create label:", err)
	me.statusLabel.SetXAlign(0.0)
}

func (me *AppWindow) makeLayout() {
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	me.container = &vbox.Container.Widget
}

func (me *AppWindow) makeConnections() {
	me.window.Connect(sigConfigure, func(_ *gtk.ApplicationWindow,
		rawEvent *gdk.Event) bool {
		event := gdk.EventConfigureNewFromEvent(rawEvent)
		me.config.updateGeometry(event.X(), event.Y(), event.Width(),
			event.Height())
		return false
	})
	me.window.Connect(sigDestroy, func(_ *gtk.ApplicationWindow) {
		me.onQuit()
	})
}

func (me *AppWindow) onQuit() {
	if me.config.dirty {
		me.config.save()
	}
	gtk.MainQuit()
}
