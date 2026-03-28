package logging

import (
	"context"
	"time"

	kratos "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

func LoggingMiddleware(logger Logger) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			startTime := time.Now()

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			var operation string
			if tr.Operation() != "" {
				operation = tr.Operation()
			} else {
				operation = "unknown"
			}

			logger.Info("request started",
				zap.String("operation", operation),
				zap.String("kind", tr.Kind().String()),
			)

			reply, err := handler(ctx, req)

			duration := time.Since(startTime)

			statusCode := "OK"
			if err != nil {
				statusCode = "ERROR"
				logger.Error("request completed",
					zap.String("operation", operation),
					zap.String("kind", tr.Kind().String()),
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				logger.Info("request completed",
					zap.String("operation", operation),
					zap.String("kind", tr.Kind().String()),
					zap.Duration("duration", duration),
					zap.String("status", statusCode),
				)
			}

			return reply, err
		}
	}
}

func LoggingMiddlewareWithConfig(logger Logger, skipPaths []string) kratos.Middleware {
	skipMap := make(map[string]struct{})
	for _, path := range skipPaths {
		skipMap[path] = struct{}{}
	}

	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			startTime := time.Now()

			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			operation := tr.Operation()
			if _, skip := skipMap[operation]; skip {
				return handler(ctx, req)
			}

			logger.Info("request started",
				zap.String("operation", operation),
				zap.String("kind", tr.Kind().String()),
			)

			reply, err := handler(ctx, req)

			duration := time.Since(startTime)

			if err != nil {
				logger.Error("request completed",
					zap.String("operation", operation),
					zap.String("kind", tr.Kind().String()),
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				logger.Info("request completed",
					zap.String("operation", operation),
					zap.String("kind", tr.Kind().String()),
					zap.Duration("duration", duration),
				)
			}

			return reply, err
		}
	}
}

func RequestLogger(logger *zap.Logger) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			startTime := time.Now()

			tr, _ := transport.FromServerContext(ctx)

			reqLogger := logger.With(
				zap.String("start_time", startTime.Format(time.RFC3339)),
			)

			if tr != nil {
				reqLogger = reqLogger.With(
					zap.String("operation", tr.Operation()),
					zap.String("kind", tr.Kind().String()),
				)
			}

			reply, err := handler(ctx, req)

			duration := time.Since(startTime)
			if err != nil {
				reqLogger.Error("request failed",
					zap.Duration("duration", duration),
					zap.Error(err),
				)
			} else {
				reqLogger.Info("request succeeded",
					zap.Duration("duration", duration),
				)
			}

			return reply, err
		}
	}
}