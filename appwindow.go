// Copyright © 2022-23 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mark-summerfield/accelhint"
	"github.com/mark-summerfield/gong"
	"golang.org/x/exp/slices"
)

type AppWindow struct {
	config        *Config
	application   *gtk.Application
	window        *gtk.ApplicationWindow
	container     *gtk.Widget
	toolFrame     *gtk.Frame
	copyButton    *gtk.ToolButton
	cutButton     *gtk.ToolButton
	pasteButton   *gtk.ToolButton
	aboutButton   *gtk.ToolButton
	quitButton    *gtk.ToolButton
	originalLabel *gtk.Label
	originalText  *gtk.TextView
	hintedLabel   *gtk.Label
	hintedText    *gtk.TextView
	alphabetLabel *gtk.Label
	alphabetEntry *gtk.Entry
	statusLabel   *gtk.Label
}

func newAppWindow(app *gtk.Application, config *Config) *AppWindow {
	appWindow := &AppWindow{application: app}
	appWindow.config = config
	appWindow.makeWidgets()
	appWindow.makeLayout()
	appWindow.makeConnections()
	appWindow.window.SetTitle(appName)
	if img := getPixbuf(icon, 0); img != nil {
		appWindow.window.SetIcon(img)
	}
	appWindow.window.SetBorderWidth(stdMargin)
	appWindow.window.Add(appWindow.container)
	appWindow.originalText.GrabFocus()
	return appWindow
}

func (me *AppWindow) makeWidgets() {
	var err error
	me.window, err = gtk.ApplicationWindowNew(me.application)
	gong.CheckError("Failed to create window:", err)
	me.makeToolbar()
	me.originalLabel, err = gtk.LabelNewWithMnemonic("_Original")
	gong.CheckError("Failed to create label:", err)
	me.originalText, err = gtk.TextViewNew()
	gong.CheckError("Failed to create text view:", err)
	me.originalText.SetAcceptsTab(false)
	me.originalLabel.SetMnemonicWidget(me.originalText)
	me.originalText.SetHExpand(true)
	me.originalText.SetVExpand(true)
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	buffer.SetText(defaultOriginal)
	start, end := buffer.GetBounds()
	buffer.SelectRange(start, end)
	me.hintedLabel, err = gtk.LabelNewWithMnemonic("_Hinted")
	gong.CheckError("Failed to create label:", err)
	me.hintedText, err = gtk.TextViewNew()
	gong.CheckError("Failed to create text view:", err)
	me.hintedLabel.SetMnemonicWidget(me.hintedText)
	me.hintedText.SetAcceptsTab(false)
	me.hintedText.SetEditable(false)
	me.hintedText.SetHExpand(true)
	me.hintedText.SetVExpand(true)
	me.alphabetLabel, err = gtk.LabelNewWithMnemonic("_Alphabet")
	gong.CheckError("Failed to create label:", err)
	me.alphabetEntry, err = gtk.EntryNew()
	me.alphabetEntry.SetHExpand(true)
	gong.CheckError("Failed to create entry:", err)
	me.alphabetLabel.SetMnemonicWidget(me.alphabetEntry)
	me.alphabetEntry.SetText(me.config.alphabet)
	me.statusLabel, err = gtk.LabelNew("")
	gong.CheckError("Failed to create label:", err)
	me.statusLabel.SetXAlign(0.0)
}

func (me *AppWindow) makeToolbar() {
	var err error
	me.toolFrame, err = gtk.FrameNew("")
	gong.CheckError("Failed to create frame:", err)
	me.copyButton = makeToolbutton(iconCopy, "Copy", "<b>Copy</b> Ctrl+C")
	me.cutButton = makeToolbutton(iconCut, "Cut", "<b>Cut</b> Ctrl+X")
	me.pasteButton = makeToolbutton(iconPaste, "Paste",
		"<b>Paste</b> Ctrl+V")
	me.aboutButton = makeToolbutton(icon, "About", "<b>About</b>")
	me.quitButton = makeToolbutton(iconQuit, "Quit",
		"<b>Quit</b> Esc <i>or</i> Ctrl+Q")
}

