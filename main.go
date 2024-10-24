package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"github.com/thejerf/suture/v4"
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/zap"
)

var (
	configPath string
)

func init() {
	flag.StringVarP(&configPath, "config", "c", "config.yaml", "start using config file")
}

func main() {
	flag.Parse()

	appConfig, err := config.LoadFromFile(configPath)
	if err != nil {
		panic("failed to load config")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sup := suture.NewSimple("echoWebServer")
	app := wireApp(appConfig)
	sup.Add(app)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	if err := sup.Serve(ctx); err != nil {
		logger.Error("Server is about to shutdown", zap.Error(err))
	}

	logger.Info("Server shutdown complete")
}
