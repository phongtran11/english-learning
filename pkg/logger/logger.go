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

func Infof(name string, template string, args ...interface{}) {
	Log.Named(name).Sugar().Infof(template, args...)
}

func Errorf(name string, template string, args ...interface{}) {
	Log.Named(name).Sugar().Errorf(template, args...)
}

func Debugf(name string, template string, args ...interface{}) {
	Log.Named(name).Sugar().Debugf(template, args...)
}

func Warnf(name string, template string, args ...interface{}) {
	Log.Named(name).Sugar().Warnf(template, args...)
}
