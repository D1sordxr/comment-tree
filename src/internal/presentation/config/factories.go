package config

import (
	api "simpleBlog/src/internal/infrastructure/api/config"
	app "simpleBlog/src/internal/infrastructure/app/config"
	load "simpleBlog/src/internal/infrastructure/config"
	storage "simpleBlog/src/internal/infrastructure/db/config"
)

func NewConfig() *Config {
	var config Config
	load.MustLoad(&config)
	return &config
}

func NewAppConfig(config *Config) app.AppConfig {
	return config.AppConfig
}

func NewStorageConfig(config *Config) storage.StorageConfig {
	return config.StorageConfig
}

func NewAPIConfig(config *Config) api.APIConfig {
	return config.APIConfig
}
