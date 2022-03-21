package utils

import (
	"fmt"
	"reflect"
)

type Metric interface {
	Group(kubePodContainer *Metric, kubePodInitContainer *Metric, metricLine MetricLine)
}

type PodMetric struct {
	Namespace              string                 `json:"namespace" group:"Namespace"`
	Pod                    string                 `json:"pod" group:"Pod"`
	KubePodInfo            map[string]interface{} `json:"kube_pod_info" get_attr:"true" get_value:"false"`
	KubePodCreated         float64                `json:"kube_pod_created" get_attr:"false" get_value:"true"`
	KubePodStartTime       float64                `json:"kube_pod_start_time" get_attr:"false" get_value:"true"`
	KubePodCompletionTime  float64                `json:"kube_pod_completion_time" get_attr:"false" get_value:"true"`
	KubePodRestartPolicy   map[string]interface{} `json:"kube_pod_restart_policy" get_attr:"true" get_value:"false"`
	KubePodOwner           map[string]interface{} `json:"kube_pod_owner" get_attr:"true" get_value:"false"`
	KubePodLabel           map[string]string      `json:"kube_pod_label" get_attr:"true" get_value:"false"`
	KubePodStatusPhase     map[string]interface{} `json:"kube_pod_status_phase" get_attr:"true" get_value:"false"`
	KubePodStatusReady     map[string]interface{} `json:"kube_pod_status_ready" get_attr:"true" get_value:"false"`
	KubePodStatusScheduled map[string]interface{} `json:"kube_pod_status_scheduled" get_attr:"true" get_value:"false"`
	KubePodContainers      map[string]Metric      `json:"kube_pod_containers" get_attr:"true" get_value:"false"`
	KubePodInitContainers  map[string]Metric      `json:"kube_pod_init_containers" get_attr:"true" get_value:"false"`
}

type KubePodContainer struct {
	Namespace                                  string                 `json:"namespace" group:"Namespace"`
	Pod                                        string                 `json:"pod" group:"Pod"`
	Container                                  string                 `json:"container" group:"Container"`
	KubePodContainerInfo                       map[string]interface{} `json:"kube_pod_container_info" get_attr:"true" get_value:"false"`
	KubePodContainerStatusWaiting              bool                   `json:"kube_pod_container_status_waiting" get_attr:"false" get_value:"true"`
	KubePodContainerStatusWaitingReason        map[string]interface{} `json:"kube_pod_container_status_waiting_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusRunning              bool                   `json:"kube_pod_container_status_running" get_attr:"false" get_value:"true"`
	KubePodContainerStateStarted               float64                `json:"kube_pod_container_state_started" get_attr:"false" get_value:"true"`
	KubePodContainerStatusTerminated           bool                   `json:"kube_pod_container_status_terminated" get_attr:"false" get_value:"true"`
	KubePodContainerStatusTerminatedReason     map[string]interface{} `json:"kube_pod_container_status_terminated_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusLastTerminatedReason map[string]interface{} `json:"kube_pod_container_status_last_terminated_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusReady                bool                   `json:"kube_pod_container_status_ready" get_attr:"false" get_value:"true"`
	KubePodContainerStatusRestartTotal         int                    `json:"kube_pod_container_status_restarts_total" get_attr:"false" get_value:"true"`
	//KubePodContainerResourceRequests           []map[string]interface{} `json:"kube_pod_container_resource_requests" get_attr:"true" get_value:"true"`
	//KubePodContainerResourceLimits             []map[string]interface{} `json:"kube_pod_container_resource_limits" get_attr:"true" get_value:"true"`
}

