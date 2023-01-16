// Copyright © 2022-23 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mark-summerfield/accelhint"
	"github.com/mark-summerfield/gong"
)

type AppWindow struct {
	config          *Config
	window          *gtk.ApplicationWindow
	container       *gtk.Widget
	originalLabel   *gtk.Label
	originalText    *gtk.TextView
	hintedLabel     *gtk.Label
	hintedText      *gtk.TextView
	underlineButton *gtk.Button
	alphabetLabel   *gtk.Label
	alphabetEntry   *gtk.Entry
	markerFrame     *gtk.Frame
	ampersandRadio  *gtk.RadioButton
	underlineRadio  *gtk.RadioButton
	statusLabel     *gtk.Label
}

func newAppWindow(app *gtk.Application, config *Config) *AppWindow {
	appWindow := &AppWindow{}
	appWindow.config = config
	appWindow.makeWidgets(app)
	appWindow.makeLayout()
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
	me.originalLabel, err = gtk.LabelNewWithMnemonic("_Original")
	gong.CheckError("Failed to create label:", err)
	me.originalText, err = gtk.TextViewNew()
	gong.CheckError("Failed to create text view:", err)
	me.originalLabel.SetMnemonicWidget(me.originalText)
	me.hintedLabel, err = gtk.LabelNewWithMnemonic("_Hinted")
	gong.CheckError("Failed to create label:", err)
	me.hintedText, err = gtk.TextViewNew()
	gong.CheckError("Failed to create text view:", err)
	me.hintedLabel.SetMnemonicWidget(me.hintedText)
	me.hintedText.SetEditable(false)
	me.underlineButton, err = gtk.ButtonNewWithMnemonic("_Underline")
	gong.CheckError("Failed to create button:", err)
	me.markerFrame, err = gtk.FrameNew("Marker")
	gong.CheckError("Failed to create frame:", err)
	me.ampersandRadio, err = gtk.RadioButtonNewWithMnemonic(nil,
		"_Ampersands")
	gong.CheckError("Failed to create radio button:", err)
	me.underlineRadio, err = gtk.RadioButtonNewWithMnemonicFromWidget(
		me.ampersandRadio, "Under_lines")
	gong.CheckError("Failed to create radio button:", err)
	if me.config.marker == accelhint.GtkMarker {
		me.underlineRadio.SetActive(true)
	}
	me.alphabetLabel, err = gtk.LabelNewWithMnemonic("_Alphabet")
	gong.CheckError("Failed to create label:", err)
	me.alphabetEntry, err = gtk.EntryNew()
	gong.CheckError("Failed to create entry:", err)
	me.alphabetLabel.SetMnemonicWidget(me.alphabetEntry)
	me.alphabetEntry.SetText(me.config.alphabet)
	me.statusLabel, err = gtk.LabelNew("Ready (" + versionInfo() + ")")
	gong.CheckError("Failed to create label:", err)
	me.statusLabel.SetXAlign(0.0)
}

func (me *AppWindow) makeLayout() {
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	left, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	left.PackStart(me.originalLabel, false, false, stdMargin)
	left.PackStart(me.originalText, true, true, stdMargin)
	right, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	right.PackStart(me.hintedLabel, false, false, stdMargin)
	right.PackStart(me.hintedText, true, true, stdMargin)
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	hbox.PackStart(left, true, true, stdMargin)
	hbox.PackStart(right, true, true, stdMargin)
	vbox.PackStart(hbox, true, true, stdMargin)
	left, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	left.PackStart(me.underlineButton, false, false, stdMargin)
	left.PackStart(me.alphabetLabel, false, false, stdMargin)
	right, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	innerBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	innerBox.PackStart(me.ampersandRadio, false, false, stdMargin)
	innerBox.PackStart(me.underlineRadio, false, false, stdMargin)
	me.markerFrame.Add(innerBox)
	right.PackStart(me.markerFrame, false, false, stdMargin)
	right.PackStart(me.alphabetEntry, true, true, stdMargin)
	hbox, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	hbox.PackStart(left, false, false, stdMargin)
	hbox.PackStart(right, true, true, stdMargin)
	vbox.PackStart(hbox, false, false, stdMargin)
	vbox.PackStart(me.statusLabel, false, false, stdMargin)
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
