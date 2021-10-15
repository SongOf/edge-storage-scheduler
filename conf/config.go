package conf

type PrometheusConf struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

type RedisConf struct {
	Address  string `yaml:"address,omitempty"`
	Password string `yaml:"password,omitempty"`
}
