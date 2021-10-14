#edge-storage-scheduler
###边缘存储节点调度

##边缘节点部署
```
交叉编译：
GOOS=windows GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=amd64 go build main.go
GOOS=linux GOARCH=arm go build main.go 用于树莓派
```

```
4节点纠删组：
nano edge-storage-start.sh

#!/bin/sh
export MINIO_ROOT_USER=admin
export MINIO_ROOT_PASSWORD=admin123
export MINIO_PROMETHEUS_AUTH_TYPE="public"
/root/edge_storage_rasp server http://192.168.1.103/export1 http://192.168.1.104/export1 http://192.168.1.101/export1 http://192.168.1.102/export1
```

```
[Unit]
Description=Minio service
Documentation=https://docs.minio.io/

[Service]
WorkingDirectory=/root/
ExecStart=/root/edge-storage-start.sh

Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```
##边缘节点启动
```
启动服务： systemctl start minio.service
查看日志： journalctl -u minio.service -f
查看服务状态： systemctl status minio.service
```