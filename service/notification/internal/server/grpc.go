package server

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewGRPCServer(cfg *Config) *grpc.Server {
	return grpc.NewServer(
		grpc.Addr(":9011"),
	)
}

func NewHTTPServer(cfg *Config) *http.Server {
	return http.NewServer(
		http.Addr(":8011"),
	)
}

func NewApp(gs *grpc.Server, hs *http.Server, logger log.Logger) *kratos.App {
	return kratos.New(
		kratos.ID("notification-service"),
		kratos.Name("notification-service"),
		kratos.Version("v1.0.0"),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}