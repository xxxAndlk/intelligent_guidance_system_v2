package server

import (
	"context"
	"time"

	v1 "intelligent-guidance-system/service/doctor/internal/biz/dto"
	"intelligent-guidance-system/service/doctor/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

func NewGRPCServer(doctorSvc *service.DoctorService, logger log.Logger) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	}

	srv := grpc.NewServer(opts...)

	RegisterDoctorGRPCServer(srv, doctorSvc)

	return srv
}

func RegisterDoctorGRPCServer(s *grpc.Server, svc *service.DoctorService) {
	s.RegisterService(&_DoctorService_ServiceDesc, &DoctorGRPCServer{service: svc})
}

type DoctorGRPCServer struct {
	service *service.DoctorService
}

func (s *DoctorGRPCServer) CreateDoctor(ctx context.Context, req *v1.DoctorCreateRequest) (*v1.DoctorDTO, error) {
	return s.service.CreateDoctor(ctx, req)
}

func (s *DoctorGRPCServer) GetDoctor(ctx context.Context, req *GetDoctorRequest) (*v1.DoctorDetailDTO, error) {
	return s.service.GetDoctor(ctx, req.DoctorId)
}

func (s *DoctorGRPCServer) UpdateDoctor(ctx context.Context, req *v1.DoctorUpdateRequest) (*v1.DoctorDTO, error) {
	return s.service.UpdateDoctor(ctx, req)
}

func (s *DoctorGRPCServer) DeleteDoctor(ctx context.Context, req *DeleteDoctorRequest) (*DeleteDoctorResponse, error) {
	err := s.service.DeleteDoctor(ctx, req.DoctorId)
	if err != nil {
		return nil, err
	}
	return &DeleteDoctorResponse{Message: "deleted"}, nil
}

func (s *DoctorGRPCServer) ListDoctors(ctx context.Context, req *v1.DoctorFilterRequest) (*v1.DoctorListResponse, error) {
	return s.service.ListDoctors(ctx, req)
}

func (s *DoctorGRPCServer) ListDoctorsByDepartment(ctx context.Context, req *ListByDepartmentRequest) (*DoctorListResponse, error) {
	result, err := s.service.ListDoctorsByDepartment(ctx, req.DepartmentId)
	if err != nil {
		return nil, err
	}
	return &DoctorListResponse{Items: result}, nil
}

func (s *DoctorGRPCServer) ListExpertDoctors(ctx context.Context, req *EmptyRequest) (*DoctorListResponse, error) {
	result, err := s.service.ListExpertDoctors(ctx)
	if err != nil {
		return nil, err
	}
	return &DoctorListResponse{Items: result}, nil
}

func (s *DoctorGRPCServer) ListDoctorsBySpecialty(ctx context.Context, req *ListBySpecialtyRequest) (*DoctorListResponse, error) {
	result, err := s.service.ListDoctorsBySpecialty(ctx, req.Specialty)
	if err != nil {
		return nil, err
	}
	return &DoctorListResponse{Items: result}, nil
}

func (s *DoctorGRPCServer) GetDoctorByEmployeeID(ctx context.Context, req *GetByEmployeeIDRequest) (*v1.DoctorDetailDTO, error) {
	return s.service.GetDoctorByEmployeeID(ctx, req.EmployeeId)
}

func (s *DoctorGRPCServer) ChangeDepartment(ctx context.Context, req *v1.DoctorDepartmentChangeRequest) (*MessageResponse, error) {
	err := s.service.ChangeDepartment(ctx, req)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "department changed"}, nil
}

func (s *DoctorGRPCServer) ChangeStatus(ctx context.Context, req *v1.DoctorStatusChangeRequest) (*MessageResponse, error) {
	err := s.service.ChangeStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "status changed"}, nil
}

func (s *DoctorGRPCServer) AssignRole(ctx context.Context, req *v1.DoctorRoleAssignRequest) (*MessageResponse, error) {
	err := s.service.AssignRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "role assigned"}, nil
}

func (s *DoctorGRPCServer) RemoveRole(ctx context.Context, req *v1.DoctorRoleRemoveRequest) (*MessageResponse, error) {
	err := s.service.RemoveRole(ctx, req)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "role removed"}, nil
}

