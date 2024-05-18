package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func displayGotoLine(app *app) (response gtk.ResponseType, line string) {
	d, _ := gtk.DialogNew()
	d.SetTitle("Goto line")
	d.SetTransientFor(app.Win)
	d.SetSizeRequest(300, 100)

	b, _ := d.GetContentArea()
	b.SetSpacing(5)
	b.SetMarginTop(10)
	b.SetMarginStart(10)
	b.SetMarginEnd(10)

	label, _ := gtk.LabelNew("Line Number:")
	label.SetHAlign(gtk.ALIGN_START)
	input, _ := gtk.EntryNew()
	input.Emit("set-focus")
	input.Connect("key-press-event", func(_ *gtk.Entry, e *gdk.Event) {
		k := gdk.EventKeyNewFromEvent(e)

		if k.KeyVal() == 65293 {
			d.Response(gtk.RESPONSE_OK)
		}

	})

	b.PackStart(label, true, true, 0)
	b.PackStart(input, true, true, 0)

	d.AddButton("OK", gtk.RESPONSE_OK)
	d.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	d.ShowAll()

	response = d.Run()
	line, _ = input.GetText()
	d.Destroy()

	return
}
