package server

import (
	"intelligent-guidance-system/service/doctor/internal/biz/usecase"
	"intelligent-guidance-system/service/doctor/internal/data/mysql"
	"intelligent-guidance-system/service/doctor/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	mysql.NewDoctorRepository,
	mysql.NewScheduleRepository,
	mysql.NewRoleRepository,
	mysql.NewDoctorRoleAssignmentRepository,
	mysql.NewDoctorEventRepository,
	usecase.NewDoctorUseCase,
	usecase.NewScheduleUseCase,
	service.NewDoctorService,
	NewHTTPServer,
	NewGRPCServer,
)

func NewLogger() log.Logger {
	return log.NewStdLogger(log.DefaultWriter)
}

type App struct {
	HttpServer *http.Server
	GrpcServer *grpc.Server
	Service    *service.DoctorService
}

func NewApp(httpSrv *http.Server, grpcSrv *grpc.Server, svc *service.DoctorService) *App {
	return &App{
		HttpServer: httpSrv,
		GrpcServer: grpcSrv,
		Service:    svc,
	}
}

func WireApp(db *gorm.DB, logger log.Logger) (*App, func(), error) {
	panic(wire.Build(ProviderSet, NewApp))
}