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

var (
	kubeClient     *KubeClient
	Namespace      string
	NODE_GAUGES    map[string]*NodeGauges
	NodePanel      *ui.List
	ContainerPanel *ui.Table
	EventsPanel    *ui.Table
	ContainerMaxes map[string]*ContainerMax
	CpuColumn      []ui.GridBufferer
	MemoryColumn   []ui.GridBufferer
)

var log = logrus.New()

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

type ContainerMax struct {
	Cpu        *CpuResource
	CpuTime    time.Time
	Memory     *MemoryResource
	MemoryTime time.Time
}

type ContainerInfo struct {
	Name      string
	Status    *ContainerStatus
	Resources *NodeResources
}

type ContainerStatus struct {
	Name     string
	Status   string
	Ready    bool
	Restarts int
	Age      string
}

type NodeGauges struct {
	Node        api_v1.Node
	CpuGauge    *ui.Gauge
	MemoryGauge *ui.Gauge
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
