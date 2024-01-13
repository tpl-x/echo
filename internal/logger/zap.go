package logger

import (
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewZapLogger(config *config.AppConfig) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.Log.FileName,
		MaxSize:    config.Log.MaxSize,
		MaxBackups: config.Log.MaxBackups,
		MaxAge:     config.Log.MaxKeepDays,
		Compress:   config.Log.Compress,
	}
	writeSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberjackLogger))
	encoder := zapcore.NewConsoleEncoder(encoderCfg)

	opt := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)
	})
	zapLogger, _ := zap.NewProduction(opt)
	return zapLogger
}
