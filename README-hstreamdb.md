##HStreamDB服务端部署
###docker版单机部署
```
db数据存储目录，由juicefs映射到的本地文件系统（后端是边缘对象存储节点）
如：/Users/lisongsong/edgenode01
```

```
Start HStream Storage
docker run -td --rm --name some-hstream-store -v /Users/lisongsong/edgenode01:/data/store --network host hstreamdb/hstream ld-dev-cluster --root /data/store --use-tcp

docker run -td --rm --name some-hstream-store -v /Users/lisongsong/dbdata:/data/store --network host hstreamdb/hstream ld-dev-cluster --root /data/store --use-tcp
```

```
Start HStreamDB Server
docker run -it --rm --name some-hstream-server -v /Users/lisongsong/edgenode01:/data/store --network host hstreamdb/hstream hstream-server --port 6570 --store-config /Users/lisongsong/logdevice.conf

docker run -it --rm --name some-hstream-server -v /Users/lisongsong/dbdata:/data/store --network host hstreamdb/hstream hstream-server --port 6570 --store-config /data/store/logdevice.conf
```

```
Start HStreamDB's interactive SQL CLI
docker run -it --rm --name some-hstream-cli -v /Users/lisongsong/edgenode01:/data/store --network host hstreamdb/hstream hstream-client --port 6570
```

