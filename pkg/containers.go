package main

import (
	"fmt"
	"sort"

	ui "github.com/dpetzold/termui"
)

func ContainersPanel() *ui.Table {
	p := ui.NewTable()
	p.Height = 30
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
		[]string{"Pod/Container", "Ready", "Status", "Restarts", "Cpu", "Memory"},
	}

	for _, c := range containers {
		ready := "0"
		if c.Status.Ready {
			ready = "1"
		}
		rows = append(rows, []string{
			c.Name,
			ready,
			c.Status.Status,
			fmt.Sprintf("%d", c.Status.Restarts),
			c.Resources.CpuUsage.String(),
			c.Resources.MemoryUsage.String(),
		})
	}

	if len(rows) > 28 {
		rows = rows[0:28]
	}

	containersPanel.Rows = rows
	containersPanel.Analysis()
	containersPanel.SetSize()
}
