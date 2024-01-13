package main

import (
	flag "github.com/spf13/pflag"
	"github.com/tpl-x/echo/internal/config"
)

var (
	configPath string
)

func init() {
	flag.StringVarP(&configPath, "config", "c", "config.yaml", "start using config file")
}
func main() {
	appConfig, err := config.LoadFromFile(configPath)
	if err != nil {
		panic("failed to load config")
	}
	app := wireApp(appConfig)
	app.start()
}
