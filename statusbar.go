package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type (
	statusbar struct {
		app          *app
		gtkStatusBar *gtk.Statusbar
	}
)

func newStatusbar(app *app) *statusbar {

	// Define statusbar
	gtkstatusBar, err := gtk.StatusbarNew()

	if err != nil {
		log.Fatal("failed setting up gtk statusbar: ", err)
	}

	msgArea, _ := gtkstatusBar.GetMessageArea()

	msgArea.SetHAlign(gtk.ALIGN_END)

	return &statusbar{
		app:          app,
		gtkStatusBar: gtkstatusBar,
	}
}

func (s *statusbar) SetText(text string) {
	s.gtkStatusBar.Push(s.gtkStatusBar.GetContextId("textView cursor position"), text)
}

func (s *statusbar) Show() {
	if s.app == nil || s.app.grid == nil {
		return
	}

	s.app.grid.Add(s.gtkStatusBar)
	s.app.grid.ShowAll()
}

func (s *statusbar) Hide() {
	if s.app == nil || s.app.grid == nil {
		return
	}

	s.app.grid.Remove(s.gtkStatusBar)
	s.app.grid.ShowAll()
}
