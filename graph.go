package main

import "sort"

type ServiceNode struct {
	ID     string
	Name   string
	Status string
	IP     string
	Deps   []string
}

func SortNodes(nodes []ServiceNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Status != nodes[j].Status {
			if nodes[i].Status != "running" {
				return true
			}
			return false
		}
		return nodes[i].Name < nodes[j].Name
	})
}
