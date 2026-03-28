package repository

import (
	"context"

	"intelligent-guidance-system/service/department/internal/domain/aggregate"
)

// DepartmentRepository defines the repository interface for department
type DepartmentRepository interface {
	// Save saves a department aggregate
	Save(ctx context.Context, department *aggregate.Department) error
	
	// FindByID finds a department by ID
	FindByID(ctx context.Context, id int64) (*aggregate.Department, error)
	
	// FindAll finds all departments with pagination
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.Department, int64, error)
	
	// FindByType finds departments by type
	FindByType(ctx context.Context, deptType int) ([]*aggregate.Department, error)
	
	// FindByStatus finds departments by status
	FindByStatus(ctx context.Context, status int) ([]*aggregate.Department, error)
	
	// FindActiveDepartments finds all active departments
	FindActiveDepartments(ctx context.Context) ([]*aggregate.Department, error)
	
	// FindByDoctorID finds departments that a doctor belongs to
	FindByDoctorID(ctx context.Context, doctorID int64) ([]*aggregate.Department, error)
	
	// Delete deletes a department by ID
	Delete(ctx context.Context, id int64) error
	
	// ExistsByName checks if a department with the name exists
	ExistsByName(ctx context.Context, name string) (bool, error)
}