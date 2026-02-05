package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var lastSelectedName string

func CreateUI(nodes []ServiceNode) (*tview.Application, *tview.TreeNode, *tview.TreeView, *tview.TextView) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	root := tview.NewTreeNode("[yellow]⬢ DHM[-]").SetColor(tcell.ColorWhite)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle("  Docker Health Monitor  ").SetTitleColor(tcell.ColorWhite)

	details := tview.NewTextView()
	details.SetDynamicColors(true).SetBorder(true).SetTitle("  Details  ")

	logsView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true)
	logsView.SetBorder(true).SetTitle("  Container Logs (ESC - back)  ")

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}
		sNode := ref.(ServiceNode)
		info := fmt.Sprintf(" [green]NAME:[-]     %s\n [green]STATUS:[-]   %s\n [green]IP ADDR:[-]  %s\n [green]ID:[-]       %s",
			sNode.Name, sNode.Status, sNode.IP, sNode.ID)
		details.SetText(info)
	})

	tree.SetChangedFunc(func(node *tview.TreeNode) {
		if node != nil {
			txt := node.GetText()
			if !strings.Contains(txt, "🔗") {
				lastSelectedName = txt
			}
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tree, 0, 2, true).
		AddItem(details, 7, 1, false)

	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", CreateHelpModal(), true, false)
	pages.AddPage("logs", logsView, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isHelpVisible(pages) || isLogsVisible(pages) {
			if event.Key() == tcell.KeyEsc {
				pages.HidePage("help")
				pages.HidePage("logs")
				app.SetFocus(tree)
				return nil
			}
			return nil
		}
		switch {
		case event.Key() == tcell.KeyF1 || event.Rune() == 'h':
			pages.ShowPage("help")
			return nil
		case event.Rune() == 'l' || event.Rune() == 'L':
			node := tree.GetCurrentNode()
			if node == nil {
				return event
			}
			ref := node.GetReference()
			if ref == nil {
				return event
			}
			sNode := ref.(ServiceNode)
			rawLogs, err := GetContainerLogs(sNode.ID)
			if err != nil {
				logsView.SetText("Error fetching logs: " + err.Error())
			} else {
				logsView.SetText(parseDockerLogs(rawLogs))
				logsView.ScrollToBeginning()
			}
			pages.ShowPage("logs")
			app.SetFocus(logsView)
			return nil
		case event.Rune() == 'r' || event.Rune() == 'R':
			node := tree.GetCurrentNode()
			if node == nil {
				return event
			}
			ref := node.GetReference()
			if ref == nil {
				return event
			}
			sNode := ref.(ServiceNode)
			details.SetText(fmt.Sprintf("\n [orange]Restarting %s...[-]", sNode.Name))
			go func() {
				err := RestartContainer(sNode.ID)
				app.QueueUpdateDraw(func() {
					if err != nil {
						details.SetText(" [red]Error: " + err.Error())
					} else {
						details.SetText(" [green]Container restarted![-]")
					}
				})
			}()
		case event.Key() == tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	updateTree(tree, root, nodes)
	app.SetRoot(pages, true).SetFocus(tree)
	return app, root, tree, details
}

func isHelpVisible(p *tview.Pages) bool {
	name, _ := p.GetFrontPage()
	return name == "help"
}

func isLogsVisible(p *tview.Pages) bool {
	name, _ := p.GetFrontPage()
	return name == "logs"
}

func StartMonitor(app *tview.Application, root *tview.TreeNode, tree *tview.TreeView) {
	for {
		time.Sleep(1 * time.Second)
		nodes, err := GetServiceNodes()
		if err != nil {
			continue
		}
		app.QueueUpdateDraw(func() {
			updateTree(tree, root, nodes)
		})
	}
}

func updateTree(tree *tview.TreeView, root *tview.TreeNode, nodes []ServiceNode) {
	SortNodes(nodes)
	root.ClearChildren()
	var nodeToFocus *tview.TreeNode
	for _, node := range nodes {
		status := strings.ToLower(strings.TrimSpace(node.Status))
		var color tcell.Color
		if status == "running" || strings.HasPrefix(status, "up") {
			color = tcell.ColorGreen
		} else if status == "" {
			color = tcell.ColorGray
		} else {
			color = tcell.ColorRed
		}
		title := fmt.Sprintf("%s [%s]", node.Name, node.Status)
		nodeView := tview.NewTreeNode(title).SetColor(color).SetSelectable(true).SetReference(node)
		if title == lastSelectedName {
			nodeToFocus = nodeView
		}
		for _, dep := range node.Deps {
			depView := tview.NewTreeNode("  🔗 " + dep).SetColor(tcell.ColorBlue).SetSelectable(false)
			nodeView.AddChild(depView)
		}
		root.AddChild(nodeView)
	}
	if nodeToFocus != nil {
		tree.SetCurrentNode(nodeToFocus)
	}
}

func parseDockerLogs(raw []byte) string {
	if len(raw) == 0 {
		return " [gray]Logs are empty...[-]"
	}
	var result bytes.Buffer
	i := 0
	for i < len(raw) {
		if len(raw[i:]) >= 8 && raw[i] <= 2 {
			size := binary.BigEndian.Uint32(raw[i+4 : i+8])
			i += 8
			if len(raw[i:]) < int(size) {
				break
			}
			result.Write(raw[i : i+int(size)])
			i += int(size)
		} else {
			result.WriteByte(raw[i])
			i++
		}
	}
	lines := strings.Split(result.String(), "\n")
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	str := strings.Join(lines, "\n")
	str = strings.ReplaceAll(str, "ERROR", "[red]ERROR[-]")
	str = strings.ReplaceAll(str, "WARN", "[yellow]WARN[-]")
	str = strings.ReplaceAll(str, "INFO", "[blue]INFO[-]")
	str = strings.ReplaceAll(str, "FAIL", "[red]FAIL[-]")
	str = strings.ReplaceAll(str, "UTC", "[gray]UTC[-]")
	return str
}
