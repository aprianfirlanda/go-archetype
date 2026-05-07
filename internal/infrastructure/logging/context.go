package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type loggerContextKey struct{}
type requestIDContextKey struct{}

func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDContextKey{}).(string)
	return requestID
}

func FromContext(ctx context.Context) *logrus.Entry {
	logger, _ := ctx.Value(loggerContextKey{}).(*logrus.Entry)
	if logger != nil {
		if requestID := RequestIDFromContext(ctx); requestID != "" {
			return logger.WithFields(logrus.Fields{
				"rid":        requestID,
				"request_id": requestID,
			})
		}
		return logger
	}

	entry := logrus.NewEntry(logrus.StandardLogger())
	if requestID := RequestIDFromContext(ctx); requestID != "" {
		return entry.WithFields(logrus.Fields{
			"rid":        requestID,
			"request_id": requestID,
		})
	}

	return entry
}
