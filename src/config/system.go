package config

type SystemConfig struct {
	DbType  string `yaml:"db-type"`  // 数据库类型 mysql/postgres
	OssType string `yaml:"oss-type"` // 存储类型
	Addr    string `yaml:"addr"`     // 服务地址
}
