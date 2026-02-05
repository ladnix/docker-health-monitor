package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

var Version = "dev"

func main() {
	vFlag := flag.Bool("v", false, "show version")
	versionFlag := flag.Bool("version", false, "show version")
	flag.Parse()

	if *vFlag || *versionFlag {
		fmt.Printf("DHM version: %s\n", Version)
		return
	}

	tview.Borders.Horizontal = '─'
	tview.Borders.Vertical = '│'
	tview.Borders.TopLeft = '┌'
	tview.Borders.TopRight = '┐'
	tview.Borders.BottomLeft = '└'
	tview.Borders.BottomRight = '┘'
	tview.Borders.LeftT = '├'
	tview.Borders.RightT = '┤'
	tview.Borders.TopT = '┬'
	tview.Borders.BottomT = '┴'
	tview.Borders.Cross = '┼'

	os.Setenv("TERM", "xterm-256color")

	nodes, err := GetServiceNodes()
	if err != nil {
		panic(err)
	}

	app, root, tree, _ := CreateUI(nodes)

	go StartMonitor(app, root, tree)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
