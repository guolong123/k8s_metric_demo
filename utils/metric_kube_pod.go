package utils

import (
	"fmt"
	"reflect"
	"strings"
)

var GLOBALNAME = "kube_pod_"

type Metric interface {
	Group(kubePodContainer *Metric, kubePodInitContainer *Metric, metricLine MetricLine)
	UpdateContainer(fieldName string)
}

type PodMetric struct {
	Namespace                string                 `json:"namespace" group:"Namespace"`
	Pod                      string                 `json:"pod" group:"Pod"`
	KubePodInfo              map[string]interface{} `json:"info" get_attr:"true" get_value:"false"`
	KubePodCreated           float64                `json:"created" get_attr:"false" get_value:"true"`
	KubePodStartTime         float64                `json:"start_time" get_attr:"false" get_value:"true"`
	KubePodCompletionTime    float64                `json:"completion_time" get_attr:"false" get_value:"true"`
	KubePodRestartPolicy     map[string]interface{} `json:"restart_policy" get_attr:"true" get_value:"false"`
	KubePodOwner             map[string]interface{} `json:"owner" get_attr:"true" get_value:"false"`
	KubePodLabel             map[string]string      `json:"label" get_attr:"true" get_value:"false"`
	KubePodStatusPhase       map[string]interface{} `json:"status_phase" get_attr:"true" get_value:"false"`
	KubePodStatusReady       map[string]interface{} `json:"status_ready" get_attr:"true" get_value:"false"`
	KubePodStatusScheduled   map[string]interface{} `json:"status_scheduled" get_attr:"true" get_value:"false"`
	KubePodContainers        map[string]Metric      `json:"containers" get_attr:"true" get_value:"false"`
	KubePodInitContainers    map[string]Metric      `json:"init_containers" get_attr:"true" get_value:"false"`
	KubePodContainerInfo     []Metric               `json:"container"`
	KubePodInitContainerInfo []Metric               `json:"init_container"`
}

type KubePodContainer struct {
	Namespace                                  string                 `json:"namespace" group:"Namespace"`
	Pod                                        string                 `json:"pod" group:"Pod"`
	Container                                  string                 `json:"container" group:"Container"`
	KubePodContainerInfo                       map[string]interface{} `json:"container_info" get_attr:"true" get_value:"false"`
	KubePodContainerStatusWaiting              bool                   `json:"container_status_waiting" get_attr:"false" get_value:"true"`
	KubePodContainerStatusWaitingReason        map[string]interface{} `json:"container_status_waiting_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusRunning              bool                   `json:"container_status_running" get_attr:"false" get_value:"true"`
	KubePodContainerStateStarted               float64                `json:"container_state_started" get_attr:"false" get_value:"true"`
	KubePodContainerStatusTerminated           bool                   `json:"container_status_terminated" get_attr:"false" get_value:"true"`
	KubePodContainerStatusTerminatedReason     map[string]interface{} `json:"container_status_terminated_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusLastTerminatedReason map[string]interface{} `json:"container_status_last_terminated_reason" get_attr:"true" get_value:"false"`
	KubePodContainerStatusReady                bool                   `json:"container_status_ready" get_attr:"false" get_value:"true"`
	KubePodContainerStatusRestartTotal         int                    `json:"container_status_restarts_total" get_attr:"false" get_value:"true"`
	//KubePodContainerResourceRequests           []map[string]interface{} `json:"container_resource_requests" get_attr:"true" get_value:"true"`
	//KubePodContainerResourceLimits             []map[string]interface{} `json:"container_resource_limits" get_attr:"true" get_value:"true"`
}

type KubePodInitContainer struct {
	Pod                                            string                 `json:"pod" group:"Pod"`
	Namespace                                      string                 `json:"namespace" group:"Namespace"`
	Container                                      string                 `json:"container" group:"Container"`
	KubePodInitContainerInfo                       map[string]interface{} `json:"init_container_info"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusWaiting              bool                   `json:"init_container_status_waiting"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusWaitingReason        map[string]interface{} `json:"init_container_status_waiting_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusRunning              bool                   `json:"init_container_status_running" get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusTerminated           bool                   `json:"init_container_status_terminated"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusTerminatedReason     map[string]interface{} `json:"init_container_status_terminated_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusLastTerminatedReason map[string]interface{} `json:"init_container_status_last_terminated_reason"  get_attr:"true" get_value:"false"`
	KubePodInitContainerStatusReady                bool                   `json:"init_container_status_ready"  get_attr:"false" get_value:"true"`
	KubePodInitContainerStatusRestartsTotal        int                    `json:"init_container_status_restarts_total" get_attr:"false" get_value:"true"`
	//KubePodInitContainerResourceRequests           []map[string]interface{} `json:"init_container_resource_requests" get_attr:"true" get_value:"false"`
	//KubePodInitContainerResourceLimits             []map[string]interface{} `json:"init_container_resource_limits" get_attr:"true" get_value:"false"`
}

