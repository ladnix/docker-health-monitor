package main

import (
	"context"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func GetServiceNodes() ([]ServiceNode, error) {
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

	var nodes []ServiceNode
	for _, c := range containers {
		name := "none"
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		inspect, err := cli.ContainerInspect(ctx, c.ID)
		cancel()

		status := c.State
		ip := ""

		if err == nil {
			status = inspect.State.Status
			ip = inspect.NetworkSettings.IPAddress
			if ip == "" && len(inspect.NetworkSettings.Networks) > 0 {
				for _, n := range inspect.NetworkSettings.Networks {
					if n.IPAddress != "" {
						ip = n.IPAddress
						break
					}
				}
			}
		}

		node := ServiceNode{
			ID:     c.ID,
			Name:   name,
			Status: status,
			IP:     ip,
		}

		if err == nil {
			for _, e := range inspect.Config.Env {
				for _, other := range allNames {
					if other != name && strings.Contains(strings.ToLower(e), strings.ToLower(other)) {
						node.Deps = append(node.Deps, other)
					}
				}
			}
		}
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

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

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       "100",
	}

	out, err := cli.ContainerLogs(context.Background(), id, options)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	return io.ReadAll(out)
}
