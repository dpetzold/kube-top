package main

import (
	ui "github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
)

var (
	kubeClient *KubeClient
	Namespace  string
)

func TopRun(k *KubeClient, namespace string) {

	kubeClient = k
	Namespace = namespace

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	nodes, err := k.Nodes()
	if err != nil {
		panic(err.Error())
	}

	nodeGauges := make(map[string]*NodeDisplay)
	var node_names []string
	var node_capacity []api_v1.ResourceList

	for _, node := range nodes {
		name := node.GetName()
		capacity := NodeCapacity(&node)
		nodeGauges[name] = &NodeDisplay{
			Node:        node,
			CpuGauge:    GaugePanel("Cpu", ui.ColorRed),
			MemoryGauge: GaugePanel("Mem", ui.ColorCyan),
		}
		node_names = append(node_names, name)
		node_capacity = append(node_capacity, capacity)
	}

	var cpu_column []ui.GridBufferer
	var mem_column []ui.GridBufferer

	for _, nd := range nodeGauges {
		cpu_column = append(cpu_column, nd.CpuGauge)
		mem_column = append(mem_column, nd.MemoryGauge)
	}

	listPanel := ListPanel(node_names, node_capacity)
	containersPanel := ContainersPanel()
	eventsPanel := EventsPanel()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, listPanel),
			ui.NewCol(4, 0, cpu_column...),
			ui.NewCol(4, 0, mem_column...),
		),
		ui.NewRow(
			ui.NewCol(12, 0, containersPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, eventsPanel),
		),
	)

	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	ui.Handle("/timer/1s", func(e ui.Event) {

		updateNodes(nodeGauges)
		updateContainers(containersPanel)
		updateEvents(eventsPanel)

		ui.Body.Align()
		ui.Render(ui.Body)
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
