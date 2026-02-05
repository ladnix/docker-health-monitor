package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/rivo/tview"
)

func CreateUI(nodes []ServiceNode) (*tview.Application, *tview.TreeNode, *tview.TreeView, *tview.TextView) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	root := tview.NewTreeNode("[yellow]⬢ DHM[-]").SetColor(tcell.ColorWhite)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle("  Docker Health Monitor  ").SetTitleColor(tcell.ColorWhite)

	details := tview.NewTextView()
	details.SetDynamicColors(true).SetBorder(true).SetTitle("  Details  ")

	logsView := tview.NewTextView().SetDynamicColors(true).SetRegions(true).SetWordWrap(true)
	logsView.SetBorder(true).SetTitle("  Logs (ESC - back)  ")

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		State.Lock()
		defer State.Unlock()
		if node.GetReference() == nil {
			State.ActiveID = "ROOT"
		} else {
			sNode := node.GetReference().(ServiceNode)
			State.ActiveID = sNode.ID
		}
	})

	tree.SetChangedFunc(func(node *tview.TreeNode) {
		if node != nil {
			txt := node.GetText()
			if !strings.Contains(txt, "🔗") {
				State.Lock()
				State.LastSelectedName = txt
				State.Unlock()
			}
		}
	})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tree, 0, 2, true).
		AddItem(details, 9, 1, false)

	pages.AddPage("main", flex, true, true)
	pages.AddPage("help", CreateHelpModal(), true, false)
	pages.AddPage("logs", logsView, true, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if isHelpVisible(pages) || isLogsVisible(pages) {
			if event.Key() == tcell.KeyEsc {
				pages.HidePage("help").HidePage("logs")
				app.SetFocus(tree)
				return nil
			}
			return nil
		}
		switch {
		case event.Key() == tcell.KeyF1 || event.Rune() == 'h':
			pages.ShowPage("help")
			return nil
		case event.Rune() == 'i' || event.Rune() == 'I':
			State.Lock()
			if State.Mode == ModeLite {
				State.Mode = ModeFull
			} else {
				State.Mode = ModeLite
			}
			State.Unlock()
			return nil
		case event.Rune() == 'l' || event.Rune() == 'L':
			node := tree.GetCurrentNode()
			if node == nil || node.GetReference() == nil {
				return event
			}
			sNode := node.GetReference().(ServiceNode)
			rawLogs, _ := GetContainerLogs(sNode.ID)
			logsView.SetText(parseDockerLogs(rawLogs)).ScrollToBeginning()
			pages.ShowPage("logs")
			app.SetFocus(logsView)
			return nil
		case event.Rune() == 'r' || event.Rune() == 'R':
			node := tree.GetCurrentNode()
			if node == nil || node.GetReference() == nil {
				return event
			}
			sNode := node.GetReference().(ServiceNode)
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
			return nil
		case event.Key() == tcell.KeyEsc:
			app.Stop()
			return nil
		}
		return event
	})

	renderTree(tree, root, nodes)
	app.SetRoot(pages, true).SetFocus(tree)
	return app, root, tree, details
}

func StartMonitor(app *tview.Application, root *tview.TreeNode, tree *tview.TreeView, details *tview.TextView) {
	for {
		State.RLock()
		mode := State.Mode
		activeID := State.ActiveID
		State.RUnlock()

		nodes, err := GetServiceNodes(mode)
		if err == nil {
			app.QueueUpdateDraw(func() {
				renderTree(tree, root, nodes)
				if activeID != "" {
					if activeID == "ROOT" {
						details.SetText(getSystemSummary(nodes, mode))
					} else {
						for _, n := range nodes {
							if n.ID == activeID {
								details.SetText(getDetailsText(n))
								break
							}
						}
					}
				}
			})
		}
		if mode == ModeLite {
			time.Sleep(1 * time.Second)
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func isHelpVisible(p *tview.Pages) bool { n, _ := p.GetFrontPage(); return n == "help" }
func isLogsVisible(p *tview.Pages) bool { n, _ := p.GetFrontPage(); return n == "logs" }
