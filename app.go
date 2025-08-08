package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/handler"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type App struct {
	config *config.AppConfig
	logger *zap.Logger
	engine *echo.Echo
}

func NewApp(lc fx.Lifecycle, config *config.AppConfig, logger *zap.Logger, handlers *handler.Handlers) *App {
	app := &App{
		config: config,
		logger: logger,
		engine: echo.New(),
	}

	// setup middleware
	app.engine.Use(echozap.ZapLogger(logger))
	app.engine.Use(middleware.Recover())
	app.engine.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// setup routes
	app.engine.GET("/", handlers.Hello.Hello)

	// User routes
	userGroup := app.engine.Group("/users")
	userGroup.GET("", handlers.User.ListUsers)
	userGroup.POST("", handlers.User.CreateUser)
	userGroup.GET("/:id", handlers.User.GetUser)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				listenAddr := fmt.Sprintf(":%d", config.Server.BindPort)
				if err := app.engine.Start(listenAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("shutting down the server")
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Warn("server will shutdown", zap.Int("in seconds", config.Server.GraceExitTimeout))
			shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Server.GraceExitTimeout)*time.Second)
			defer cancel()
			return app.engine.Shutdown(shutdownCtx)
		},
	})

	return app
}
