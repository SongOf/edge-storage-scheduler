package api

import (
	"context"
	"edge-storage-scheduler/internal/contains"
	"edge-storage-scheduler/internal/globals"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	"runtime"
	"sync"
	"time"
)

var TIMEOUT = 10 * time.Second

type EdgeNode struct {
	EdgeSetName string
	JobName     string
	//cpu idle状态的百分比
	NodeCpuIdle model.Vector
	//cpu system状态的百分比
	NodeCpuSystem model.Vector
	//cpu user状态的百分比
	NodeCpuUser model.Vector
	//cpu 1m负载
	NodeLoad1 model.Vector
	//cpu 5m负载
	NodeLoad5 model.Vector
	//cpu 15m负载
	NodeLoad15 model.Vector
	//可用内存 1m内均值 单位MB
	NodeMemoryAvailable model.Vector
	//总内存 1m内均值 单位MB
	NodeMemoryTotal model.Vector
	//Buffer缓存 1m内均值 单位MB
	NodeMemoryBuffer model.Vector
	//Cached缓存 1m内均值 单位MB
	NodeMemoryCached model.Vector
	//Swap总容量 1m内均值 单位MB
	NodeMemorySwapTotal model.Vector
	//Swap可用容量 1m内均值 单位MB
	NodeMemorySwapFree model.Vector
	//device="/dev/root" fstype="ext4" mountpoint="/" 的文件系统总容量 1m内均值 单位MB
	NodeFileSystemSize model.Vector
	//device="/dev/root" fstype="ext4" mountpoint="/" 的文件系统可用容量 1m内均值 单位MB
	NodeFileSystemAvail model.Vector
	//磁盘读取速率 单位KB/s
	NodeDiskIORead model.Vector
	//磁盘写入速率 单位KB/s
	NodeDiskIOWrite model.Vector
	//device="eth0" 以太网接口的上传速率 1m内的速率 单位KB/s
	NodeNetworkTransmitTotal model.Vector
	//device="eth0" 以太网接口的下载速率 1m内的速率 单位KB/s
	NodeNetworkReceiveTotal model.Vector

	NodeOnline map[string]void
	NodeScore  map[string]float64
}

func NewEdgeNode(edgeSetName string) *EdgeNode {
	return &EdgeNode{
		EdgeSetName: edgeSetName,
		JobName:     edgeSetName + contains.EDGE_NODE_JOB_NAME_SUB_TAIL,
		NodeOnline:  make(map[string]void),
		NodeScore:   make(map[string]float64),
	}
}

