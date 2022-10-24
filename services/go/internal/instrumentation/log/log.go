package log

import (
	"context"
	"strings"

	"go.uber.org/zap"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"
)

type Logger = logur.Logger
type Fields = logur.Fields
type contextKey int

var WithFields = logur.WithFields

const (
	logCtxKey contextKey = iota
)
const (
	TagError string = "err"
)

func NewLogger(ctx context.Context, cfg *Config, opts ...zap.Option) (logur.Logger, error) {

	inner, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	if !strings.Contains(strings.ToLower(cfg.Environment), "dev") {
		inner, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	}

	inner = inner.WithOptions(opts...)
	return zapadapter.New(inner), nil
}

func WithLogger(ctx context.Context, logger logur.Logger) context.Context {
	return context.WithValue(ctx, logCtxKey, logger)
}

// GetLogger tries to get the logger out of context.
// If one cannot be found, it will wrap the global zap logger with the interface.
func GetLogger(ctx context.Context) logur.Logger {

	if l, ok := ctx.Value(logCtxKey).(logur.Logger); ok {
		return l
	}
	l, _ := zap.NewDevelopment()
	return zapadapter.New(l)
}
