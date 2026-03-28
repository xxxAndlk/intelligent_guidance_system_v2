package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"intelligent-guidance-system/service/doctor/internal/domain/aggregate"
	"intelligent-guidance-system/service/doctor/internal/domain/entity"
	"intelligent-guidance-system/service/doctor/internal/domain/event"
	"intelligent-guidance-system/service/doctor/internal/domain/repository"
	"intelligent-guidance-system/service/doctor/internal/domain/vo"

	"gorm.io/gorm"
)

type DoctorRepoImpl struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) repository.DoctorRepository {
	return &DoctorRepoImpl{db: db}
}

func (r *DoctorRepoImpl) Save(ctx context.Context, doctor *aggregate.Doctor) error {
	po := r.toDoctorPO(doctor)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *DoctorRepoImpl) Update(ctx context.Context, doctor *aggregate.Doctor) error {
	po := r.toDoctorPO(doctor)
	return r.db.WithContext(ctx).Save(po).Error
}

func (r *DoctorRepoImpl) Delete(ctx context.Context, doctorID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("doctor_id = ?", doctorID).Delete(&DoctorRoleAssignmentPO{}).Error; err != nil {
			return err
		}
		if err := tx.Where("doctor_id = ?", doctorID).Delete(&SchedulePO{}).Error; err != nil {
			return err
		}
		if err := tx.Where("doctor_id = ?", doctorID).Delete(&TimeSlotPO{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", doctorID).Delete(&DoctorPO{}).Error
	})
}

func (r *DoctorRepoImpl) FindByID(ctx context.Context, doctorID string) (*aggregate.Doctor, error) {
	var po DoctorPO
	err := r.db.WithContext(ctx).Where("id = ?", doctorID).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDoctorDomain(&po, ctx)
}

func (r *DoctorRepoImpl) FindByEmployeeID(ctx context.Context, employeeID string) (*aggregate.Doctor, error) {
	var po DoctorPO
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDoctorDomain(&po, ctx)
}

func (r *DoctorRepoImpl) FindByDepartment(ctx context.Context, departmentID string) ([]*aggregate.Doctor, error) {
	var poList []DoctorPO
	err := r.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(&poList).Error
	if err != nil {
		return nil, err
	}
	return r.toDoctorDomains(poList, ctx)
}

func (r *DoctorRepoImpl) FindActive(ctx context.Context) ([]*aggregate.Doctor, error) {
	var poList []DoctorPO
	err := r.db.WithContext(ctx).Where("status = ?", "ACTIVE").Find(&poList).Error
	if err != nil {
		return nil, err
	}
	return r.toDoctorDomains(poList, ctx)
}

func (r *DoctorRepoImpl) FindExperts(ctx context.Context) ([]*aggregate.Doctor, error) {
	var poList []DoctorPO
	err := r.db.WithContext(ctx).Where("position IN ?", []string{"CHIEF", "ASSOCIATE_CHIEF"}).
		Where("status = ?", "ACTIVE").Find(&poList).Error
	if err != nil {
		return nil, err
	}
	return r.toDoctorDomains(poList, ctx)
}

func (r *DoctorRepoImpl) FindBySpecialty(ctx context.Context, specialty string) ([]*aggregate.Doctor, error) {
	var poList []DoctorPO
	err := r.db.WithContext(ctx).Where("specialties LIKE ?", "%"+specialty+"%").
		Where("status = ?", "ACTIVE").Find(&poList).Error
	if err != nil {
		return nil, err
	}
	return r.toDoctorDomains(poList, ctx)
}

func (r *DoctorRepoImpl) FindByPosition(ctx context.Context, positionCode string) ([]*aggregate.Doctor, error) {
	var poList []DoctorPO
	err := r.db.WithContext(ctx).Where("position = ?", positionCode).Find(&poList).Error
	if err != nil {
		return nil, err
	}
	return r.toDoctorDomains(poList, ctx)
}

