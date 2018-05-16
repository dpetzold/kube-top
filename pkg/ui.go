package main

import (
	"fmt"
	"time"

	ui "github.com/dpetzold/termui"
)

func UpdatePanels() {
	updateNodes(NODE_PANEL)
	updateContainers(CONTAINER_PANEL)
	updateEvents(EVENTS_PANEL)
}

func TopRun(k *KubeClient, namespace string) {

	kubeClient = k
	Namespace = namespace

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	var cpu_column []ui.GridBufferer
	var mem_column []ui.GridBufferer

	NODE_PANEL, cpu_column, mem_column = NodePanel()
	CONTAINER_PANEL = ContainersPanel()
	EVENTS_PANEL = EventsPanel()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, NODE_PANEL),
			ui.NewCol(4, 0, cpu_column...),
			ui.NewCol(4, 0, mem_column...),
		),
		ui.NewRow(
			ui.NewCol(12, 0, CONTAINER_PANEL),
		),
		ui.NewRow(
			ui.NewCol(12, 0, EVENTS_PANEL),
		),
	)

	UpdatePanels()

	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	refresh_duration := time.Second * REFRESH_SECONDS
	timer_path := fmt.Sprintf("/timer/%s", refresh_duration)
	ui.DefaultEvtStream.Merge("timer", ui.NewTimerCh(refresh_duration))

	ui.Handle(timer_path, func(ui.Event) {
		UpdatePanels()
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
