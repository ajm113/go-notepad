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
	AppName             string = "Notepad"
	DefaultFilename            = "Untitled"
	DefaultWindowWidth  int    = 900
	DefaultWindowHeight int    = 500
)

type (
	App struct {
		openedFilename  string
		hasChanges      bool
		isFileOpened    bool
		lineCount       int
		lineOffsetCount int

		Win        *gtk.Window
		TextView   *TextView
		menu       *Menu
		accelGroup *gtk.AccelGroup
		statusBar  *Statusbar
		grid       *gtk.Grid
	}
)

func (app *App) UpdateTitle() {
	title := filepath.Base(app.openedFilename) + " - " + AppName

	if app.hasChanges {
		title = "*" + title
	}

	app.Win.SetTitle(title)
}

func (app *App) LoadFile(filename string) {
	app.openedFilename = filename
	err := app.TextView.LoadSource(filename)

	if err != nil {
		d := gtk.MessageDialogNew(app.Win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "")
		d.FormatSecondaryText("Unexpected error saving file: %s", err)
		d.SetTitle(AppName)
		d.Run()
		d.Destroy()
	}

	app.hasChanges = false
	app.isFileOpened = true
	app.UpdateTitle()
}

func (app *App) Init(args []string) {
	if len(args) > 1 && fileExist(args[1]) {
		app.LoadFile(args[1])
	}
}

func (app *App) SetupWindow() {
	var err error

	app.grid, err = gtk.GridNew()
	if err != nil {
		log.Fatal("unable to create grid:", err)
	}

	app.Win.Add(app.grid)
	app.grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	app.menu = NewMenu(app)
	app.TextView = NewTextView(app)
	app.statusBar = NewStatusbar(app)

	app.UpdateTitle()
	app.Win.SetBorderWidth(2)
	app.Win.SetDefaultSize(DefaultWindowWidth, DefaultWindowHeight)
	app.Win.SetPosition(gtk.WIN_POS_CENTER)
	app.Win.ShowAll()
}

func (app *App) displayUnsavedChangesMessagedialog() (response gtk.ResponseType) {
	d := gtk.MessageDialogNew(
		app.Win,
		gtk.DIALOG_DESTROY_WITH_PARENT,
		gtk.MESSAGE_WARNING,
		gtk.BUTTONS_YES_NO,
		"The text in the %s file has changed.",
		app.openedFilename,
	)
	d.FormatSecondaryText("Do you want to save the changes?")
	d.SetTitle(AppName)
	response = d.Run()
	d.Destroy()

	return
}

func (app *App) SetupEvents() {
	tb, _ := app.TextView.GTKtextView.GetBuffer()
	tb.Connect("mark-set", func(tb *gtk.TextBuffer, itr *gtk.TextIter) {
		if tb.GetHasSelection() {
			app.menu.cutMenuItem.SetSensitive(true)
			app.menu.copyMenuItem.SetSensitive(true)
			app.menu.deleteMenuItem.SetSensitive(true)
		} else {
			app.menu.cutMenuItem.SetSensitive(false)
			app.menu.copyMenuItem.SetSensitive(false)
			app.menu.deleteMenuItem.SetSensitive(false)
		}

		app.lineCount = itr.GetLine() + 1
		app.lineOffsetCount = itr.GetLineOffset()

		if app.menu.statusBarMenuItem.GetActive() {
			app.statusBar.SetText(fmt.Sprintf("col: %d | line: %d", app.lineOffsetCount, app.lineCount))
		}
	})

	tb.Connect("changed", func(tb *gtk.TextBuffer) {
		app.hasChanges = true
		app.UpdateTitle()
	})

	app.menu.openMenuItem.Connect("activate", func() {
		file := gtk.OpenFileChooserNative("Open File", app.Win)

		if file != nil {
			app.LoadFile(*file)
		}
	})

	app.menu.newMenuItem.Connect("activate", func() {
		if app.hasChanges {
			response := app.displayUnsavedChangesMessagedialog()

			if response == gtk.RESPONSE_CANCEL || response == gtk.RESPONSE_DELETE_EVENT {
				return
			}
		}

		app.openedFilename = DefaultFilename
		app.TextView.Clear()
		app.hasChanges = false
		app.isFileOpened = false
		app.UpdateTitle()
	})

	app.menu.saveAsMenuItem.Connect("activate", func() {
		fc, _ := gtk.FileChooserNativeDialogNew("Save As...", app.Win, gtk.FILE_CHOOSER_ACTION_SAVE, "Save", "Cancel")
		response := fc.Run()
		fc.Destroy()

		filename := fc.GetFilename()

		if response == int(gtk.RESPONSE_ACCEPT) {
			if fileExist(filename) {
				d := gtk.MessageDialogNew(app.Win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK_CANCEL, "")
				d.FormatSecondaryText("You are about to write to a already saved file! Are you sure you wish to do this?")
				d.SetTitle(AppName)
				response := d.Run()
				d.Destroy()

				if response == gtk.RESPONSE_CANCEL || response == gtk.RESPONSE_DELETE_EVENT {
					return
				}
			}

			app.TextView.SaveSource(filename)
			app.openedFilename = filename
			app.hasChanges = false
			app.isFileOpened = true
			app.UpdateTitle()
		}
	})

	app.menu.saveMenuItem.Connect("activate", func() {
		if !app.isFileOpened {
			app.menu.saveAsMenuItem.Emit("activate")
			return
		}

		app.TextView.SaveSource(app.openedFilename)
		app.hasChanges = false
		app.UpdateTitle()
	})

	app.menu.wordWrapMenuItem.Connect("activate", func() {
		if app.menu.wordWrapMenuItem.GetActive() {
			app.TextView.WrapText(true)
		} else {
			app.TextView.WrapText(false)
		}
	})

	app.menu.statusBarMenuItem.Connect("activate", func() {
		if app.menu.statusBarMenuItem.GetActive() {
			app.statusBar.Show()
			app.statusBar.SetText(fmt.Sprintf("col: %d | line: %d", app.lineOffsetCount, app.lineCount))
		} else {
			app.statusBar.Hide()
		}
	})

	app.menu.fontMenuItem.Connect("activate", func() {

	})

	app.menu.timedateMenuItem.Connect("activate", func() {
		app.TextView.InsertTimestamp()
	})

	// Handle on-close events.
	app.menu.exitMenuItem.Connect("activate", func() {
		app.Win.Close()
	})

	app.accelGroup.Connect(gdk.KEY_W, gdk.CONTROL_MASK, 0, func() {
		app.Win.Close()
	})
}

func main() {
	gtk.Init(nil)
	var err error

	app := &App{
		openedFilename: DefaultFilename,
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