func (r *DoctorRepoImpl) List(ctx context.Context, filters *repository.DoctorFilter, page, pageSize int) ([]*aggregate.Doctor, int64, error) {
	query := r.db.WithContext(ctx).Model(&DoctorPO{})

	if filters.DepartmentID != "" {
		query = query.Where("department_id = ?", filters.DepartmentID)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.Position != "" {
		query = query.Where("position = ?", filters.Position)
	}
	if filters.Specialty != "" {
		query = query.Where("specialties LIKE ?", "%"+filters.Specialty+"%")
	}
	if filters.Name != "" {
		query = query.Where("name LIKE ?", "%"+filters.Name+"%")
	}
	if filters.EmployeeID != "" {
		query = query.Where("employee_id = ?", filters.EmployeeID)
	}
	if filters.IsExpert {
		query = query.Where("position IN ?", []string{"CHIEF", "ASSOCIATE_CHIEF"})
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var poList []DoctorPO
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&poList).Error; err != nil {
		return nil, 0, err
	}

	doctors, err := r.toDoctorDomains(poList, ctx)
	if err != nil {
		return nil, 0, err
	}

	return doctors, total, nil
}

func (r *DoctorRepoImpl) ExistsByEmployeeID(ctx context.Context, employeeID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&DoctorPO{}).Where("employee_id = ?", employeeID).Count(&count).Error
	return count > 0, err
}

func (r *DoctorRepoImpl) ExistsByLicenseNumber(ctx context.Context, licenseNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&DoctorPO{}).Where("license_number = ?", licenseNumber).Count(&count).Error
	return count > 0, err
}

func (r *DoctorRepoImpl) toDoctorPO(d *aggregate.Doctor) *DoctorPO {
	specialties, _ := json.Marshal(d.Professional().Specialties())
	certifications, _ := json.Marshal(d.Professional().Certifications())

	return &DoctorPO{
		ID:              d.ID(),
		EmployeeID:      d.EmployeeID(),
		Name:            d.BasicInfo().Name(),
		Gender:          d.BasicInfo().Gender(),
		BirthDate:       d.BasicInfo().BirthDate(),
		Phone:           d.BasicInfo().Phone(),
		Email:           d.BasicInfo().Email(),
		IDCard:          d.BasicInfo().IDCard(),
		Address:         d.BasicInfo().Address(),
		AvatarURL:       d.BasicInfo().AvatarURL(),
		LicenseNumber:   d.Professional().LicenseNumber(),
		Position:        d.Professional().Position().Code(),
		Title:           d.Professional().Title().Code(),
		Specialties:     string(specialties),
		Education:       d.Professional().Education().Code(),
		GraduationSchool: d.Professional().GraduationSchool(),
		GraduationYear:  d.Professional().GraduationYear(),
		YearsOfExp:      d.Professional().YearsOfExperience(),
		Certifications:  string(certifications),
		Introduction:    d.Professional().Introduction(),
		DepartmentID:    d.DepartmentID(),
		Status:          d.Status().Code(),
		CreatedAt:       d.CreatedAt(),
		UpdatedAt:       d.UpdatedAt(),
		Version:         d.Version(),
	}
}

