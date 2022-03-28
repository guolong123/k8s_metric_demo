package kube_pod

import (
	"fmt"
	"k8s_metric/utils2"
	"strconv"
)

type PodContainer struct {
	Container                                  string            `json:"name" group:"Container"`
	KubePodContainerInfo                       map[string]string `json:"info"`
	KubePodContainerStatusWaiting              bool              `json:"status_waiting"`
	KubePodContainerStatusWaitingReason        string            `json:"status_waiting_reason"`
	KubePodContainerStatusRunning              bool              `json:"status_running"`
	KubePodContainerStateStarted               float64           `json:"state_started"`
	KubePodContainerStatusTerminated           bool              `json:"status_terminated"`
	KubePodContainerStatusTerminatedReason     string            `json:"status_terminated_reason"`
	KubePodContainerStatusLastTerminatedReason string            `json:"status_last_terminated_reason"`
	KubePodContainerStatusReady                bool              `json:"status_ready"`
	KubePodContainerStatusRestartTotal         int               `json:"status_restarts_total"`
	KubePodContainerResourceRequests           map[string]string `json:"resource_requests"`
	KubePodContainerResourceLimits             map[string]string `json:"resource_limits"`
}

func (m *PodMetric) GetContainer(line utils2.MetricLine, groupFields string) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_container_info",
		"kube_pod_container_status_waiting",
		"kube_pod_container_status_waiting_reason",
		"kube_pod_container_status_running",
		"kube_pod_container_state_started",
		"kube_pod_container_status_terminated",
		"kube_pod_container_status_terminated_reason",
		"kube_pod_container_status_last_terminated_reason",
		"kube_pod_container_status_ready",
		"kube_pod_container_status_restarts_total",
		"kube_pod_container_resource_limits",
		"kube_pod_container_resource_requests",
	}) {
		return
	}
	pod := m.Pods[groupFields]
	containerGroupField := groupFields + "_" + line.Attribute["container"]

	container, ok := m.Pods[groupFields].KubePodContainerMap[containerGroupField]
	if !ok {
		container = PodContainer{
			Container: line.Attribute["container"]}
	}
	if pod.KubePodContainerMap == nil {
		pod.KubePodContainerMap = make(map[string]PodContainer)
	}
	switch line.Type {
	case "kube_pod_container_info":
		container.GetContainerInfo(line)
	case "kube_pod_container_status_waiting":
		container.GetContainerStatus(line)
	case "kube_pod_container_status_running":
		container.GetContainerStatus(line)
	case "kube_pod_container_status_terminated":
		container.GetContainerStatus(line)
	case "kube_pod_container_status_waiting_reason":
		container.GetContainerStatusReason(line)
	case "kube_pod_container_status_terminated_reason":
		container.GetContainerStatusReason(line)
	case "kube_pod_container_status_last_terminated_reason":
		container.GetContainerStatusReason(line)
	case "kube_pod_container_status_restarts_total":
		container.GetContainerRestartCount(line)
	case "kube_pod_container_resource_limits":
		container.GetContainerResource(line)
	case "kube_pod_container_resource_requests":
		container.GetContainerResource(line)
	}
	pod.KubePodContainerMap[containerGroupField] = container
	m.Pods[groupFields] = pod
}

func (c *PodContainer) GetContainerInfo(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_container_info"}) {
		return
	}
	if c.KubePodContainerInfo == nil {
		c.KubePodContainerInfo = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if utils2.In(key, []interface{}{"image", "image_id", "container_id"}) {
			c.KubePodContainerInfo[key] = value
		}
	}
}

func (c *PodContainer) GetContainerStateStarted(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"container_state_started"}) {
		return
	}
	var newNum float64
	value, ok := line.Value.(string)
	if !ok {
		fmt.Printf("%v not convert to int", line.Value)
	}
	_, err := fmt.Sscanf(value, "%e", &newNum)
	if err != nil {
		fmt.Printf("%v not convert to int", line.Value)
	}
	c.KubePodContainerStateStarted = newNum
}

func (c *PodContainer) GetContainerStatus(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_container_status_waiting",
		"kube_pod_container_status_running",
		"kube_pod_container_status_terminated",
		"kube_pod_container_status_ready",
	}) {
		return
	}
	switch line.Type {
	case "kube_pod_container_status_waiting":
		if line.Value == "0" {
			c.KubePodContainerStatusWaiting = false
		} else if line.Value == "1" {
			c.KubePodContainerStatusWaiting = true
		}
	case "kube_pod_container_status_running":
		if line.Value == "0" {
			c.KubePodContainerStatusRunning = false
		} else if line.Value == "1" {
			c.KubePodContainerStatusRunning = true
		}
	case "kube_pod_container_status_terminated":
		if line.Value == "0" {
			c.KubePodContainerStatusTerminated = false
		} else if line.Value == "1" {
			c.KubePodContainerStatusTerminated = true
		}
	case "kube_pod_container_status_ready":
		if line.Value == "0" {
			c.KubePodContainerStatusReady = false
		} else if line.Value == "1" {
			c.KubePodContainerStatusReady = true
		}
	}
}

func (c *PodContainer) GetContainerStatusReason(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_container_status_waiting_reason",
		"kube_pod_container_status_terminated_reason",
		"kube_pod_container_status_last_terminated_reason"}) {
		return
	}
	for key, value := range line.Attribute {
		switch line.Type {
		case "kube_pod_container_status_waiting_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodContainerStatusWaitingReason = value
			}
		case "kube_pod_container_status_terminated_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodContainerStatusTerminatedReason = value
			}
		case "kube_pod_container_status_last_terminated_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodContainerStatusLastTerminatedReason = value
			}
		}
	}
}

func (c *PodContainer) GetContainerRestartCount(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_container_status_restarts_total"}) {
		return
	}
	value := line.Value.(string)
	number, err := strconv.Atoi(value)
	if err == nil {
		c.KubePodContainerStatusRestartTotal = number
	}
}

func (c *PodContainer) GetContainerResource(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_container_resource_requests", "kube_pod_container_resource_limits"}) {
		return
	}
	if c.KubePodContainerResourceLimits == nil {
		c.KubePodContainerResourceLimits = make(map[string]string)
	}
	if c.KubePodContainerResourceRequests == nil {
		c.KubePodContainerResourceRequests = make(map[string]string)
	}

	switch line.Type {
	case "kube_pod_container_resource_requests":
		for k, v := range line.Attribute {

			if k == "resource" {
				if v == "memory" {
					number := utils2.ENum2float64(line.Value)
					c.KubePodContainerResourceRequests[v] = fmt.Sprintf("%.2f (%s)", number, line.Attribute["unit"])
				} else if v == "cpu" {
					value := line.Value.(string)
					c.KubePodContainerResourceRequests[v] = fmt.Sprintf("%s (%s)", value, line.Attribute["unit"])
				}
			}
		}
	case "kube_pod_container_resource_limits":
		for k, v := range line.Attribute {
			if k == "resource" {
				if v == "memory" {
					number := utils2.ENum2float64(line.Value)
					c.KubePodContainerResourceLimits[v] = fmt.Sprintf("%.2f (%s)", number, line.Attribute["unit"])
				} else if v == "cpu" {
					value := line.Value.(string)
					c.KubePodContainerResourceLimits[v] = fmt.Sprintf("%s (%s)", value, line.Attribute["unit"])
				}
			}
		}
	}
}
