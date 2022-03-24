package utils2

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type Metric interface {
	Group()
	Add(line MetricLine)
	Sender(func(jsonData string))
}

type MetricLine struct {
	Type      string
	Attribute map[string]string
	Value     interface{}
}

func ReadFromUrl(url string, todoFunc func(jsonData string)) (err error) {
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
		if metricLine.Type == "" {
			continue
		}
		for k, metric := range Constructor {
			if strings.HasPrefix(metricLine.Type, k) {
				metric.Add(metricLine)
			}
		}
	}
	for _, metric := range Constructor {
		metric.Group()
		metric.Sender(todoFunc)
	}
	return nil
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
	var labelMap = make(map[string]string)
	for i := 0; i < len(labels); i++ {
		label := strings.Split(labels[i], "=")
		labelMap[label[0]] = strings.ReplaceAll(label[1], "\"", "")
	}
	metric.Attribute = labelMap
	return
}
