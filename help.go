package main

import (
	"fmt"

	"github.com/rivo/tview"
)

func CreateHelpModal() *tview.Grid {
	helpText := fmt.Sprintf(`
  [yellow]Docker Health Monitor[-]
  [gray]Version: %s
  Author: ladnix[-]

  [yellow]Hotkeys[-]
  Tree Navigation         : [white]Arrows[-]
  Select for Monitoring   : [white]Enter[-]
  Container Logs          : [white]L[-]
  Restart Service         : [white]R[-] 
  Toggle Info Mode        : [white]I[-]
  Show Help               : [white]H / F1[-]
  Exit                    : [white]ESC / Ctrl+C[-]

  [yellow]Info Modes (Toggle with 'I')[-]
  [green]Lite[-] : Fast (1s), name & status only
  [red]Full[-] : Slow (3s), adds CPU & RAM stats

  [yellow]Container Status[-]
  Running (Healthy)       : [green]Green[-]
  Stopped (Manual/Clean)  : [gray]Gray[-] 
  Failed (Exit Code != 0) : [red]Red[-]
  Dependency Link         : [blue]Blue[-]`, Version)

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(helpText).
		SetTextAlign(tview.AlignLeft)

	textView.SetBorder(true).
		SetTitle("  Help (ESC - back)  ")

	return tview.NewGrid().
		SetColumns(0, 44, 0).
		SetRows(0, 26, 0).
		AddItem(textView, 1, 1, 1, 1, 0, 0, true)
}
