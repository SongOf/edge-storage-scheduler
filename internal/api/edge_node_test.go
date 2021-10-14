package api

import (
	"edge-storage-scheduler/internal/globals"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"sync"
	"testing"
)

func TestEdgeNode_LoadEdgeNodeInfo(t *testing.T) {
	client, err := api.NewClient(api.Config{
		Address: "http://127.0.0.1:9090",
	})
	if err != nil {
		fmt.Println(err)
	}
	globals.PrometheusClient = &client
	edgeNode := NewEdgeNode()
	wg := sync.WaitGroup{}
	wg.Add(18)
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeCpuIdleInfo()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeCpuSystemInfo()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeCpuUserInfo()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeLoad1Info()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeLoad5Info()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeLoad15Info()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemoryAvailable()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemoryTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeFileSystemAvail()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeFileSystemSize()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeNetworkTransmitTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeNetworkReceiveTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemoryBuffer()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemoryCached()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemorySwapFree()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeMemorySwapTotal()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeDiskIORead()
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		defer wg.Done()
		err := edgeNode.RefreshNodeDiskIOWrite()
		if err != nil {
			fmt.Println(err)
		}
	}()
	wg.Wait()
	fmt.Println(edgeNode.NodeCpuUser.String())
	fmt.Println(edgeNode.NodeCpuSystem.String())
	fmt.Println(edgeNode.NodeCpuIdle.String())
	fmt.Println(edgeNode.NodeLoad1.String())
	fmt.Println(edgeNode.NodeLoad5.String())
	fmt.Println(edgeNode.NodeLoad15.String())
	fmt.Println(edgeNode.NodeMemoryAvailable.String())
	fmt.Println(edgeNode.NodeMemoryTotal.String())
	fmt.Println(edgeNode.NodeFileSystemAvail.String())
	fmt.Println(edgeNode.NodeFileSystemSize.String())
	fmt.Println(edgeNode.NodeNetworkTransmitTotal.String())
	fmt.Println(edgeNode.NodeNetworkReceiveTotal.String())
	fmt.Println(edgeNode.NodeMemoryBuffer.String())
	fmt.Println(edgeNode.NodeMemoryCached.String())
	fmt.Println(edgeNode.NodeMemorySwapFree.String())
	fmt.Println(edgeNode.NodeMemorySwapTotal.String())
	fmt.Println(edgeNode.NodeDiskIORead.String())
	fmt.Println(edgeNode.NodeDiskIOWrite.String())
}
