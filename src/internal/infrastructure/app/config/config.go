package config

type AppConfig struct {
	Mode string `yaml:"env" env-default:"local"`
}
