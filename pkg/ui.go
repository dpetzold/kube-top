package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	ui "github.com/dpetzold/termui"
)

func Init() {
	Globals.SortField = "CpuUsage"
	Globals.NodePanel = NewNodePanel()
	Globals.ContainerMaxes = make(map[string]*ContainerMaxes)
	containers_height := ui.TermHeight() - EVENTS_PANEL_HEIGHT - Globals.NodePanel.Height
	Globals.ContainerPanel = NewContainersPanel(containers_height)
	Globals.EventsPanel = NewEventsPanel(EVENTS_PANEL_HEIGHT)
}

func UpdateResources() {

	containers, err := Globals.KubeClient.Containers(Globals.Namespace)
	if err != nil {
		panic(err.Error())
	}

	for _, c := range containers {

		if m, ok := Globals.ContainerMaxes[c.Name]; ok {
			if m.CpuMax.Quantity.Cmp(*c.CpuUsage.Quantity) < 0 {
				m.CpuMax = c.CpuUsage
				m.CpuMaxTime = time.Now()
			}

			if m.MemoryMax.Quantity.Cmp(*c.MemoryUsage.Quantity) < 0 {
				m.MemoryMax = c.MemoryUsage
				m.MemoryMaxTime = time.Now()
			}
		} else {
			Globals.ContainerMaxes[c.Name] = &ContainerMaxes{
				CpuMax:        c.CpuUsage,
				CpuMaxTime:    time.Now(),
				MemoryMax:     c.MemoryUsage,
				MemoryMaxTime: time.Now(),
			}
		}

		maxes := Globals.ContainerMaxes[c.Name]

		c.CpuMax = maxes.CpuMax
		c.CpuMaxTime = maxes.CpuMaxTime
		c.MemoryMax = maxes.MemoryMax
		c.MemoryMaxTime = maxes.MemoryMaxTime

	}
	Globals.Containers = containers

	events, err := Globals.KubeClient.Events(Globals.Namespace)
	if err != nil {
		panic(err.Error())
	}
	Globals.Events = events

	nodes, err := Globals.KubeClient.Nodes()
	if err != nil {
		panic(err.Error())
	}
	Globals.Nodes = nodes

	var nodeResources []*NodeResources
	for _, node := range nodes {

		resources, err := Globals.KubeClient.NodeResources(&node)
		if err != nil {
			panic(err.Error())
		}

		nodeResources = append(nodeResources, resources)
	}

	Globals.NodeResources = nodeResources
}

func UpdatePanels() {
	updateNodes(Globals.NodePanel)
	updateContainers(Globals.ContainerPanel)
	updateEvents(Globals.EventsPanel)
	ui.Body.Align()
	ui.Render(ui.Body)
}

func CenterText(text string) string {
	padding := ui.TermWidth()/2 - len(text)/2
	return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), text)
}

func Footer(text string) *ui.Par {
	p := ui.NewPar(CenterText(text))
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
	UpdatePanels()
	ui.Clear()
	ui.Body.Align()
	ui.Render(ui.Body)
}

func ContainersHandler(e ui.EvtKbd) {

	ContainerKeyMapping := map[string]string{
		"a": "Age",
		"e": "MemoryUsage",
		"E": "MemoryMax",
		"c": "CpuUsage",
		"C": "CpuMax",
		"p": "Name",
	}

	key := e.KeyStr

	if key == "b" {
		showWindow(ShowDashboard)
		return
	} else {
		if field, ok := ContainerKeyMapping[key]; ok {
			if Globals.SortField == field {
				Globals.SortOrder = !Globals.SortOrder
			} else {
				Globals.SortField = field
			}
		}
	}

	ui.Clear()
	UpdatePanels()
}

func DashboardHandler(e ui.EvtKbd) {
	switch e.KeyStr {
	case "c":
		showWindow(ShowContainers)
	case "d":
		showWindow(ShowDashboard)
	case "e":
		showWindow(ShowEvents)
	case "n":
		showWindow(ShowNodeWindow)
	}
}

func DefaultHandler(e ui.Event) {

	if e.Type != "keyboard" {
		return
	}

	k := e.Data.(ui.EvtKbd)

	log.Infof("%s", spew.Sdump(e))
	switch Globals.ActiveWindow {
	case ContainersWindow:
		ContainersHandler(k)
	default:
		DashboardHandler(k)
	}
}

func KubeTop() {

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	Init()

	showWindow(ShowDashboard)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	timer_path := createTimer(REFRESH_SECONDS)
	ui.Handle(timer_path, func(ui.Event) {
		go func() {
			UpdateResources()
			UpdatePanels()
		}()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	ui.Handle("/", func(e ui.Event) {
		DefaultHandler(e)
	})

	ui.Loop()
}
