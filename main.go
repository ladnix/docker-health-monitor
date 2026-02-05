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

	flag.Usage = func() {
		fmt.Println(

			"DHM - Docker Health Monitor")
		fmt.Println("\nUsage: dhm [options] [command]")
		fmt.Println("\nOptions:")
		fmt.Println("  -v, --version  Show version info")
		fmt.Println("  -h, --help     Show this help message")
		fmt.Println("\nCommands:")
		fmt.Println("  help           Show interactive help (same as -h)")
		fmt.Println("\nControls (inside app):")
		fmt.Println("  Arrows         Navigate tree")
		fmt.Println("  Enter          Select for monitoring")
		fmt.Println("  L              View logs")
		fmt.Println("  I              Toggle Lite/Full mode")
	}

	flag.Parse()
	args := flag.Args()

	if *vFlag || *versionFlag {
		fmt.Printf("DHM version: %s\n", Version)
		return
	}

	if len(args) > 0 && args[0] == "help" {
		flag.Usage()
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

	nodes, err := GetServiceNodes(State.Mode)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Docker: %v", err))
	}

	app, root, tree, details := CreateUI(nodes)

	go StartMonitor(app, root, tree, details)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
