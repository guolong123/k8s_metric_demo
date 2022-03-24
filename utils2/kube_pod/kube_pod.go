package kube_pod

import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_pod_", &PodMetric{groupField: [2]string{"namespace", "pod"}})
}

type PodMetric struct {
	lines      []utils2.MetricLine
	groupField [2]string
	Pods       map[string]Pod
}

type Pod struct {
	Namespace               string                      `json:"Namespace"`
	Pod                     string                      `json:"PodName"`
	KubePodInfo             PodInfo                     `json:"PodInfo"`
	KubePodCreated          float64                     `json:"PodCreated"`
	KubePodStartTime        float64                     `json:"PodStartTime"`
	KubePodCompletionTime   float64                     `json:"PodCompletionTime"`
	KubePodRestartPolicy    string                      `json:"PodRestartPolicy"`
	KubePodOwner            PodOwner                    `json:"PodOwner"`
	KubePodLabels           []string                    `json:"PodLabels"`
	KubePodStatusPhase      string                      `json:"PodStatusPhase"`
	KubePodStatusReady      string                      `json:"PodStatusReady"`
	KubePodStatusScheduled  string                      `json:"PodStatusScheduled"`
	KubePodContainerMap     map[string]PodContainer     `json:"-"`
	KubePodInitContainerMap map[string]PodInitContainer `json:"-"`
	KubePodContainers       []PodContainer              `json:"PodContainers"`
	KubePodInitContainers   []PodInitContainer          `json:"PodInitContainers"`
}

type PodInfo struct {
	HostIp        string
	PodIp         string
	uid           string
	Node          string
	CreatedByKind string
	CreatedByName string
	PriorityClass string
}

type PodOwner struct {
	OwnerKind         string
	OwnerName         string
	OwnerIsController bool
}

func (m *PodMetric) Group() {
	if m.Pods == nil {
		m.Pods = make(map[string]Pod)
	}
	for _, line := range m.lines {
		var groupFieldList []string
		for _, v := range m.groupField {
			groupField, ok := line.Attribute[v]
			if !ok {
				continue
			}
			groupFieldList = append(groupFieldList, groupField)
		}
		groupFields := strings.Join(groupFieldList, "_")
		_, ok := m.Pods[groupFields]
		if !ok {
			m.Pods[groupFields] = Pod{Namespace: groupFieldList[0], Pod: groupFieldList[1]}
		}

		m.GetPodInfo(line, groupFields)
		m.GetTimeField(line, groupFields)
		m.GetPodOwner(line, groupFields)
		m.GetContainer(line, groupFields)
		m.GetInitContainer(line, groupFields)
		m.getPodStatus(line, groupFields)
		m.getPodLabel(line, groupFields)

	}
	m.ConvertContainer()
	m.ConvertInitContainer()
}

func (m *PodMetric) ConvertContainer() {
	for k, pod := range m.Pods {
		for _, container := range pod.KubePodContainerMap {
			pod.KubePodContainers = append(pod.KubePodContainers, container)
		}
		m.Pods[k] = pod
	}

}

func (m *PodMetric) ConvertInitContainer() {
	for k, pod := range m.Pods {
		for _, container := range pod.KubePodInitContainerMap {
			pod.KubePodInitContainers = append(pod.KubePodInitContainers, container)
		}
		m.Pods[k] = pod
	}
}
func (m *PodMetric) GetPodInfo(line utils2.MetricLine, groupFields string) {
	if line.Type != "kube_pod_info" {
		return
	}
	kubePodInfo := PodInfo{}
	for key, value := range line.Attribute {
		switch key {
		case "host_ip":
			kubePodInfo.HostIp = value
		case "pod_ip":
			kubePodInfo.PodIp = value
		case "uid":
			kubePodInfo.uid = value
		case "node":
			kubePodInfo.Node = value
		case "created_by_kind":
			kubePodInfo.CreatedByKind = value
		case "created_by_name":
			kubePodInfo.CreatedByName = value
		case "priority_class":
			kubePodInfo.PriorityClass = value
		}
	}

	pod := m.Pods[groupFields]
	pod.KubePodInfo = kubePodInfo
	m.Pods[groupFields] = pod
}

