package kube_Namespace

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/namespace-metrics.md

import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_namespace_", &NamespaceMetric{groupField: [1]string{"namespace"}})
}

type NamespaceMetric struct {
	lines      []utils2.MetricLine
	groupField [1]string
	Namespaces map[string]Namespace
}

type Namespace struct {
	Timestamp                int64             `json:"timestamp"`
	Type                     string            `json:"type"`
	Namespace                string            `json:"namespace"`
	NamespaceCreated         float64           `json:"created"`
	NamespaceLabels          []string          `json:"labels"`
	NamespaceStatusCondition map[string]string `json:"status_condition"`
	NamespaceStatusPhase     string            `json:"status_phase"`
}

func (m *NamespaceMetric) Group() {
	if m.Namespaces == nil {
		m.Namespaces = make(map[string]Namespace)
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
		_, ok := m.Namespaces[groupFields]
		if !ok {
			m.Namespaces[groupFields] = Namespace{Timestamp: utils2.Timestamp, Type: "namespace", Namespace: groupFieldList[0]}
		}
		m.GetNamespaceLabels(line, groupFields)
		m.GetNamespaceCreated(line, groupFields)
		m.GetNamespaceStatusCondition(line, groupFields)
		m.GetNamespaceStatusPhase(line, groupFields)
	}
}

func (m *NamespaceMetric) GetNamespaceCreated(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_namespace_created" {
		return
	}
	ns := m.Namespaces[groupField]
	number := utils2.ENum2float64(line.Value)
	ns.NamespaceCreated = number
	m.Namespaces[groupField] = ns
}

func (m *NamespaceMetric) GetNamespaceLabels(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_namespace_labels" {
		return
	}
	ns := m.Namespaces[groupField]
	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			ns.NamespaceLabels = append(ns.NamespaceLabels, key+"="+value)
		}
	}
	m.Namespaces[groupField] = ns
}

func (m *NamespaceMetric) GetNamespaceStatusCondition(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_namespace_status_condition" || line.Value != "1" {
		return
	}
	ns := m.Namespaces[groupField]
	if ns.NamespaceStatusCondition == nil {
		ns.NamespaceStatusCondition = make(map[string]string)
	}

	ns.NamespaceStatusCondition[line.Attribute["condition"]] = line.Attribute["status"]
	m.Namespaces[groupField] = ns
}

func (m *NamespaceMetric) GetNamespaceStatusPhase(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_namespace_status_phase" || line.Value != "1" {
		return
	}
	ns := m.Namespaces[groupField]
	ns.NamespaceStatusPhase = line.Attribute["phase"]
	m.Namespaces[groupField] = ns
}

func (m *NamespaceMetric) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m NamespaceMetric) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Namespace
	for _, v := range m.Namespaces {
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
