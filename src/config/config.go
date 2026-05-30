package config

type Config struct {
	System SystemConfig `yaml:"system"` // 系统配置
	DB     DBConfig     `yaml:"db"`     // 数据库配置
}
