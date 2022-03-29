package kube_daemonSet

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/daemonset-metrics.md
import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_daemonset_", &DaemonSetMetrics{groupField: [2]string{"namespace", "daemonset"}})
}

type DaemonSetMetrics struct {
	lines      []utils2.MetricLine
	groupField [2]string
	DaemonSets map[string]DaemonSet
}

type DaemonSet struct {
	Namespace                          string
	DaemonSetName                      string
	DaemonSetCreated                   float64
	DaemonSetStatusCurrentNumber       int
	DaemonSetStatusDesiredNumber       int
	DaemonSetStatusNumberAvailable     int
	DaemonSetStatusNumberMissScheduled int
	DaemonSetStatusNumberReady         int
	DaemonSetStatusNumberUnavailable   int
	DaemonSetMetadataGeneration        int
	DaemonSetLabels                    []string
}

func (m *DaemonSetMetrics) Group() {
	if m.DaemonSets == nil {
		m.DaemonSets = make(map[string]DaemonSet)
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
		_, ok := m.DaemonSets[groupFields]
		if !ok {
			m.DaemonSets[groupFields] = DaemonSet{Namespace: groupFieldList[0], DaemonSetName: groupFieldList[1]}
		}
		m.GetDaemonSetCreated(line, groupFields)
		m.GetValue(line, groupFields)
		m.GetDaemonSetLabel(line, groupFields)
	}
}

func (m *DaemonSetMetrics) GetDaemonSetCreated(line utils2.MetricLine, groupFields string) {
	if !utils2.In(line.Type, []interface{}{"kube_daemonset_created"}) {
		return
	}
	daemonset := m.DaemonSets[groupFields]
	newNum := utils2.ENum2float64(line.Value)
	daemonset.DaemonSetCreated = newNum
	m.DaemonSets[groupFields] = daemonset
}

func (m *DaemonSetMetrics) GetValue(line utils2.MetricLine, groupFields string) {
	if !utils2.In(line.Type, []interface{}{
		"kube_daemonset_status_current_number_scheduled",
		"kube_daemonset_status_desired_number_scheduled",
		"kube_daemonset_status_number_available",
		"kube_daemonset_status_number_misscheduled",
		"kube_daemonset_status_number_ready",
		"kube_daemonset_status_number_unavailable",
		"kube_daemonset_metadata_generation",
	}) {
		return
	}
	daemonset := m.DaemonSets[groupFields]
	newNum := utils2.String2int(line.Value)
	switch line.Type {
	case "kube_daemonset_status_current_number_scheduled":
		daemonset.DaemonSetStatusCurrentNumber = newNum
	case "kube_daemonset_status_desired_number_scheduled":
		daemonset.DaemonSetStatusDesiredNumber = newNum
	case "kube_daemonset_status_number_available":
		daemonset.DaemonSetStatusNumberAvailable = newNum
	case "kube_daemonset_status_number_misscheduled":
		daemonset.DaemonSetStatusNumberMissScheduled = newNum
	case "kube_daemonset_status_number_ready":
		daemonset.DaemonSetStatusNumberReady = newNum
	case "kube_daemonset_status_number_unavailable":
		daemonset.DaemonSetStatusNumberUnavailable = newNum
	case "kube_daemonset_metadata_generation":
		daemonset.DaemonSetMetadataGeneration = newNum
	}

	m.DaemonSets[groupFields] = daemonset
}

func (m *DaemonSetMetrics) GetDaemonSetLabel(line utils2.MetricLine, groupFields string) {
	if line.Type != "kube_daemonset_labels" {
		return
	}

	daemonset := m.DaemonSets[groupFields]

	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			daemonset.DaemonSetLabels = append(daemonset.DaemonSetLabels, key+"="+value)
		}
	}
	m.DaemonSets[groupFields] = daemonset
}

func (m *DaemonSetMetrics) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *DaemonSetMetrics) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []DaemonSet
	for _, v := range m.DaemonSets {
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
