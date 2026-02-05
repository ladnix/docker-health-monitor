package main

import "sync"

type AppMode int

const (
	ModeLite AppMode = iota // 1s interval, no stats
	ModeFull                // 3s interval, with CPU/RAM stats
)

type AppState struct {
	sync.RWMutex
	Mode             AppMode
	LastSelectedName string
	ActiveID         string
}

type ServiceNode struct {
	ID       string
	Name     string
	Status   string
	IP       string
	Deps     []string
	ExitCode int
	CPU      float64
	MemUsage uint64
	MemLimit uint64
}

var State = &AppState{Mode: ModeLite}
