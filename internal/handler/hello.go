package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
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

func (h *HelloHandler) Hello(c *echo.Context) error {
	h.logger.Info("Hello endpoint called")
	return c.String(http.StatusOK, "hello,world!")
}
