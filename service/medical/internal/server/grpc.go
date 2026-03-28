package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	v1 "intelligent_guidance_system_v2/api/medical/v1"
	"intelligent_guidance_system_v2/service/medical/internal/service"
)

func NewGRPCServer(medicalSvc *service.MedicalService) *grpc.Server {
	srv := grpc.NewServer(
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	v1.RegisterMedicalServiceServer(srv, medicalSvc)

	return srv
}

func NewHTTPServer(grpcSrv *grpc.Server) *http.Server {
	srv := http.NewServer(
		http.Middleware(
			recovery.Recovery(),
		),
	)

	v1.RegisterMedicalServiceHTTPServer(srv, grpcSrv)

	return srv
}