func (s *DoctorGRPCServer) SetSchedule(ctx context.Context, req *v1.ScheduleSetRequest) (*v1.ScheduleDTO, error) {
	return s.service.SetSchedule(ctx, req)
}

func (s *DoctorGRPCServer) AddScheduleSlot(ctx context.Context, req *v1.ScheduleAddSlotRequest) (*v1.ScheduleDTO, error) {
	return s.service.AddScheduleSlot(ctx, req)
}

func (s *DoctorGRPCServer) GetSchedule(ctx context.Context, req *GetScheduleRequest) (*v1.ScheduleDTO, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.BadRequest("INVALID_DATE", "invalid date format")
	}
	return s.service.GetSchedule(ctx, req.DoctorId, date)
}

func (s *DoctorGRPCServer) GetSchedules(ctx context.Context, req *GetSchedulesRequest) (*v1.ScheduleListResponse, error) {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)
	if startDate.IsZero() {
		startDate = time.Now()
	}
	if endDate.IsZero() {
		endDate = startDate.AddDate(0, 0, 30)
	}
	return s.service.GetSchedules(ctx, req.DoctorId, startDate, endDate)
}

func (s *DoctorGRPCServer) PublishSchedule(ctx context.Context, req *v1.SchedulePublishRequest) (*v1.ScheduleDTO, error) {
	return s.service.PublishSchedule(ctx, req)
}

func (s *DoctorGRPCServer) CancelSchedule(ctx context.Context, req *v1.ScheduleCancelRequest) (*MessageResponse, error) {
	err := s.service.CancelSchedule(ctx, req)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "schedule cancelled"}, nil
}

func (s *DoctorGRPCServer) BookSlot(ctx context.Context, req *BookSlotRequest) (*MessageResponse, error) {
	err := s.service.BookSlot(ctx, req.ScheduleId, req.SlotId)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "slot booked"}, nil
}

func (s *DoctorGRPCServer) CancelBooking(ctx context.Context, req *CancelBookingRequest) (*MessageResponse, error) {
	err := s.service.CancelBooking(ctx, req.ScheduleId, req.SlotId)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "booking cancelled"}, nil
}

func (s *DoctorGRPCServer) GetAvailableSlots(ctx context.Context, req *GetAvailableSlotsRequest) (*v1.AvailableSlotsResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.BadRequest("INVALID_DATE", "invalid date format")
	}
	return s.service.GetAvailableSlots(ctx, req.DoctorId, date)
}

func (s *DoctorGRPCServer) GetPublishedSchedulesByDate(ctx context.Context, req *GetPublishedSchedulesRequest) (*ScheduleListResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.BadRequest("INVALID_DATE", "invalid date format")
	}
	result, err := s.service.GetPublishedSchedulesByDate(ctx, date)
	if err != nil {
		return nil, err
	}
	return &ScheduleListResponse{Items: result}, nil
}

func (s *DoctorGRPCServer) GetSchedulesByDepartment(ctx context.Context, req *GetSchedulesByDepartmentRequest) (*ScheduleListResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.BadRequest("INVALID_DATE", "invalid date format")
	}
	result, err := s.service.GetSchedulesByDepartment(ctx, req.DepartmentId, date)
	if err != nil {
		return nil, err
	}
	return &ScheduleListResponse{Items: result}, nil
}

func (s *DoctorGRPCServer) RemoveScheduleSlot(ctx context.Context, req *RemoveScheduleSlotRequest) (*MessageResponse, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.BadRequest("INVALID_DATE", "invalid date format")
	}
	err = s.service.RemoveScheduleSlot(ctx, req.DoctorId, date, req.SlotId)
	if err != nil {
		return nil, err
	}
	return &MessageResponse{Message: "slot removed"}, nil
}

