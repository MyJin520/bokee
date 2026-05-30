package config

type SystemConfig struct {
	DbType  string `yaml:"db-type"`
	OssType string `yaml:"oss-type"`
	Addr    string `yaml:"addr"`
}
