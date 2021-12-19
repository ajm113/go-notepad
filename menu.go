package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func (app *App) setupMenu(grid *gtk.Grid) {
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	grid.Add(vbox)

	menubar, err := gtk.MenuBarNew()
	if err != nil {
		log.Fatal("unable to create menubar:", err)
	}

	// Define accel group to define our hotkeys.
	app.accelGroup, err = gtk.AccelGroupNew()

	if err != nil {
		log.Fatal("unable to create accelGroup:", err)
	}

	app.Win.AddAccelGroup(app.accelGroup)

	app.setupFileMenu(menubar)
	app.setupEditMenu(menubar)
	app.setupViewMenu(menubar)
	app.setupFormatMenu(menubar)
	app.setupHelpMenu(menubar)

	menubar.SetHExpand(true)
	vbox.PackStart(menubar, true, true, 0)
}

func (app *App) setupFileMenu(menubar *gtk.MenuBar) {
	fileMenu, _ := gtk.MenuNew()
	fileMain, _ := gtk.MenuItemNewWithLabel("File")

	newMi, _ := gtk.MenuItemNewWithLabel("New")
	key, mod := gtk.AcceleratorParse("<Control>N")
	newMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	openMi, _ := gtk.MenuItemNewWithLabel("Open...")
	key, mod = gtk.AcceleratorParse("<Control>O")
	openMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	saveMi, _ := gtk.MenuItemNewWithLabel("Save")
	key, mod = gtk.AcceleratorParse("<Control>S")
	saveMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	saveAsMi, _ := gtk.MenuItemNewWithLabel("Save As...")

	pageSetupMi, _ := gtk.MenuItemNewWithLabel("Page Setup...")
	printMi, _ := gtk.MenuItemNewWithLabel("Print...")
	key, mod = gtk.AcceleratorParse("<Control>P")
	printMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	exitMi, _ := gtk.MenuItemNewWithLabel("Exit")

	sepMi1, _ := gtk.SeparatorMenuItemNew()
	sepMi2, _ := gtk.SeparatorMenuItemNew()

	fileMain.SetSubmenu(fileMenu)
	fileMenu.Append(newMi)
	fileMenu.Append(openMi)
	fileMenu.Append(saveMi)
	fileMenu.Append(saveAsMi)
	fileMenu.Append(sepMi1)
	fileMenu.Append(pageSetupMi)
	fileMenu.Append(printMi)
	fileMenu.Append(sepMi2)
	fileMenu.Append(exitMi)

	menubar.Append(fileMain)

	exitMi.Connect("button-press-event", func() {
		app.Win.Emit("destroy")
	})

}

func (app *App) setupEditMenu(menubar *gtk.MenuBar) {
	editMenu, _ := gtk.MenuNew()
	editMain, _ := gtk.MenuItemNewWithLabel("Edit")

	undoMi, _ := gtk.MenuItemNewWithLabel("Undo")
	key, mod := gtk.AcceleratorParse("<Control>Z")
	undoMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	sepMi1, _ := gtk.SeparatorMenuItemNew()
	cutMi, _ := gtk.MenuItemNewWithLabel("Cut")
	key, mod = gtk.AcceleratorParse("<Control>X")
	cutMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	cutMi.Connect("activate", func() {
		app.textView.Cut()
	})

	copyMi, _ := gtk.MenuItemNewWithLabel("Copy")
	key, mod = gtk.AcceleratorParse("<Control>C")
	copyMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	copyMi.Connect("activate", func() {
		app.textView.Copy()
	})

	pasteMi, _ := gtk.MenuItemNewWithLabel("Paste")
	key, mod = gtk.AcceleratorParse("<Control>V")
	pasteMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	pasteMi.Connect("activate", func() {
		app.textView.Paste()
	})

	deleteMi, _ := gtk.MenuItemNewWithLabel("Delete")
	deleteMi.Connect("activate", func() {
		app.textView.Backspace()
	})

	sepMi2, _ := gtk.SeparatorMenuItemNew()

	findMi, _ := gtk.MenuItemNewWithLabel("Find...")
	key, mod = gtk.AcceleratorParse("<Control>F")
	findMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	findNextMi, _ := gtk.MenuItemNewWithLabel("Find Next")
	key, mod = gtk.AcceleratorParse("F3")
	findNextMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	replaceMi, _ := gtk.MenuItemNewWithLabel("Replace...")
	key, mod = gtk.AcceleratorParse("<Control>H")
	replaceMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	goToMi, _ := gtk.MenuItemNewWithLabel("Go To...")
	key, mod = gtk.AcceleratorParse("<Control>G")
	goToMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	sepMi3, _ := gtk.SeparatorMenuItemNew()

	selectAllMi, _ := gtk.MenuItemNewWithLabel("Select All")
	key, mod = gtk.AcceleratorParse("<Control>A")
	selectAllMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	selectAllMi.Connect("activate", func() {
		app.textView.SelectAll()
	})

	timeDateMi, _ := gtk.MenuItemNewWithLabel("Time/Date")

	editMain.SetSubmenu(editMenu)
	editMenu.Append(undoMi)
	editMenu.Append(sepMi1)
	editMenu.Append(cutMi)
	editMenu.Append(copyMi)
	editMenu.Append(pasteMi)
	editMenu.Append(deleteMi)
	editMenu.Append(sepMi2)

	editMenu.Append(findMi)
	editMenu.Append(findNextMi)
	editMenu.Append(replaceMi)
	editMenu.Append(goToMi)
	editMenu.Append(sepMi3)
	editMenu.Append(selectAllMi)
	editMenu.Append(timeDateMi)

	menubar.Append(editMain)
}

func (app *App) setupFormatMenu(menubar *gtk.MenuBar) {
	formatMenu, _ := gtk.MenuNew()
	formatMain, _ := gtk.MenuItemNewWithLabel("Format")

	wordWrapMi, _ := gtk.MenuItemNewWithLabel("Word Wrap")
	fontMi, _ := gtk.MenuItemNewWithLabel("Font...")

	formatMain.SetSubmenu(formatMenu)
	formatMenu.Append(wordWrapMi)
	formatMenu.Append(fontMi)

	menubar.Append(formatMain)
}

func (app *App) setupViewMenu(menubar *gtk.MenuBar) {
	viewMenu, _ := gtk.MenuNew()
	viewMain, _ := gtk.MenuItemNewWithLabel("View")

	statusBarMi, _ := gtk.MenuItemNewWithLabel("Status Bar")

	viewMain.SetSubmenu(viewMenu)
	viewMenu.Append(statusBarMi)

	menubar.Append(viewMain)
}

func (app *App) setupHelpMenu(menubar *gtk.MenuBar) {
	helpMenu, _ := gtk.MenuNew()
	helpMain, _ := gtk.MenuItemNewWithLabel("Help")

	aboutMi, _ := gtk.MenuItemNewWithLabel("About")
	key, mod := gtk.AcceleratorParse("F1")
	aboutMi.AddAccelerator("activate", app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	helpMain.SetSubmenu(helpMenu)
	helpMenu.Append(aboutMi)

	menubar.Append(helpMain)
}
