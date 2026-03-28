package server

import (
	pb "intelligent_guidance_system_v2/service/ai/api"

	"intelligent_guidance_system_v2/service/ai/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(c *Config, aiSvc *service.AIService, logger log.Logger) *grpc.Server {
	srv := grpc.NewServer(
		grpc.Address(c.Server.GRPC.Addr),
		grpc.Timeout(c.Server.GRPC.Timeout),
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.Validator(),
		),
	)

	pb.RegisterAIServiceServer(srv, aiSvc)

	return srv
}