package main

import (
	ui "github.com/dpetzold/termui"
)

func TopRun(k *KubeClient, namespace string) {

	kubeClient = k
	Namespace = namespace

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	nodePanel, cpu_column, mem_column := NodePanel()
	containersPanel := ContainersPanel()
	eventsPanel := EventsPanel()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, nodePanel),
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

		updateNodes(nodePanel)
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
