// Copyright © 2022 Mark Summerfield. All rights reserved.
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
	icon           = "gaccelhints.svg"
	configFilename = "gaccelhints.ini"

	stdMargin     = 6
	defaultWidth  = 640
	defaultHeight = 480

	sigActivate  = "activate"
	sigConfigure = "configure-event"
	sigClicked   = "clicked"
	sigDestroy   = "destroy"
	sigToggled   = "toggled"

	propActive = "active"
)
