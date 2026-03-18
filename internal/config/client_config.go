// Package config содержит конфигурацию клиента и сервера.
package config

import (
	"github.com/spf13/viper"

	"github.com/pkg/errors"
)

// ClientConfig содержит конфигурацию клиента.
type ClientConfig struct {
	GRPCAddr    string
	TLSCertFile string
}

// NewClientConfig загружает конфигурацию клиента из файла.
func NewClientConfig() (*ClientConfig, error) {
	viper.SetConfigFile(GetConfigPath())

	cfg := ClientConfig{
		GRPCAddr:    ":8081",
		TLSCertFile: "certs/server.crt",
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "не удалось прочитать конфигурацию из файла")
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "не удалось распарсить конфигурацию в структуру")
	}

	return &cfg, nil
}
