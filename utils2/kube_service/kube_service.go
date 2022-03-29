package kube_Service

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/service-metrics.md

import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_service_", &ServiceMetrics{groupField: [2]string{"namespace", "service"}})
}

type ServiceMetrics struct {
	lines      []utils2.MetricLine
	groupField [2]string
	Services   map[string]Service
}

type Service struct {
	Namespace             string
	ServiceName           string
	ServiceInfo           map[string]string // Information about service
	ServiceLabels         []string          //  Kubernetes labels converted to Prometheus labels
	ServiceCreated        float64           // Unix creation timestamp
	ServiceSpecType       string            // Type about service
	ServiceSpecExternalIp string            // Service external ips. One series for each ip
}

func (m *ServiceMetrics) Group() {
	if m.Services == nil {
		m.Services = make(map[string]Service)
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
		_, ok := m.Services[groupFields]
		if !ok {
			m.Services[groupFields] = Service{Namespace: groupFieldList[0], ServiceName: groupFieldList[1]}
		}
		m.GetServiceInfo(line, groupFields)
		m.GetServiceLabels(line, groupFields)
		m.GetServiceCreated(line, groupFields)
		m.GetServiceSpecExternalIp(line, groupFields)
		m.GetServiceSpecType(line, groupFields)
	}
}

func (m *ServiceMetrics) GetServiceInfo(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_service_info" {
		return
	}
	service := m.Services[groupField]
	if service.ServiceInfo == nil {
		service.ServiceInfo = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if key == "namespace" || key == "service" {
			continue
		}
		service.ServiceInfo[key] = value
	}
	m.Services[groupField] = service
}

func (m *ServiceMetrics) GetServiceLabels(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_service_labels" {
		return
	}
	service := m.Services[groupField]
	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			service.ServiceLabels = append(service.ServiceLabels, key+"="+value)
		}
	}
	m.Services[groupField] = service
}

func (m *ServiceMetrics) GetServiceCreated(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_service_created" {
		return
	}
	service := m.Services[groupField]
	number := utils2.ENum2float64(line.Value)
	service.ServiceCreated = number
	m.Services[groupField] = service
}

func (m *ServiceMetrics) GetServiceSpecType(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_service_spec_type" {
		return
	}
	service := m.Services[groupField]
	service.ServiceSpecType = line.Attribute["type"]
	m.Services[groupField] = service
}

func (m *ServiceMetrics) GetServiceSpecExternalIp(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_service_spec_external_ip" {
		return
	}
	service := m.Services[groupField]
	service.ServiceSpecExternalIp = line.Attribute["external_ip"]
	m.Services[groupField] = service
}

func (m *ServiceMetrics) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *ServiceMetrics) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Service
	for _, v := range m.Services {
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
