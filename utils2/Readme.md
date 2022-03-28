### 指标开发流程
新增一个指标类型，需要实现metric_state.go文件中的Metric interface中的方法，
分别是（Group，Add，Sender）

* Add 
```go
// 将读取到的行信息添加到该指标采集器的待处理列表中
func (m *PodMetric) Add(line utils2.MetricLine) {
	m.lines = append(m.lines, line)
}
```

* Group
```go
// 处理所有行信息，编写逻辑实现分组聚合
func (m *PodMetric) Group() {
	if m.Pods == nil {
		m.Pods = make(map[string]Pod)
	}
	...
```

* Sender
```go
// 将处理好的数据以指定的方法发送出去
func (m *PodMetric) Sender(todoFunc func(jsonData string)) {
	index := 0
	var listData []Pod
	for _, v := range m.Pods {
		listData = append(listData, v)
	}
	...
}


```