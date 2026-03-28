package recovery

import (
	"context"
	"fmt"
	"runtime"

	"intelligent_guidance_system_v2/common/pkg/errors"
	"intelligent_guidance_system_v2/common/pkg/response"

	kratos "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.uber.org/zap"
)

type RecoveryOption func(*recoveryOptions)

type recoveryOptions struct {
	logger   *zap.Logger
	handlers []RecoveryHandler
}

type RecoveryHandler func(ctx context.Context, req, reply interface{}, err error, stack []byte) error

func WithLogger(logger *zap.Logger) RecoveryOption {
	return func(o *recoveryOptions) {
		o.logger = logger
	}
}

func WithHandler(handler RecoveryHandler) RecoveryOption {
	return func(o *recoveryOptions) {
		o.handlers = append(o.handlers, handler)
	}
}

func RecoveryMiddleware(opts ...RecoveryOption) kratos.Middleware {
	options := &recoveryOptions{
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					stack := make([]byte, 4096)
					n := runtime.Stack(stack, false)
					stack = stack[:n]

					err := fmt.Errorf("panic recovered: %v", r)

					options.logger.Error("panic recovered",
						zap.Error(err),
						zap.String("stack", string(stack)),
					)

					for _, h := range options.handlers {
						if handlerErr := h(ctx, req, nil, err, stack); handlerErr != nil {
							options.logger.Error("recovery handler error", zap.Error(handlerErr))
						}
					}
				}
			}()

			return handler(ctx, req)
		}
	}
}

func RecoveryMiddlewareWithResponse(opts ...RecoveryOption) kratos.Middleware {
	options := &recoveryOptions{
		logger: zap.NewNop(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				if r := recover(); r != nil {
					stack := make([]byte, 4096)
					n := runtime.Stack(stack, false)
					stack = stack[:n]

					err := fmt.Errorf("panic recovered: %v", r)

					tr, _ := transport.FromServerContext(ctx)
					operation := "unknown"
					if tr != nil {
						operation = tr.Operation()
					}

					options.logger.Error("panic recovered",
						zap.Error(err),
						zap.String("operation", operation),
						zap.String("stack", string(stack)),
					)

					for _, h := range options.handlers {
						if handlerErr := h(ctx, req, nil, err, stack); handlerErr != nil {
							options.logger.Error("recovery handler error", zap.Error(handlerErr))
						}
					}
				}
			}()

			reply, err := handler(ctx, req)
			if err != nil {
				return reply, err
			}

			return reply, err
		}
	}
}

func DefaultRecoveryHandler() RecoveryHandler {
	return func(ctx context.Context, req, reply interface{}, err error, stack []byte) error {
		return errors.NewInternal("internal server error", errors.WithInternal(err))
	}
}

func RecoveryResponse() *response.Response {
	return response.Error(errors.CodeInternal, "internal server error")
}