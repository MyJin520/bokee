package config

type LogConfig struct {
	Level      string `yaml:"level"`       // 日志级别 debug/info/warn/error
	Format     string `yaml:"format"`      // 输出格式 json/console
	Output     string `yaml:"output"`      // 输出路径 stdout/file
	FilePath   string `yaml:"file-path"`   // 文件路径
	MaxSize    int    `yaml:"max-size"`    // 单个日志文件最大大小(MB)
	MaxBackups int    `yaml:"max-backups"` // 保留的旧日志文件最大数量
	MaxAge     int    `yaml:"max-age"`     // 保留旧日志文件的最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}
