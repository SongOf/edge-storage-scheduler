package api

import (
	"edge-storage-scheduler/internal/globals"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"sync"
	"testing"
)

func TestEdgeSet_RefreshSet(t *testing.T) {
	client, err := api.NewClient(api.Config{
		Address: "http://127.0.0.1:9090",
	})
	if err != nil {
		fmt.Println(err)
	}
	globals.PrometheusClient = &client

	edgeSet := NewEdgeSet()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		err := edgeSet.RefreshSetDiskOnlineTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeSet.RefreshSetDiskOfflineTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeSet.RefreshSetDiskFree()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeSet.RefreshSetDiskTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
	fmt.Println(edgeSet.SetDiskOnlineTotal.String())
	fmt.Println(edgeSet.SetDiskOfflineTotal.String())
	fmt.Println(edgeSet.SetDiskFree.String())
	fmt.Println(edgeSet.SetDiskTotal.String())

	for _, set := range edgeSet.SetDiskOnlineTotal {
		fmt.Println(set.Metric)
	}
}
