package main

import (
	"fmt"
	"sort"

	api_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/kubectl/metricsutil"
)

func NewKubeClient(
	clientset *kubernetes.Clientset,
	heapsterClient *metricsutil.HeapsterMetricsClient,
) *KubeClient {
	return &KubeClient{
		clientset:      clientset,
		heapsterClient: heapsterClient,
	}
}

// Return all active pods in the specified namespace. Can also limit to a node if specified.
func (k *KubeClient) ActivePods(namespace, nodeName string) ([]api_v1.Pod, error) {

	// selector := fmt.Sprintf("status.phase!=%s,status.phase!=%s", string(api_v1.PodSucceeded), string(api_v1.PodFailed))
	selector := fmt.Sprintf("status.phase!=%s", string(api_v1.PodSucceeded))
	if nodeName != "" {
		selector += fmt.Sprintf(",spec.nodeName=%s", nodeName)
	}

	fieldSelector, err := fields.ParseSelector(selector)
	if err != nil {
		return nil, err
	}

	activePods, err := k.clientset.Core().Pods(
		namespace,
	).List(
		metav1.ListOptions{FieldSelector: fieldSelector.String()},
	)
	if err != nil {
		return nil, err
	}

	return activePods.Items, err
}

// Return the containers request and limit QoS.
func containerRequestsAndLimits(container *api_v1.Container) (reqs api_v1.ResourceList, limits api_v1.ResourceList) {
	reqs, limits = api_v1.ResourceList{}, api_v1.ResourceList{}

	for name, quantity := range container.Resources.Requests {
		if _, ok := reqs[name]; ok {
			panic(fmt.Sprintf("Duplicate key: %s", name))
		} else {
			reqs[name] = *quantity.Copy()
		}
	}

	for name, quantity := range container.Resources.Limits {
		if _, ok := limits[name]; ok {
			panic(fmt.Sprintf("Duplicate key: %s", name))
		} else {
			limits[name] = *quantity.Copy()
		}
	}
	return
}

// Return the nodes capacity.
func NodeCapacity(node *api_v1.Node) api_v1.ResourceList {
	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}
	return allocatable
}

// Return NodeResources struct for the specified object.
func (k *KubeClient) NodeResources(node *api_v1.Node) (*NodeResources, error) {

	metricsList, err := k.heapsterClient.GetNodeMetrics(node.GetName(), labels.Everything().String())
	if err != nil {
		return nil, err
	}

	if len(metricsList.Items) != 1 {
		return nil, fmt.Errorf("Got bad number of results from client.GetNodeMetrics")
	}

	pods, err := k.ActivePods(Globals.Namespace, node.GetName())
	if err != nil {
		return nil, err
	}

	metrics := metricsList.Items[0]

	capacity := NodeCapacity(node)

	cpuCapacity := capacity[api_v1.ResourceCPU]
	memoryCapacity := capacity[api_v1.ResourceMemory]

	cpuQuantity := metrics.Usage[api_v1.ResourceCPU]
	memoryQuantity := metrics.Usage[api_v1.ResourceMemory]

	cpuUsage := NewCpuResource(cpuQuantity.MilliValue())
	memoryUsage := NewMemoryResource(memoryQuantity.Value())

	return &NodeResources{
		Name: node.GetName(),
		Pods: len(pods),

		CpuUsage:    NewCpuResource(cpuQuantity.MilliValue()),
		CpuCapacity: NewCpuResource(cpuCapacity.MilliValue()),
		PercentCpu:  cpuUsage.calcPercentage(capacity.Cpu()),

		MemoryUsage:    NewMemoryResource(memoryQuantity.Value()),
		MemoryCapacity: NewMemoryResource(memoryCapacity.Value()),
		PercentMemory:  memoryUsage.calcPercentage(capacity.Memory()),
	}, nil
}

// Returns the resource usage of the pods in the specified namespace.
func (k *KubeClient) ContainerNodeResources(namespace string) (map[string]*ContainerUsage, error) {

	allNamespaces := false
	if namespace == "" {
		allNamespaces = true
	}

	metricsList, err := k.heapsterClient.GetPodMetrics(namespace, "", allNamespaces, labels.Everything())
	if err != nil {
		return nil, err
	}

	usage := make(map[string]*ContainerUsage)

	for _, item := range metricsList.Items {

		for _, metrics := range item.Containers {

			cpuQuantity := metrics.Usage[api_v1.ResourceCPU]
			memoryQuantity := metrics.Usage[api_v1.ResourceMemory]
			name := fmt.Sprintf("%s/%s", item.ObjectMeta.Name, metrics.Name)

			usage[name] = &ContainerUsage{
				Name:        name,
				CpuUsage:    NewCpuResource(cpuQuantity.MilliValue()),
				MemoryUsage: NewMemoryResource(memoryQuantity.Value()),
			}
		}

	}

	return usage, nil
}

