package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type textView struct {
	app         *app
	GTKtextView *gtk.TextView
}

func newTextView(app *app) *textView {
	tv, err := gtk.TextViewNew()

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

	scrolled.Add(tv)
	app.grid.Add(scrolled)

	tv.SetMonospace(true)

	target, err := gtk.TargetEntryNew("text/uri-list", gtk.TargetFlags(0), 0)
	if err != nil {
		log.Fatal("failed creating target for textView", err)
	}

	tv.DragDestSet(gtk.DEST_DEFAULT_ALL, []gtk.TargetEntry{*target}, gdk.ACTION_COPY)

	return &textView{
		app:         app,
		GTKtextView: tv,
	}
}

func (t *textView) SetFont(font string, size int64, style string) error {
	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatal("Unable to create CSS provider:", err)
	}
	err = cssProvider.LoadFromData(`
        textview {
            font-family: "` + font + `";
            font-size: ` + strconv.FormatInt(size, 10) + `pt;
			font-style: "` + style + `";
		}
    `)

	if err != nil {
		return err
	}

	// Get the default screen for the GTK application.
	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		log.Fatal("Unable to get default screen:", err)
	}

	// Add the CSS provider to the screen's style context.
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_USER)

	return nil
}

func (t *textView) SetText(text string) {
	b, _ := t.GTKtextView.GetBuffer()
	b.SetText(text)
}

func (t *textView) WrapText(wrap bool) {
	if wrap {
		t.GTKtextView.SetWrapMode(gtk.WRAP_WORD)
	} else {
		t.GTKtextView.SetWrapMode(gtk.WRAP_NONE)
	}
}

func (t *textView) Copy() {

	b, _ := t.GTKtextView.GetBuffer()

	if !t.GTKtextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.GTKtextView.Emit("copy-clipboard")
}

func (t *textView) Cut() {
	b, _ := t.GTKtextView.GetBuffer()
	if !t.GTKtextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.GTKtextView.Emit("cut-clipboard")
}

func (t *textView) Paste() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("paste-clipboard")
}

func (t *textView) SelectAll() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("select-all")
}

func (t *textView) Backspace() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("backspace")
}

func (t *textView) LoadSource(filename string) (err error) {
	src, err := os.ReadFile(filename)

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

func (t *textView) SaveSource(filename string) (err error) {
	buff, _ := t.GTKtextView.GetBuffer()

	source, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), true)

	err = ioutil.WriteFile(filename, []byte(source), 0666)

	return
}

func (t *textView) Clear() {
	buff, _ := t.GTKtextView.GetBuffer()

	buff.Delete(buff.GetStartIter(), buff.GetEndIter())
}

func (t *textView) InsertTimestamp() {
	timestamp := time.Now().Format("1:04 PM 02/01/2006")

	buff, _ := t.GTKtextView.GetBuffer()
	buff.InsertAtCursor(timestamp)
}

func (t *textView) GoToLine(i int) {
	buff, _ := t.GTKtextView.GetBuffer()
	buff.PlaceCursor(buff.GetIterAtLine(i))
}
