package config

import (
	api "simpleBlog/src/internal/infrastructure/api/config"
	app "simpleBlog/src/internal/infrastructure/app/config"
	storage "simpleBlog/src/internal/infrastructure/db/config"
)

type Config struct {
	app.AppConfig
	storage.StorageConfig
	api.APIConfig
}
