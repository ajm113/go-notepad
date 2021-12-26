package main

import (
	"log"
	"strconv"

	"github.com/gotk3/gotk3/gtk"
)

type (
	menu struct {
		app        *app
		gtkmenuBar *gtk.MenuBar

		newMenuItem    *gtk.MenuItem
		openMenuItem   *gtk.MenuItem
		saveMenuItem   *gtk.MenuItem
		saveAsMenuItem *gtk.MenuItem
		exitMenuItem   *gtk.MenuItem

		undoMenuItem     *gtk.MenuItem
		cutMenuItem      *gtk.MenuItem
		copyMenuItem     *gtk.MenuItem
		pasteMenuItem    *gtk.MenuItem
		deleteMenuItem   *gtk.MenuItem
		timedateMenuItem *gtk.MenuItem

		wordWrapMenuItem  *gtk.CheckMenuItem
		statusBarMenuItem *gtk.CheckMenuItem

		fontMenuItem *gtk.MenuItem

		aboutMenuItem *gtk.MenuItem
	}
)

func newMenu(app *app) *menu {
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	app.grid.Add(vbox)

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

	m := &menu{
		app:        app,
		gtkmenuBar: menubar,
	}

	m.setupFileMenu()
	m.setupEditMenu()
	m.setupViewMenu()
	m.setupFormatMenu()
	m.setupHelpMenu()

	menubar.SetHExpand(true)
	vbox.PackStart(menubar, true, true, 0)

	return m
}

func (m *menu) setupFileMenu() {
	fileMenu, _ := gtk.MenuNew()
	fileMain, _ := gtk.MenuItemNewWithLabel("File")

	m.newMenuItem, _ = gtk.MenuItemNewWithLabel("New")
	key, mod := gtk.AcceleratorParse("<Control>N")
	m.newMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	m.openMenuItem, _ = gtk.MenuItemNewWithLabel("Open...")
	key, mod = gtk.AcceleratorParse("<Control>O")
	m.openMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	m.saveMenuItem, _ = gtk.MenuItemNewWithLabel("Save")
	key, mod = gtk.AcceleratorParse("<Control>S")
	m.saveMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	m.saveAsMenuItem, _ = gtk.MenuItemNewWithLabel("Save As...")

	pageSetupMi, _ := gtk.MenuItemNewWithLabel("Page Setup...")
	printMi, _ := gtk.MenuItemNewWithLabel("Print...")
	key, mod = gtk.AcceleratorParse("<Control>P")
	printMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	m.exitMenuItem, _ = gtk.MenuItemNewWithLabel("Exit")

	sepMi1, _ := gtk.SeparatorMenuItemNew()
	sepMi2, _ := gtk.SeparatorMenuItemNew()

	fileMain.SetSubmenu(fileMenu)
	fileMenu.Append(m.newMenuItem)
	fileMenu.Append(m.openMenuItem)
	fileMenu.Append(m.saveMenuItem)
	fileMenu.Append(m.saveAsMenuItem)
	fileMenu.Append(sepMi1)
	fileMenu.Append(pageSetupMi)
	fileMenu.Append(printMi)
	fileMenu.Append(sepMi2)
	fileMenu.Append(m.exitMenuItem)

	m.gtkmenuBar.Append(fileMain)
}

