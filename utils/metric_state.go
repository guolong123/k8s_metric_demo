package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

//
//var (
//	PodGroupFields       = []string{"namespace", "pod"}
//	ContainerGroupFields = []string{"namespace", "pod", "container"}
//)

type MetricInterface interface {
	Group(metric MetricLine)
}

type MetricLine struct {
	Type      string
	Attribute map[string]interface{}
	Value     interface{}
}

func GetGroupFields(metric *Metric, metricLine MetricLine) (ret string) {
	// 获取分组字段值
	m := *metric
	st := reflect.TypeOf(m).Elem()
	sv := reflect.ValueOf(m).Elem()
	lenGroupField := int(0)
	var retList []string
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tagField := field.Tag.Get("json")
		groupField := field.Tag.Get("group")
		if groupField != "<nil>" && groupField != "" {
			lenGroupField++
			groupValue, ok := metricLine.Attribute[tagField]
			if !ok {
				continue
			}
			groupValue2 := fmt.Sprint(groupValue)
			metricAttr := reflect.ValueOf(&groupValue2).Elem()
			sv.FieldByName(groupField).Set(metricAttr)
			groupFieldValue, ok := metricLine.Attribute[tagField].(string)
			if ok {
				retList = append(retList, groupFieldValue)
			}

		}
	}
	if lenGroupField == len(retList) {
		ret = strings.Join(retList, "_")
	}

	return
}

func ReadFromUrl(url string, todoFunc func(jsonData string)) error {
	// 从给定的kube_state_metric服务地址中按行取数据，舍弃注释行；
	// 将单行数据解析为MetricLine结构体，并使用转换后的结构体去和Metric类型的结构体中的json tag进行匹配，
	// 匹配一致，则将metricLine中的值赋给metric
	groupPodHandler := map[string]Metric{}
	//GroupContainerHandler := map[string]Metric{}
	//GroupInitContainerHandler := map[string]Metric{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if resp.StatusCode == 200 {
		fmt.Println("ok")
	}
	defer resp.Body.Close()
	//var listBody []string
	buf := bufio.NewReader(resp.Body)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		//fmt.Println(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return err
			}
		}
		// 将单行字符串转换为metricLine类型
		metricLine := readLine(line)

		if metricLine.Type != "" {
			var podMetric Metric = &PodMetric{}
			var kubePodInitContainer Metric = &KubePodInitContainer{}
			var kubePodContainer Metric = &KubePodContainer{}

			// 获取分组字段，作为pod信息的key（分组字段值为metric结构体中tag为group字段的字符串拼接）
			groupField := GetGroupFields(&podMetric, metricLine) // pod的分组字段为namespace + "_" + pod_name
			if groupField == "" {
				continue
			}
			// 从map中获取key为分组字段的对象，不存在则创建
			_, ok := groupPodHandler[groupField]
			if ok != true {
				groupPodHandler[groupField] = podMetric
			}
			// 将单行数据metricLine按照规则放到指定的分组map中
			groupPodHandler[groupField].Group(&kubePodContainer, &kubePodInitContainer, metricLine)

		}
	}
	for _, v := range groupPodHandler {
		v.UpdateContainer("KubePodContainers")

		v.UpdateContainer("KubePodInitContainers")

		//var mapContainer = map[string]Metric{}
		//var mapInitContainer = map[string]Metric{}
		var listData = []Metric{v}
		jsonData, err := json.Marshal(listData)
		if err == nil {
			todoFunc(string(jsonData))
			//listBody = append(listBody, string(jsonData))
		}
	}

	return nil
}

func Distribution(metricLine MetricLine, metric Metric) bool {
	// 判断给定的metric中是否包含metricLine.Type类型
	st := reflect.TypeOf(metric).Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tag := GLOBALNAME + field.Tag.Get("json")
		if tag == metricLine.Type {
			return true
		}
	}
	return false
}

func readLine(line string) (metric MetricLine) {

	// 忽略#开头的注释信息
	if strings.HasPrefix(line, "#") {
		return metric
	}

	// 使用正则匹配解析，将单行字符串解析为MetricLine结构体对象
	reg1 := regexp.MustCompile(`^(\w+)\{(\w+=".+")+\}\s(.*)$`)
	if reg1 == nil {
		fmt.Println("MustCompile err")
		return
	}
	result := reg1.FindAllStringSubmatch(line, -1)
	metric.Value = result[0][len(result[0])-1]
	metric.Type = result[0][1]
	labels := strings.Split(result[0][2], ",")
	var labelMap = make(map[string]interface{})
	for i := 0; i < len(labels); i++ {
		label := strings.Split(labels[i], "=")
		labelMap[label[0]] = strings.ReplaceAll(label[1], "\"", "")
	}
	metric.Attribute = labelMap
	return
}
