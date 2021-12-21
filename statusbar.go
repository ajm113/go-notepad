package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type (
	Statusbar struct {
		app          *App
		gtkStatusBar *gtk.Statusbar
	}
)

func NewStatusbar(app *App) (statusbar *Statusbar) {

	// Define statusbar
	gtkstatusBar, err := gtk.StatusbarNew()

	if err != nil {
		log.Fatal("failed setting up gtk statusbar: ", err)
	}

	msgArea, _ := gtkstatusBar.GetMessageArea()

	msgArea.SetHAlign(gtk.ALIGN_END)

	statusbar = &Statusbar{
		app:          app,
		gtkStatusBar: gtkstatusBar,
	}

	return
}

func (s *Statusbar) SetText(text string) {
	s.gtkStatusBar.Push(s.gtkStatusBar.GetContextId("textView cursor position"), text)
}

func (s *Statusbar) Show() {
	if s.app == nil || s.app.grid == nil {
		return
	}

	s.app.grid.Add(s.gtkStatusBar)
	s.app.grid.ShowAll()
}

func (s *Statusbar) Hide() {
	if s.app == nil || s.app.grid == nil {
		return
	}

	s.app.grid.Remove(s.gtkStatusBar)
	s.app.grid.ShowAll()
}