func (me *AppWindow) makeLayout() {
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	me.makeToolbarLayout(vbox)
	left, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	left.PackStart(me.originalLabel, false, false, stdMargin)
	scroller, err := gtk.ScrolledWindowNew(nil, nil)
	gong.CheckError("Failed to create scroller:", err)
	scroller.Add(me.originalText)
	left.PackStart(scroller, true, true, stdMargin)
	right, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, stdMargin)
	gong.CheckError("Failed to create vbox:", err)
	right.PackStart(me.hintedLabel, false, false, stdMargin)
	scroller, err = gtk.ScrolledWindowNew(nil, nil)
	gong.CheckError("Failed to create scroller:", err)
	scroller.Add(me.hintedText)
	right.PackStart(scroller, true, true, stdMargin)
	hbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	hbox.PackStart(left, true, true, stdMargin)
	hbox.PackStart(right, true, true, stdMargin)
	vbox.PackStart(hbox, true, true, stdMargin)
	hbox, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	hbox.PackStart(me.alphabetLabel, false, false, stdMargin)
	hbox.PackStart(me.alphabetEntry, true, true, stdMargin)
	vbox.PackStart(hbox, false, false, stdMargin)
	vbox.PackStart(me.statusLabel, false, false, stdMargin)
	me.container = &vbox.Container.Widget
}

func (me *AppWindow) makeToolbarLayout(vbox *gtk.Box) {
	toolbox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, stdMargin)
	gong.CheckError("Failed to create hbox:", err)
	toolbox.PackStart(me.copyButton, false, false, 0)
	toolbox.PackStart(me.cutButton, false, false, 0)
	toolbox.PackStart(me.pasteButton, false, false, 0)
	toolbox.PackStart(me.aboutButton, false, false, 0)
	toolbox.PackEnd(me.quitButton, false, false, 0)
	me.toolFrame.Add(toolbox)
	vbox.PackStart(me.toolFrame, false, false, stdMargin)
}

func (me *AppWindow) makeConnections() {
	me.window.Connect(sigConfigure, func(_ *gtk.ApplicationWindow,
		rawEvent *gdk.Event) bool {
		event := gdk.EventConfigureNewFromEvent(rawEvent)
		me.config.updateGeometry(event.X(), event.Y(), event.Width(),
			event.Height())
		return false
	})
	me.window.Connect(sigKeyPress, func(_ *gtk.ApplicationWindow,
		event *gdk.Event) {
		me.onKeyPress(event)
	})
	me.window.Connect(sigDestroy, func(_ *gtk.ApplicationWindow) {
		me.onQuit()
	})
	me.copyButton.Connect(sigClicked, func() { me.onCopy() })
	me.cutButton.Connect(sigClicked, func() { me.onCut() })
	me.pasteButton.Connect(sigClicked, func() { me.onPaste() })
	me.aboutButton.Connect(sigClicked, func() { me.onAbout() })
	me.quitButton.Connect(sigClicked, func() { me.onQuit() })
	me.alphabetEntry.Connect(sigChanged, func(_ *gtk.Entry) {
		me.onTextChanged()
	})
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	buffer.Connect(sigChanged, func(_ *gtk.TextBuffer) {
		me.onTextChanged()
	})
}

func (me *AppWindow) onTextChanged() {
	alphabet, err := me.alphabetEntry.GetText()
	gong.CheckError("Failed to get text:", err)
	if len(alphabet) < 5 {
		alphabet = accelhint.Alphabet
	} else {
		alphabet = strings.ToUpper(strings.TrimSpace(alphabet))
	}
	me.alphabetEntry.SetText(alphabet)
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	start, end := buffer.GetBounds()
	text, err := buffer.GetText(start, end, false)
	gong.CheckError("Failed to get text buffer's text:", err)
	buffer, err = me.hintedText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	start, end = buffer.GetBounds()
	buffer.Delete(start, end)
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if text == "" {
		me.statusLabel.SetMarkup("Enter items…")
		return
	}
	hinted, n, err := accelhint.HintedX(lines, '&', alphabet)
	if err != nil {
		me.statusLabel.SetMarkup(fmt.Sprintf(
			"<span color='darkred'>Failed to set accelerators:</span> "+
				"<span color='red'>%s</span>", err))
	} else {
		me.displayHinted(lines, hinted, n, buffer)
	}
}

