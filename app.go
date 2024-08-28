package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	config *config.AppConfig
	logger *zap.Logger
	engine *echo.Echo
}

func newApp(config *config.AppConfig, logger *zap.Logger) *app {
	return &app{
		config: config,
		logger: logger,
		engine: echo.New(),
	}
}

func (a *app) start() {
	// setup middleware
	a.engine.Use(echozap.ZapLogger(a.logger))
	a.engine.Use(middleware.Recover())
	a.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	// setup routers here

	// setup base hello world
	a.engine.GET("/", func(c echo.Context) error {
		return c.String(200, "hello,world!")
	})

	// start serve
	go func() {
		// start listen
		listenAddr := fmt.Sprintf(":%d", a.config.Server.BindPort)
		if err := a.engine.Start(listenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal("shutting down the server")
		}
	}()

	// wait signal to exit
	quit := make(chan os.Signal, 1)
	// capture SIGINT（Ctrl+C）and SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// set timeout to exit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Server.GraceExitTimeout)*time.Second)
	defer cancel()
	if err := a.engine.Shutdown(ctx); err != nil {
		a.logger.Fatal(err.Error())
	}

	a.logger.Info("Server gracefully stopped")
}
