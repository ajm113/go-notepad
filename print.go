package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

func printText(buffer *gtk.TextBuffer, exportName string) {
	printOp, _ := gtk.PrintOperationNew()
	printOp.SetExportFilename(exportName + ".pdf")
	printOp.SetUnit(gtk.GTK_UNIT_POINTS)

	printOp.Connect("begin-print", func() {
		printTextBegin(printOp)
	})

	printOp.Connect("draw-page", func(operation *gtk.PrintOperation, ctx *gtk.PrintContext, pageNr int) {
		printTextDraw(buffer, operation, ctx, pageNr)
	})

	printOp.Connect("end-print", func() {
		// Clean up or notify when printing is finished
	})

	printOp.Run(gtk.PRINT_OPERATION_ACTION_PRINT_DIALOG, nil)
}

func printTextBegin(operation *gtk.PrintOperation) {

	operation.SetNPages(1) // Assuming only one page for simplicity
	operation.SetHasSelection(false)
	operation.SetUnit(gtk.GTK_UNIT_POINTS)
	operation.SetEmbedPageSetup(true)
	operation.SetUseFullPage(true)

	defaultPageSetup, _ := operation.GetDefaultPageSetup()
	printSettings, _ := operation.GetPrintSettings(defaultPageSetup)
	operation.SetPrintSettings(printSettings)
	operation.SetUnit(gtk.GTK_UNIT_POINTS)
}

func printTextDraw(buffer *gtk.TextBuffer, operation *gtk.PrintOperation, ctx *gtk.PrintContext, pageNr int) {
	startIter := buffer.GetStartIter()
	endIter := buffer.GetEndIter()
	text, _ := buffer.GetText(startIter, endIter, true)

	carioContext := ctx.GetCairoContext()
	pangoContext := pango.CairoCreateContext(carioContext)
	layout := pango.LayoutNew(pangoContext)
	layout.SetText(text, len(text))

	defaultPageSetup, _ := operation.GetDefaultPageSetup()
	printSettings, _ := operation.GetPrintSettings(defaultPageSetup)

	printSettings.SetResolution(72)

	// width, height := operation.GetPageSize()
	layout.SetWidth(int(printSettings.GetPaperWidth(gtk.GTK_UNIT_POINTS)) * pango.SCALE)
	layout.SetHeight(int(printSettings.GetPaperHeight(gtk.GTK_UNIT_POINTS)) * pango.SCALE)

	carioContext.SetSourceRGB(0, 0, 0) // Set text color
	carioContext.MoveTo(0, 0)
	pango.CairoShowLayout(carioContext, layout)
}
