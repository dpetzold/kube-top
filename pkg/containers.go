package main

import (
	"fmt"
	"sort"
	"time"

	ui "github.com/dpetzold/termui"
)

func NewContainersPanel(height int) *ui.Table {
	p := ui.NewTable()
	p.Height = height
	p.BorderLabel = "Containers"
	p.TextAlign = ui.AlignLeft
	p.Separator = false
	p.Headers = true
	p.Analysis()
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
		[]string{
			"Pod/Container",
			"Ready",
			"Status",
			"Restarts",
			"Cpu",
			"CpuMax",
			"Memory",
			"MemoryMax",
			"Age",
		},
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

			if maxes, ok := ContainerMaxes[c.Name]; ok {
				if maxes.Cpu.ToQuantity().Cmp(*c.Resources.CpuUsage.ToQuantity()) < 0 {
					maxes.Cpu = c.Resources.CpuUsage
					maxes.CpuTime = time.Now()
				}

				if maxes.Memory.ToQuantity().Cmp(*c.Resources.MemoryUsage.ToQuantity()) < 0 {
					maxes.Memory = c.Resources.MemoryUsage
					maxes.MemoryTime = time.Now()
				}
			} else {
				ContainerMaxes[c.Name] = &ContainerMax{
					Cpu:        c.Resources.CpuUsage,
					CpuTime:    time.Now(),
					Memory:     c.Resources.MemoryUsage,
					MemoryTime: time.Now(),
				}
			}

		}

		var cpuString, memoryString string
		if maxes, ok := ContainerMaxes[c.Name]; ok {
			cpuString = maxes.Cpu.String()
			memoryString = maxes.Memory.String()
		}

		rows = append(rows, []string{
			c.Name,
			ready,
			c.Status.Status,
			fmt.Sprintf("%d", c.Status.Restarts),
			cpuString,
			cpuResources,
			memoryResources,
			memoryString,
			c.Status.Age,
		})
	}

	max_rows := ContainerPanel.Height - 2
	if len(rows) > max_rows {
		rows = rows[0:max_rows]
	}

	containersPanel.Rows = rows
	containersPanel.Analysis()
}
