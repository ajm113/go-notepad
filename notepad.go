package main

import (
	"log"

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
		openFilename string
		hasChanges   bool

		Win        *gtk.Window
		textView   *TextView
		accelGroup *gtk.AccelGroup
		statusBar  *Statusbar
		grid       *gtk.Grid
	}
)

func (app *App) UpdateTitle() {
	title := app.openFilename + " - " + AppName

	if app.hasChanges {
		title = "*" + title
	}

	app.Win.SetTitle(title)
}

func (app *App) SetupWindow() {
	var err error

	app.grid, err = gtk.GridNew()
	if err != nil {
		log.Fatal("unable to create grid:", err)
	}

	app.Win.Add(app.grid)
	app.grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	app.setupMenu(app.grid)
	app.textView = NewTextView(app)
	app.statusBar = NewStatusbar(app)

	app.UpdateTitle()
	app.Win.SetBorderWidth(2)
	app.Win.SetDefaultSize(DefaultWindowWidth, DefaultWindowHeight)
	app.Win.SetPosition(gtk.WIN_POS_CENTER)
	app.Win.ShowAll()
}

func main() {
	gtk.Init(nil)
	var err error

	app := &App{
		openFilename: DefaultFilename,
	}

	app.Win, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	if err != nil {
		log.Fatal("failed creating window:", err)
	}

	app.Win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	app.SetupWindow()

	gtk.Main()
}
