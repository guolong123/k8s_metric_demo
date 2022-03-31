package kube_node

// https://github.com/kubernetes/kube-state-metrics/blob/master/docs/node-metrics.md
import (
	"encoding/json"
	"fmt"
	"k8s_metric/utils2"
	"strings"
)

func init() {
	utils2.Register("kube_node_", &NodeMetrics{groupField: [1]string{"node"}})
}

type NodeMetrics struct {
	lines      []utils2.MetricLine
	groupField [1]string
	Nodes      map[string]Node
}

type Node struct {
	Timestamp             int64             `json:"timestamp"`
	Type                  string            `json:"type"`
	Node                  string            `json:"node"`
	NodeInfo              map[string]string `json:"info"`               // Information about a cluster node
	NodeLabels            []string          `json:"labels"`             // Kubernetes labels converted to Prometheus labels
	NodeRole              string            `json:"role"`               // The role of a cluster node
	NodeSpecUnschedulable bool              `json:"spec_unschedulable"` // Whether a node can schedule new pods
	NodeSpecTaint         map[string]string `json:"spec_taint"`         // The taint of a cluster node
	NodeStatusCapacity    map[string]string `json:"status_capacity"`    // The capacity for different resources of a node
	NodeStatusAllocatable map[string]string `json:"status_allocatable"` // The allocatable for different resources of a node that are available for scheduling.
	NodeStatusCondition   map[string]string `json:"status_condition"`   // The condition of a cluster node
	NodeCreated           float64           `json:"created"`            // Unix creation timestamp
}

func (m *NodeMetrics) Group() {
	if m.Nodes == nil {
		m.Nodes = make(map[string]Node)
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
		_, ok := m.Nodes[groupFields]
		if !ok {
			m.Nodes[groupFields] = Node{Timestamp: utils2.Timestamp, Type: "node", Node: groupFieldList[0]}
		}
		m.GetNodeInfo(line, groupFields)
		m.GetNodeLabels(line, groupFields)
		m.GetNodeCreated(line, groupFields)
		m.GetNodeRole(line, groupFields)
		m.GetNodeSpecTaint(line, groupFields)
		m.GetNodeSpecUnscheduled(line, groupFields)
		m.GetNodeStatus(line, groupFields)
		m.GetNodeStatusCondition(line, groupFields)
	}
}

func (m *NodeMetrics) GetNodeInfo(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_info" {
		return
	}
	node := m.Nodes[groupField]
	if node.NodeInfo == nil {
		node.NodeInfo = make(map[string]string)
	}

	for key, value := range line.Attribute {
		if key == "node" {
			continue
		}
		node.NodeInfo[key] = value
	}
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeCreated(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_created" {
		return
	}
	node := m.Nodes[groupField]
	number := utils2.ENum2float64(line.Value)
	node.NodeCreated = number
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeLabels(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_labels" {
		return
	}
	node := m.Nodes[groupField]
	for key, value := range line.Attribute {
		if strings.HasPrefix(key, "label_") {
			node.NodeLabels = append(node.NodeLabels, key+"="+value)
		}
	}
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeRole(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_role" {
		return
	}
	node := m.Nodes[groupField]
	for key, value := range line.Attribute {
		if key == "role" {
			node.NodeRole = value
		}
	}
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeSpecUnscheduled(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_spec_unschedulable" {
		return
	}
	node := m.Nodes[groupField]
	if line.Value == "0" {
		node.NodeSpecUnschedulable = false
	} else if line.Value == "1" {
		node.NodeSpecUnschedulable = true
	}
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeSpecTaint(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_spec_taint" {
		return
	}

	node := m.Nodes[groupField]
	if node.NodeSpecTaint == nil {
		node.NodeSpecTaint = make(map[string]string)
	}
	for key, value := range line.Attribute {
		if key == "node" {
			continue
		}
		node.NodeSpecTaint[key] = value
	}
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeStatus(line utils2.MetricLine, groupField string) {
	if !utils2.In(line.Type, []interface{}{"kube_node_status_capacity", "kube_node_status_allocatable"}) {
		return
	}
	node := m.Nodes[groupField]
	strValue := line.Value.(string)
	if node.NodeStatusCapacity == nil {
		node.NodeStatusCapacity = make(map[string]string)
	}
	if node.NodeStatusAllocatable == nil {
		node.NodeStatusAllocatable = make(map[string]string)
	}
	switch line.Type {
	case "kube_node_status_capacity":
		node.NodeStatusCapacity[line.Attribute["resource"]] = strValue + "(" + line.Attribute["unit"] + ")"

	case "kube_node_status_allocatable":
		node.NodeStatusAllocatable[line.Attribute["resource"]] = strValue + "(" + line.Attribute["unit"] + ")"
	}

	m.Nodes[groupField] = node
}

func (m *NodeMetrics) GetNodeStatusCondition(line utils2.MetricLine, groupField string) {
	if line.Type != "kube_node_status_condition" || line.Value != "1" {
		return
	}
	node := m.Nodes[groupField]
	if node.NodeStatusCondition == nil {
		node.NodeStatusCondition = make(map[string]string)
	}
	node.NodeStatusCondition[line.Attribute["condition"]] = line.Attribute["status"]
	m.Nodes[groupField] = node
}

func (m *NodeMetrics) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}

func (m *NodeMetrics) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Node
	for _, v := range m.Nodes {
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
