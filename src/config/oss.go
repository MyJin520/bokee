package config

type OssConfig struct {
	LocalOss LocalOss `yaml:"local-oss"` // 本地存储
	MinioOss MinioOss `yaml:"minio-oss"` // Minio 存储
	AliOss   AliOss   `yaml:"ali-oss"`   // 阿里云 OSS
}

type LocalOss struct {
	Path string `yaml:"path"` // 存储路径
}

type MinioOss struct {
	Url string `yaml:"url"` // 服务地址
}

type AliOss struct {
	Url string `yaml:"url"` // 服务地址
}
