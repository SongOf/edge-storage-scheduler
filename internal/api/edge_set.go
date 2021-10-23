package api

import (
	"context"
	"edge-storage-scheduler/internal/contains"
	"edge-storage-scheduler/internal/globals"
	"edge-storage-scheduler/internal/timer"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	"runtime"
	"sync"
	"time"
)

type void struct{}

var member void

//func Set(set string) {
//	EdgeSetOnline[set] = member
//}
//func Delete(set string) {
//	delete(EdgeSetOnline, set)
//}
//func Size() int {
//	return len(EdgeSetOnline)
//}
//func Exists(set string) bool {
//	_, exists := EdgeSetOnline[set]
//	return exists
//}
type EdgeSet struct {
	SetDiskOnlineTotal  model.Vector
	SetDiskOfflineTotal model.Vector
	SetDiskFree         model.Vector
	SetDiskTotal        model.Vector

	SetOnline      map[string]void
	SetScore       map[string]float64
	EdgeNodesOfSet map[string]*EdgeNode
}

func NewEdgeSet() *EdgeSet {
	return &EdgeSet{
		SetOnline:      make(map[string]void),
		SetScore:       make(map[string]float64),
		EdgeNodesOfSet: make(map[string]*EdgeNode),
	}
}

//v := result.(model.Vector)
//
//for _,s := range v {
//	fmt.Println(s.Metric)
//	fmt.Println(s.Value)
//	fmt.Println(s.Timestamp)
//}

func (es *EdgeSet) Run() {
	//刷新数据
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		err := es.RefreshSetDiskOnlineTotal()
		if err != nil {
			klog.Error("RefreshSetDiskOnlineTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := es.RefreshSetDiskOfflineTotal()
		if err != nil {
			klog.Error("RefreshSetDiskOfflineTotal fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := es.RefreshSetDiskFree()
		if err != nil {
			klog.Error("RefreshSetDiskFree fail", err)
		}
	}()
	go func() {
		defer wg.Done()
		err := es.RefreshSetDiskTotal()
		if err != nil {
			klog.Error("RefreshSetDiskTotal fail", err)
		}
	}()
	wg.Wait()

	//解析在线纠删集
	for _, set := range es.SetDiskOnlineTotal {
		setNameBytes, err := json.Marshal(set.Metric)
		if err != nil {
			klog.Error("map to json fail for ", set.Metric.String())
		}
		setName := string(setNameBytes)
		_, exists := es.SetOnline[setName]
		if !exists {
			//新增纠删集
			//为每个纠删集添加一个定时任务
			es.SetOnline[setName] = member
			es.EdgeNodesOfSet[setName] = NewEdgeNode(string(set.Metric["job"]))
			go func() {
				edgeNodeTimer := timer.Timer{
					Function: es.EdgeNodesOfSet[setName].Run,
					Duration: 60 * 1 * time.Second,
					Times:    0,
					Shutdown: make(chan string),
				}
				edgeNodeTimer.Start()
			}()
		}
	}

	//纠删集打分
	EdgeSetScheduler(es)
	//同步到redis zset
	ctx, cacel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cacel()
	for key, value := range es.SetScore {
		//err := globals.RedisClient.GetClient().ZIncrBy(ctx, contains.EDGE_SET_REDIS_KEY, value, key)
		cmd := globals.RedisClient.GetClient().ZAdd(ctx, contains.EDGE_SET_REDIS_KEY, &redis.Z{
			Score:  value,
			Member: key,
		})
		if cmd != nil {
			klog.Info("redis zset operation ", cmd)
		}
	}
}

func (es *EdgeSet) RefreshSetDiskOnlineTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, "avg(minio_cluster_disk_online_total) by (job)", time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("SetDiskOnlineTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	es.SetDiskOnlineTotal = v
	return nil
}

func (es *EdgeSet) RefreshSetDiskOfflineTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, "avg(minio_cluster_disk_offline_total) by (job)", time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("SetDiskOfflineTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	es.SetDiskOfflineTotal = v
	return nil
}

func (es *EdgeSet) RefreshSetDiskFree() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, "min(minio_node_disk_free_bytes) by (job) / 1024 / 1024", time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("SetDiskFree to Vector error")
		return &runtime.TypeAssertionError{}
	}
	es.SetDiskFree = v
	return nil
}

func (es *EdgeSet) RefreshSetDiskTotal() error {
	v1api := v1.NewAPI(*globals.PrometheusClient)
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, "min(minio_node_disk_total_bytes) by (job) / 1024 / 1024", time.Now())
	if err != nil {
		klog.Error("Error querying Prometheus: %v\n", err)
		return err
	}
	if len(warnings) > 0 {
		klog.Info("Warnings: %v\n", warnings)
	}
	v, ok := result.(model.Vector)
	if !ok {
		klog.Error("SetDiskTotal to Vector error")
		return &runtime.TypeAssertionError{}
	}
	es.SetDiskTotal = v
	return nil
}
