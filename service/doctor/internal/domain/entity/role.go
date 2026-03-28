package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidRoleName   = errors.New("invalid role name")
	ErrInvalidPermission = errors.New("invalid permission")
	ErrRoleNotFound      = errors.New("role not found")
)

// RoleType represents the type of role
type RoleType int

const (
	RoleTypeSystem RoleType = iota
	RoleTypeDepartment
	RoleTypeCustom
)

func (r RoleType) String() string {
	switch r {
	case RoleTypeSystem:
		return "系统角色"
	case RoleTypeDepartment:
		return "科室角色"
	case RoleTypeCustom:
		return "自定义角色"
	default:
		return "未知"
	}
}

// Permission represents a specific permission
type Permission struct {
	code        string
	name        string
	description string
	resource    string
	action      string
}

// NewPermission creates a new permission
func NewPermission(code, name, description, resource, action string) (*Permission, error) {
	if code == "" || name == "" {
		return nil, ErrInvalidPermission
	}
	return &Permission{
		code:        code,
		name:        name,
		description: description,
		resource:    resource,
		action:      action,
	}, nil
}

func (p *Permission) Code() string        { return p.code }
func (p *Permission) Name() string        { return p.name }
func (p *Permission) Description() string { return p.description }
func (p *Permission) Resource() string    { return p.resource }
func (p *Permission) Action() string       { return p.action }

// Role represents a doctor's role with permissions
type Role struct {
	id          string
	name        string
	code        string
	description string
	roleType    RoleType
	permissions []*Permission
	createdAt   time.Time
	updatedAt   time.Time
}

// NewRole creates a new role
func NewRole(name, code, description string, roleType RoleType) (*Role, error) {
	if name == "" || code == "" {
		return nil, ErrInvalidRoleName
	}

	return &Role{
		id:          uuid.New().String(),
		name:        name,
		code:        code,
		description: description,
		roleType:    roleType,
		permissions: make([]*Permission, 0),
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}, nil
}

// ReconstructRole reconstructs a role from persistence
func ReconstructRole(
	id, name, code, description string,
	roleType RoleType,
	permissions []*Permission,
	createdAt, updatedAt time.Time,
) *Role {
	return &Role{
		id:          id,
		name:        name,
		code:        code,
		description: description,
		roleType:    roleType,
		permissions: permissions,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (r *Role) ID() string                   { return r.id }
func (r *Role) Name() string                  { return r.name }
func (r *Role) Code() string                  { return r.code }
func (r *Role) Description() string           { return r.description }
func (r *Role) RoleType() RoleType            { return r.roleType }
func (r *Role) Permissions() []*Permission     { return r.permissions }
func (r *Role) CreatedAt() time.Time          { return r.createdAt }
func (r *Role) UpdatedAt() time.Time          { return r.updatedAt }

// AddPermission adds a permission to the role
func (r *Role) AddPermission(permission *Permission) {
	for _, p := range r.permissions {
		if p.Code() == permission.Code() {
			return
		}
	}
	r.permissions = append(r.permissions, permission)
	r.updatedAt = time.Now()
}

// RemovePermission removes a permission by code
func (r *Role) RemovePermission(code string) {
	for i, p := range r.permissions {
		if p.Code() == code {
			r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
			r.updatedAt = time.Now()
			return
		}
	}
}

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(code string) bool {
	for _, p := range r.permissions {
		if p.Code() == code {
			return true
		}
	}
	return false
}

// HasResourceAction checks if the role has permission for a resource action
func (r *Role) HasResourceAction(resource, action string) bool {
	for _, p := range r.permissions {
		if p.Resource() == resource && p.Action() == action {
			return true
		}
	}
	return false
}

// UpdateDescription updates the role description
func (r *Role) UpdateDescription(description string) {
	r.description = description
	r.updatedAt = time.Now()
}

// DoctorRoleAssignment represents a role assigned to a doctor
type DoctorRoleAssignment struct {
	id         string
	doctorID   string
	roleID     string
	departmentID string
	assignedBy string
	assignedAt time.Time
	expiresAt  *time.Time
	isActive   bool
}

// NewDoctorRoleAssignment creates a new role assignment
func NewDoctorRoleAssignment(doctorID, roleID, departmentID, assignedBy string, expiresAt *time.Time) *DoctorRoleAssignment {
	return &DoctorRoleAssignment{
		id:           uuid.New().String(),
		doctorID:     doctorID,
		roleID:       roleID,
		departmentID: departmentID,
		assignedBy:   assignedBy,
		assignedAt:   time.Now(),
		expiresAt:    expiresAt,
		isActive:     true,
	}
}

// ReconstructDoctorRoleAssignment reconstructs an assignment from persistence
func ReconstructDoctorRoleAssignment(
	id, doctorID, roleID, departmentID, assignedBy string,
	assignedAt time.Time,
	expiresAt *time.Time,
	isActive bool,
) *DoctorRoleAssignment {
	return &DoctorRoleAssignment{
		id:           id,
		doctorID:     doctorID,
		roleID:       roleID,
		departmentID: departmentID,
		assignedBy:   assignedBy,
		assignedAt:   assignedAt,
		expiresAt:    expiresAt,
		isActive:     isActive,
	}
}

func (d *DoctorRoleAssignment) ID() string            { return d.id }
func (d *DoctorRoleAssignment) DoctorID() string       { return d.doctorID }
func (d *DoctorRoleAssignment) RoleID() string         { return d.roleID }
func (d *DoctorRoleAssignment) DepartmentID() string   { return d.departmentID }
func (d *DoctorRoleAssignment) AssignedBy() string     { return d.assignedBy }
func (d *DoctorRoleAssignment) AssignedAt() time.Time { return d.assignedAt }
func (d *DoctorRoleAssignment) ExpiresAt() *time.Time  { return d.expiresAt }
func (d *DoctorRoleAssignment) IsActive() bool          { return d.isActive }

// Deactivate deactivates the role assignment
func (d *DoctorRoleAssignment) Deactivate() {
	d.isActive = false
}

// Activate activates the role assignment
func (d *DoctorRoleAssignment) Activate() {
	d.isActive = true
}

// IsExpired checks if the role assignment has expired
func (d *DoctorRoleAssignment) IsExpired() bool {
	if d.expiresAt == nil {
		return false
	}
	return time.Now().After(*d.expiresAt)
}

// IsValid checks if the role assignment is valid (active and not expired)
func (d *DoctorRoleAssignment) IsValid() bool {
	return d.isActive && !d.IsExpired()
}