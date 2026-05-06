package logger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZapLogger(lc fx.Lifecycle, logConfig *config.LogConfig) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logConfig.FileName,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxKeepDays,
		Compress:   logConfig.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(lumberjackLogger),
	)
	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewCore(
		encoder,
		writeSyncer,
		zap.NewAtomicLevelAt(parseLevel(logConfig.Level)),
	)
	zapLogger := zap.New(core, zap.AddCaller())
	restoreGlobals := zap.ReplaceGlobals(zapLogger)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			restoreGlobals()
			return sync(zapLogger)
		},
	})
	return zapLogger
}

func RequestLogger(log *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:     true,
		LogLatency:      true,
		LogProtocol:     true,
		LogRemoteIP:     true,
		LogHost:         true,
		LogMethod:       true,
		LogURI:          true,
		LogRoutePath:    true,
		LogRequestID:    true,
		LogUserAgent:    true,
		LogStatus:       true,
		LogResponseSize: true,
		LogValuesFunc: func(_ *echo.Context, v middleware.RequestLoggerValues) error {
			fields := []zap.Field{
				zap.String("request_id", v.RequestID),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("host", v.Host),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.String("route", v.RoutePath),
				zap.String("protocol", v.Protocol),
				zap.Int("status", v.Status),
				zap.Int64("response_size", v.ResponseSize),
				zap.Duration("latency", v.Latency),
				zap.String("user_agent", v.UserAgent),
			}
			if v.Error != nil {
				fields = append(fields, zap.Error(v.Error))
			}

			switch {
			case v.Status >= http.StatusInternalServerError:
				log.Error("request completed", fields...)
			case v.Status >= http.StatusBadRequest:
				log.Warn("request completed", fields...)
			default:
				log.Info("request completed", fields...)
			}
			return nil
		},
	})
}

func parseLevel(level string) zapcore.Level {
	var parsed zapcore.Level
	if err := parsed.UnmarshalText([]byte(level)); err != nil {
		return zapcore.InfoLevel
	}
	return parsed
}

func sync(log *zap.Logger) error {
	if err := log.Sync(); err != nil &&
		!errors.Is(err, os.ErrInvalid) &&
		!errors.Is(err, syscall.EINVAL) &&
		!errors.Is(err, syscall.ENOTTY) {
		return fmt.Errorf("sync logger: %w", err)
	}
	return nil
}
