package config

import (
	"go.uber.org/fx"
	"simpleBlog/src/internal/presentation/config"
)

var Module = fx.Module(
	"presentation.config",
	fx.Provide(
		config.NewConfig,
		config.NewStorageConfig,
		config.NewAPIConfig,
		config.NewAppConfig,
	),
)