func (r *DoctorRepoImpl) toDoctorDomain(po *DoctorPO, ctx context.Context) (*aggregate.Doctor, error) {
	basicInfo, err := vo.NewDoctorBasicInfo(
		po.Name, po.Gender, po.BirthDate,
		po.Phone, po.Email, po.IDCard,
		po.Address, po.AvatarURL,
	)
	if err != nil {
		return nil, err
	}

	var specialties []string
	if po.Specialties != "" {
		json.Unmarshal([]byte(po.Specialties), &specialties)
	}

	var certifications []string
	if po.Certifications != "" {
		json.Unmarshal([]byte(po.Certifications), &certifications)
	}

	professional, err := vo.NewProfessionalInfo(
		po.LicenseNumber,
		vo.PositionFromCode(po.Position),
		vo.TitleFromCode(po.Title),
		specialties,
		vo.EducationLevelFromCode(po.Education),
		po.GraduationSchool,
		po.GraduationYear,
		po.YearsOfExp,
		certifications,
		po.Introduction,
	)
	if err != nil {
		return nil, err
	}

	schedules, err := r.loadSchedules(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	roles, err := r.loadRoleAssignments(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	return aggregate.ReconstructDoctor(
		po.ID, po.EmployeeID,
		basicInfo, professional,
		po.DepartmentID,
		aggregate.DoctorStatusFromCode(po.Status),
		schedules, roles,
		po.CreatedAt, po.UpdatedAt,
		po.Version,
	), nil
}

func (r *DoctorRepoImpl) toDoctorDomains(poList []DoctorPO, ctx context.Context) ([]*aggregate.Doctor, error) {
	doctors := make([]*aggregate.Doctor, 0, len(poList))
	for _, po := range poList {
		d, err := r.toDoctorDomain(&po, ctx)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, d)
	}
	return doctors, nil
}

func (r *DoctorRepoImpl) loadSchedules(ctx context.Context, doctorID string) ([]*entity.Schedule, error) {
	var schedulePOs []SchedulePO
	err := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&schedulePOs).Error
	if err != nil {
		return nil, err
	}

	schedules := make([]*entity.Schedule, 0, len(schedulePOs))
	for _, spo := range schedulePOs {
		var slotPOs []TimeSlotPO
		err := r.db.WithContext(ctx).Where("schedule_id = ?", spo.ID).Find(&slotPOs).Error
		if err != nil {
			return nil, err
		}

		slots := make([]*entity.TimeSlot, 0, len(slotPOs))
		for _, tpo := range slotPOs {
			slots = append(slots, entity.ReconstructTimeSlot(
				tpo.ID, tpo.StartTime, tpo.EndTime,
				tpo.IsAvailable, tpo.MaxPatients, tpo.BookedCount,
			))
		}

		var status entity.ScheduleStatus
		switch spo.Status {
		case "DRAFT":
			status = entity.ScheduleStatusDraft
		case "PUBLISHED":
			status = entity.ScheduleStatusPublished
		case "CANCELLED":
			status = entity.ScheduleStatusCancelled
		}

		schedules = append(schedules, entity.ReconstructSchedule(
			spo.ID, spo.DoctorID, spo.Date,
			slots, status, spo.CreatedAt, spo.UpdatedAt,
		))
	}

	return schedules, nil
}

func (r *DoctorRepoImpl) loadRoleAssignments(ctx context.Context, doctorID string) ([]*entity.DoctorRoleAssignment, error) {
	var assignmentPOs []DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&assignmentPOs).Error
	if err != nil {
		return nil, err
	}

	assignments := make([]*entity.DoctorRoleAssignment, 0, len(assignmentPOs))
	for _, apo := range assignmentPOs {
		assignments = append(assignments, entity.ReconstructDoctorRoleAssignment(
			apo.ID, apo.DoctorID, apo.RoleID, apo.DepartmentID, apo.AssignedBy,
			apo.AssignedAt, apo.ExpiresAt, apo.IsActive,
		))
	}

	return assignments, nil
}

type ScheduleRepoImpl struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) repository.ScheduleRepository {
	return &ScheduleRepoImpl{db: db}
}

