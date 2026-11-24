package logger

import (
	"log"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/pavel-pj/auth/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globbalLogger *zap.Logger

func Init(cfg config.LoggerConfig) {

	core := getCore(cfg)

	//core zapcore.Core, options ...zap.Option
	globbalLogger = zap.New(core)
}

func Debug(msg string, fields ...zap.Field) {
	globbalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globbalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globbalLogger.Error(msg, fields...)
}

func getAtomicLevel(logLevel string) zap.AtomicLevel {

	var level zapcore.Level

	if err := level.Set(logLevel); err != nil {
		log.Fatal("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func getCore(cfg config.LoggerConfig) zapcore.Core {

	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename(),
		MaxSize:    cfg.MaxSize(),
		MaxBackups: cfg.MaxBackups(),
		MaxAge:     cfg.MaxAge(),
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCgf := zap.NewDevelopmentEncoderConfig()
	developmentCgf.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCgf)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	level := getAtomicLevel(cfg.Level())

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

}
