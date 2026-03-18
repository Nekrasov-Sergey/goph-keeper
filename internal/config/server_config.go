// Package config содержит конфигурацию клиента и сервера.
package config

import (
	"encoding/base64"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// ServerConfig содержит конфигурацию сервера.
type ServerConfig struct {
	GRPCAddr    string
	DatabaseDSN string
	JWTSecret   []byte
	MasterKey   []byte
	TLSCertFile string
	TLSKeyFile  string
}

// rawServerConfig содержит конфигурацию сервера до декодирования ключей.
type rawServerConfig struct {
	GRPCAddr    string
	DatabaseDSN string
	JWTSecret   string
	MasterKey   string
	TLSCertFile string
	TLSKeyFile  string
}

// GetConfigPath возвращает путь к файлу конфигурации.
func GetConfigPath() string {
	c := os.Getenv("CONFIG_PATH")
	if c == "" {
		c = "./config/local.yml"
	}
	return c
}

// NewServerConfig загружает конфигурацию сервера из файла.
func NewServerConfig(logger zerolog.Logger) (*ServerConfig, error) {
	viper.SetConfigFile(GetConfigPath())

	raw := rawServerConfig{
		GRPCAddr:    ":8081",
		TLSCertFile: "certs/server.crt",
		TLSKeyFile:  "certs/server.key",
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "не удалось прочитать конфигурацию из файла")
	}

	if err := viper.Unmarshal(&raw); err != nil {
		return nil, errors.Wrap(err, "не удалось распарсить конфигурацию в структуру")
	}

	jwtSecret, err := mustDecodeKey(raw.JWTSecret)
	if err != nil {
		return nil, err
	}

	masterKey, err := mustDecodeKey(raw.MasterKey)
	if err != nil {
		return nil, err
	}

	cfg := ServerConfig{
		GRPCAddr:    raw.GRPCAddr,
		DatabaseDSN: raw.DatabaseDSN,
		JWTSecret:   jwtSecret,
		MasterKey:   masterKey,
		TLSCertFile: raw.TLSCertFile,
		TLSKeyFile:  raw.TLSKeyFile,
	}

	logger.Info().
		Str("GRPCAddr", cfg.GRPCAddr).
		Str("DatabaseDSN", cfg.DatabaseDSN).
		Str("TLSCertFile", cfg.TLSCertFile).
		Str("TLSKeyFile", cfg.TLSKeyFile).
		Msg("Загружена конфигурация сервера")

	return &cfg, nil
}

// mustDecodeKey декодирует ключ из base64 и проверяет его длину.
func mustDecodeKey(key string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, errors.Wrap(err, "не удалось декодировать ключ")
	}

	if len(b) != 32 {
		return nil, errors.New("ключ должен быть 32 байт")
	}

	return b, nil
}
