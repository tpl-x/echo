package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type HelloHandler struct {
	logger *zap.Logger
}

func NewHelloHandler(logger *zap.Logger) *HelloHandler {
	return &HelloHandler{
		logger: logger,
	}
}

func (h *HelloHandler) Hello(c echo.Context) error {
	h.logger.Info("Hello endpoint called")
	return c.String(200, "hello,world!")
}
