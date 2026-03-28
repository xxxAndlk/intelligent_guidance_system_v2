// service/patient/internal/server/http.go
package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"

	pb "intelligent-guidance-system/api/patient/v1"
	"intelligent-guidance-system/service/patient/internal/service"
)

func NewHTTPServer(c *Config, patientSvc *service.PatientService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.Validator(),
		),
	}

	if c.HTTP.Network != "" {
		opts = append(opts, http.Network(c.HTTP.Network))
	}
	if c.HTTP.Addr != "" {
		opts = append(opts, http.Address(c.HTTP.Addr))
	}
	if c.HTTP.Timeout > 0 {
		opts = append(opts, http.Timeout(c.HTTP.Timeout))
	}

	srv := http.NewServer(opts...)
	pb.RegisterPatientServiceHTTPServer(srv, patientSvc)

	return srv
}