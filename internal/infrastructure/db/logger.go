package db

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// GormLogrusLogger adapts logrus.Entry to satisfy gorm/logger.Interface.
type GormLogrusLogger struct {
	entry                *logrus.Entry
	level                logger.LogLevel
	slowThreshold        time.Duration
	ignoreRecordNotFound bool
	parameterizedQueries bool
}

type GormLogrusOption func(*GormLogrusLogger)

func WithSlowThreshold(d time.Duration) GormLogrusOption {
	return func(g *GormLogrusLogger) { g.slowThreshold = d }
}

func WithIgnoreRecordNotFound(v bool) GormLogrusOption {
	return func(g *GormLogrusLogger) { g.ignoreRecordNotFound = v }
}

func WithParameterizedQueries(v bool) GormLogrusOption {
	return func(g *GormLogrusLogger) { g.parameterizedQueries = v }
}

// NewGormLogrusLogger creates a GORM logger backed by a logrus.Entry.
// gormLevel mirrors config.Database.LogLevel: "silent"|"error"|"warn"|"info".
func NewGormLogrusLogger(entry *logrus.Entry, gormLevel string, opts ...GormLogrusOption) logger.Interface {
	var lvl logger.LogLevel
	switch gormLevel {
	case "silent":
		lvl = logger.Silent
	case "error":
		lvl = logger.Error
	case "info":
		lvl = logger.Info
	default: // "warn" or empty
		lvl = logger.Warn
	}

	g := &GormLogrusLogger{
		entry:                entry,
		level:                lvl,
		slowThreshold:        200 * time.Millisecond,
		ignoreRecordNotFound: true,
		parameterizedQueries: false,
	}
	for _, o := range opts {
		o(g)
	}
	return g
}

// LogMode implements logger.Interface.
func (g *GormLogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *g
	clone.level = level
	return &clone
}

// Info implements logger.Interface.
func (g *GormLogrusLogger) Info(ctx context.Context, msg string, args ...any) {
	if g.level >= logger.Info {
		g.entry.WithContext(ctx).
			WithField("component", "gorm").
			Infof(msg, args...)
	}
}

// Warn implements logger.Interface.
func (g *GormLogrusLogger) Warn(ctx context.Context, msg string, args ...any) {
	if g.level >= logger.Warn {
		g.entry.WithContext(ctx).
			WithField("component", "gorm").
			Warnf(msg, args...)
	}
}

// Error implements logger.Interface.
func (g *GormLogrusLogger) Error(ctx context.Context, msg string, args ...any) {
	if g.level >= logger.Error {
		g.entry.WithContext(ctx).
			WithField("component", "gorm").
			Errorf(msg, args...)
	}
}

// Trace implements logger.Interface — called for every SQL statement.
func (g *GormLogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if g.level <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := logrus.Fields{
		"component":     "gorm",
		"elapsed_ms":    float64(elapsed.Nanoseconds()) / 1e6,
		"rows_affected": rows,
	}

	if !g.parameterizedQueries {
		fields["sql"] = sql
	}

	entry := g.entry.WithContext(ctx).WithFields(fields)

	switch {
	case err != nil && g.level >= logger.Error &&
		!(errors.Is(err, logger.ErrRecordNotFound) && g.ignoreRecordNotFound):
		entry.WithError(err).Error("gorm query error")

	case g.slowThreshold > 0 && elapsed > g.slowThreshold && g.level >= logger.Warn:
		entry.Warnf("gorm slow query >= %s", g.slowThreshold)

	case g.level >= logger.Info:
		entry.Info("gorm query")
	}
}
