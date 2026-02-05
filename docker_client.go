package main

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetServiceNodes(mode AppMode) ([]ServiceNode, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var allNames []string
	for _, c := range containers {
		if len(c.Names) > 0 {
			allNames = append(allNames, strings.TrimPrefix(c.Names[0], "/"))
		}
	}

	var wg sync.WaitGroup
	nodeChan := make(chan ServiceNode, len(containers))

	for _, c := range containers {
		wg.Add(1)
		go func(c types.Container) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			inspect, err := cli.ContainerInspect(ctx, c.ID)
			cancel()

			node := ServiceNode{
				ID:     c.ID,
				Name:   strings.TrimPrefix(c.Names[0], "/"),
				Status: c.State,
				IP:     "",
			}

			if err == nil {
				node.Status = inspect.State.Status
				node.ExitCode = inspect.State.ExitCode
				node.IP = inspect.NetworkSettings.IPAddress
				if node.IP == "" && len(inspect.NetworkSettings.Networks) > 0 {
					for _, n := range inspect.NetworkSettings.Networks {
						if n.IPAddress != "" {
							node.IP = n.IPAddress
							break
						}
					}
				}

				if mode == ModeFull && node.Status == "running" {
					sCtx, sCancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
					statsResp, sErr := cli.ContainerStats(sCtx, c.ID, false)
					if sErr == nil {
						var v struct {
							CPUStats struct {
								CPUUsage struct {
									TotalUsage  uint64   `json:"total_usage"`
									PercpuUsage []uint64 `json:"percpu_usage"`
								} `json:"cpu_usage"`
								SystemUsage uint64 `json:"system_usage"`
							} `json:"cpu_stats"`
							PreCPUStats struct {
								CPUUsage    struct{ TotalUsage uint64 } `json:"cpu_usage"`
								SystemUsage uint64                      `json:"system_usage"`
							} `json:"precpu_stats"`
							MemoryStats struct {
								Usage uint64 `json:"usage"`
								Limit uint64 `json:"limit"`
							} `json:"memory_stats"`
						}
						if json.NewDecoder(statsResp.Body).Decode(&v) == nil {
							cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
							systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
							numCPUs := float64(len(v.CPUStats.CPUUsage.PercpuUsage))
							if numCPUs == 0 {
								numCPUs = 1.0
							}
							if systemDelta > 0 && cpuDelta > 0 {
								node.CPU = (cpuDelta / systemDelta) * numCPUs * 100.0
							}
							node.MemUsage = v.MemoryStats.Usage
							node.MemLimit = v.MemoryStats.Limit
						}
						statsResp.Body.Close()
					}
					sCancel()
				}

				for _, e := range inspect.Config.Env {
					for _, other := range allNames {
						if other != node.Name && strings.Contains(strings.ToLower(e), strings.ToLower(other)) {
							node.Deps = append(node.Deps, other)
						}
					}
				}
			}
			nodeChan <- node
		}(c)
	}

	go func() {
		wg.Wait()
		close(nodeChan)
	}()

	var nodes []ServiceNode
	for n := range nodeChan {
		nodes = append(nodes, n)
	}
	SortNodes(nodes)
	return nodes, nil
}

func RestartContainer(id string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	timeout := 10
	return cli.ContainerRestart(context.Background(), id, container.StopOptions{Timeout: &timeout})
}

func GetContainerLogs(id string) ([]byte, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()
	out, err := cli.ContainerLogs(context.Background(), id, container.LogsOptions{ShowStdout: true, ShowStderr: true, Timestamps: true, Tail: "100"})
	if err != nil {
		return nil, err
	}
	defer out.Close()
	return io.ReadAll(out)
}
