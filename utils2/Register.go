package utils2

var Constructor = map[string]Metric{}

func Register(metricType string, metric Metric) {
	Constructor[metricType] = metric
}
