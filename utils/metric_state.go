package utils

import (
	"bufio"
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
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tagField := field.Tag.Get("json")
		groupField := field.Tag.Get("group")
		if groupField != "<nil>" && groupField != "" {
			groupValue, ok := metricLine.Attribute[tagField]
			if !ok {
				continue
			}
			groupValue2 := fmt.Sprint(groupValue)
			metricAttr := reflect.ValueOf(&groupValue2).Elem()
			sv.FieldByName(groupField).Set(metricAttr)
			ret = fmt.Sprintf("%s_%s", ret, metricLine.Attribute[tagField])
		}
	}
	return
}

func ReadFromUrl(url string) ([]string, error) {
	groupPodHandler := map[string]Metric{}
	GroupContainerHandler := map[string]Metric{}
	GroupInitContainerHandler := map[string]Metric{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if resp.StatusCode == 200 {
		fmt.Println("ok")
	}
	defer resp.Body.Close()
	var listBody []string
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
				return nil, err
			}
		}
		metricLine := readLine(line)

		if metricLine.Type != "" {
			var metric Metric = &PodMetric{}
			var kubePodInitContainer Metric = &KubePodInitContainer{}
			var kubePodContainer Metric = &KubePodContainer{}
			groupField := GetGroupFields(&metric, metricLine)
			containerGroupField := GetGroupFields(&kubePodContainer, metricLine)
			initContainerGroupField := GetGroupFields(&kubePodContainer, metricLine)

			_, ok := GroupContainerHandler[containerGroupField]
			if !ok {
				GroupContainerHandler[containerGroupField] = kubePodContainer
			}
			_, ok = GroupInitContainerHandler[initContainerGroupField]
			if !ok {
				GroupInitContainerHandler[initContainerGroupField] = kubePodInitContainer
			}

			_, ok = groupPodHandler[groupField]
			if ok != true {
				groupPodHandler[groupField] = metric
			}
			a := GroupContainerHandler[containerGroupField]
			b := GroupContainerHandler[containerGroupField]
			groupPodHandler[groupField].Group(&a, &b, metricLine)

		}
	}

	return listBody, nil
}

func Distribution(metricLine MetricLine, metric Metric) bool {
	// 判断给定的metric中是否包含metricLine中的类型
	st := reflect.TypeOf(metric).Elem()
	for i := 0; i < st.NumField(); i++ {
		//fmt.Println(st.Field(i).Tag) //将tag输出出来
		field := st.Field(i)
		tag := field.Tag.Get("json")

		if tag == metricLine.Type {
			return true
		}
	}
	return false
}

func readLine(line string) (metric MetricLine) {
	if strings.HasPrefix(line, "#") {
		return metric
	}
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