func evaluatePod(pod *api_v1.Pod) map[string]*ContainerStatus {

	containers := make(map[string]*ContainerStatus)

	status := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		status = pod.Status.Reason
	}

	// why is this in reverse?
	for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
		ready := false
		container := pod.Status.ContainerStatuses[i]
		restarts := int(container.RestartCount)

		if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
			status = container.State.Waiting.Reason
		} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
			status = container.State.Terminated.Reason
		} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
			if container.State.Terminated.Signal != 0 {
				status = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
			} else {
				status = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
			}
		} else if container.Ready && container.State.Running != nil {
			ready = true
		}

		containers[container.Name] = &ContainerStatus{
			Name:     container.Name,
			Status:   status,
			Ready:    ready,
			Restarts: restarts,
			Age:      pod.CreationTimestamp.Time,
		}
	}

	return containers

}

func (k *KubeClient) Containers(namespace string) ([]*ContainerInfo, error) {

	resource_map, err := k.ContainerNodeResources(namespace)
	if err != nil {
		return nil, err
	}

	active, err := k.ActivePods(namespace, "")
	if err != nil {
		return nil, err
	}

	var rows []*ContainerInfo
	for _, pod := range active {
		containers := evaluatePod(&pod)

		for container, status := range containers {
			name := fmt.Sprintf("%s/%s", pod.GetName(), container)
			resources := resource_map[name]

			var containerUsage *ContainerUsage

			if resources != nil {
				containerUsage = &ContainerUsage{
					CpuUsage:    resources.CpuUsage,
					MemoryUsage: resources.MemoryUsage,
				}
			} else {
				containerUsage = &ContainerUsage{
					CpuUsage:    NewCpuResource(0),
					MemoryUsage: NewMemoryResource(0),
				}
			}

			rows = append(rows, &ContainerInfo{
				Name:           name,
				ContainerUsage: containerUsage,
				ContainerStatus: &ContainerStatus{
					Status:   status.Status,
					Ready:    status.Ready,
					Restarts: status.Restarts,
					Age:      status.Age,
				},
			})
		}
	}
	return rows, nil
}

// Return a list of container resources for all containers running on the specified node
func (k *KubeClient) NodeContainerResources(namespace, nodeName string) (resources []*ContainerResources, err error) {

	mc := k.clientset.Core().Nodes()
	node, err := mc.Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	activePodsList, err := k.ActivePods(namespace, nodeName)
	if err != nil {
		return nil, err
	}

	capacity := NodeCapacity(node)

	// https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/describe.go#L2970
	for _, pod := range activePodsList {
		for _, container := range pod.Spec.Containers {
			req, limit := containerRequestsAndLimits(&container)

			_cpuReq := req[api_v1.ResourceCPU]
			cpuReq := NewCpuResource(_cpuReq.MilliValue())

			_cpuLimit := limit[api_v1.ResourceCPU]
			cpuLimit := NewCpuResource(_cpuLimit.MilliValue())

			_memoryReq := req[api_v1.ResourceMemory]
			memoryReq := NewMemoryResource(_memoryReq.Value())

			_memoryLimit := limit[api_v1.ResourceMemory]
			memoryLimit := NewMemoryResource(_memoryLimit.Value())

			resources = append(resources, &ContainerResources{
				Name:               fmt.Sprintf("%s/%s", pod.GetName(), container.Name),
				Namespace:          pod.GetNamespace(),
				CpuReq:             cpuReq,
				CpuLimit:           cpuLimit,
				PercentCpuReq:      cpuReq.calcPercentage(capacity.Cpu()),
				PercentCpuLimit:    cpuLimit.calcPercentage(capacity.Cpu()),
				MemReq:             memoryReq,
				MemLimit:           memoryLimit,
				PercentMemoryReq:   memoryReq.calcPercentage(capacity.Memory()),
				PercentMemoryLimit: memoryLimit.calcPercentage(capacity.Memory()),
			})
		}
	}

	return resources, nil
}

// Return the resources in use by containers in the cluster as list of ContainerResources
func (k *KubeClient) ContainerResources(namespace string) (resources []*ContainerResources, err error) {
	nodes, err := k.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		nodeUsage, err := k.NodeContainerResources(namespace, node.GetName())
		if err != nil {
			return nil, err
		}
		resources = append(resources, nodeUsage...)
	}

	return resources, nil
}

// Return the total cluster capacity as a ResourceList
func (k *KubeClient) ClusterCapacity() (capacity api_v1.ResourceList, err error) {

	nodes, err := k.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	capacity = api_v1.ResourceList{}

	for _, node := range nodes.Items {

		allocatable := NodeCapacity(&node)

		for name, quantity := range allocatable {
			if value, ok := capacity[name]; ok {
				value.Add(quantity)
				capacity[name] = value
			} else {
				capacity[name] = *quantity.Copy()
			}
		}

	}

	return capacity, nil
}

// Returns events from the specified namespace sorted by timestamp
func (k *KubeClient) Events(namespace string) ([]api_v1.Event, error) {

	eventList, err := k.clientset.Core().Events(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	events := eventList.Items
	sort.Slice(events, func(i, j int) bool {
		return events[i].LastTimestamp.Time.After(events[j].LastTimestamp.Time)
	})

	return events, nil
}

// Returns nodes in the cluster
func (k *KubeClient) Nodes() ([]api_v1.Node, error) {
	nodeList, err := k.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodeList.Items, nil
}
