package main

import (
	"sort"
	"strings"
)

func SortNodes(nodes []ServiceNode) {
	sort.Slice(nodes, func(i, j int) bool {
		getPriority := func(n ServiceNode) int {
			status := strings.ToLower(n.Status)

			if status == "dead" || (status == "exited" && n.ExitCode != 0 && n.ExitCode != 137 && n.ExitCode != 143) {
				return 0
			}

			if status == "exited" {
				return 1
			}

			return 2
		}

		p1, p2 := getPriority(nodes[i]), getPriority(nodes[j])
		if p1 != p2 {
			return p1 < p2
		}
		return nodes[i].Name < nodes[j].Name
	})
}
