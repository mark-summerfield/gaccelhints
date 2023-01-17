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

	stdMargin     = 6
	defaultWidth  = 640
	defaultHeight = 480

	sigActivate  = "activate"
	sigChanged   = "changed"
	sigConfigure = "configure-event"
	sigClicked   = "clicked"
	sigDestroy   = "destroy"
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
