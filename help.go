package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func CreateHelpModal() *tview.Grid {
	helpText := fmt.Sprintf(`
    [yellow]Docker Health Monitor[-]
    [gray]Version: %s | Author: ladnix[-]

    [yellow]Hotkeys[-]
    [white]Arrows[-]       : Tree Navigation
    [white]Enter[-]        : Container Info
    [white]L[-]            : Container Logs
    [white]R[-]            : Restart Service
    [white]H / F1[-]       : Show Help
    [white]ESC / Ctrl+C[-] : Exit

    [yellow]Container Status[-]
    [green]Green[-]  : Running
    [red]Red[-]    : Stopped / Exited
    [blue]Blue[-]   : Relation / Dependency`, Version)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpText).
		SetTextAlign(tview.AlignLeft)

	textView.SetBorder(true).
		SetTitle("  Help (ESC - back)  ")

	return tview.NewGrid().
		SetColumns(0, 40, 0).
		SetRows(0, 19, 0).
		AddItem(textView, 1, 1, 1, 1, 0, 0, true)
}
