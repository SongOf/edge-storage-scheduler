package cmd

import (
	"edge-storage-scheduler/conf"
	"edge-storage-scheduler/internal/globals"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
	"strings"
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

}

func main() {

}
