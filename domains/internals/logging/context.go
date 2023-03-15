package logging

import (
	"context"
	"log"
)

type contextKey string

const contextKeyLogger contextKey = "logger"

// CtxWithLogger returns child context with passed logger.
func CtxWithLogger(baseCtx context.Context, logger *log.Logger) context.Context {
	return context.WithValue(baseCtx, contextKeyLogger, logger)
}

// LoggerFromCtx returns logger from given context.
func LoggerFromCtx(ctx context.Context) *log.Logger {
	value := ctx.Value(contextKeyLogger)

	logger, ok := value.(*log.Logger)
	if !ok {
		return nil
	}

	return logger
}
