package main

import "k8s_metric/utils"

func main() {
	_, err := utils.ReadFromUrl("http://10.202.185.50:8080/metrics")
	if err != nil {
		return
	}
}
