package logger

import (
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/fx"
)

var Module = fx.Module("logger",
	fx.Provide(
		func(appConfig *config.AppConfig) *config.LogConfig {
			return &appConfig.Log
		},
		NewZapLogger,
	),
)
