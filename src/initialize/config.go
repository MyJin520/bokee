package initialize

import (
	"bokee/config"
	"gopkg.in/yaml.v3"
	"os"
)

func initLoadConfig(configPath string) (*config.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