func (r *ScheduleRepoImpl) Save(ctx context.Context, schedule *entity.Schedule) error {
	spo := &SchedulePO{
		ID:       schedule.ID(),
		DoctorID: schedule.DoctorID(),
		Date:     schedule.Date(),
		Status:   r.statusToCode(schedule.Status()),
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(spo).Error; err != nil {
			return err
		}

		for _, slot := range schedule.TimeSlots() {
			tpo := &TimeSlotPO{
				ID:          slot.ID(),
				ScheduleID:  schedule.ID(),
				StartTime:   slot.StartTime(),
				EndTime:     slot.EndTime(),
				IsAvailable: slot.IsAvailable(),
				MaxPatients: slot.MaxPatients(),
				BookedCount: slot.BookedCount(),
			}
			if err := tx.Create(tpo).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ScheduleRepoImpl) Update(ctx context.Context, schedule *entity.Schedule) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		spo := &SchedulePO{
			ID:       schedule.ID(),
			DoctorID: schedule.DoctorID(),
			Date:     schedule.Date(),
			Status:   r.statusToCode(schedule.Status()),
			UpdatedAt: schedule.UpdatedAt(),
		}
		if err := tx.Save(spo).Error; err != nil {
			return err
		}

		for _, slot := range schedule.TimeSlots() {
			tpo := &TimeSlotPO{
				ID:          slot.ID(),
				ScheduleID:  schedule.ID(),
				StartTime:   slot.StartTime(),
				EndTime:     slot.EndTime(),
				IsAvailable: slot.IsAvailable(),
				MaxPatients: slot.MaxPatients(),
				BookedCount: slot.BookedCount(),
				UpdatedAt:   time.Now(),
			}
			if err := tx.Save(tpo).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ScheduleRepoImpl) Delete(ctx context.Context, scheduleID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("schedule_id = ?", scheduleID).Delete(&TimeSlotPO{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", scheduleID).Delete(&SchedulePO{}).Error
	})
}

func (r *ScheduleRepoImpl) FindByID(ctx context.Context, scheduleID string) (*entity.Schedule, error) {
	var spo SchedulePO
	err := r.db.WithContext(ctx).Where("id = ?", scheduleID).First(&spo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toScheduleDomain(&spo, ctx)
}

func (r *ScheduleRepoImpl) FindByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) (*entity.Schedule, error) {
	var spo SchedulePO
	err := r.db.WithContext(ctx).Where("doctor_id = ? AND date = ?", doctorID, date.Truncate(24*time.Hour)).First(&spo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toScheduleDomain(&spo, ctx)
}

func (r *ScheduleRepoImpl) FindByDoctor(ctx context.Context, doctorID string, startDate, endDate time.Time) ([]*entity.Schedule, error) {
	var spoList []SchedulePO
	err := r.db.WithContext(ctx).Where("doctor_id = ? AND date >= ? AND date <= ?", doctorID, startDate, endDate).Find(&spoList).Error
	if err != nil {
		return nil, err
	}
	return r.toScheduleDomains(spoList, ctx)
}

func (r *ScheduleRepoImpl) FindPublishedByDate(ctx context.Context, date time.Time) ([]*entity.Schedule, error) {
	var spoList []SchedulePO
	err := r.db.WithContext(ctx).Where("date = ? AND status = ?", date.Truncate(24*time.Hour), "PUBLISHED").Find(&spoList).Error
	if err != nil {
		return nil, err
	}
	return r.toScheduleDomains(spoList, ctx)
}

func (r *ScheduleRepoImpl) FindByDepartmentAndDate(ctx context.Context, departmentID string, date time.Time) ([]*entity.Schedule, error) {
	var spoList []SchedulePO
	err := r.db.WithContext(ctx).
		Joins("JOIN doctors ON doctors.id = doctor_schedules.doctor_id").
		Where("doctors.department_id = ? AND doctor_schedules.date = ? AND doctor_schedules.status = ?", departmentID, date.Truncate(24*time.Hour), "PUBLISHED").
		Find(&spoList).Error
	if err != nil {
		return nil, err
	}
	return r.toScheduleDomains(spoList, ctx)
}

func (r *ScheduleRepoImpl) BookSlot(ctx context.Context, scheduleID, slotID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tpo TimeSlotPO
		err := tx.Where("id = ? AND schedule_id = ?", slotID, scheduleID).First(&tpo).Error
		if err != nil {
			return err
		}

		if tpo.BookedCount >= tpo.MaxPatients {
			return entity.ErrTimeSlotOverlap
		}

		tpo.BookedCount++
		tpo.IsAvailable = tpo.BookedCount < tpo.MaxPatients
		return tx.Save(&tpo).Error
	})
}

func (r *ScheduleRepoImpl) CancelSlotBooking(ctx context.Context, scheduleID, slotID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tpo TimeSlotPO
		err := tx.Where("id = ? AND schedule_id = ?", slotID, scheduleID).First(&tpo).Error
		if err != nil {
			return err
		}

		if tpo.BookedCount > 0 {
			tpo.BookedCount--
		}
		tpo.IsAvailable = true
		return tx.Save(&tpo).Error
	})
}

