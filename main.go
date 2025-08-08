package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"github.com/tpl-x/echo/internal/biz"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/data"
	"github.com/tpl-x/echo/internal/handler"
	"github.com/tpl-x/echo/internal/pkg/logger"
	"go.uber.org/fx"
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

	fxApp := fx.New(
		fx.Provide(config.ProvideConfig(configPath)),
		logger.Module,
		data.Module,
		biz.Module,
		handler.Module,
		fx.Provide(NewApp),
		fx.NopLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zap.L().Info("Received signal", zap.String("signal", sig.String()))
		cancel()
	}()

	if err := fxApp.Start(ctx); err != nil {
		panic("failed to start fx app")
	}

	<-ctx.Done()

	if err := fxApp.Stop(context.Background()); err != nil {
		zap.L().Error("Server shutdown error", zap.Error(err))
	}

	zap.L().Info("Server shutdown complete")
}
