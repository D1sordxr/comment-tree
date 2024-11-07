package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

const BasicConfigPath = "./src/configs/dev.yaml"

func MustLoad(val interface{}) {
	path := fetchConfigPath()

	if path == "" {
		path = BasicConfigPath
	}
	fmt.Println(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("app file does not exist: " + path)
	}

	if err := cleanenv.ReadConfig(path, val); err != nil {
		panic("failed to read app" + err.Error())
	}

}

func fetchConfigPath() string {
	var res string

	// go run ./src/cmd/main.go --config="./src/configs/dev.yaml"
	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
