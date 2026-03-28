package server

import (
	"intelligent_guidance_system_v2/service/ai/internal/biz/dto"
	"intelligent_guidance_system_v2/service/ai/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(c *Config, aiSvcHTTP *service.AIServiceHTTP, logger log.Logger) *khttp.Server {
	srv := khttp.NewServer(
		khttp.Address(c.Server.HTTP.Addr),
		khttp.Timeout(c.Server.HTTP.Timeout),
		khttp.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)

	registerHTTPRoutes(srv, aiSvcHTTP)

	return srv
}

func registerHTTPRoutes(srv *khttp.Server, svc *service.AIServiceHTTP) {
	srv.Route("/api/v1").POST("/sessions", func(ctx khttp.Context) error {
		var req dto.CreateSessionRequest
		if err := ctx.BindVars(&req); err != nil {
			return ctx.JSON(400, map[string]string{"error": "invalid request"})
		}

		resp, err := svc.DiagnosisUC().CreateSession(ctx, &req)
		if err != nil {
			return ctx.JSON(500, map[string]string{"error": err.Error()})
		}

		return ctx.JSON(200, resp)
	})

	srv.Route("/api/v1").GET("/health", func(ctx khttp.Context) error {
		resp, err := svc.HealthCheck(ctx)
		if err != nil {
			return ctx.JSON(500, map[string]string{"error": err.Error()})
		}

		return ctx.JSON(200, resp)
	})
}