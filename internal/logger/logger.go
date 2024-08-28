package logger

import (
	"github.com/google/wire"
	"github.com/tpl-x/echo/internal/config"
)

var ProviderSet = wire.NewSet(
	wire.FieldsOf(new(*config.AppConfig), "Log"),
	NewZapLogger,
)
