package main

import (
	"io/ioutil"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type TextView struct {
	app         *App
	GTKtextView *gtk.TextView
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
		GTKtextView: textView,
	}
}

func (t *TextView) SetText(text string) {
	b, _ := t.GTKtextView.GetBuffer()
	b.SetText(text)
}

func (t *TextView) WrapText(wrap bool) {
	if wrap {
		t.GTKtextView.SetWrapMode(gtk.WRAP_WORD)
	} else {
		t.GTKtextView.SetWrapMode(gtk.WRAP_NONE)
	}
}

func (t *TextView) Copy() {

	b, _ := t.GTKtextView.GetBuffer()

	if !t.GTKtextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.GTKtextView.Emit("copy-clipboard")
}

func (t *TextView) Cut() {
	b, _ := t.GTKtextView.GetBuffer()
	if !t.GTKtextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.GTKtextView.Emit("cut-clipboard")
}

func (t *TextView) Paste() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("paste-clipboard")
}

func (t *TextView) SelectAll() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("select-all")
}

func (t *TextView) Backspace() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("backspace")
}

func (t *TextView) LoadSource(filename string) (err error) {
	src, err := ioutil.ReadFile(filename)

	if err != nil {
		return
	}

	buff, err := t.GTKtextView.GetBuffer()

	if err != nil {
		return
	}

	t.Clear()
	buff.Insert(buff.GetStartIter(), string(src))

	return
}

func (t *TextView) SaveSource(filename string) (err error) {
	buff, _ := t.GTKtextView.GetBuffer()

	source, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), true)

	err = ioutil.WriteFile(filename, []byte(source), 0666)

	return
}

func (t *TextView) Clear() {
	buff, _ := t.GTKtextView.GetBuffer()

	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
}
