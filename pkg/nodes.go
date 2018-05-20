package main

import (
	"fmt"
	"strconv"

	ui "github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
)

func NewNodePanel() *ui.List {

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

	for _, nd := range NODE_GAUGES {
		CpuColumn = append(CpuColumn, nd.CpuGauge)
		MemoryColumn = append(MemoryColumn, nd.MemoryGauge)
	}

	p := ui.NewList()
	p.Border = false
	p.Height = NODE_PANEL_HEIGHT
	return p
}

func nodeRow() *ui.Row {
	return ui.NewRow(
		ui.NewCol(4, 0, NodePanel),
		ui.NewCol(4, 0, CpuColumn...),
		ui.NewCol(4, 0, MemoryColumn...),
	)
}

func ShowNodes() {
	NodePanel.Height = ui.TermHeight() - 1
	ui.Body.Rows = []*ui.Row{
		nodeRow(),
		ui.NewRow(ui.NewCol(12, 0, Footer())),
	}
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

	max_rows := nodePanel.Height

	if len(items) > max_rows {
		items = items[0:max_rows]
	}

	nodePanel.Items = items

	for _, nd := range NODE_GAUGES {
		r, err := kubeClient.NodeNodeResources(&nd.Node)
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
