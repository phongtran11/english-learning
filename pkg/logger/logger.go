package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(env string) {
	var config zap.Config

	if env == "prod" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		os.Exit(1)
	}
}

func Infof(template string, args ...interface{}) {
	Log.Sugar().Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Log.Sugar().Errorf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	Log.Sugar().Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Log.Sugar().Warnf(template, args...)
}
