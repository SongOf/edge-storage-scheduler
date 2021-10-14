package conf

type PrometheusConf struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}
