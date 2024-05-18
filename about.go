package main

import "github.com/gotk3/gotk3/gtk"

func displayAboutDialog(app *app) {
	d, _ := gtk.DialogNew()
	d.SetTitle("About " + appName)
	d.SetTransientFor(app.Win)
	d.SetSizeRequest(400, 400)
	d.SetResizable(false)

	b, _ := d.GetContentArea()
	b.SetSpacing(5)
	b.SetMarginStart(5)
	b.SetMarginTop(10)

	label, _ := gtk.LabelNew("")
	label.SetMarkup(
		appName + "\nVersion 1.0\nMIT License 2021-2024 <a href=\"https://github.com/ajm113\">@ajm113</a>\n" +
			"Visit our <a href=\"https://github.com/ajm113/go-notepad\">GitHub Page</a> for more information or support!",
	)
	label.SetHAlign(gtk.ALIGN_START)

	b.PackStart(label, true, true, 0)

	d.AddButton("OK", gtk.RESPONSE_OK)
	d.ShowAll()
	d.Run()
	d.Destroy()

}
