package main

import (
	"fmt"
	"io/ioutil"
	"k8s_metric/utils2"
	_ "k8s_metric/utils2/kube_configmap"
	_ "k8s_metric/utils2/kube_daemonset"
	_ "k8s_metric/utils2/kube_ingress"
	_ "k8s_metric/utils2/kube_node"
	_ "k8s_metric/utils2/kube_pod"
	"net/http"
	"strings"
	"time"
)

func main() {
	startTime := time.Now().Unix()
	err := utils2.ReadFromUrl("http://10.202.185.50:8080/metrics", uploadToPandora)
	if err != nil {
		return
	}
	endTime := time.Now().Unix()
	dustion := endTime - startTime
	fmt.Printf("总共耗时：%d", dustion)
}

func uploadToPandora(jsonData string) {
	defer func() {
		if err := recover(); err != nil {
			// 打印异常，关闭资源，退出此函数
			fmt.Println(err)
			uploadToPandora(jsonData)
		}
	}()
	token := "eyJhbGciOiJIUzUxMiIsInppcCI6IkRFRiJ9.eJwVy0sOgyAQANC7zBoaPjPUsvIqjEBCY9GqNE2Mdy_dv3dC-q7gtSOHD2UUCXgeBTw4ooBkkiTFWqKdouQhosRkORmKOdsMAirn_8b7QH1rAXvjvtdQ47KF8V1qabdpeXX6KdvRwgw-h3lP1w_XdyNd.13uwdKHaXXduQi2FhuHnb-FZxtOBKsqCG_bPbDDA4wzFHPxaYsFfoTI7Q5ZdP2RnKhdPAgtydFokG6DL4vxasA"
	client := &http.Client{}

	req, err := http.NewRequest("POST", "http://pandora-web-svc.pandora-jks-guolong.qa.qiniu.io/api/v1/data?repo=k8s_metrics2&sourcetype=json", strings.NewReader(jsonData))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
