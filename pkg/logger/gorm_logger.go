package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses zap
type GormLogger struct {
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

// NewGormLogger creates a new GormLogger
func NewGormLogger(zapLogger *zap.Logger, slowThreshold time.Duration) gormlogger.Interface {
	return &GormLogger{
		ZapLogger:     zapLogger.Named("gorm"),
		SlowThreshold: slowThreshold,
	}
}

// LogMode implements gormlogger.Interface
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info implements gormlogger.Interface
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Infof(msg, data...)
}

// Warn implements gormlogger.Interface
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Warnf(msg, data...)
}

// Error implements gormlogger.Interface
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Sugar().Errorf(msg, data...)
}

// Trace implements gormlogger.Interface
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
		zap.String("sql", sql),
	}

	if err != nil && !errors.Is(err, gormlogger.ErrRecordNotFound) {
		l.ZapLogger.Error("SQL execution error", append(fields, zap.Error(err))...)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.ZapLogger.Warn("Slow SQL query detected", fields...)
		return
	}

	l.ZapLogger.Debug("SQL query executed", fields...)
}