var _DoctorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "doctor.DoctorService",
	HandlerType: (*DoctorGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateDoctor", Handler: _DoctorService_CreateDoctor_Handler},
		{MethodName: "GetDoctor", Handler: _DoctorService_GetDoctor_Handler},
		{MethodName: "UpdateDoctor", Handler: _DoctorService_UpdateDoctor_Handler},
		{MethodName: "DeleteDoctor", Handler: _DoctorService_DeleteDoctor_Handler},
		{MethodName: "ListDoctors", Handler: _DoctorService_ListDoctors_Handler},
		{MethodName: "ListDoctorsByDepartment", Handler: _DoctorService_ListDoctorsByDepartment_Handler},
		{MethodName: "ListExpertDoctors", Handler: _DoctorService_ListExpertDoctors_Handler},
		{MethodName: "ListDoctorsBySpecialty", Handler: _DoctorService_ListDoctorsBySpecialty_Handler},
		{MethodName: "GetDoctorByEmployeeID", Handler: _DoctorService_GetDoctorByEmployeeID_Handler},
		{MethodName: "ChangeDepartment", Handler: _DoctorService_ChangeDepartment_Handler},
		{MethodName: "ChangeStatus", Handler: _DoctorService_ChangeStatus_Handler},
		{MethodName: "AssignRole", Handler: _DoctorService_AssignRole_Handler},
		{MethodName: "RemoveRole", Handler: _DoctorService_RemoveRole_Handler},
		{MethodName: "SetSchedule", Handler: _DoctorService_SetSchedule_Handler},
		{MethodName: "AddScheduleSlot", Handler: _DoctorService_AddScheduleSlot_Handler},
		{MethodName: "GetSchedule", Handler: _DoctorService_GetSchedule_Handler},
		{MethodName: "GetSchedules", Handler: _DoctorService_GetSchedules_Handler},
		{MethodName: "PublishSchedule", Handler: _DoctorService_PublishSchedule_Handler},
		{MethodName: "CancelSchedule", Handler: _DoctorService_CancelSchedule_Handler},
		{MethodName: "BookSlot", Handler: _DoctorService_BookSlot_Handler},
		{MethodName: "CancelBooking", Handler: _DoctorService_CancelBooking_Handler},
		{MethodName: "GetAvailableSlots", Handler: _DoctorService_GetAvailableSlots_Handler},
		{MethodName: "GetPublishedSchedulesByDate", Handler: _DoctorService_GetPublishedSchedulesByDate_Handler},
		{MethodName: "GetSchedulesByDepartment", Handler: _DoctorService_GetSchedulesByDepartment_Handler},
		{MethodName: "RemoveScheduleSlot", Handler: _DoctorService_RemoveScheduleSlot_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "doctor.proto",
}

func _DoctorService_CreateDoctor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).CreateDoctor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/CreateDoctor"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).CreateDoctor(ctx, req.(*v1.DoctorCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetDoctor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDoctorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetDoctor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetDoctor"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetDoctor(ctx, req.(*GetDoctorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_UpdateDoctor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).UpdateDoctor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/UpdateDoctor"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).UpdateDoctor(ctx, req.(*v1.DoctorUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_DeleteDoctor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDoctorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).DeleteDoctor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/DeleteDoctor"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).DeleteDoctor(ctx, req.(*DeleteDoctorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ListDoctors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ListDoctors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ListDoctors"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ListDoctors(ctx, req.(*v1.DoctorFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ListDoctorsByDepartment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListByDepartmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ListDoctorsByDepartment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ListDoctorsByDepartment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ListDoctorsByDepartment(ctx, req.(*ListByDepartmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ListExpertDoctors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ListExpertDoctors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ListExpertDoctors"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ListExpertDoctors(ctx, req.(*EmptyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ListDoctorsBySpecialty_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListBySpecialtyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ListDoctorsBySpecialty(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ListDoctorsBySpecialty"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ListDoctorsBySpecialty(ctx, req.(*ListBySpecialtyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetDoctorByEmployeeID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetByEmployeeIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetDoctorByEmployeeID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetDoctorByEmployeeID"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetDoctorByEmployeeID(ctx, req.(*GetByEmployeeIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ChangeDepartment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorDepartmentChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ChangeDepartment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ChangeDepartment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ChangeDepartment(ctx, req.(*v1.DoctorDepartmentChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_ChangeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorStatusChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).ChangeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/ChangeStatus"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).ChangeStatus(ctx, req.(*v1.DoctorStatusChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_AssignRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorRoleAssignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).AssignRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/AssignRole"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).AssignRole(ctx, req.(*v1.DoctorRoleAssignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_RemoveRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.DoctorRoleRemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).RemoveRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/RemoveRole"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).RemoveRole(ctx, req.(*v1.DoctorRoleRemoveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_SetSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ScheduleSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).SetSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/SetSchedule"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).SetSchedule(ctx, req.(*v1.ScheduleSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_AddScheduleSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ScheduleAddSlotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).AddScheduleSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/AddScheduleSlot"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).AddScheduleSlot(ctx, req.(*v1.ScheduleAddSlotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScheduleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetSchedule"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetSchedule(ctx, req.(*GetScheduleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetSchedules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchedulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetSchedules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetSchedules"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetSchedules(ctx, req.(*GetSchedulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_PublishSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.SchedulePublishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).PublishSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/PublishSchedule"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).PublishSchedule(ctx, req.(*v1.SchedulePublishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_CancelSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ScheduleCancelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).CancelSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/CancelSchedule"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).CancelSchedule(ctx, req.(*v1.ScheduleCancelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_BookSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BookSlotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).BookSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/BookSlot"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).BookSlot(ctx, req.(*BookSlotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_CancelBooking_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelBookingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).CancelBooking(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/CancelBooking"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).CancelBooking(ctx, req.(*CancelBookingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetAvailableSlots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAvailableSlotsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetAvailableSlots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetAvailableSlots"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetAvailableSlots(ctx, req.(*GetAvailableSlotsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetPublishedSchedulesByDate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublishedSchedulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetPublishedSchedulesByDate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetPublishedSchedulesByDate"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetPublishedSchedulesByDate(ctx, req.(*GetPublishedSchedulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_GetSchedulesByDepartment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSchedulesByDepartmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).GetSchedulesByDepartment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/GetSchedulesByDepartment"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).GetSchedulesByDepartment(ctx, req.(*GetSchedulesByDepartmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoctorService_RemoveScheduleSlot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveScheduleSlotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(*DoctorGRPCServer).RemoveScheduleSlot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/doctor.DoctorService/RemoveScheduleSlot"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(*DoctorGRPCServer).RemoveScheduleSlot(ctx, req.(*RemoveScheduleSlotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

type GetDoctorRequest struct {
	DoctorId string `json:"doctor_id"`
}

type DeleteDoctorRequest struct {
	DoctorId string `json:"doctor_id"`
}

type DeleteDoctorResponse struct {
	Message string `json:"message"`
}

type ListByDepartmentRequest struct {
	DepartmentId string `json:"department_id"`
}

type ListBySpecialtyRequest struct {
	Specialty string `json:"specialty"`
}

type GetByEmployeeIDRequest struct {
	EmployeeId string `json:"employee_id"`
}

type EmptyRequest struct{}

type MessageResponse struct {
	Message string `json:"message"`
}

type DoctorListResponse struct {
	Items []*v1.DoctorDTO `json:"items"`
}

type ScheduleListResponse struct {
	Items []*v1.ScheduleDTO `json:"items"`
}

type GetScheduleRequest struct {
	DoctorId string `json:"doctor_id"`
	Date     string `json:"date"`
}

type GetSchedulesRequest struct {
	DoctorId  string `json:"doctor_id"`
 StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type BookSlotRequest struct {
	ScheduleId string `json:"schedule_id"`
	SlotId     string `json:"slot_id"`
}

type CancelBookingRequest struct {
	ScheduleId string `json:"schedule_id"`
	SlotId     string `json:"slot_id"`
}

type GetAvailableSlotsRequest struct {
	DoctorId string `json:"doctor_id"`
	Date     string `json:"date"`
}

type GetPublishedSchedulesRequest struct {
	Date string `json:"date"`
}

type GetSchedulesByDepartmentRequest struct {
	DepartmentId string `json:"department_id"`
	Date         string `json:"date"`
}

type RemoveScheduleSlotRequest struct {
	DoctorId string `json:"doctor_id"`
	Date     string `json:"date"`
	SlotId   string `json:"slot_id"`
}