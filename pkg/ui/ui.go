package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/dpetzold/kube-top/pkg/kube"
	"github.com/dpetzold/termui"
)

func UpdateResources() {

	containers, err := KubeClient.Containers(Namespace)
	if err != nil {
		panic(err.Error())
	}

	for _, c := range containers {

		if m, ok := ContainerMaxes[c.Name]; ok {
			if m.CpuMax.Quantity.Cmp(*c.CpuUsage.Quantity) < 0 {
				m.CpuMax = c.CpuUsage
				m.CpuMaxTime = time.Now()
			}

			if m.MemoryMax.Quantity.Cmp(*c.MemoryUsage.Quantity) < 0 {
				m.MemoryMax = c.MemoryUsage
				m.MemoryMaxTime = time.Now()
			}
		} else {
			ContainerMaxes[c.Name] = &kube.ContainerMaxes{
				CpuMax:        c.CpuUsage,
				CpuMaxTime:    time.Now(),
				MemoryMax:     c.MemoryUsage,
				MemoryMaxTime: time.Now(),
			}
		}

		maxes := ContainerMaxes[c.Name]

		c.CpuMax = maxes.CpuMax
		c.CpuMaxTime = maxes.CpuMaxTime
		c.MemoryMax = maxes.MemoryMax
		c.MemoryMaxTime = maxes.MemoryMaxTime

	}
	Containers = containers

	events, err := KubeClient.Events(Namespace)
	if err != nil {
		panic(err.Error())
	}
	Events = events

	nodes, err := KubeClient.Nodes()
	if err != nil {
		panic(err.Error())
	}
	Nodes = nodes

	var nodeResources []*kube.NodeResources
	for _, node := range nodes {

		resources, err := KubeClient.NodeResources(&node, Namespace)
		if err != nil {
			panic(err.Error())
		}

		nodeResources = append(nodeResources, resources)
	}

	NodeResources = nodeResources
}

func UpdatePanels() {
	updateNodes(NodePanel)
	updateContainers(ContainerPanel)
	updateEvents(EventsPanel)
	termui.Body.Align()
	termui.Render(termui.Body)
}

func CenterText(text string) string {
	padding := termui.TermWidth()/2 - len(text)/2
	return fmt.Sprintf("%s%s", strings.Repeat(" ", padding), text)
}

func Footer(text string) *termui.Par {
	p := termui.NewPar(CenterText(text))
	p.Border = false
	p.Height = 1
	p.TextFgColor = termui.ColorYellow
	return p
}

func createTimer(seconds time.Duration) string {
	refresh_duration := time.Second * seconds
	termui.DefaultEvtStream.Merge("timer", termui.NewTimerCh(refresh_duration))
	return fmt.Sprintf("/timer/%s", refresh_duration)
}

func showWindow(displayFunc func()) {
	displayFunc()
	UpdatePanels()
	termui.Clear()
	termui.Body.Align()
	termui.Render(termui.Body)
}

func ContainersHandler(e termui.EvtKbd) {

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
			if SortField == field {
				SortOrder = !SortOrder
			} else {
				SortField = field
			}
		}
	}

	termui.Clear()
	UpdatePanels()
}

func DashboardHandler(e termui.EvtKbd) {
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

func DefaultHandler(e termui.Event) {

	if e.Type != "keyboard" {
		return
	}

	k := e.Data.(termui.EvtKbd)

	switch ActiveWindow {
	case ContainersWindow:
		ContainersHandler(k)
	default:
		DashboardHandler(k)
	}
}

func KubeTop(kubeClient *kube.KubeClient, namespace string) {

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()

	SortField = "CpuUsage"
	NodePanel = NewNodePanel()
	ContainerMaxes = make(map[string]*kube.ContainerMaxes)
	containers_height := termui.TermHeight() - EVENTS_PANEL_HEIGHT - NodePanel.Height
	ContainerPanel = NewContainersPanel(containers_height)
	EventsPanel = NewEventsPanel(EVENTS_PANEL_HEIGHT)
	Namespace = namespace

	showWindow(ShowDashboard)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	timer_path := createTimer(REFRESH_SECONDS)
	termui.Handle(timer_path, func(termui.Event) {
		go func() {
			UpdateResources()
			UpdatePanels()
		}()
	})

	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	termui.Handle("/", func(e termui.Event) {
		DefaultHandler(e)
	})

	termui.Loop()
}
