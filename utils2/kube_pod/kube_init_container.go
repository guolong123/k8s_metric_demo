package kube_pod

import (
	"fmt"
	"k8s_metric/utils2"
	"strconv"
)

type PodInitContainer struct {
	Container                                      string            `json:"init_container"`
	KubePodInitContainerInfo                       map[string]string `json:"info"`
	KubePodInitContainerStatusWaiting              bool              `json:"status_waiting"`
	KubePodInitContainerStatusWaitingReason        string            `json:"status_waiting_reason" `
	KubePodInitContainerStatusRunning              bool              `json:"status_running"`
	KubePodInitContainerStatusTerminated           bool              `json:"status_terminated"`
	KubePodInitContainerStatusTerminatedReason     string            `json:"status_terminated_reason" `
	KubePodInitContainerStatusLastTerminatedReason string            `json:"status_last_terminated_reason" `
	KubePodInitContainerStatusReady                bool              `json:"status_ready"`
	KubePodInitContainerStatusRestartsTotal        int               `json:"status_restarts_total"`
	KubePodInitContainerResourceRequests           map[string]string `json:"resource_requests"`
	KubePodInitContainerResourceLimits             map[string]string `json:"resource_limits"`
}

func (m *PodMetric) GetInitContainer(line utils2.MetricLine, groupFields string) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_init_container_info",
		"kube_pod_init_container_status_waiting",
		"kube_pod_init_container_status_waiting_reason",
		"kube_pod_init_container_status_running",
		"kube_pod_init_container_state_started",
		"kube_pod_init_container_status_terminated",
		"kube_pod_init_container_status_terminated_reason",
		"kube_pod_init_container_status_last_terminated_reason",
		"kube_pod_init_container_status_ready",
		"kube_pod_init_container_status_restarts_total",
		"kube_pod_init_container_resource_limits",
		"kube_pod_init_container_resource_requests",
	}) {
		return
	}
	pod := m.Pods[groupFields]
	containerGroupField := groupFields + "_" + line.Attribute["container"]

	container, ok := m.Pods[groupFields].KubePodInitContainerMap[containerGroupField]
	if !ok {
		container = PodInitContainer{
			Container: line.Attribute["container"]}
	}
	if pod.KubePodInitContainerMap == nil {
		pod.KubePodInitContainerMap = make(map[string]PodInitContainer)
	}
	switch line.Type {
	case "kube_pod_init_container_info":
		container.GetInitContainerInfo(line)
	case "kube_pod_init_container_status_waiting":
		container.GetInitContainerStatus(line)
	case "kube_pod_init_container_status_running":
		container.GetInitContainerStatus(line)
	case "kube_pod_init_container_status_terminated":
		container.GetInitContainerStatus(line)
	case "kube_pod_init_container_status_waiting_reason":
		container.GetInitContainerStatusReason(line)
	case "kube_pod_init_container_status_terminated_reason":
		container.GetInitContainerStatusReason(line)
	case "kube_pod_init_container_status_last_terminated_reason":
		container.GetInitContainerStatusReason(line)
	case "kube_pod_init_container_status_restarts_total":
		container.GetInitContainerRestartCount(line)
	case "kube_pod_init_container_resource_limits":
		container.GetInitContainerResource(line)
	case "kube_pod_init_container_resource_requests":
		container.GetInitContainerResource(line)
	}
	pod.KubePodInitContainerMap[containerGroupField] = container
	m.Pods[groupFields] = pod
}

func (c *PodInitContainer) GetInitContainerInfo(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_init_container_info"}) {
		return
	}
	if c.KubePodInitContainerInfo == nil {
		c.KubePodInitContainerInfo = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if utils2.In(key, []interface{}{"image", "image_id", "container_id"}) {
			c.KubePodInitContainerInfo[key] = value
		}
	}
}

func (c *PodInitContainer) GetInitContainerStatus(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_init_container_status_waiting",
		"kube_pod_init_container_status_running",
		"kube_pod_init_container_status_terminated",
		"kube_pod_init_container_status_ready",
	}) {
		return
	}
	switch line.Type {
	case "kube_pod_init_container_status_waiting":
		if line.Value == "0" {
			c.KubePodInitContainerStatusWaiting = false
		} else if line.Value == "1" {
			c.KubePodInitContainerStatusWaiting = true
		}
	case "kube_pod_init_container_status_running":
		if line.Value == "0" {
			c.KubePodInitContainerStatusRunning = false
		} else if line.Value == "1" {
			c.KubePodInitContainerStatusRunning = true
		}
	case "kube_pod_init_container_status_terminated":
		if line.Value == "0" {
			c.KubePodInitContainerStatusTerminated = false
		} else if line.Value == "1" {
			c.KubePodInitContainerStatusTerminated = true
		}
	case "kube_pod_init_container_status_ready":
		if line.Value == "0" {
			c.KubePodInitContainerStatusReady = false
		} else if line.Value == "1" {
			c.KubePodInitContainerStatusReady = true
		}
	}
}

func (c *PodInitContainer) GetInitContainerStatusReason(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{
		"kube_pod_init_container_status_waiting_reason",
		"kube_pod_init_container_status_terminated_reason",
		"kube_pod_init_container_status_last_terminated_reason"}) {
		return
	}
	for key, value := range line.Attribute {
		switch line.Type {
		case "kube_pod_init_container_status_waiting_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodInitContainerStatusWaitingReason = value
			}
		case "kube_pod_init_container_status_terminated_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodInitContainerStatusTerminatedReason = value
			}
		case "kube_pod_init_container_status_last_terminated_reason":
			if line.Value == "1" && key == "reason" {
				c.KubePodInitContainerStatusLastTerminatedReason = value
			}
		}
	}
}

func (c *PodInitContainer) GetInitContainerRestartCount(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_container_status_restarts_total"}) {
		return
	}
	value := line.Value.(string)
	number, err := strconv.Atoi(value)
	if err == nil {
		c.KubePodInitContainerStatusRestartsTotal = number
	}
}

func (c *PodInitContainer) GetInitContainerResource(line utils2.MetricLine) {
	if !utils2.In(line.Type, []interface{}{"kube_pod_init_container_resource_requests", "kube_pod_init_container_resource_limits"}) {
		return
	}
	if c.KubePodInitContainerResourceLimits == nil {
		c.KubePodInitContainerResourceLimits = make(map[string]string)
	}
	if c.KubePodInitContainerResourceRequests == nil {
		c.KubePodInitContainerResourceRequests = make(map[string]string)
	}

	switch line.Type {
	case "kube_pod_init_container_resource_requests":
		for k, v := range line.Attribute {

			if k == "resource" {
				if v == "memory" {
					number := utils2.ENum2float64(line.Value)
					c.KubePodInitContainerResourceRequests[v] = fmt.Sprintf("%.2f (%s)", number, line.Attribute["unit"])
				} else if v == "cpu" {
					value := line.Value.(string)
					c.KubePodInitContainerResourceRequests[v] = fmt.Sprintf("%s (%s)", value, line.Attribute["unit"])
				}
			}
		}
	case "kube_pod_init_container_resource_limits":
		for k, v := range line.Attribute {
			if k == "resource" {
				if v == "memory" {
					number := utils2.ENum2float64(line.Value)
					c.KubePodInitContainerResourceLimits[v] = fmt.Sprintf("%.2f (%s)", number, line.Attribute["unit"])
				} else if v == "cpu" {
					value := line.Value.(string)
					c.KubePodInitContainerResourceLimits[v] = fmt.Sprintf("%s (%s)", value, line.Attribute["unit"])
				}
			}
		}
	}
}
