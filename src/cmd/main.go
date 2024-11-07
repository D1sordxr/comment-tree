package main

import (
	"fmt"
	"go.uber.org/fx"
	config2 "simpleBlog/src/internal/presentation/config"
	"simpleBlog/src/internal/presentation/di/config"
)

func main() {
	cfg := config2.NewConfig()
	fmt.Println(cfg)

	fx.New(config.Module).Run()
}
