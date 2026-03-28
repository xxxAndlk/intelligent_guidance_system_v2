package server

import (
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	v1 "intelligent_guidance_system_v2/api/medical/v1"
	"intelligent_guidance_system_v2/service/medical/internal/service"
)

func NewHTTPServerDirect(medicalSvc *service.MedicalService) *http.Server {
	srv := http.NewServer(
		http.Middleware(
			recovery.Recovery(),
		),
	)

	v1.RegisterMedicalServiceHTTPServer(srv, medicalSvc)

	return srv
}