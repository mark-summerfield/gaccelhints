// Copyright © 2022-23 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"embed"
)

//go:embed Version.dat
var Version string

//go:embed images/*.svg
var Images embed.FS

const (
	appID          = "eu.qtrac.gaccelhints"
	appName        = "gAccelHints"
	configFilename = "gaccelhints.ini"
	icon           = "images/accelhint.svg"
	iconUndo       = "images/edit-undo.svg"
	iconRedo       = "images/edit-redo.svg"
	iconCopy       = "images/edit-copy.svg"
	iconCut        = "images/edit-cut.svg"
	iconPaste      = "images/edit-paste.svg"
	iconQuit       = "images/shutdown.svg"
	accelTag       = "accel"
	escMarker      = "&&"
	placeholder    = "||"

	stdMargin     = 6
	defaultWidth  = 640
	defaultHeight = 480
	iconSize      = 32

	sigActivate  = "activate"
	sigChanged   = "changed"
	sigConfigure = "configure-event"
	sigClicked   = "clicked"
	sigDestroy   = "destroy"
	sigKeyPress  = "key-press-event"
	sigToggled   = "toggled"

	defaultOriginal = `Undo
Redo
Copy
Cu&t
Paste
Find
Find Again
Find && Replace`
)