type KubePodInitContainer struct {
	Pod                                            string                 `json:"pod" group:"Pod"`
	Namespace                                      string                 `json:"namespace" group:"Namespace"`
	Container                                      string                 `json:"container" group:"Container"`
	KubePodInitContainerInfo                       map[string]interface{} `json:"kube_pod_init_container_info"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusWaiting              bool                   `json:"kube_pod_init_container_status_waiting"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusWaitingReason        map[string]interface{} `json:"kube_pod_init_container_status_waiting_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusRunning              bool                   `json:"kube_pod_init_container_status_running" get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusTerminated           bool                   `json:"kube_pod_init_container_status_terminated"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusTerminatedReason     map[string]interface{} `json:"kube_pod_init_container_status_terminated_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusLastTerminatedReason map[string]interface{} `json:"kube_pod_init_container_status_last_terminated_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusReady                bool                   `json:"kube_pod_init_container_status_ready"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusRestartsTotal        int                    `json:"kube_pod_init_container_status_restarts_total" get_attr:"false" get_value:"true"`
	//KubePodInitContainerResourceRequests           []map[string]interface{} `json:"kube_pod_init_container_resource_requests" get_attr:"true" get_value:"false"`
	//KubePodInitContainerResourceLimits             []map[string]interface{} `json:"kube_pod_init_container_resource_limits" get_attr:"true" get_value:"false"`
}

func (m *PodMetric) Group(kubePodContainer *Metric, kubePodInitContainer *Metric, metricLine MetricLine) {
	var valueFalse interface{} = "0"
	st := reflect.TypeOf(*m)

	sv := reflect.ValueOf(m).Elem()
	if m.KubePodInitContainers == nil {
		m.KubePodInitContainers = make(map[string]Metric)
	}
	if m.KubePodContainers == nil {
		m.KubePodContainers = make(map[string]Metric)
	}

	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		jsonTag := field.Tag.Get("json")
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")
		if jsonTag == "kube_pod_containers" && Distribution(metricLine, *kubePodContainer) {
			groupField := GetGroupFields(kubePodContainer, metricLine)
			kubePodContainerValue, ok := m.KubePodContainers[groupField]
			if !ok {
				kubePodContainerValue = *kubePodContainer
			}
			kubePodContainerValue.Group(kubePodContainer, kubePodInitContainer, metricLine)
			m.KubePodContainers[groupField] = kubePodContainerValue
		} else if jsonTag == "kube_pod_init_containers" && Distribution(metricLine, *kubePodInitContainer) {
			groupField := GetGroupFields(kubePodInitContainer, metricLine)
			kubePodInitContainerValue, ok := m.KubePodContainers[groupField]
			if !ok {
				kubePodInitContainerValue = *kubePodInitContainer
			}
			kubePodInitContainerValue.Group(kubePodContainer, kubePodInitContainer, metricLine)
			m.KubePodInitContainers[groupField] = kubePodInitContainerValue

		} else if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			if getAttr == "true" {
				metricAttr := reflect.ValueOf(&metricLine.Attribute).Elem()
				sv.Field(i).Set(metricAttr)
			} else if getValue == "true" {
				// 根据struct原定义的类型进行类型转换
				value, ok := metricLine.Value.(string)
				if !ok {
					fmt.Printf("%v not convert to int", metricLine.Value)
				}
				typeName := field.Type.Name()
				switch typeName {
				case "float64":
					var newNum float64
					_, err := fmt.Sscanf(value, "%e", &newNum)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}
					metricValue := reflect.ValueOf(&newNum).Elem()
					sv.Field(i).Set(metricValue)
				case "string":
					metricValue := reflect.ValueOf(&value).Elem()
					sv.Field(i).Set(metricValue)
				case "int":
					var newValue int
					_, err := fmt.Sscanf(value, "%d", &newValue)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				case "bool":
					var newValue bool
					if value == "0" {
						newValue = false
					} else {
						newValue = true
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				}
			}
		}
	}
}

func (m *KubePodInitContainer) Group(kubePodContainer *Metric, kubePodInitContainer *Metric, metricLine MetricLine) {
	var valueFalse interface{} = "0"
	st := reflect.TypeOf(*m)

	sv := reflect.ValueOf(m).Elem()

	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		jsonTag := field.Tag.Get("json")
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")

		if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			if getAttr == "true" {

				//if jsonTag == "kube_pod_init_container_resource_requests" && Distribution(metricLine, *kubePodInitContainer) {
				//	typeName := field.Type.Name()
				//	fmt.Println(typeName)
				//	//m.KubePodInitContainerResourceLimits = append(m.KubePodInitContainerResourceLimits, )
				//}
				metricAttr := reflect.ValueOf(&metricLine.Attribute).Elem()
				sv.Field(i).Set(metricAttr)
			} else if getValue == "true" {
				// 根据struct原定义的类型进行类型转换
				value, ok := metricLine.Value.(string)
				if !ok {
					fmt.Printf("%v not convert to int", metricLine.Value)
				}
				typeName := field.Type.Name()
				switch typeName {
				case "float64":
					var newNum float64
					_, err := fmt.Sscanf(value, "%e", &newNum)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}

					metricValue := reflect.ValueOf(&newNum).Elem()
					sv.Field(i).Set(metricValue)
				case "string":
					metricValue := reflect.ValueOf(&value).Elem()
					sv.Field(i).Set(metricValue)
				case "int":
					var newValue int
					_, err := fmt.Sscanf(value, "%d", &newValue)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				case "bool":
					var newValue bool
					if value == "0" {
						newValue = false
					} else {
						newValue = true
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				}
			}
		}

	}
}

func (m *KubePodContainer) Group(kubePodContainer *Metric, kubePodInitContainer *Metric, metricLine MetricLine) {
	var valueFalse interface{} = "0"
	st := reflect.TypeOf(*m)

	sv := reflect.ValueOf(m).Elem()

	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		jsonTag := field.Tag.Get("json")
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")
		if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			if getAttr == "true" {
				//if jsonTag == "kube_pod_container_resource_requests" && Distribution(metricLine, *kubePodInitContainer) {
				//	typeName := field.Type.Name()
				//	fmt.Println(typeName)
				//	//m.KubePodInitContainerResourceLimits = append(m.KubePodInitContainerResourceLimits, )
				//}
				metricAttr := reflect.ValueOf(&metricLine.Attribute).Elem()
				sv.Field(i).Set(metricAttr)
			} else if getValue == "true" {
				// 根据struct原定义的类型进行类型转换
				value, ok := metricLine.Value.(string)
				if !ok {
					fmt.Printf("%v not convert to int", metricLine.Value)
				}
				typeName := field.Type.Name()
				switch typeName {
				case "float64":
					var newNum float64
					_, err := fmt.Sscanf(value, "%e", &newNum)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}

					metricValue := reflect.ValueOf(&newNum).Elem()
					sv.Field(i).Set(metricValue)
				case "string":
					metricValue := reflect.ValueOf(&value).Elem()
					sv.Field(i).Set(metricValue)
				case "int":
					var newValue int
					_, err := fmt.Sscanf(value, "%d", &newValue)
					if err != nil {
						fmt.Printf("%v not convert to int", metricLine.Value)
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				case "bool":
					var newValue bool
					if value == "0" {
						newValue = false
					} else {
						newValue = true
					}
					metricValue := reflect.ValueOf(&newValue).Elem()
					sv.Field(i).Set(metricValue)
				}
			}
		}
	}
}
