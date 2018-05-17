package main

import (
	"fmt"
	"os"
	"time"

	ui "github.com/dpetzold/termui"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func UpdatePanels() {
	updateNodes(NodePanel)
	updateContainers(ContainerPanel)
	updateEvents(EventsPanel)
}

var log = logrus.New()

func TopRun(k *KubeClient, namespace string) {

	filenameHook := filename.NewHook()
	filenameHook.Field = "source"
	log.AddHook(filenameHook)

	file, err := os.OpenFile("/tmp/kube-top.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	kubeClient = k
	Namespace = namespace

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	var cpu_column []ui.GridBufferer
	var mem_column []ui.GridBufferer

	ContainerMaxes = make(map[string]*ContainerMax)

	NodePanel, cpu_column, mem_column = NewNodePanel()

	containers_height := ui.TermHeight() - EVENTS_PANEL_HEIGHT - NodePanel.Height

	log.Infof("%d %d %d %d", ui.TermHeight(), EVENTS_PANEL_HEIGHT, NodePanel.Height, containers_height)

	ContainerPanel = NewContainersPanel(containers_height)
	EventsPanel = NewEventsPanel()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(4, 0, NodePanel),
			ui.NewCol(4, 0, cpu_column...),
			ui.NewCol(4, 0, mem_column...),
		),
		ui.NewRow(
			ui.NewCol(12, 0, ContainerPanel),
		),
		ui.NewRow(
			ui.NewCol(12, 0, EventsPanel),
		),
	)

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
