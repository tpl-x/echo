package logger

import (
	"github.com/tpl-x/echo/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewZapLogger(logConfig *config.LogConfig) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
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
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	)
	zapLogger := zap.New(core, zap.AddCaller())
	return zapLogger
}
