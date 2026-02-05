package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func renderTree(tree *tview.TreeView, root *tview.TreeNode, nodes []ServiceNode) {
	SortNodes(nodes)
	root.ClearChildren()
	var nodeToFocus *tview.TreeNode

	State.RLock()
	mode := State.Mode
	lastSelected := State.LastSelectedName
	State.RUnlock()

	statusMap := make(map[string]ServiceNode)
	for _, n := range nodes {
		statusMap[n.Name] = n
	}

	for _, node := range nodes {
		status := strings.ToLower(strings.TrimSpace(node.Status))
		var color tcell.Color
		switch {
		case status == "running" || strings.HasPrefix(status, "up"):
			color = tcell.ColorGreen
		case status == "exited" && (node.ExitCode == 0 || node.ExitCode == 137 || node.ExitCode == 143):
			color = tcell.ColorGray
		default:
			color = tcell.ColorRed
		}

		statsInfo := ""
		if mode == ModeFull && (status == "running" || strings.HasPrefix(status, "up")) {
			cpuC, memC := "white", "white"
			if node.CPU > 80 {
				cpuC = "red"
			} else if node.CPU > 40 {
				cpuC = "yellow"
			}
			memP := 0.0
			if node.MemLimit > 0 {
				memP = (float64(node.MemUsage) / float64(node.MemLimit)) * 100
			}
			if memP > 90 {
				memC = "red"
			} else if memP > 70 {
				memC = "yellow"
			}
			statsInfo = fmt.Sprintf(" [%s]%.1f%%[-] | [%s]%.1fMB[-] ", cpuC, node.CPU, memC, float64(node.MemUsage)/1024/1024)
		}

		title := node.Name
		if statsInfo != "" {
			title = fmt.Sprintf("%s %s", node.Name, statsInfo)
		}

		nodeView := tview.NewTreeNode(title).SetColor(color).SetSelectable(true).SetReference(node)
		if node.Name == lastSelected || title == lastSelected {
			nodeToFocus = nodeView
		}

		for _, depName := range node.Deps {
			depNode, exists := statusMap[depName]
			depColor := tcell.ColorDarkRed
			if exists {
				ds := strings.ToLower(depNode.Status)
				if ds == "running" || strings.HasPrefix(ds, "up") {
					depColor = tcell.ColorBlue
				} else {
					depColor = tcell.ColorGray
				}
			}
			nodeView.AddChild(tview.NewTreeNode("  🔗 " + depName).SetColor(depColor).SetSelectable(false))
		}
		root.AddChild(nodeView)
	}
	if nodeToFocus != nil {
		tree.SetCurrentNode(nodeToFocus)
	}
}

func getDetailsText(node ServiceNode) string {
	status := strings.ToLower(node.Status)
	exitPart := ""
	if status == "exited" {
		exitPart = fmt.Sprintf("\n [green]EXIT CODE:[-] %d (%s)", node.ExitCode, getExitCodeDescription(node.ExitCode))
	}
	resPart := ""
	if (status == "running" || strings.HasPrefix(status, "up")) && node.CPU >= 0 {
		resPart = fmt.Sprintf("\n [green]CPU USAGE:[-] %.2f%%\n [green]MEM USAGE:[-] %.1f MB / %.1f MB",
			node.CPU, float64(node.MemUsage)/1024/1024, float64(node.MemLimit)/1024/1024)
	}
	return fmt.Sprintf(" [green]NAME:[-] %s\n [green]STATUS:[-] %s\n [green]IP ADDR:[-] %s\n [green]ID:[-] %s%s%s",
		node.Name, node.Status, node.IP, node.ID, exitPart, resPart)
}

func getSystemSummary(nodes []ServiceNode, mode AppMode) string {
	total := len(nodes)
	running := 0
	var tCPU float64
	var tMem uint64
	for _, n := range nodes {
		if strings.ToLower(n.Status) == "running" || strings.HasPrefix(strings.ToLower(n.Status), "up") {
			running++
			tCPU += n.CPU
			tMem += n.MemUsage
		}
	}
	mStr := "Lite"
	if mode == ModeFull {
		mStr = "Full"
	}
	return fmt.Sprintf(" [yellow]SYSTEM OVERVIEW[-]\n\n"+
		" [white]Containers:[-] %d total, [green]%d running[-]\n"+
		" [white]Total CPU:[-] [blue]%.2f%%[-]\n"+
		" [white]Total RAM:[-] [blue]%.1f MB[-]\n\n"+
		" [gray]Mode: %s (Press 'I' to toggle)[-]",
		total, running, tCPU, float64(tMem)/1024/1024, mStr)
}

func parseDockerLogs(raw []byte) string {
	if len(raw) == 0 {
		return " [gray]Logs are empty...[-]"
	}
	var r bytes.Buffer
	for i := 0; i < len(raw); {
		if len(raw[i:]) >= 8 && raw[i] <= 2 {
			size := binary.BigEndian.Uint32(raw[i+4 : i+8])
			i += 8
			if len(raw[i:]) < int(size) {
				break
			}
			r.Write(raw[i : i+int(size)])
			i += int(size)
		} else {
			r.WriteByte(raw[i])
			i++
		}
	}
	l := strings.Split(r.String(), "\n")
	for i, j := 0, len(l)-1; i < j; i, j = i+1, j-1 {
		l[i], l[j] = l[j], l[i]
	}
	s := strings.Join(l, "\n")
	s = strings.ReplaceAll(s, "ERROR", "[red]ERROR[-]")
	s = strings.ReplaceAll(s, "WARN", "[yellow]WARN[-]")
	return s
}

func getExitCodeDescription(code int) string {
	switch code {
	case 0:
		return "Clean Exit"
	case 137:
		return "Manual Stop (SIGKILL)"
	case 143:
		return "Graceful Stop (SIGTERM)"
	default:
		return "Error"
	}
}