func (m *menu) setupEditMenu() {
	editMenu, _ := gtk.MenuNew()
	editMain, _ := gtk.MenuItemNewWithLabel("Edit")

	m.undoMenuItem, _ = gtk.MenuItemNewWithLabel("Undo")
	key, mod := gtk.AcceleratorParse("<Control>Z")
	m.undoMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	sepMi1, _ := gtk.SeparatorMenuItemNew()
	m.cutMenuItem, _ = gtk.MenuItemNewWithLabel("Cut")
	key, mod = gtk.AcceleratorParse("<Control>X")
	m.cutMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	m.cutMenuItem.Connect("activate", func() {
		m.app.textView.Cut()
	})

	m.copyMenuItem, _ = gtk.MenuItemNewWithLabel("Copy")
	key, mod = gtk.AcceleratorParse("<Control>C")
	m.copyMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	m.copyMenuItem.Connect("activate", func() {
		m.app.textView.Copy()
	})

	m.pasteMenuItem, _ = gtk.MenuItemNewWithLabel("Paste")
	key, mod = gtk.AcceleratorParse("<Control>V")
	m.pasteMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	m.pasteMenuItem.Connect("activate", func() {
		m.app.textView.Paste()
	})

	m.deleteMenuItem, _ = gtk.MenuItemNewWithLabel("Delete")
	m.deleteMenuItem.Connect("activate", func() {
		m.app.textView.Backspace()
	})

	sepMi2, _ := gtk.SeparatorMenuItemNew()

	findMi, _ := gtk.MenuItemNewWithLabel("Find...")
	key, mod = gtk.AcceleratorParse("<Control>F")
	findMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	findNextMi, _ := gtk.MenuItemNewWithLabel("Find Next")
	key, mod = gtk.AcceleratorParse("F3")
	findNextMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	replaceMi, _ := gtk.MenuItemNewWithLabel("Replace...")
	key, mod = gtk.AcceleratorParse("<Control>H")
	replaceMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	goToMi, _ := gtk.MenuItemNewWithLabel("Go To...")
	key, mod = gtk.AcceleratorParse("<Control>G")
	goToMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	goToMi.Connect("activate", func() {
		response, line := displayGotoLine(m.app)

		if response == gtk.RESPONSE_OK {
			i, _ := strconv.ParseInt(line, 10, 64)
			m.app.textView.GoToLine(int(i - 1))
		}
	})

	sepMi3, _ := gtk.SeparatorMenuItemNew()

	selectAllMi, _ := gtk.MenuItemNewWithLabel("Select All")
	key, mod = gtk.AcceleratorParse("<Control>A")
	selectAllMi.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)
	selectAllMi.Connect("activate", func() {
		m.app.textView.SelectAll()
	})

	m.timedateMenuItem, _ = gtk.MenuItemNewWithLabel("Time/Date")
	key, mod = gtk.AcceleratorParse("F5")
	m.timedateMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	editMain.SetSubmenu(editMenu)
	editMenu.Append(m.undoMenuItem)
	editMenu.Append(sepMi1)
	editMenu.Append(m.cutMenuItem)
	editMenu.Append(m.copyMenuItem)
	editMenu.Append(m.pasteMenuItem)
	editMenu.Append(m.deleteMenuItem)
	editMenu.Append(sepMi2)

	editMenu.Append(findMi)
	editMenu.Append(findNextMi)
	editMenu.Append(replaceMi)
	editMenu.Append(goToMi)
	editMenu.Append(sepMi3)
	editMenu.Append(selectAllMi)
	editMenu.Append(m.timedateMenuItem)

	m.gtkmenuBar.Append(editMain)

	// Setup signals from our textView
	m.cutMenuItem.SetSensitive(false)
	m.copyMenuItem.SetSensitive(false)
	m.deleteMenuItem.SetSensitive(false)

}

func (m *menu) setupFormatMenu() {
	formatMenu, _ := gtk.MenuNew()
	formatMain, _ := gtk.MenuItemNewWithLabel("Format")

	m.fontMenuItem, _ = gtk.MenuItemNewWithLabel("Font...")
	m.wordWrapMenuItem, _ = gtk.CheckMenuItemNewWithLabel("Word Wrap")

	formatMain.SetSubmenu(formatMenu)
	formatMenu.Append(m.wordWrapMenuItem)
	formatMenu.Append(m.fontMenuItem)

	m.gtkmenuBar.Append(formatMain)
}

func (m *menu) setupViewMenu() {
	viewMenu, _ := gtk.MenuNew()
	viewMain, _ := gtk.MenuItemNewWithLabel("View")

	m.statusBarMenuItem, _ = gtk.CheckMenuItemNewWithLabel("Status Bar")

	viewMain.SetSubmenu(viewMenu)
	viewMenu.Append(m.statusBarMenuItem)

	m.gtkmenuBar.Append(viewMain)
}

func (m *menu) setupHelpMenu() {
	helpMenu, _ := gtk.MenuNew()
	helpMain, _ := gtk.MenuItemNewWithLabel("Help")

	m.aboutMenuItem, _ = gtk.MenuItemNewWithLabel("About")
	key, mod := gtk.AcceleratorParse("F1")
	m.aboutMenuItem.AddAccelerator("activate", m.app.accelGroup, key, mod, gtk.ACCEL_VISIBLE)

	helpMain.SetSubmenu(helpMenu)
	helpMenu.Append(m.aboutMenuItem)

	m.gtkmenuBar.Append(helpMain)
}
