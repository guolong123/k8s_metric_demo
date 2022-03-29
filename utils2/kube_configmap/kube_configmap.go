package kube_configmap

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/configmap-metrics.md
import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_configmap_", &ConfigMapMetrics{groupField: [2]string{"namespace", "configmap"}})
}

type ConfigMapMetrics struct {
	lines      []utils2.MetricLine
	groupField [2]string
	ConfigMaps map[string]ConfigMap
}

type ConfigMap struct {
	Namespace                        string
	ConfigMapName                    string
	ConfigMapCreated                 float64
	ConfigMapMetaDataResourceVersion float64
}

func (m *ConfigMapMetrics) Group() {
	if m.ConfigMaps == nil {
		m.ConfigMaps = make(map[string]ConfigMap)
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
		_, ok := m.ConfigMaps[groupFields]
		if !ok {
			m.ConfigMaps[groupFields] = ConfigMap{Namespace: groupFieldList[0], ConfigMapName: groupFieldList[1]}
		}
		m.GetCreated(line, groupFields)
	}
}

func (m *ConfigMapMetrics) GetCreated(line utils2.MetricLine, groupFields string) {
	if !utils2.In(line.Type, []interface{}{"kube_configmap_created", "kube_configmap_metadata_resource_version"}) {
		return
	}
	configMap := m.ConfigMaps[groupFields]
	newNum := utils2.ENum2float64(line.Value)
	switch line.Type {
	case "kube_configmap_created":
		configMap.ConfigMapCreated = newNum
	case "kube_configmap_metadata_resource_version":
		configMap.ConfigMapMetaDataResourceVersion = newNum
	}
	m.ConfigMaps[groupFields] = configMap
}

func (m *ConfigMapMetrics) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *ConfigMapMetrics) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []ConfigMap
	for _, v := range m.ConfigMaps {
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
