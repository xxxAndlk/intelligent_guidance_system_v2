package server

import (
	"context"
	"time"

	v1 "intelligent-guidance-system/service/doctor/internal/biz/dto"
	"intelligent-guidance-system/service/doctor/internal/service"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewHTTPServer(doctorSvc *service.DoctorService, logger log.Logger) *http.Server {
	opts := []http.ServerOption{
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
	}

	srv := http.NewServer(opts...)

	registerHTTPRoutes(srv, doctorSvc)

	return srv
}

func registerHTTPRoutes(srv *http.Server, svc *service.DoctorService) {
	srv.Route("/").POST("/doctors", _HTTP_CreateDoctor_Handler(svc))
	srv.Route("/").GET("/doctors/{id}", _HTTP_GetDoctor_Handler(svc))
	srv.Route("/").PUT("/doctors/{id}", _HTTP_UpdateDoctor_Handler(svc))
	srv.Route("/").DELETE("/doctors/{id}", _HTTP_DeleteDoctor_Handler(svc))
	srv.Route("/").GET("/doctors", _HTTP_ListDoctors_Handler(svc))
	srv.Route("/").GET("/doctors/department/{dept_id}", _HTTP_ListByDept_Handler(svc))
	srv.Route("/").GET("/doctors/experts", _HTTP_ListExperts_Handler(svc))
	srv.Route("/").GET("/doctors/specialty/{specialty}", _HTTP_ListBySpecialty_Handler(svc))
	srv.Route("/").GET("/doctors/employee/{employee_id}", _HTTP_GetByEmployeeID_Handler(svc))
	srv.Route("/").PUT("/doctors/{id}/department", _HTTP_ChangeDepartment_Handler(svc))
	srv.Route("/").PUT("/doctors/{id}/status", _HTTP_ChangeStatus_Handler(svc))
	srv.Route("/").POST("/doctors/{id}/roles", _HTTP_AssignRole_Handler(svc))
	srv.Route("/").DELETE("/doctors/{id}/roles/{role_id}", _HTTP_RemoveRole_Handler(svc))

	srv.Route("/").POST("/schedules", _HTTP_SetSchedule_Handler(svc))
	srv.Route("/").POST("/schedules/slots", _HTTP_AddScheduleSlot_Handler(svc))
	srv.Route("/").GET("/schedules/{doctor_id}/{date}", _HTTP_GetSchedule_Handler(svc))
	srv.Route("/").GET("/schedules/{doctor_id}", _HTTP_GetSchedules_Handler(svc))
	srv.Route("/").PUT("/schedules/publish", _HTTP_PublishSchedule_Handler(svc))
	srv.Route("/").PUT("/schedules/cancel", _HTTP_CancelSchedule_Handler(svc))
	srv.Route("/").POST("/schedules/book", _HTTP_BookSlot_Handler(svc))
	srv.Route("/").DELETE("/schedules/book", _HTTP_CancelBooking_Handler(svc))
	srv.Route("/").GET("/schedules/available/{doctor_id}/{date}", _HTTP_GetAvailableSlots_Handler(svc))
	srv.Route("/").GET("/schedules/published/{date}", _HTTP_GetPublishedSchedules_Handler(svc))
	srv.Route("/").GET("/schedules/department/{dept_id}/{date}", _HTTP_GetSchedulesByDept_Handler(svc))
	srv.Route("/").DELETE("/schedules/slots/{doctor_id}/{date}/{slot_id}", _HTTP_RemoveScheduleSlot_Handler(svc))
}

func _HTTP_CreateDoctor_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.DoctorCreateRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		result, err := svc.CreateDoctor(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetDoctor_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		result, err := svc.GetDoctor(ctx.Context(), id)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_UpdateDoctor_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		var req v1.DoctorUpdateRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		req.ID = id
		result, err := svc.UpdateDoctor(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_DeleteDoctor_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		err := svc.DeleteDoctor(ctx.Context(), id)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "deleted"})
	}
}

func _HTTP_ListDoctors_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.DoctorFilterRequest
		if err := ctx.BindQuery(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		result, err := svc.ListDoctors(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_ListByDept_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		deptID := ctx.Vars()["dept_id"]
		result, err := svc.ListDoctorsByDepartment(ctx.Context(), deptID)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_ListExperts_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		result, err := svc.ListExpertDoctors(ctx.Context())
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_ListBySpecialty_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		specialty := ctx.Vars()["specialty"]
		result, err := svc.ListDoctorsBySpecialty(ctx.Context(), specialty)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetByEmployeeID_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		employeeID := ctx.Vars()["employee_id"]
		result, err := svc.GetDoctorByEmployeeID(ctx.Context(), employeeID)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_ChangeDepartment_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		var req v1.DoctorDepartmentChangeRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		req.DoctorID = id
		err := svc.ChangeDepartment(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "department changed"})
	}
}

