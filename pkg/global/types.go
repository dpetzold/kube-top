package global

import (
	"github.com/dpetzold/kube-top/pkg/kube"
	"github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
)

type Window string

const (
	REFRESH_SECONDS     = 3
	EVENTS_PANEL_HEIGHT = 15
	NODE_PANEL_HEIGHT   = 9

	NodesWindow      Window = "Nodes"
	DashboardWindow  Window = "Dashboard"
	EventsWindow     Window = "Events"
	ContainersWindow Window = "Containers"
)

var (
	ActiveWindow   Window
	ContainerPanel *termui.Table
	Containers     []*kube.ContainerInfo
	ContainerMaxes map[string]*kube.ContainerMaxes
	CpuColumn      []termui.GridBufferer
	Events         []api_v1.Event
	EventsPanel    *termui.Table
	KubeClient     *kube.KubeClientType
	MemoryColumn   []termui.GridBufferer
	Namespace      string
	NodePanel      *termui.List
	NodeResources  []*kube.NodeResources
	Nodes          []api_v1.Node
	SortField      string
	SortOrder      bool
)
