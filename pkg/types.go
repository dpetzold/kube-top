package main

import (
	"time"

	ui "github.com/dpetzold/termui"
	"github.com/sirupsen/logrus"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
)

const (
	REFRESH_SECONDS     = 3
	EVENTS_PANEL_HEIGHT = 15
	NODE_PANEL_HEIGHT   = 9
)

type GlobalsType struct {
	ActiveWindow   Window
	ContainerPanel *ui.Table
	Containers     []*ContainerInfo
	ContainerMaxes map[string]*ContainerMaxes
	CpuColumn      []ui.GridBufferer
	Events         []api_v1.Event
	EventsPanel    *ui.Table
	KubeClient     *KubeClient
	MemoryColumn   []ui.GridBufferer
	Namespace      string
	NodePanel      *ui.List
	NodeResources  []*NodeResources
	Nodes          []api_v1.Node
	SortField      string
	SortOrder      bool
}

var Globals GlobalsType

var log = logrus.New()

type Window string

const (
	NodesWindow      Window = "Nodes"
	DashboardWindow  Window = "Dashboard"
	EventsWindow     Window = "Events"
	ContainersWindow Window = "Containers"
)

type KubeClient struct {
	clientset      *kubernetes.Clientset
	heapsterClient *metricsutil.HeapsterMetricsClient
}

type CpuResource struct {
	*resource.Quantity
}

type MemoryResource struct {
	*resource.Quantity
}

type ContainerMaxes struct {
	CpuMax        *CpuResource
	CpuMaxTime    time.Time
	MemoryMax     *MemoryResource
	MemoryMaxTime time.Time
}

type ContainerInfo struct {
	Name string

	CpuMax        *CpuResource
	CpuMaxTime    time.Time
	MemoryMax     *MemoryResource
	MemoryMaxTime time.Time

	*ContainerStatus
	*ContainerUsage
}

type ContainerUsage struct {
	Name        string
	CpuUsage    *CpuResource
	MemoryUsage *MemoryResource
}

type ContainerResources struct {
	Name               string
	Namespace          string
	CpuReq             *CpuResource
	CpuLimit           *CpuResource
	PercentCpuReq      int
	PercentCpuLimit    int
	MemReq             *MemoryResource
	MemLimit           *MemoryResource
	PercentMemoryReq   int
	PercentMemoryLimit int
}

type ContainerStatus struct {
	Name     string
	Status   string
	Ready    bool
	Restarts int
	Age      time.Time
}

type NodeResources struct {
	Name           string
	Pods           int
	CpuUsage       *CpuResource
	CpuCapacity    *CpuResource
	PercentCpu     int
	MemoryUsage    *MemoryResource
	MemoryCapacity *MemoryResource
	PercentMemory  int
}
