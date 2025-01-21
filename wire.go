//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/google/wire"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/pkg/logger"
)

// wireApp init for builder backend
func wireApp(conf *config.AppConfig) *app {
	panic(wire.Build(
		logger.ProviderSet,
		newApp))
}
