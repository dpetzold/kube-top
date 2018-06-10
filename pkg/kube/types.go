package kube

import (
	"time"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
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
