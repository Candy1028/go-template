package logger

import (
	"context"
	"go-template/log"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	LogLevel logger.LogLevel
}

func NewGormLogger(level logger.LogLevel) *GormLogger {
	return &GormLogger{LogLevel: level}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Logger.WithContext(ctx).Infof(msg, data...)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Logger.WithContext(ctx).Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Logger.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	entry := log.Logger.WithContext(ctx).WithFields(map[string]interface{}{
		"duration": elapsed.String(),
		"rows":     rows,
		"sql":      sql,
	})

	if err != nil && l.LogLevel >= logger.Error {
		entry.WithField("error", err).Error("GORM SQL 错误")
	} else if elapsed > 200*time.Millisecond && l.LogLevel >= logger.Warn {
		entry.Warn("GORM 慢查询")
	} else if l.LogLevel >= logger.Info {
		entry.Info("GORM 执行SQL")
	}
}