func (m *PodMetric) UpdateContainer(fieldName string) {
	st := reflect.TypeOf(*m)
	sv := reflect.ValueOf(m).Elem()
	var containerList []Metric
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		if field.Name == fieldName {
			t, ok := sv.FieldByName(fieldName).Interface().(map[string]Metric)
			if ok {
				for _, value := range t {
					containerList = append(containerList, value)
				}
			}
		}
	}
	if len(containerList) > 0 {
		if fieldName == "KubePodContainers" {
			fmt.Printf("%T,%v", containerList, containerList)
			sv.FieldByName("KubePodContainerInfo").Set(reflect.ValueOf(containerList))
			zeroMapMetric := map[string]Metric{}
			sv.FieldByName("KubePodContainers").Set(reflect.ValueOf(&zeroMapMetric).Elem())

		} else if fieldName == "KubePodInitContainers" {
			sv.FieldByName("KubePodInitContainerInfo").Set(reflect.ValueOf(containerList))
			zeroMapMetric := map[string]Metric{}
			sv.FieldByName("KubePodInitContainers").Set(reflect.ValueOf(&zeroMapMetric).Elem())

		}
	}

}

func (m KubePodInitContainer) UpdateContainer(fieldName string) {}
func (m KubePodContainer) UpdateContainer(fieldName string)     {}

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
	var groupFields []string
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		jsonTag := GLOBALNAME + field.Tag.Get("json")
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")
		groupTag := field.Tag.Get("group")
		if groupTag != "" {
			groupFields = append(groupFields, groupTag)
		}
		if jsonTag == "kube_pod_containers" && Distribution(metricLine, *kubePodContainer) {
			// 当jsonTag为"kube_pod_containers"并且kubePodContainer对象属性中包含metricLine.type时
			// 获取并更新KubePodContainers[groupField]的内容
			groupField := GetGroupFields(kubePodContainer, metricLine)
			kubePodContainerValue, ok := m.KubePodContainers[groupField]
			if !ok {
				kubePodContainerValue = *kubePodContainer
			}
			kubePodContainerValue.Group(kubePodContainer, kubePodInitContainer, metricLine)
			m.KubePodContainers[groupField] = kubePodContainerValue
		} else if jsonTag == "kube_pod_init_containers" && Distribution(metricLine, *kubePodInitContainer) {
			// 同上
			groupField := GetGroupFields(kubePodInitContainer, metricLine)
			kubePodInitContainerValue, ok := m.KubePodInitContainers[groupField]
			if !ok {
				kubePodInitContainerValue = *kubePodInitContainer
			}
			kubePodInitContainerValue.Group(kubePodContainer, kubePodInitContainer, metricLine)
			m.KubePodInitContainers[groupField] = kubePodInitContainerValue

		} else if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			// 处理其它类型的对象；
			//metricLine.Value != valueFalse，这里有个特殊处理，当metricLine.value等于0时将舍弃此行；算是一个坑

			if getAttr == "true" {
				// 当getAttr为true时，将metricLine.Attribute直接赋给m对象
				// 另外一个坑，这里getAttr和getValue不能同时为true；因为字段名只有一个；
				attrs := map[string]interface{}{}
				for k, v := range metricLine.Attribute {
					t := 0
					for _, g := range groupFields {
						if strings.ToUpper(g) == strings.ToUpper(k) {
							t++
						}
					}
					if t == 0 {
						attrs[k] = v
					}
				}
				metricAttr := reflect.ValueOf(&attrs).Elem()
				sv.Field(i).Set(metricAttr)
			} else if getValue == "true" {
				// 当getValue为true时，将针对原定义的字段类型来转换metricLine.value；
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
	var groupFields []string
	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		jsonTag := GLOBALNAME + field.Tag.Get("json")
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")
		groupTag := field.Tag.Get("group")
		if groupTag != "" {
			groupFields = append(groupFields, groupTag)
		}
		if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			if getAttr == "true" {

				//if jsonTag == "kube_pod_init_container_resource_requests" && Distribution(metricLine, *kubePodInitContainer) {
				//	typeName := field.Type.Name()
				//	fmt.Println(typeName)
				//	//m.KubePodInitContainerResourceLimits = append(m.KubePodInitContainerResourceLimits, )
				//}
				attrs := map[string]interface{}{}
				for k, v := range metricLine.Attribute {
					t := 0
					for _, g := range groupFields {
						if strings.ToUpper(g) == strings.ToUpper(k) {
							t++
						}
					}
					if t == 0 {
						attrs[k] = v
					}
				}
				metricAttr := reflect.ValueOf(&attrs).Elem()
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

	var groupFields []string

	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		jsonTag := GLOBALNAME + field.Tag.Get("json")
		groupTag := field.Tag.Get("group")
		if groupTag != "" {
			groupFields = append(groupFields, groupTag)
		}
		getAttr := field.Tag.Get("get_attr")
		getValue := field.Tag.Get("get_value")
		if jsonTag == metricLine.Type && metricLine.Value != valueFalse {
			if getAttr == "true" {
				//if jsonTag == "kube_pod_container_resource_requests" && Distribution(metricLine, *kubePodInitContainer) {
				//	typeName := field.Type.Name()
				//	fmt.Println(typeName)
				//	//m.KubePodInitContainerResourceLimits = append(m.KubePodInitContainerResourceLimits, )
				//}
				attrs := map[string]interface{}{}
				for k, v := range metricLine.Attribute {
					t := 0
					for _, g := range groupFields {
						if strings.ToUpper(g) == strings.ToUpper(k) {
							t++
						}
					}
					if t == 0 {
						attrs[k] = v
					}
				}
				metricAttr := reflect.ValueOf(&attrs).Elem()
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