func (me *AppWindow) displayHinted(lines, hinted []string, n int,
	buffer *gtk.TextBuffer) {
	for i := 0; i < len(hinted); i++ {
		chars := []rune(strings.ReplaceAll(hinted[i], escMarker,
			placeholder))
		j := slices.Index(chars, rune('&'))
		if j > -1 && j+1 < len(chars) {
			left := strings.ReplaceAll(string(chars[:j]), placeholder,
				escMarker)
			accel := string(chars[j+1])
			right := ""
			if j+2 < len(chars) {
				right = strings.ReplaceAll(string(chars[j+2:]),
					placeholder, escMarker)
			}
			start := buffer.GetEndIter()
			buffer.InsertMarkup(start, fmt.Sprintf(
				"<span color='gray'>%d</span>\t", j))
			buffer.InsertAtCursor(left)
			start = buffer.GetEndIter()
			buffer.InsertMarkup(start, fmt.Sprintf(
				"<span color='blue' underline='single'>%s</span>",
				accel))
			buffer.InsertAtCursor(right)
			buffer.InsertAtCursor("\n")
		} else {
			buffer.InsertAtCursor("\t" + strings.ReplaceAll(
				string(chars), placeholder, escMarker) + "\n")
		}
	}
	me.statusLabel.SetMarkup(fmt.Sprintf(
		"<span color='darkgreen'>Hinted — %d/%d — %.0f%%</span>", n,
		len(lines), (float64(n) / float64(len(lines)) * 100.0)))
}

func (me *AppWindow) onKeyPress(event *gdk.Event) {
	keyEvent := &gdk.EventKey{Event: event}
	keyVal := keyEvent.KeyVal()
	if (keyEvent.State() & gdk.CONTROL_MASK) != 0 {
		switch keyVal {
		case gdk.KEY_Q, gdk.KEY_q:
			me.onQuit()
		}
	} else if keyVal == gdk.KEY_Escape {
		me.onQuit()
	}
}

func (me *AppWindow) onCopy() {
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	if clipboard := getClipboard(); clipboard != nil {
		buffer.CopyClipboard(clipboard)
	}
}

func (me *AppWindow) onCut() {
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	if clipboard := getClipboard(); clipboard != nil {
		buffer.CutClipboard(clipboard, true)
	}
}

func (me *AppWindow) onPaste() {
	buffer, err := me.originalText.GetBuffer()
	gong.CheckError("Failed to get text buffer:", err)
	if clipboard := getClipboard(); clipboard != nil {
		pos := buffer.GetIterAtMark(buffer.GetInsert())
		buffer.PasteClipboard(clipboard, pos, true)
	}
}

func (me *AppWindow) onAbout() {
	dlg, err := gtk.AboutDialogNew()
	gong.CheckError("Failed to create about dialog:", err)
	dlg.SetTitle(appName) // doesn't work
	dlg.SetProgramName(appName)
	dlg.SetVersion(strings.TrimSpace(Version))
	dlg.SetWebsite("https://github.com/mark-summerfield/gaccelhints")
	dlg.SetWebsiteLabel("github gaccelhints")
	dlg.SetComments("An application for efficiently calculating keyboard " +
		"accelerators.\n" + libInfo())
	year := time.Now().Year()
	yearStr := "2022"
	if year > 2022 {
		yearStr = fmt.Sprintf("%s-%d", yearStr, year-2000)
	}
	dlg.SetCopyright(fmt.Sprintf(
		"Copyright © %s Mark Summerfield. All rights reserved.", yearStr))
	dlg.SetLicenseType(gtk.LICENSE_GPL_3_0)
	if img := getPixbuf(icon, 128); img != nil {
		dlg.SetLogo(img)
	}
	dlg.ShowAll()
}

func (me *AppWindow) onQuit() {
	if me.config.dirty {
		me.config.save()
	}
	me.application.Quit()
}
