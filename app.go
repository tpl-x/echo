package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tpl-x/echo/internal/config"
	"github.com/tpl-x/echo/internal/handler"
	applogger "github.com/tpl-x/echo/internal/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type App struct {
	config       *config.AppConfig
	logger       *zap.Logger
	engine       *echo.Echo
	serverCancel context.CancelFunc
	serverDone   chan error
}

func NewApp(lc fx.Lifecycle, config *config.AppConfig, logger *zap.Logger, handlers *handler.Handlers) *App {
	app := &App{
		config: config,
		logger: logger,
		engine: echo.New(),
	}

	// setup middleware
	app.engine.Use(middleware.RequestID())
	app.engine.Use(applogger.RequestLogger(logger))
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
			listenAddr := fmt.Sprintf(":%d", config.Server.BindPort)
			serverCtx, cancel := context.WithCancel(context.Background())
			app.serverCancel = cancel
			app.serverDone = make(chan error, 1)
			go func() {
				err := (echo.StartConfig{
					Address:         listenAddr,
					HideBanner:      true,
					HidePort:        true,
					GracefulTimeout: time.Duration(config.Server.GraceExitTimeout) * time.Second,
					OnShutdownError: func(err error) {
						logger.Error("server shutdown error", zap.Error(err))
					},
				}).Start(serverCtx, app.engine)
				if err != nil {
					logger.Error("server exited", zap.Error(err))
				}
				app.serverDone <- err
			}()
			logger.Info("server started", zap.String("address", listenAddr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Warn("server will shutdown", zap.Int("in_seconds", config.Server.GraceExitTimeout))
			if app.serverCancel == nil || app.serverDone == nil {
				return nil
			}
			app.serverCancel()
			shutdownCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Server.GraceExitTimeout)*time.Second)
			defer cancel()
			select {
			case err := <-app.serverDone:
				if err != nil && err != http.ErrServerClosed {
					return err
				}
				logger.Info("server shutdown complete")
				return nil
			case <-shutdownCtx.Done():
				return shutdownCtx.Err()
			}
		},
	})

	return app
}
