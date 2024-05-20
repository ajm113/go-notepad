package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	appName             string = "Go Notepad"
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

		config *ConfigSchema
	}
)

func (a *app) LoadConfig() {
	c, err := searchAndLoadConfig()
	if err != nil {
		a.UnexpectedErrorMessageBox("Unexpected error parsing config file: %s\n\nUsing defaults", err)
		c = &DefaultConfig
	}

	a.config = c
}

func (a *app) UnexpectedErrorMessageBox(format string, args ...interface{}) {
	d := gtk.MessageDialogNew(a.Win, gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "")
	d.FormatSecondaryText(format, args...)
	d.SetTitle(appName)
	d.Run()
	d.Destroy()
}

func (a *app) updateStatusBar() {
	if a.statusBar == nil || !a.menu.statusBarMenuItem.GetActive() {
		return
	}

	a.statusBar.SetText(fmt.Sprintf("col: %d | line: %d", a.lineOffsetCount+1, a.lineCount))
}

func (a *app) UpdateTitle() {
	title := filepath.Base(a.openedFilename) + " - " + appName
	a.Win.SetTitle(title)
}

func (a *app) LoadFile(filename string) {
	a.openedFilename = filename
	err := a.textView.LoadSource(filename)

	if err != nil {
		a.UnexpectedErrorMessageBox("Unexpected error loading file: %s\n\n%s", filename, err)
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

	a.textView.SetFont(a.config.Font.Family, a.config.Font.Size)

	if a.config.StatusBar.Enable {
		a.menu.statusBarMenuItem.SetActive(true)
		a.statusBar.Show()
		a.updateStatusBar()
	}

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

		a.updateStatusBar()
	})

	a.textView.GTKtextView.Connect("drag-data-received", func(tv *gtk.TextView, ctx *gdk.DragContext, x, y int, data *gtk.SelectionData, info uint, time uint32) {
		if a.hasChanges {
			response := a.displayUnsavedChangesMessagedialog()

			if response == gtk.RESPONSE_CANCEL || response == gtk.RESPONSE_DELETE_EVENT {
				return
			}
		}

		uris := string(data.GetData())
		// Split by new lines to handle multiple files
		for _, uri := range strings.Split(uris, "\n") {
			if uri != "" {
				filePath := strings.TrimPrefix(uri, "file://")
				filePath = filePath[:len(filePath)-1]
				a.LoadFile(filePath)
				break
			}
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

			err := a.textView.SaveSource(filename)

			if err != nil {
				a.UnexpectedErrorMessageBox("Unexpected error saving the file to disk!\n\n%s", err)
				return
			}

			a.openedFilename = filename
			a.hasChanges = false
			a.isFileOpened = true
			a.UpdateTitle()
		}
	})

	a.menu.saveMenuItem.Connect("activate", func() {
		if !a.isFileOpened {
			a.menu.saveAsMenuItem.Emit("activate", glib.TYPE_NONE)
			return
		}

		err := a.textView.SaveSource(a.openedFilename)
		if err != nil {
			a.UnexpectedErrorMessageBox("Unexpected error saving the file to disk!\n\n%s", err)
			return
		}

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
			a.updateStatusBar()
		} else {
			a.statusBar.Hide()
		}
	})

	a.menu.fontMenuItem.Connect("activate", func() {
		fd, err := gtk.FontChooserDialogNew(appName, a.Win)

		if err != nil {
			a.UnexpectedErrorMessageBox("Unexpected error creating font chooser dialog:\n\n%s", err)
			fmt.Printf("failed creating font chooser dialog: %s\n", err)
		}

		fd.SetFont(fmt.Sprintf("%s %d", a.config.Font.Family, a.config.Font.Size))

		fd.ShowAll()
		response := fd.Run()

		if response == gtk.RESPONSE_OK {
			fontText := fd.GetFont()

			fontTokens := strings.Split(fontText, " ")
			fontSize, err := strconv.Atoi(fontTokens[len(fontTokens)-1])

			if err != nil {
				a.UnexpectedErrorMessageBox("Unexpected error extracting font size:\n\n%s", err)
				fmt.Printf("failed selecting font: %s\n", err)
			}

			fontFamily := strings.Join(fontTokens[:len(fontTokens)-1], " ")
			fontFamily = strings.Trim(fontFamily, ",")

			err = a.textView.SetFont(fontFamily, int64(fontSize))
			if err != nil {
				a.UnexpectedErrorMessageBox("Unexpected error choosing font:\n\n%s", err)
			}
		}

		fd.Destroy()

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
				app.menu.saveMenuItem.Emit("activate", glib.TYPE_NONE)
				return true
			}
		}

		return false
	})

	app.Win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	app.LoadConfig()
	app.SetupWindow()
	app.SetupEvents()
	app.Init(os.Args)

	gtk.Main()
}
