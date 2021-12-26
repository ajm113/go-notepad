package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const (
	appName             string = "Notepad"
	defaultFilename            = "Untitled"
	defaultWindowWidth  int    = 900
	defaultWindowHeight int    = 500
)

type (
	app struct {
		openedFilename  string
		hasChanges      bool
		isFileOpened    bool
		lineCount       int
		lineOffsetCount int

		Win        *gtk.Window
		textView   *textView
		menu       *menu
		accelGroup *gtk.AccelGroup
		statusBar  *statusbar
		grid       *gtk.Grid
	}
)

func (a *app) UpdateTitle() {
	title := filepath.Base(a.openedFilename) + " - " + appName

	if a.hasChanges {
		title = "*" + title
	}

	a.Win.SetTitle(title)
}

func (a *app) LoadFile(filename string) {
	a.openedFilename = filename
	err := a.textView.LoadSource(filename)

	if err != nil {
		d := gtk.MessageDialogNew(a.Win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "")
		d.FormatSecondaryText("Unexpected error saving file: %s", err)
		d.SetTitle(appName)
		d.Run()
		d.Destroy()
	}

	a.hasChanges = false
	a.isFileOpened = true
	a.UpdateTitle()
}

func (a *app) Init(args []string) {
	if len(args) > 1 && fileExist(args[1]) {
		a.LoadFile(args[1])
	}
}

func (a *app) SetupWindow() {
	var err error

	a.grid, err = gtk.GridNew()
	if err != nil {
		log.Fatal("unable to create grid:", err)
	}

	a.Win.Add(a.grid)
	a.grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	a.menu = newMenu(a)
	a.textView = newTextView(a)
	a.statusBar = newStatusbar(a)

	a.UpdateTitle()
	a.Win.SetBorderWidth(2)
	a.Win.SetDefaultSize(defaultWindowWidth, defaultWindowHeight)
	a.Win.SetPosition(gtk.WIN_POS_CENTER)
	a.Win.ShowAll()
}

func (a *app) displayUnsavedChangesMessagedialog() (response gtk.ResponseType) {
	d := gtk.MessageDialogNew(
		a.Win,
		gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_WARNING,
		gtk.BUTTONS_YES_NO,
		"The text in the %s file has changed.",
		a.openedFilename,
	)
	d.FormatSecondaryText("Do you want to save the changes?")
	d.SetTitle(appName)
	response = d.Run()
	d.Destroy()

	return
}

func (a *app) SetupEvents() {
	tb, _ := a.textView.GTKtextView.GetBuffer()
	tb.Connect("mark-set", func(tb *gtk.TextBuffer, itr *gtk.TextIter) {
		if tb.GetHasSelection() {
			a.menu.cutMenuItem.SetSensitive(true)
			a.menu.copyMenuItem.SetSensitive(true)
			a.menu.deleteMenuItem.SetSensitive(true)
		} else {
			a.menu.cutMenuItem.SetSensitive(false)
			a.menu.copyMenuItem.SetSensitive(false)
			a.menu.deleteMenuItem.SetSensitive(false)
		}

		a.lineCount = itr.GetLine() + 1
		a.lineOffsetCount = itr.GetLineOffset()

		if a.menu.statusBarMenuItem.GetActive() {
			a.statusBar.SetText(fmt.Sprintf("col: %d | line: %d", a.lineOffsetCount, a.lineCount))
		}
	})

	tb.Connect("changed", func(tb *gtk.TextBuffer) {
		a.hasChanges = true
		a.UpdateTitle()
	})

	a.menu.openMenuItem.Connect("activate", func() {
		file := gtk.OpenFileChooserNative("Open File", a.Win)

		if file != nil {
			a.LoadFile(*file)
		}
	})

	a.menu.newMenuItem.Connect("activate", func() {
		if a.hasChanges {
			response := a.displayUnsavedChangesMessagedialog()

			if response == gtk.RESPONSE_CANCEL || response == gtk.RESPONSE_DELETE_EVENT {
				return
			}
		}

		a.openedFilename = defaultFilename
		a.textView.Clear()
		a.hasChanges = false
		a.isFileOpened = false
		a.UpdateTitle()
	})

	a.menu.saveAsMenuItem.Connect("activate", func() {
		fc, _ := gtk.FileChooserNativeDialogNew("Save As...", a.Win, gtk.FILE_CHOOSER_ACTION_SAVE, "Save", "Cancel")
		response := fc.Run()
		fc.Destroy()

		filename := fc.GetFilename()

		if response == int(gtk.RESPONSE_ACCEPT) {
			if fileExist(filename) {
				d := gtk.MessageDialogNew(a.Win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK_CANCEL, "")
				d.FormatSecondaryText("You are about to write to a already saved file! Are you sure you wish to do this?")
				d.SetTitle(appName)
				response := d.Run()
				d.Destroy()

				if response == gtk.RESPONSE_CANCEL || response == gtk.RESPONSE_DELETE_EVENT {
					return
				}
			}

			a.textView.SaveSource(filename)
			a.openedFilename = filename
			a.hasChanges = false
			a.isFileOpened = true
			a.UpdateTitle()
		}
	})

	a.menu.saveMenuItem.Connect("activate", func() {
		if !a.isFileOpened {
			a.menu.saveAsMenuItem.Emit("activate")
			return
		}

		a.textView.SaveSource(a.openedFilename)
		a.hasChanges = false
		a.UpdateTitle()
	})

	a.menu.wordWrapMenuItem.Connect("activate", func() {
		if a.menu.wordWrapMenuItem.GetActive() {
			a.textView.WrapText(true)
		} else {
			a.textView.WrapText(false)
		}
	})

	a.menu.statusBarMenuItem.Connect("activate", func() {
		if a.menu.statusBarMenuItem.GetActive() {
			a.statusBar.Show()
			a.statusBar.SetText(fmt.Sprintf("col: %d | line: %d", a.lineOffsetCount, a.lineCount))
		} else {
			a.statusBar.Hide()
		}
	})

	a.menu.fontMenuItem.Connect("activate", func() {

	})

	a.menu.timedateMenuItem.Connect("activate", func() {
		a.textView.InsertTimestamp()
	})

	// Handle on-close events.
	a.menu.exitMenuItem.Connect("activate", func() {
		a.Win.Close()
	})

	a.accelGroup.Connect(gdk.KEY_W, gdk.CONTROL_MASK, 0, func() {
		a.Win.Close()
	})

	a.menu.aboutMenuItem.Connect("activate", func() {
		displayAboutDialog(a)
	})
}

func main() {
	gtk.Init(nil)
	var err error

	app := &app{
		openedFilename: defaultFilename,
	}

	app.Win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	if err != nil {
		log.Fatal("failed creating window:", err)
	}

	app.Win.Connect("delete-event", func() bool {
		if app.hasChanges {
			switch app.displayUnsavedChangesMessagedialog() {
			case gtk.RESPONSE_NO, gtk.RESPONSE_DELETE_EVENT:
				return false
			case gtk.RESPONSE_YES:
				app.menu.saveMenuItem.Emit("activate")
				return true
			}
		}

		return false
	})

	app.Win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	app.SetupWindow()
	app.SetupEvents()
	app.Init(os.Args)

	gtk.Main()
}
