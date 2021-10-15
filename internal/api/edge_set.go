package api

import (
	"context"
	"edge-storage-scheduler/internal/globals"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog/v2"
	"runtime"
	"time"
)

type EdgeSet struct {
	SetDiskOnlineTotal  model.Vector
	SetDiskOfflineTotal model.Vector
	SetDiskFree         model.Vector
	SetDiskTotal        model.Vector
}

func NewEdgeSet() *EdgeSet {
	return new(EdgeSet)
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
