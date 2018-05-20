package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	ui "github.com/dpetzold/termui"
	"github.com/onrik/logrus/filename"
)

func UpdatePanels() {
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
	timer_path := fmt.Sprintf("/timer/%s", refresh_duration)
	ui.DefaultEvtStream.Merge("timer", ui.NewTimerCh(refresh_duration))
	return timer_path
}

func showWindow(displayFunc func()) {
	displayFunc()
	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

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

	ContainerMaxes = make(map[string]*ContainerMax)
	NodePanel = NewNodePanel()
	containers_height := ui.TermHeight() - EVENTS_PANEL_HEIGHT - NodePanel.Height
	ContainerPanel = NewContainersPanel(containers_height)
	EventsPanel = NewEventsPanel(EVENTS_PANEL_HEIGHT)

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
		UpdatePanels()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Loop()
}
