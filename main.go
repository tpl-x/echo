package main

import (
	"context"
	flag "github.com/spf13/pflag"
	"github.com/thejerf/suture/v4"
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

	sup := suture.NewSimple("echoWebServer")
	app := wireApp(appConfig)
	sup.Add(app)
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	if err = sup.Serve(ctx); err != nil {
		panic(err)
	}
}
