package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func Load(appName string) (*Config, error) {
	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.AppName = appName

	return &cfg, nil
}
