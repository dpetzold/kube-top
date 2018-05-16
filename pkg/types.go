package main

import (
	ui "github.com/dpetzold/termui"
	api_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
)

const (
	NODE_DISPLAY_COUNT = 3
)

var (
	kubeClient  *KubeClient
	Namespace   string
	NODE_GAUGES map[string]*NodeGauges
)

type KubeClient struct {
	clientset      *kubernetes.Clientset
	heapsterClient *metricsutil.HeapsterMetricsClient
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
}

type CpuResource struct {
	*resource.Quantity
}

type MemoryResource struct {
	*resource.Quantity
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
