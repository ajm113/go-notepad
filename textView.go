package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type TextView struct {
	app         *App
	gtktextView *gtk.TextView
}

func NewTextView(app *App) *TextView {
	textView, err := gtk.TextViewNew()

	if err != nil {
		log.Fatal("failed setting up gtk textview: ", err)
	}

	adjv, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)
	adjh, _ := gtk.AdjustmentNew(0, 0, 0, 0, 0, 0)

	scrolled, err := gtk.ScrolledWindowNew(adjv, adjh)

	if err != nil {
		log.Fatal("failed setting up gtk scroll: ", err)
	}

	scrolled.SetHExpand(true)
	scrolled.SetVExpand(true)
	scrolled.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	scrolled.Add(textView)
	app.grid.Add(scrolled)

	return &TextView{
		app:         app,
		gtktextView: textView,
	}
}

func (t *TextView) SetText(text string) {
	b, _ := t.gtktextView.GetBuffer()
	b.SetText(text)
}

func (t *TextView) WrapText(wrap bool) {
	if wrap {
		t.gtktextView.SetWrapMode(gtk.WRAP_WORD)
	} else {
		t.gtktextView.SetWrapMode(gtk.WRAP_NONE)
	}
}

func (t *TextView) Copy() {

	b, _ := t.gtktextView.GetBuffer()

	if !t.gtktextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.gtktextView.Emit("copy-clipboard")
}

func (t *TextView) Cut() {
	b, _ := t.gtktextView.GetBuffer()
	if !t.gtktextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.gtktextView.Emit("cut-clipboard")
}

func (t *TextView) Paste() {
	if !t.gtktextView.IsFocus() {
		return
	}

	t.gtktextView.Emit("paste-clipboard")
}

func (t *TextView) SelectAll() {
	if !t.gtktextView.IsFocus() {
		return
	}

	t.gtktextView.Emit("select-all")
}

func (t *TextView) Backspace() {
	if !t.gtktextView.IsFocus() {
		return
	}

	t.gtktextView.Emit("backspace")
}
