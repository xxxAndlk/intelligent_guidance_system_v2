//go:build ignore

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"

	"intelligent_guidance_system_v2/service/medical/internal/acl"
	"intelligent_guidance_system_v2/service/medical/internal/biz/usecase"
	"intelligent_guidance_system_v2/service/medical/internal/data/mysql"
	"intelligent_guidance_system_v2/service/medical/internal/server"
	"intelligent_guidance_system_v2/service/medical/internal/service"
)

var infrastructureSet = wire.NewSet(
	server.ProviderSet,
)

func wireApp(cfg *server.Config, logger log.Logger) (*kratos.App, func(), error) {
	wire.Build(
		infrastructureSet,
		newApp,
	)
	return nil, nil, nil
}

func newApp(gs *grpc.Server, hs *http.Server, cfg *server.Config, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID("medical-service"),
		kratos.Name(cfg.Server.Name),
		kratos.Version("v1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}