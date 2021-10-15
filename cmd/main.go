package main

import (
	"edge-storage-scheduler/conf"
	api2 "edge-storage-scheduler/internal/api"
	"edge-storage-scheduler/internal/globals"
	"edge-storage-scheduler/internal/model"
	"edge-storage-scheduler/internal/timer"
	"edge-storage-scheduler/internal/utils/redis"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
	"strings"
	"sync"
	"time"
)

func initViper(confFilePath, configFileFullName string) *viper.Viper {
	oviper := viper.New()
	r := strings.Split(configFileFullName, ".")
	if len(r) != 2 {
		klog.Fatal("invalid config file to init viper")
	}
	oviper.SetConfigName(r[0])
	oviper.SetConfigType(r[1])
	oviper.AddConfigPath(confFilePath)
	err := oviper.ReadInConfig()
	if err != nil {
		klog.Fatal("parse config file error", err)
	}
	return oviper
}

func loadPrometheusConfig(v *viper.Viper) *conf.PrometheusConf {
	prometheusConf := conf.PrometheusConf{}
	if err := v.UnmarshalKey("prometheus", &prometheusConf); err != nil {
		klog.Fatal("load prometheus config error", err)
	}
	return &prometheusConf
}

func loadRedisConfig(v *viper.Viper) *conf.RedisConf {
	redisConf := conf.RedisConf{}
	if err := v.UnmarshalKey("redis", &redisConf); err != nil {
		klog.Fatal("load redis config error", err)
	}
	return &redisConf
}

func init() {
	commonViper := initViper("./conf", "common_config.yaml")
	//初始化prometheus客户端
	prometheusConf := loadPrometheusConfig(commonViper)
	prometheusAddr := fmt.Sprintf("http://%s:%d", prometheusConf.Host, prometheusConf.Port)
	client, err := api.NewClient(api.Config{
		Address: prometheusAddr,
	})
	if err != nil {
		fmt.Printf("Error creating prometheus client: %v\n", err)
		os.Exit(1)
	}
	globals.PrometheusClient = &client

	//初始化redis客户端
	redisConf := loadRedisConfig(commonViper)
	redisClient := redis.NewRedisCache(redis.RedisOption{
		Address:  redisConf.Address,
		Password: redisConf.Password,
	})
	globals.RedisClient = redisClient
}

func main() {
	//初始化在线就删集
	model.NewEdgeSetOnline()

	scanCycle := 60 * 1 * time.Second
	edgeSet := api2.NewEdgeSet()
	edgeSetTimer := timer.Timer{
		Function: edgeSet.Run,
		Duration: scanCycle,
		Times:    0,
		Shutdown: make(chan string),
	}
	defer edgeSetTimer.Terminated()
	go func() {
		edgeSetTimer.Start()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