func (r *ScheduleRepoImpl) ListAvailableSlots(ctx context.Context, doctorID string, date time.Time) ([]*entity.TimeSlot, error) {
	schedule, err := r.FindByDoctorAndDate(ctx, doctorID, date)
	if err != nil || schedule == nil {
		return nil, err
	}

	slots := make([]*entity.TimeSlot, 0)
	for _, slot := range schedule.TimeSlots() {
		if slot.IsAvailable() {
			slots = append(slots, slot)
		}
	}
	return slots, nil
}

func (r *ScheduleRepoImpl) statusToCode(status entity.ScheduleStatus) string {
	switch status {
	case entity.ScheduleStatusDraft:
		return "DRAFT"
	case entity.ScheduleStatusPublished:
		return "PUBLISHED"
	case entity.ScheduleStatusCancelled:
		return "CANCELLED"
	default:
		return "DRAFT"
	}
}

func (r *ScheduleRepoImpl) toScheduleDomain(spo *SchedulePO, ctx context.Context) (*entity.Schedule, error) {
	var slotPOs []TimeSlotPO
	err := r.db.WithContext(ctx).Where("schedule_id = ?", spo.ID).Find(&slotPOs).Error
	if err != nil {
		return nil, err
	}

	slots := make([]*entity.TimeSlot, 0, len(slotPOs))
	for _, tpo := range slotPOs {
		slots = append(slots, entity.ReconstructTimeSlot(
			tpo.ID, tpo.StartTime, tpo.EndTime,
			tpo.IsAvailable, tpo.MaxPatients, tpo.BookedCount,
		))
	}

	var status entity.ScheduleStatus
	switch spo.Status {
	case "DRAFT":
		status = entity.ScheduleStatusDraft
	case "PUBLISHED":
		status = entity.ScheduleStatusPublished
	case "CANCELLED":
		status = entity.ScheduleStatusCancelled
	}

	return entity.ReconstructSchedule(
		spo.ID, spo.DoctorID, spo.Date,
		slots, status, spo.CreatedAt, spo.UpdatedAt,
	), nil
}

func (r *ScheduleRepoImpl) toScheduleDomains(spoList []SchedulePO, ctx context.Context) ([]*entity.Schedule, error) {
	schedules := make([]*entity.Schedule, 0, len(spoList))
	for _, spo := range spoList {
		s, err := r.toScheduleDomain(&spo, ctx)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

type RoleRepoImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &RoleRepoImpl{db: db}
}

func (r *RoleRepoImpl) Save(ctx context.Context, role *entity.Role) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		rpo := &RolePO{
			ID:          role.ID(),
			Name:        role.Name(),
			Code:        role.Code(),
			Description: role.Description(),
			RoleType:    r.roleTypeToCode(role.RoleType()),
		}
		if err := tx.Create(rpo).Error; err != nil {
			return err
		}

		for _, perm := range role.Permissions() {
			ppo := &RolePermissionPO{
				ID:       perm.Code(),
				RoleID:   role.ID(),
				PermCode: perm.Code(),
				PermName: perm.Name(),
				Resource: perm.Resource(),
				Action:   perm.Action(),
			}
			if err := tx.Create(ppo).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RoleRepoImpl) Update(ctx context.Context, role *entity.Role) error {
	rpo := &RolePO{
		ID:          role.ID(),
		Name:        role.Name(),
		Code:        role.Code(),
		Description: role.Description(),
		RoleType:    r.roleTypeToCode(role.RoleType()),
	}
	return r.db.WithContext(ctx).Save(rpo).Error
}

func (r *RoleRepoImpl) Delete(ctx context.Context, roleID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&RolePermissionPO{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", roleID).Delete(&RolePO{}).Error
	})
}

