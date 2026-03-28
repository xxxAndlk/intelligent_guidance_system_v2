package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/department/internal/domain/aggregate"
	"intelligent-guidance-system/service/department/internal/domain/entity"
	"intelligent-guidance-system/service/department/internal/domain/repository"
	"intelligent-guidance-system/service/department/internal/domain/vo"
)

type DepartmentRepoImpl struct {
	db *gorm.DB
}

func NewDepartmentRepoImpl(db *gorm.DB) repository.DepartmentRepository {
	return &DepartmentRepoImpl{db: db}
}

func (r *DepartmentRepoImpl) Save(ctx context.Context, department *aggregate.Department) error {
	po := r.toPO(department)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if po.ID == 0 {
			if err := tx.Create(&po).Error; err != nil {
				return err
			}
			department.SetID(po.ID)
		} else {
			if err := tx.Save(&po).Error; err != nil {
				return err
			}
		}

		for _, assignment := range department.DoctorAssignments() {
			assignmentPO := DoctorAssignmentPO{
				DepartmentID: department.ID(),
				DoctorID:     assignment.DoctorID(),
				AssignedAt:   assignment.AssignedAt(),
				IsDuty:       assignment.IsDuty(),
				IsActive:     assignment.IsActive(),
			}
			if assignment.IsActive() {
				if err := tx.Where("department_id = ? AND doctor_id = ?", department.ID(), assignment.DoctorID()).
					Assign(assignmentPO).FirstOrCreate(&assignmentPO).Error; err != nil {
					return err
				}
			} else {
				tx.Where("department_id = ? AND doctor_id = ?", department.ID(), assignment.DoctorID()).
					Update("is_active", false)
			}
		}

		return nil
	})
}

func (r *DepartmentRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.Department, error) {
	var po DepartmentPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var assignments []DoctorAssignmentPO
	if err := r.db.WithContext(ctx).Where("department_id = ? AND is_active = ?", id, true).Find(&assignments).Error; err != nil {
		return nil, err
	}

	return r.toAggregate(&po, assignments), nil
}

func (r *DepartmentRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Department, int64, error) {
	var pos []DepartmentPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&DepartmentPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	departments := make([]*aggregate.Department, 0, len(pos))
	for _, po := range pos {
		var assignments []DoctorAssignmentPO
		if err := r.db.WithContext(ctx).Where("department_id = ? AND is_active = ?", po.ID, true).Find(&assignments).Error; err != nil {
			return nil, 0, err
		}
		departments = append(departments, r.toAggregate(&po, assignments))
	}

	return departments, total, nil
}

func (r *DepartmentRepoImpl) FindByType(ctx context.Context, deptType int) ([]*aggregate.Department, error) {
	var pos []DepartmentPO
	typeCode := entity.DepartmentType(deptType).Code()
	if err := r.db.WithContext(ctx).Where("type = ?", typeCode).Find(&pos).Error; err != nil {
		return nil, err
	}

	departments := make([]*aggregate.Department, 0, len(pos))
	for _, po := range pos {
		var assignments []DoctorAssignmentPO
		if err := r.db.WithContext(ctx).Where("department_id = ? AND is_active = ?", po.ID, true).Find(&assignments).Error; err != nil {
			return nil, err
		}
		departments = append(departments, r.toAggregate(&po, assignments))
	}

	return departments, nil
}

func (r *DepartmentRepoImpl) FindByStatus(ctx context.Context, status int) ([]*aggregate.Department, error) {
	var pos []DepartmentPO
	statusCode := entity.DepartmentStatus(status).Code()
	if err := r.db.WithContext(ctx).Where("status = ?", statusCode).Find(&pos).Error; err != nil {
		return nil, err
	}

	departments := make([]*aggregate.Department, 0, len(pos))
	for _, po := range pos {
		var assignments []DoctorAssignmentPO
		if err := r.db.WithContext(ctx).Where("department_id = ? AND is_active = ?", po.ID, true).Find(&assignments).Error; err != nil {
			return nil, err
		}
		departments = append(departments, r.toAggregate(&po, assignments))
	}

	return departments, nil
}

func (r *DepartmentRepoImpl) FindActiveDepartments(ctx context.Context) ([]*aggregate.Department, error) {
	return r.FindByStatus(ctx, int(entity.DepartmentStatusActive))
}

func (r *DepartmentRepoImpl) FindByDoctorID(ctx context.Context, doctorID int64) ([]*aggregate.Department, error) {
	var assignmentPOs []DoctorAssignmentPO
	if err := r.db.WithContext(ctx).Where("doctor_id = ? AND is_active = ?", doctorID, true).Find(&assignmentPOs).Error; err != nil {
		return nil, err
	}

	departmentIDs := make([]int64, 0, len(assignmentPOs))
	for _, a := range assignmentPOs {
		departmentIDs = append(departmentIDs, a.DepartmentID)
	}

	if len(departmentIDs) == 0 {
		return []*aggregate.Department{}, nil
	}

	var pos []DepartmentPO
	if err := r.db.WithContext(ctx).Where("id IN ?", departmentIDs).Find(&pos).Error; err != nil {
		return nil, err
	}

	departments := make([]*aggregate.Department, 0, len(pos))
	for _, po := range pos {
		departments = append(departments, r.toAggregate(&po, nil))
	}

	return departments, nil
}

func (r *DepartmentRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("department_id = ?", id).Delete(&DoctorAssignmentPO{}).Error; err != nil {
			return err
		}
		return tx.Delete(&DepartmentPO{}, id).Error
	})
}

func (r *DepartmentRepoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&DepartmentPO{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DepartmentRepoImpl) toPO(dept *aggregate.Department) *DepartmentPO {
	return &DepartmentPO{
		ID:                  dept.ID(),
		Name:                dept.Name(),
		Introduction:        dept.Introduction(),
		PersonInChargeID:    dept.PersonInChargeID(),
		Phone:               dept.Phone(),
		Address:             dept.Address(),
		Type:                dept.Type().Code(),
		StaffCount:          dept.StaffCount(),
		Status:              dept.Status().Code(),
		RegistrationFee:     dept.Fee().RegistrationFee(),
		ExpertRegistrationFee: dept.Fee().ExpertRegistrationFee(),
		DutyDoctorID:        dept.DutyDoctorID(),
		CreatedAt:           dept.CreatedAt(),
		UpdatedAt:           dept.UpdatedAt(),
	}
}

func (r *DepartmentRepoImpl) toAggregate(po *DepartmentPO, assignments []DoctorAssignmentPO) *aggregate.Department {
	fee, _ := vo.NewDepartmentFee(po.RegistrationFee, po.ExpertRegistrationFee, "CNY")

	doctorAssignments := make([]*entity.DoctorAssignment, 0, len(assignments))
	for _, a := range assignments {
		doctorAssignments = append(doctorAssignments, entity.ReconstructDoctorAssignment(
			a.DoctorID,
			a.AssignedAt,
			a.IsDuty,
			a.IsActive,
		))
	}

	return aggregate.ReconstructDepartment(
		po.ID,
		po.Name,
		po.Introduction,
		po.PersonInChargeID,
		po.Phone,
		po.Address,
		entity.DepartmentTypeFromCode(po.Type),
		po.StaffCount,
		entity.DepartmentStatusFromCode(po.Status),
		fee,
		po.DutyDoctorID,
		doctorAssignments,
		po.CreatedAt,
		po.UpdatedAt,
	)
}