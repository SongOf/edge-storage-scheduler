##juicefs客户端部署
```
交叉编译：
GOOS=windows GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=amd64 go build main.go
GOOS=linux GOARCH=arm go build main.go 用于树莓派
```

```
format:基于某个边缘节点的对象存储建立文件系统，如下：

$ ./juicefs format \
	--storage minio \
	--bucket http://192.168.1.102:9000/edgenode01 \
	--access-key admin \
	--secret-key admin123 \
	redis://:songof123@127.0.0.1:6379/5 \
	edgenode01
```

```
注意：mac需要安装macFUSE 4.1.2 windows也需要安装相应的fuse
mount:挂载到指定目录上

$ sudo ./juicefs mount -d redis://:songof123@127.0.0.1:6379/5 ~/edgenode01

$ sudo ./juicefs umount ~/edgenode01
```

```
status:查看文件系统状态
./juicefs status redis://:songof123@127.0.0.1:6379/5
```

```
验证：
$ df -Th
```