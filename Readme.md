### 仓库说明

该仓库代码用来解析k8s_metric，从[kube_state_metric](https://github.com/kubernetes/kube-state-metrics/tree/master/docs)
的URL地址获取指标原数据（该数据示例文件：[metrics.txt](example/metrics.txt)），按照[kube_state_metric](https://github.com/kubernetes/kube-state-metrics/tree/master/docs) 的分组情况
对指标进行分组聚合成为json格式数据；

比如：
[Pod Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/pod-metrics.md) 信息聚合后：
```json
[
    {
        "Namespace": "kodo-staging",
        "PodName": "kodo-stat-logexporter-r24tj",
        "PodInfo": {
            "HostIp": "10.200.20.xx",
            "PodIp": "10.202.107.xx",
            "Node": "cs34",
            "CreatedByKind": "DaemonSet",
            "CreatedByName": "kodo-stat-logexporter",
            "PriorityClass": ""
        },
        "PodCreated": 1647755616,
        "PodStartTime": 1647755616,
        "PodCompletionTime": 0,
        "PodRestartPolicy": "",
        "PodOwner": {
            "OwnerKind": "DaemonSet",
            "OwnerName": "kodo-stat-logexporter",
            "OwnerIsController": true
        },
        "PodLabels": [
            "label_app_kubernetes_io_name=stat",
            "label_controller_revision_hash=586758dff4",
            "label_kodo_qiniu_com_profile=true",
            "label_pod_template_generation=1",
            "label_app_kubernetes_io_component=logexporter",
            "label_app_kubernetes_io_instance=kodo-stat"
        ],
        "PodStatusPhase": "Running",
        "PodStatusReady": "true",
        "PodStatusScheduled": "true",
        "PodContainers": [
            {
                "name": "logexporter",
                "info": {
                    "container_id": "docker://751f9a473865acb5e67481fc4f07b63424a28e805037cc9333961f41f783aa9a",
                    "image": "xx-xx-xx.xx.io/kodo/qboxlogexporter.v2:enterprise_ec16p4-5d4935b4-1647754697",
                    "image_id": "docker-pullable://xx-xx-xx.xx.io/kodo/qboxlogexporter.v2@sha256:6aff0b6bee75c8360c91c5d4ecc62cd2c053d82ad31198103ecd09b495f422e1"
                },
                "status_waiting": false,
                "status_waiting_reason": "",
                "status_running": true,
                "state_started": 0,
                "status_terminated": false,
                "status_terminated_reason": "",
                "status_last_terminated_reason": "",
                "status_ready": false,
                "status_restarts_total": 0,
                "resource_requests": null,
                "resource_limits": null
            }
        ],
        "PodInitContainers": [
            {
                "name": "waiting",
                "info": {
                    "container_id": "docker://f3fc65b5b372a76bf0efb954bbfcc179174af1e8c894dcdfadfdac0ab475ef79",
                    "image": "xx-xx-xx.xx.io/qa/mongo:3.6.13",
                    "image_id": "docker-pullable://xx-xx-xx.xx.io/qa/mongo@sha256:d6541adc0c65cd9adf8690830e8dc8b916f82d4067a8b5e32c4f1be143e462e9"
                },
                "status_waiting": false,
                "status_waiting_reason": "",
                "status_running": false,
                "status_terminated": true,
                "status_terminated_reason": "Completed",
                "status_last_terminated_reason": "",
                "status_ready": false,
                "status_restarts_total": 0,
                "resource_requests": null,
                "resource_limits": null
            },
            {
                "name": "init-config",
                "info": {
                    "container_id": "docker://44af45c59d2c747744b58ebefd7137718b9a89f6e2fbc5086dc89b20c16c713d",
                    "image": "xx-xx-xx.xx.io/kodo_pub/vans:v1.0.0",
                    "image_id": "docker-pullable://xx-xx-xx.xx.io/kodo_pub/vans@sha256:42d8a53d92fae7b6a69d09c7057b49b9d11479777742d551fb72a3169866ba28"
                },
                "status_waiting": false,
                "status_waiting_reason": "",
                "status_running": false,
                "status_terminated": true,
                "status_terminated_reason": "Completed",
                "status_last_terminated_reason": "",
                "status_ready": false,
                "status_restarts_total": 0,
                "resource_requests": null,
                "resource_limits": null
            }
        ]
    }
]
```

[Daemonset Metrics](https://github.com/kubernetes/kube-state-metrics/blob/master/docs/daemonset-metrics.md) 信息聚合后：
```json
{
	"Namespace": "kube-system",
	"DaemonSetName": "kube-proxy",
	"DaemonSetCreated": 1589489617,
	"DaemonSetStatusCurrentNumber": 35,
	"DaemonSetStatusDesiredNumber": 35,
	"DaemonSetStatusNumberAvailable": 35,
	"DaemonSetStatusNumberMissScheduled": 0,
	"DaemonSetStatusNumberReady": 35,
	"DaemonSetStatusNumberUnavailable": 0,
	"DaemonSetMetadataGeneration": 1,
	"DaemonSetLabels": [
		"label_k8s_app=kube-proxy"
	]
}

```