func (r *RoleRepoImpl) FindByID(ctx context.Context, roleID string) (*entity.Role, error) {
	var rpo RolePO
	err := r.db.WithContext(ctx).Where("id = ?", roleID).First(&rpo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toRoleDomain(&rpo, ctx)
}

func (r *RoleRepoImpl) FindByCode(ctx context.Context, code string) (*entity.Role, error) {
	var rpo RolePO
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&rpo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toRoleDomain(&rpo, ctx)
}

func (r *RoleRepoImpl) FindByType(ctx context.Context, roleType entity.RoleType) ([]*entity.Role, error) {
	var rpoList []RolePO
	err := r.db.WithContext(ctx).Where("role_type = ?", r.roleTypeToCode(roleType)).Find(&rpoList).Error
	if err != nil {
		return nil, err
	}
	return r.toRoleDomains(rpoList, ctx)
}

func (r *RoleRepoImpl) FindAll(ctx context.Context) ([]*entity.Role, error) {
	var rpoList []RolePO
	err := r.db.WithContext(ctx).Find(&rpoList).Error
	if err != nil {
		return nil, err
	}
	return r.toRoleDomains(rpoList, ctx)
}

func (r *RoleRepoImpl) List(ctx context.Context, page, pageSize int) ([]*entity.Role, int64, error) {
	var total int64
	if err := r.db.WithContext(ctx).Model(&RolePO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rpoList []RolePO
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&rpoList).Error
	if err != nil {
		return nil, 0, err
	}

	roles, err := r.toRoleDomains(rpoList, ctx)
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *RoleRepoImpl) roleTypeToCode(roleType entity.RoleType) string {
	switch roleType {
	case entity.RoleTypeSystem:
		return "SYSTEM"
	case entity.RoleTypeDepartment:
		return "DEPARTMENT"
	case entity.RoleTypeCustom:
		return "CUSTOM"
	default:
		return "SYSTEM"
	}
}

func (r *RoleRepoImpl) toRoleDomain(rpo *RolePO, ctx context.Context) (*entity.Role, error) {
	var permPOs []RolePermissionPO
	err := r.db.WithContext(ctx).Where("role_id = ?", rpo.ID).Find(&permPOs).Error
	if err != nil {
		return nil, err
	}

	perms := make([]*entity.Permission, 0, len(permPOs))
	for _, ppo := range permPOs {
		perm, err := entity.NewPermission(ppo.PermCode, ppo.PermName, "", ppo.Resource, ppo.Action)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	var roleType entity.RoleType
	switch rpo.RoleType {
	case "SYSTEM":
		roleType = entity.RoleTypeSystem
	case "DEPARTMENT":
		roleType = entity.RoleTypeDepartment
	case "CUSTOM":
		roleType = entity.RoleTypeCustom
	}

	return entity.ReconstructRole(
		rpo.ID, rpo.Name, rpo.Code, rpo.Description,
		roleType, perms, rpo.CreatedAt, rpo.UpdatedAt,
	), nil
}

func (r *RoleRepoImpl) toRoleDomains(rpoList []RolePO, ctx context.Context) ([]*entity.Role, error) {
	roles := make([]*entity.Role, 0, len(rpoList))
	for _, rpo := range rpoList {
		role, err := r.toRoleDomain(&rpo, ctx)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

type DoctorEventRepoImpl struct {
	db *gorm.DB
}

func NewDoctorEventRepository(db *gorm.DB) repository.DoctorEventRepository {
	return &DoctorEventRepoImpl{db: db}
}

func (r *DoctorEventRepoImpl) Save(ctx context.Context, evt interface{}) error {
	var epo *DoctorEventPO

	switch v := evt.(type) {
	case *event.DoctorEvent:
		payload, _ := json.Marshal(v.Payload())
		epo = &DoctorEventPO{
			ID:           v.EventID(),
			EventType:    string(v.EventType()),
			DoctorID:     v.DoctorID(),
			EmployeeID:   v.EmployeeID(),
			DepartmentID: v.DepartmentID(),
			Payload:      string(payload),
			OccurredAt:   v.OccurredAt(),
		}
	case *event.DoctorCreatedEvent:
		payload, _ := json.Marshal(v.Payload())
		epo = &DoctorEventPO{
			ID:           v.EventID(),
			EventType:    string(v.EventType()),
			DoctorID:     v.DoctorID(),
			EmployeeID:   v.EmployeeID(),
			DepartmentID: v.DepartmentID,
			Payload:      string(payload),
			OccurredAt:   v.OccurredAt(),
		}
	default:
		return errors.New("unsupported event type")
	}

	return r.db.WithContext(ctx).Create(epo).Error
}

func (r *DoctorEventRepoImpl) FindByDoctor(ctx context.Context, doctorID string, limit int) ([]interface{}, error) {
	var epoList []DoctorEventPO
	err := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Order("occurred_at DESC").Limit(limit).Find(&epoList).Error
	if err != nil {
		return nil, err
	}

	events := make([]interface{}, 0, len(epoList))
	for _, epo := range epoList {
		events = append(events, r.toEventDomain(&epo))
	}
	return events, nil
}

func (r *DoctorEventRepoImpl) FindByType(ctx context.Context, eventType string, limit int) ([]interface{}, error) {
	var epoList []DoctorEventPO
	err := r.db.WithContext(ctx).Where("event_type = ?", eventType).Order("occurred_at DESC").Limit(limit).Find(&epoList).Error
	if err != nil {
		return nil, err
	}

	events := make([]interface{}, 0, len(epoList))
	for _, epo := range epoList {
		events = append(events, r.toEventDomain(&epo))
	}
	return events, nil
}

func (r *DoctorEventRepoImpl) FindRecent(ctx context.Context, limit int) ([]interface{}, error) {
	var epoList []DoctorEventPO
	err := r.db.WithContext(ctx).Order("occurred_at DESC").Limit(limit).Find(&epoList).Error
	if err != nil {
		return nil, err
	}

	events := make([]interface{}, 0, len(epoList))
	for _, epo := range epoList {
		events = append(events, r.toEventDomain(&epo))
	}
	return events, nil
}

func (r *DoctorEventRepoImpl) toEventDomain(epo *DoctorEventPO) *event.DoctorEvent {
	var payload map[string]interface{}
	if epo.Payload != "" {
		json.Unmarshal([]byte(epo.Payload), &payload)
	}

	return event.NewDoctorEvent(
		event.DoctorEventType(epo.EventType),
		epo.DoctorID, epo.EmployeeID, payload,
	)
}

type DoctorRoleAssignmentRepoImpl struct {
	db *gorm.DB
}

func NewDoctorRoleAssignmentRepository(db *gorm.DB) repository.DoctorRoleAssignmentRepository {
	return &DoctorRoleAssignmentRepoImpl{db: db}
}

func (r *DoctorRoleAssignmentRepoImpl) Save(ctx context.Context, assignment *entity.DoctorRoleAssignment) error {
	apo := &DoctorRoleAssignmentPO{
		ID:           assignment.ID(),
		DoctorID:     assignment.DoctorID(),
		RoleID:       assignment.RoleID(),
		DepartmentID: assignment.DepartmentID(),
		AssignedBy:   assignment.AssignedBy(),
		AssignedAt:   assignment.AssignedAt(),
		ExpiresAt:    assignment.ExpiresAt(),
		IsActive:     assignment.IsActive(),
	}
	return r.db.WithContext(ctx).Create(apo).Error
}

func (r *DoctorRoleAssignmentRepoImpl) Update(ctx context.Context, assignment *entity.DoctorRoleAssignment) error {
	apo := &DoctorRoleAssignmentPO{
		ID:           assignment.ID(),
		DoctorID:     assignment.DoctorID(),
		RoleID:       assignment.RoleID(),
		DepartmentID: assignment.DepartmentID(),
		AssignedBy:   assignment.AssignedBy(),
		AssignedAt:   assignment.AssignedAt(),
		ExpiresAt:    assignment.ExpiresAt(),
		IsActive:     assignment.IsActive(),
	}
	return r.db.WithContext(ctx).Save(apo).Error
}

func (r *DoctorRoleAssignmentRepoImpl) Delete(ctx context.Context, assignmentID string) error {
	return r.db.WithContext(ctx).Where("id = ?", assignmentID).Delete(&DoctorRoleAssignmentPO{}).Error
}

func (r *DoctorRoleAssignmentRepoImpl) FindByID(ctx context.Context, assignmentID string) (*entity.DoctorRoleAssignment, error) {
	var apo DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("id = ?", assignmentID).First(&apo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&apo), nil
}

func (r *DoctorRoleAssignmentRepoImpl) FindByDoctor(ctx context.Context, doctorID string) ([]*entity.DoctorRoleAssignment, error) {
	var apoList []DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("doctor_id = ?", doctorID).Find(&apoList).Error
	if err != nil {
		return nil, err
	}
	return r.toDomains(apoList), nil
}

func (r *DoctorRoleAssignmentRepoImpl) FindActiveByDoctor(ctx context.Context, doctorID string) ([]*entity.DoctorRoleAssignment, error) {
	var apoList []DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("doctor_id = ? AND is_active = ?", doctorID, true).Find(&apoList).Error
	if err != nil {
		return nil, err
	}
	return r.toDomains(apoList), nil
}

func (r *DoctorRoleAssignmentRepoImpl) FindByRole(ctx context.Context, roleID string) ([]*entity.DoctorRoleAssignment, error) {
	var apoList []DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&apoList).Error
	if err != nil {
		return nil, err
	}
	return r.toDomains(apoList), nil
}

func (r *DoctorRoleAssignmentRepoImpl) FindByDoctorAndRole(ctx context.Context, doctorID, roleID string) (*entity.DoctorRoleAssignment, error) {
	var apo DoctorRoleAssignmentPO
	err := r.db.WithContext(ctx).Where("doctor_id = ? AND role_id = ?", doctorID, roleID).First(&apo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(&apo), nil
}

func (r *DoctorRoleAssignmentRepoImpl) Deactivate(ctx context.Context, assignmentID string) error {
	return r.db.WithContext(ctx).Model(&DoctorRoleAssignmentPO{}).Where("id = ?", assignmentID).Update("is_active", false).Error
}

func (r *DoctorRoleAssignmentRepoImpl) toDomain(apo *DoctorRoleAssignmentPO) *entity.DoctorRoleAssignment {
	return entity.ReconstructDoctorRoleAssignment(
		apo.ID, apo.DoctorID, apo.RoleID, apo.DepartmentID, apo.AssignedBy,
		apo.AssignedAt, apo.ExpiresAt, apo.IsActive,
	)
}

func (r *DoctorRoleAssignmentRepoImpl) toDomains(apoList []DoctorRoleAssignmentPO) []*entity.DoctorRoleAssignment {
	result := make([]*entity.DoctorRoleAssignment, 0, len(apoList))
	for _, apo := range apoList {
		result = append(result, r.toDomain(&apo))
	}
	return result
}