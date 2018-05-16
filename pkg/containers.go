package main

import (
	"fmt"
	"sort"

	ui "github.com/dpetzold/termui"
)

func ContainersPanel() *ui.Table {
	p := ui.NewTable()
	p.Height = CONTAINER_PANEL_HEIGHT
	p.BorderLabel = "Containers"
	p.TextAlign = ui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
	p.SetSize()
	return p
}

func updateContainers(containersPanel *ui.Table) {
	containers, err := kubeClient.Containers(Namespace)
	if err != nil {
		panic(err.Error())
	}

	sort.Slice(containers, func(i, j int) bool {
		if containers[j].Resources == nil || containers[i].Resources == nil {
			return false
		}
		q := containers[j].Resources.CpuUsage.ToQuantity()
		return containers[i].Resources.CpuUsage.ToQuantity().Cmp(*q) > 0
	})

	rows := [][]string{
		[]string{"Pod/Container", "Ready", "Status", "Restarts", "Cpu", "Memory", "Age"},
	}

	for _, c := range containers {
		ready := "0"
		if c.Status.Ready {
			ready = "1"
		}

		cpuResources := ""
		memoryResources := ""
		if c.Resources != nil {
			cpuResources = c.Resources.CpuUsage.String()
			memoryResources = c.Resources.MemoryUsage.String()
		}

		rows = append(rows, []string{
			c.Name,
			ready,
			c.Status.Status,
			fmt.Sprintf("%d", c.Status.Restarts),
			cpuResources,
			memoryResources,
			c.Status.Age,
		})
	}

	max_rows := CONTAINER_PANEL_HEIGHT - 2
	if len(rows) > max_rows {
		rows = rows[0:max_rows]
	}

	containersPanel.Rows = rows
	containersPanel.Analysis()
	containersPanel.SetSize()
}
