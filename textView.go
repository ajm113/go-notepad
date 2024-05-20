package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
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

func (t *textView) SetFont(font string, size int64) error {
	styleContext, err := t.GTKtextView.GetStyleContext()

	if err != nil {
		return err
	}

	cssProvider, err := gtk.CssProviderNew()
	if err != nil {
		return err
	}

	css := `
	textview {
		padding-top: 2px;
		padding-left: 2px;
		font-family: "` + font + `", "Lucida Console";
		font-size: ` + strconv.FormatInt(size, 10) + `pt;
	}
	`

	err = cssProvider.LoadFromData(css)

	if err != nil {
		return err
	}

	// Get the default screen for the GTK application.
	screen, err := styleContext.GetScreen()
	if err != nil {
		return err
	}

	// Add the CSS provider to the screen's style context.
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	t.GTKtextView.ShowAll()

	fmt.Printf("setting font: '%s' %d\n", font, size)

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

	t.GTKtextView.Emit("copy-clipboard", glib.TYPE_NONE)
}

func (t *textView) Cut() {
	b, _ := t.GTKtextView.GetBuffer()
	if !t.GTKtextView.IsFocus() || !b.GetHasSelection() {
		return
	}

	t.GTKtextView.Emit("cut-clipboard", glib.TYPE_NONE)
}

func (t *textView) Paste() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("paste-clipboard", glib.TYPE_NONE)
}

func (t *textView) SelectAll() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("select-all", glib.TYPE_NONE)
}

func (t *textView) Backspace() {
	if !t.GTKtextView.IsFocus() {
		return
	}

	t.GTKtextView.Emit("backspace", glib.TYPE_NONE)
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

func (t *textView) SaveSource(filename string) error {
	buff, _ := t.GTKtextView.GetBuffer()

	// TODO: Add a write file error
	source, err := buff.GetText(buff.GetStartIter(), buff.GetEndIter(), true)

	if err != nil {
		return err
	}

	err = os.WriteFile(filename, []byte(source), 0666)

	return err
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
