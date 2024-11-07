package config

type StorageConfig struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

// DB: "./src/internal/infrastructure/db/storage.db"
