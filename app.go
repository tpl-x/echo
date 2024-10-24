package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thejerf/suture/v4"
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var _ suture.Service = (*app)(nil)

type app struct {
	config *config.AppConfig
	logger *zap.Logger
	engine *echo.Echo
}

func (a *app) Serve(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return a.stop()
		default:
			a.start()
		}
	}
}

func newApp(config *config.AppConfig, logger *zap.Logger) *app {
	return &app{
		config: config,
		logger: logger,
		engine: echo.New(),
	}
}

func (a *app) stop() error {
	if a.engine != nil {
		a.logger.Warn("server will  shutdown", zap.Int("in seconds", a.config.Server.GraceExitTimeout))
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Server.GraceExitTimeout)*time.Second)
		defer cancel()
		a.engine.Shutdown(ctx)
	}
	return nil
}

func (a *app) start() error {
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

	// start listen
	listenAddr := fmt.Sprintf(":%d", a.config.Server.BindPort)
	if err := a.engine.Start(listenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Fatal("shutting down the server")
		return err
	}
	return nil
}