func (m *PodMetric) GetPodRestartPolicy(line utils2.MetricLine, groupFields string) {
	if line.Type != "kube_pod_restart_policy" {
		return
	}
	pod := m.Pods[groupFields]
	for key, value := range line.Attribute {
		if key == "type" {
			pod.KubePodRestartPolicy = value
		}
	}
	m.Pods[groupFields] = pod
}

func (m *PodMetric) GetTimeField(line utils2.MetricLine, groupFields string) {
	timeFields := []interface{}{"kube_pod_start_time", "kube_pod_created", "kube_pod_completion_time"}
	if !in(line.Type, timeFields) {
		return
	}
	pod := m.Pods[groupFields]
	var newNum float64
	value, ok := line.Value.(string)
	if !ok {
		fmt.Printf("%v not convert to int", line.Value)
	}
	_, err := fmt.Sscanf(value, "%e", &newNum)
	if err != nil {
		fmt.Printf("%v not convert to int", line.Value)
	}
	switch line.Type {
	case "kube_pod_start_time":
		pod.KubePodStartTime = newNum
	case "kube_pod_created":
		pod.KubePodCreated = newNum
	case "kube_pod_completion_time":
		pod.KubePodCompletionTime = newNum
	}
	m.Pods[groupFields] = pod
}

func (m *PodMetric) GetPodOwner(line utils2.MetricLine, groupFields string) {
	if line.Type != "kube_pod_owner" {
		return
	}
	owner := PodOwner{}
	pod := m.Pods[groupFields]
	for key, value := range line.Attribute {
		switch key {
		case "owner_name":
			owner.OwnerName = value
		case "owner_kind":
			owner.OwnerKind = value
		case "owner_is_controller":
			if value == "true" {
				owner.OwnerIsController = true
			} else if value == "false" {
				owner.OwnerIsController = false
			}
		}
	}
	pod.KubePodOwner = owner
	m.Pods[groupFields] = pod
}

func (m *PodMetric) getPodLabel(line utils2.MetricLine, groupFields string) {
	if line.Type != "kube_pod_labels" {
		return
	}

	pod := m.Pods[groupFields]

	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			pod.KubePodLabels = append(pod.KubePodLabels, key+"="+value)
		}
	}
	m.Pods[groupFields] = pod
}

func (m *PodMetric) getPodStatus(line utils2.MetricLine, groupFields string) {
	if !in(line.Type, []interface{}{"kube_pod_status_scheduled", "kube_pod_status_ready", "kube_pod_status_phase"}) {
		return
	}
	pod := m.Pods[groupFields]
	for key, value := range line.Attribute {
		if (key == "condition" || key == "phase") && line.Value == "1" {
			if line.Type == "kube_pod_status_scheduled" {
				pod.KubePodStatusScheduled = value
			} else if line.Type == "kube_pod_status_ready" {
				pod.KubePodStatusReady = value
			} else if line.Type == "kube_pod_status_phase" {
				pod.KubePodStatusPhase = value
			}
		}
	}
	m.Pods[groupFields] = pod
}

func (m *PodMetric) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *PodMetric) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Pod
	for _, v := range m.Pods {
		listData = append(listData, v)
	}
	total := len(listData)
	var jsonData []byte
	var err error
	for index < len(listData) {
		if total < 500 {
			jsonData, err = json.Marshal(listData[index : index+total])
		} else {
			jsonData, err = json.Marshal(listData[index : index+500])
		}

		if err == nil {
			todoFunc(string(jsonData))
		} else {
			fmt.Println(err)
		}
		index += 500
		total -= 500
	}

}

func eNum2float64(enum interface{}) float64 {
	var newNum float64
	value := enum.(string)

	_, err := fmt.Sscanf(value, "%e", &newNum)
	if err != nil {
		fmt.Printf("%v not convert to int", enum)
	}
	return newNum
}

func in(obj interface{}, objList []interface{}) bool {
	for _, v := range objList {
		if obj == v {
			return true
		}
	}
	return false
}
