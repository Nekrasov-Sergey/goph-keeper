package config

import (
	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

type ClientConfig struct {
	GRPCAddr string
}

func NewClientConfig() (*ClientConfig, error) {
	viper.SetConfigFile(GetConfigPath())

	cfg := ClientConfig{
		GRPCAddr: ":8081",
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "не удалось прочитать конфигурацию из файла")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "не удалось распарсить конфигурацию в структуру")
	}

	return &cfg, nil
}
