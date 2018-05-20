package main

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/dpetzold/termui"
)

func Init() {
	ContainerMaxes = make(map[string]*ContainerMax)
	NodePanel = NewNodePanel()
	containers_height := ui.TermHeight() - EVENTS_PANEL_HEIGHT - NodePanel.Height
	ContainerPanel = NewContainersPanel(containers_height)
	EventsPanel = NewEventsPanel(EVENTS_PANEL_HEIGHT)
}

func updatePanels() {
	updateNodes(NodePanel)
	updateContainers(ContainerPanel)
	updateEvents(EventsPanel)
	ui.Body.Align()
	ui.Render(ui.Body)
}

func Footer() *ui.Par {
	text := "(D)ashboard    (C)ontainers    (E)vents    (N)odes    (Q)uit"
	padding := ui.TermWidth()/2 - len(text)/2

	p := ui.NewPar(fmt.Sprintf("%s%s", strings.Repeat(" ", padding), text))
	p.Border = false
	p.Height = 1
	p.TextFgColor = ui.ColorYellow
	return p
}

func createTimer(seconds time.Duration) string {
	refresh_duration := time.Second * seconds
	ui.DefaultEvtStream.Merge("timer", ui.NewTimerCh(refresh_duration))
	return fmt.Sprintf("/timer/%s", refresh_duration)
}

func showWindow(displayFunc func()) {
	displayFunc()
	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func KubeTop() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	Init()

	showWindow(ShowDashboard)

	ui.Handle("/sys/kbd/c", func(ui.Event) {
		showWindow(ShowContainers)
	})

	ui.Handle("/sys/kbd/d", func(ui.Event) {
		showWindow(ShowDashboard)
	})

	ui.Handle("/sys/kbd/e", func(ui.Event) {
		showWindow(ShowEvents)
	})

	ui.Handle("/sys/kbd/n", func(ui.Event) {
		showWindow(ShowNodes)
	})

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	timer_path := createTimer(REFRESH_SECONDS)
	ui.Handle(timer_path, func(ui.Event) {
		updatePanels()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
