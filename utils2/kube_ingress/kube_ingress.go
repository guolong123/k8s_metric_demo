package kube_ingress

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/ingress-metrics.md

import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_ingress_", &IngressMetrics{groupField: [2]string{"namespace", "ingress"}})
}

type IngressMetrics struct {
	lines      []utils2.MetricLine
	groupField [2]string
	Ingresses  map[string]Ingress
}

type Ingress struct {
	Namespace                      string
	IngressName                    string
	IngressLabels                  []string
	IngressCreated                 float64
	IngressMetadataResourceVersion float64
	IngressPath                    map[string]string
	IngressTls                     map[string]string
}

func (m *IngressMetrics) Group() {
	if m.Ingresses == nil {
		m.Ingresses = make(map[string]Ingress)
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
		_, ok := m.Ingresses[groupFields]
		if !ok {
			m.Ingresses[groupFields] = Ingress{Namespace: groupFieldList[0], IngressName: groupFieldList[1]}
		}
		m.GetTime(line, groupFields)
		m.GetIngressPath(line, groupFields)
		m.GetIngressLabel(line, groupFields)
		m.GetIngressTls(line, groupFields)
	}
}

func (m *IngressMetrics) GetIngressLabel(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_ingress_labels" {
		return
	}
	ingress := m.Ingresses[groupField]
	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			ingress.IngressLabels = append(ingress.IngressLabels, key+"="+value)
		}
	}
	m.Ingresses[groupField] = ingress
}

func (m *IngressMetrics) GetTime(line utils2.MetricLine, groupField string) {
	if !utils2.In(line.Type, []interface{}{"kube_ingress_created", "kube_ingress_metadata_resource_version"}) {
		return
	}
	ingress := m.Ingresses[groupField]
	number := utils2.ENum2float64(line.Value)
	switch line.Type {
	case "kube_ingress_created":
		ingress.IngressCreated = number
	case "kube_ingress_metadata_resource_version":
		ingress.IngressMetadataResourceVersion = number
	}
	m.Ingresses[groupField] = ingress
}

func (m *IngressMetrics) GetIngressPath(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_ingress_path" {
		return
	}

	ingress := m.Ingresses[groupField]
	if ingress.IngressPath == nil {
		ingress.IngressPath = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if !utils2.In(key, []interface{}{"namespace", "ingress"}) {
			ingress.IngressPath[key] = value
		}
	}
	m.Ingresses[groupField] = ingress
}
func (m *IngressMetrics) GetIngressTls(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_ingress_tls" {
		return
	}
	ingress := m.Ingresses[groupField]
	if ingress.IngressTls == nil {
		ingress.IngressTls = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if !utils2.In(key, []interface{}{"namespace", "ingress"}) {
			ingress.IngressTls[key] = value
		}
	}
	m.Ingresses[groupField] = ingress
}

func (m *IngressMetrics) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *IngressMetrics) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Ingress
	for _, v := range m.Ingresses {
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
