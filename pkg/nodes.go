package main

import (
	"fmt"

	ui "github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
)

type NodeDisplay struct {
	Node        api_v1.Node
	CpuGauge    *ui.Gauge
	MemoryGauge *ui.Gauge
}

func ListPanel(names []string, capacities []api_v1.ResourceList) *ui.List {

	var items []string
	for i, name := range names {

		capacity := capacities[i]

		cpuR := capacity[api_v1.ResourceCPU]
		cpu := NewCpuResource(cpuR.MilliValue())

		memR := capacity[api_v1.ResourceMemory]
		mem := NewMemoryResource(memR.Value())

		items = append(items, []string{
			fmt.Sprintf(" %s", name),
			fmt.Sprintf("   Cpu: %s  Mem: %s", cpu.String(), mem.String()),
			"",
		}...)
	}

	p := ui.NewList()
	p.Border = false
	p.Items = items
	p.Height = len(names) * 3
	return p
}

func GaugePanel(label string, barColor ui.Attribute) *ui.Gauge {

	p := ui.NewGauge()
	p.BarColor = barColor
	p.BorderFg = ui.ColorWhite
	p.BorderLabelFg = ui.ColorCyan
	p.BorderLabel = label
	p.Height = 3
	p.LabelAlign = ui.AlignRight
	p.PaddingBottom = 0
	p.Percent = 0
	return p
}

func updateNodes(nodeGauges map[string]*NodeDisplay) {
	for _, nd := range nodeGauges {
		r, _ := kubeClient.NodeResourceUsage(&nd.Node)
		nd.MemoryGauge.Percent = r.PercentMemory
		nd.MemoryGauge.Label = fmt.Sprintf("%d%% (%s)", r.PercentMemory, r.MemoryUsage.String())
		nd.CpuGauge.Percent = r.PercentCpu
		nd.CpuGauge.Label = fmt.Sprintf("%d%% (%s)", r.PercentCpu, r.CpuUsage.String())
	}
}