func _HTTP_ChangeStatus_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		var req v1.DoctorStatusChangeRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		req.DoctorID = id
		err := svc.ChangeStatus(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "status changed"})
	}
}

func _HTTP_AssignRole_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		var req v1.DoctorRoleAssignRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		req.DoctorID = id
		err := svc.AssignRole(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "role assigned"})
	}
}

func _HTTP_RemoveRole_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		id := ctx.Vars()["id"]
		roleID := ctx.Vars()["role_id"]
		req := &v1.DoctorRoleRemoveRequest{
			DoctorID: id,
			RoleID:   roleID,
		}
		err := svc.RemoveRole(ctx.Context(), req)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "role removed"})
	}
}

func _HTTP_SetSchedule_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.ScheduleSetRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		result, err := svc.SetSchedule(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_AddScheduleSlot_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.ScheduleAddSlotRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		result, err := svc.AddScheduleSlot(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetSchedule_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		doctorID := ctx.Vars()["doctor_id"]
		dateStr := ctx.Vars()["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return errors.BadRequest("INVALID_DATE", "invalid date format")
		}
		result, err := svc.GetSchedule(ctx.Context(), doctorID, date)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetSchedules_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		doctorID := ctx.Vars()["doctor_id"]
		startDateStr := ctx.Query().Get("start_date")
		endDateStr := ctx.Query().Get("end_date")
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		if startDate.IsZero() {
			startDate = time.Now()
		}
		if endDate.IsZero() {
			endDate = startDate.AddDate(0, 0, 30)
		}
		result, err := svc.GetSchedules(ctx.Context(), doctorID, startDate, endDate)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_PublishSchedule_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.SchedulePublishRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		result, err := svc.PublishSchedule(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_CancelSchedule_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		var req v1.ScheduleCancelRequest
		if err := ctx.BindVars(&req); err != nil {
			return errors.BadRequest("BIND_ERROR", err.Error())
		}
		err := svc.CancelSchedule(ctx.Context(), &req)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "schedule cancelled"})
	}
}

func _HTTP_BookSlot_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		scheduleID := ctx.Query().Get("schedule_id")
		slotID := ctx.Query().Get("slot_id")
		err := svc.BookSlot(ctx.Context(), scheduleID, slotID)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "slot booked"})
	}
}

func _HTTP_CancelBooking_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		scheduleID := ctx.Query().Get("schedule_id")
		slotID := ctx.Query().Get("slot_id")
		err := svc.CancelBooking(ctx.Context(), scheduleID, slotID)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "booking cancelled"})
	}
}

func _HTTP_GetAvailableSlots_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		doctorID := ctx.Vars()["doctor_id"]
		dateStr := ctx.Vars()["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return errors.BadRequest("INVALID_DATE", "invalid date format")
		}
		result, err := svc.GetAvailableSlots(ctx.Context(), doctorID, date)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetPublishedSchedules_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		dateStr := ctx.Vars()["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return errors.BadRequest("INVALID_DATE", "invalid date format")
		}
		result, err := svc.GetPublishedSchedulesByDate(ctx.Context(), date)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_GetSchedulesByDept_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		deptID := ctx.Vars()["dept_id"]
		dateStr := ctx.Vars()["date"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return errors.BadRequest("INVALID_DATE", "invalid date format")
		}
		result, err := svc.GetSchedulesByDepartment(ctx.Context(), deptID, date)
		if err != nil {
			return err
		}
		return ctx.Result(200, result)
	}
}

func _HTTP_RemoveScheduleSlot_Handler(svc *service.DoctorService) http.HandlerFunc {
	return func(ctx http.Context) error {
		doctorID := ctx.Vars()["doctor_id"]
		dateStr := ctx.Vars()["date"]
		slotID := ctx.Vars()["slot_id"]
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return errors.BadRequest("INVALID_DATE", "invalid date format")
		}
		err = svc.RemoveScheduleSlot(ctx.Context(), doctorID, date, slotID)
		if err != nil {
			return err
		}
		return ctx.Result(200, map[string]string{"message": "slot removed"})
	}
}