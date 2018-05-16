package main

import (
	"fmt"
	"strconv"

	ui "github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
)

func NodePanel() (*ui.List, []ui.GridBufferer, []ui.GridBufferer) {

	nodes, err := kubeClient.Nodes()
	if err != nil {
		panic(err.Error())
	}

	NODE_GAUGES = make(map[string]*NodeGauges)
	var node_names []string
	var capacities []api_v1.ResourceList

	for _, node := range nodes {
		name := node.GetName()
		capacity := NodeCapacity(&node)
		NODE_GAUGES[name] = &NodeGauges{
			Node:        node,
			CpuGauge:    GaugePanel("Cpu", ui.ColorRed),
			MemoryGauge: GaugePanel("Mem", ui.ColorCyan),
		}
		node_names = append(node_names, name)
		capacities = append(capacities, capacity)
	}

	var cpu_column []ui.GridBufferer
	var mem_column []ui.GridBufferer

	for _, nd := range NODE_GAUGES {
		cpu_column = append(cpu_column, nd.CpuGauge)
		mem_column = append(mem_column, nd.MemoryGauge)
	}

	p := ui.NewList()
	p.Border = false
	p.Height = len(node_names) * NODE_DISPLAY_COUNT
	return p, cpu_column, mem_column
}

func GaugePanel(label string, barColor ui.Attribute) *ui.Gauge {

	p := ui.NewGauge()
	p.BarColor = barColor
	p.BorderFg = ui.ColorWhite
	p.BorderLabelFg = ui.ColorCyan
	p.BorderLabel = label
	p.Height = 3
	p.LabelAlign = ui.AlignRight
	return p
}

func updateNodes(nodePanel *ui.List) error {

	nodes, err := kubeClient.Nodes()
	if err != nil {
		panic(err.Error())
	}

	nodes = nodes[0:NODE_DISPLAY_COUNT]

	var node_resources []*NodeResources

	for _, node := range nodes {
		name := node.GetName()
		capacity := NodeCapacity(&node)

		pods, err := kubeClient.ActivePods("", name)
		if err != nil {
			return err
		}

		cpuQuantity := capacity[api_v1.ResourceCPU]
		memQuantity := capacity[api_v1.ResourceMemory]

		node_resources = append(node_resources, &NodeResources{
			Name:           name,
			Pods:           len(pods),
			CpuCapacity:    NewCpuResource(cpuQuantity.MilliValue()),
			MemoryCapacity: NewMemoryResource(memQuantity.Value()),
		})
	}

	var items []string
	resource_fmt := "%s[%s:](fg-cyan) [%s](fg-cyan,fg-bold)%s"

	for _, n := range node_resources {

		resources := fmt.Sprintf(resource_fmt, "  ", "Cpu", n.CpuCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Mem", n.MemoryCapacity.String(), "  ")
		resources += fmt.Sprintf(resource_fmt, "", "Pods", strconv.FormatInt(int64(n.Pods), 10), "")

		items = append(items, []string{
			fmt.Sprintf(" [%s](fg-white,fg-bold)", n.Name),
			resources,
			"",
		}...)
	}

	nodePanel.Items = items

	for _, nd := range NODE_GAUGES {
		r, _ := kubeClient.NodeNodeResources(&nd.Node)
		if err != nil {
			return err
		}

		nd.MemoryGauge.Percent = r.PercentMemory
		nd.MemoryGauge.Label = fmt.Sprintf("%d%% (%s)", r.PercentMemory, r.MemoryUsage.String())
		nd.CpuGauge.Percent = r.PercentCpu
		nd.CpuGauge.Label = fmt.Sprintf("%d%% (%s)", r.PercentCpu, r.CpuUsage.String())
	}

	return nil
}