func (en *EdgeNode) Run() {
	//刷新数据
	wg := sync.WaitGroup{}
	wg.Add(18)
	go func() {
		defer wg.Done()
		err := en.RefreshNodeCpuIdleInfo()
		if err != nil {
			klog.Error("RefreshNodeCpuIdleInfo fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeCpuSystemInfo()
		if err != nil {
			klog.Error("RefreshNodeCpuSystemInfo fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeCpuUserInfo()
		if err != nil {
			klog.Error("RefreshNodeCpuUserInfo fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeLoad1Info()
		if err != nil {
			klog.Error("RefreshNodeLoad1Info fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeLoad5Info()
		if err != nil {
			klog.Error("RefreshNodeLoad5Info fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeLoad15Info()
		if err != nil {
			klog.Error("RefreshNodeLoad15Info fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemoryAvailable()
		if err != nil {
			klog.Error("RefreshNodeMemoryAvailable fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemoryTotal()
		if err != nil {
			klog.Error("RefreshNodeMemoryTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeFileSystemAvail()
		if err != nil {
			klog.Error("RefreshNodeFileSystemAvail fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeFileSystemSize()
		if err != nil {
			klog.Error("RefreshNodeFileSystemSize fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeNetworkTransmitTotal()
		if err != nil {
			klog.Error("RefreshNodeNetworkTransmitTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeNetworkReceiveTotal()
		if err != nil {
			klog.Error("RefreshNodeNetworkReceiveTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemoryBuffer()
		if err != nil {
			klog.Error("RefreshNodeMemoryBuffer fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemoryCached()
		if err != nil {
			klog.Error("RefreshNodeMemoryCached fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemorySwapFree()
		if err != nil {
			klog.Error("RefreshNodeMemorySwapFree fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeMemorySwapTotal()
		if err != nil {
			klog.Error("RefreshNodeMemorySwapTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeDiskIORead()
		if err != nil {
			klog.Error("RefreshNodeDiskIORead fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := en.RefreshNodeDiskIOWrite()
		if err != nil {
			klog.Error("RefreshNodeDiskIOWrite fail", err)
		}
	}()
	wg.Wait()

	//解析在线节点
	for _, node := range en.NodeCpuIdle {
		nodeNameBytes, err := json.Marshal(node.Metric)
		if err != nil {
			klog.Error("map to json fail for ", node.Metric.String())
		}
		nodeName := string(nodeNameBytes)
		_, exists := en.NodeOnline[nodeName]
		if !exists {
			//新增节点
			en.NodeOnline[nodeName] = member
		}
	}

	//边缘节点打分
	EdgeNodeScheduler(en)
	//同步到redis zset
	ctx, cacel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cacel()
	for key, value := range en.NodeScore {
		//globals.RedisClient.GetClient().ZIncrBy(ctx, en.EdgeSetName, value, key)
		cmd := globals.RedisClient.GetClient().ZAdd(ctx, en.EdgeSetName, &redis.Z{
			Score:  value,
			Member: key,
		})
		if cmd != nil {
			klog.Info("redis zset operation ", cmd)
		}
	}
}

func (en *EdgeNode) RefreshNodeCpuIdleInfo() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("(sum(increase(node_cpu_seconds_total{mode='idle',job='%s'}[1m]))by(instance)) / (sum(increase(node_cpu_seconds_total{job='%s'}[1m]))by(instance)) * 100", en.JobName, en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeCpuIdle to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeCpuIdle = v
	return nil
}

func (en *EdgeNode) RefreshNodeCpuSystemInfo() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("(sum(increase(node_cpu_seconds_total{mode='system',job='%s'}[1m]))by(instance)) / (sum(increase(node_cpu_seconds_total{job='%s'}[1m]))by(instance))  *100", en.JobName, en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeCpuSystem to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeCpuSystem = v
	return nil
}

func (en *EdgeNode) RefreshNodeCpuUserInfo() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("(sum(increase(node_cpu_seconds_total{mode='user',job='%s'}[1m]))by(instance)) / (sum(increase(node_cpu_seconds_total{job='%s'}[1m]))by(instance))  *100", en.JobName, en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeCpuUser to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeCpuUser = v
	return nil
}

func (en *EdgeNode) RefreshNodeLoad1Info() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("node_load1{job='%s'}", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeLoad1 to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeLoad1 = v
	return nil
}

func (en *EdgeNode) RefreshNodeLoad5Info() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("node_load5{job='%s'}", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeLoad5 to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeLoad5 = v
	return nil
}

func (en *EdgeNode) RefreshNodeLoad15Info() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("node_load15{job='%s'}", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeLoad15 to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeLoad15 = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemoryAvailable() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_MemAvailable_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemoryAvailable to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemoryAvailable = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemoryTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_MemTotal_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemoryTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemoryTotal = v
	return nil
}

func (en *EdgeNode) RefreshNodeFileSystemSize() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_filesystem_size_bytes{mountpoint=\"/\", job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeFileSystemAvail to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeFileSystemSize = v
	return nil
}

func (en *EdgeNode) RefreshNodeFileSystemAvail() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_filesystem_avail_bytes{mountpoint=\"/\", job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeFileSystemSize to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeFileSystemAvail = v
	return nil
}

func (en *EdgeNode) RefreshNodeNetworkTransmitTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg(irate(node_network_transmit_bytes_total{device=\"eth0\", job='%s'}[1m]) / 1024) by (instance)", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeNetworkTransmitTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeNetworkTransmitTotal = v
	return nil
}

func (en *EdgeNode) RefreshNodeNetworkReceiveTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg(irate(node_network_receive_bytes_total{device=\"eth0\", job='%s'}[1m]) / 1024) by (instance)", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeNetworkReceiveTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeNetworkReceiveTotal = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemoryBuffer() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_Buffers_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemoryBuffer to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemoryBuffer = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemoryCached() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_Cached_bytes{job='%s'}[1m]) / 1024 / 1024 + avg_over_time(node_memory_Slab_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName, en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemoryCached to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemoryCached = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemorySwapFree() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_SwapFree_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemorySwapFree to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemorySwapFree = v
	return nil
}

func (en *EdgeNode) RefreshNodeMemorySwapTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("avg_over_time(node_memory_SwapTotal_bytes{job='%s'}[1m]) / 1024 / 1024", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeMemorySwapTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeMemorySwapTotal = v
	return nil
}

func (en *EdgeNode) RefreshNodeDiskIORead() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("sum(irate(node_disk_reads_completed_total{job='%s'}[5m])) by (instance)", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeDiskIORead to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeDiskIORead = v
	return nil
}

func (en *EdgeNode) RefreshNodeDiskIOWrite() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	queryStr := fmt.Sprintf("sum(irate(node_disk_writes_completed_total{job='%s'}[5m])) by (instance)", en.JobName)
	result, warnings, err := v1api.Query(ctx, queryStr, time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("NodeDiskIOWrite to Vector error")
		return &runtime.TypeAssertionError{}
	}
	en.NodeDiskIOWrite = v
	return nil
}

//v := result.(model.Vector)
//
//for _,s := range v {
//	fmt.Println(s.Metric)
//	fmt.Println(s.Value)
//	fmt.Println(s.Timestamp)
//}
