package logger

import (
	"item-service/config"
	"item-service/utils/constants"
	"log"
	"os"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var LoggerModule = fx.Provide(NewZapLogger)

type _zap struct{}

func (z *_zap) GetEncoder(config config.Config) zapcore.Encoder {
	if config.Logger.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig(config))
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig(config))
}

func (z *_zap) GetEncoderConfig(config config.Config) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     z.CustomTimeEncoder(config.Logger.Prefix),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

func (z *_zap) CustomTimeEncoder(prefix string) zapcore.TimeEncoder {
	return func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(prefix + t.Format("2006/01/02 - 15:04:05.000"))
	}
}

func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool {
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool {
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool {
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool {
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool {
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool {
			return level == zap.DebugLevel
		}
	}
}

func NewZapLogger(config config.Config) *zap.Logger {
	var (
		zapLogger *zap.Logger
		zapLogLvl zapcore.Level
	)
	var Zap = new(_zap)

	err := zapLogLvl.Set(config.Logger.Level)
	if err != nil {
		log.Println("cannot parse logLevel, err:", err.Error())
		zapLogLvl = zap.WarnLevel
	}

	switch config.Server.Env {
	case constants.Dev, constants.Prod:
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		consoleEncoder := Zap.GetEncoder(config)
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleErrors, Zap.GetLevelPriority(zapcore.ErrorLevel)),
			zapcore.NewCore(consoleEncoder, consoleInfos, Zap.GetLevelPriority(zapcore.InfoLevel)),
		)

		zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
		zap.RedirectStdLog(zapLogger)

		return zapLogger
	default:
		c := zap.NewDevelopmentConfig()
		c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		c.Level = zap.NewAtomicLevelAt(zapLogLvl)
		zapLogger, err = c.Build()
		if err != nil {
			log.Println("cannot build logger, err:", err.Error())
		}
		return zapLogger
	}
}
