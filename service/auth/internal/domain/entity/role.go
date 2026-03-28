package entity

import "time"

type Role struct {
	id          int64
	code        string
	name        string
	dataScope   DataScope
	permissions []*Permission
	createdAt   time.Time
	updatedAt   time.Time
}

func NewRole(code, name string, dataScope DataScope) (*Role, error) {
	if code == "" {
		return nil, ErrInvalidRole
	}
	if name == "" {
		return nil, ErrInvalidRole
	}

	now := time.Now()
	return &Role{
		id:          0,
		code:        code,
		name:        name,
		dataScope:   dataScope,
		permissions: make([]*Permission, 0),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructRole(id int64, code, name string, dataScope DataScope, permissions []*Permission, createdAt, updatedAt time.Time) *Role {
	return &Role{
		id:          id,
		code:        code,
		name:        name,
		dataScope:   dataScope,
		permissions: permissions,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (r *Role) ID() int64 { return r.id }
func (r *Role) Code() string { return r.code }
func (r *Role) Name() string { return r.name }
func (r *Role) DataScope() DataScope { return r.dataScope }
func (r *Role) Permissions() []*Permission { return r.permissions }
func (r *Role) CreatedAt() time.Time { return r.createdAt }
func (r *Role) UpdatedAt() time.Time { return r.updatedAt }

func (r *Role) SetID(id int64) { r.id = id }

func (r *Role) AddPermission(permission *Permission) error {
	for _, p := range r.permissions {
		if p.Code() == permission.Code() {
			return ErrPermissionDenied
		}
	}
	r.permissions = append(r.permissions, permission)
	r.updatedAt = time.Now()
	return nil
}

func (r *Role) RemovePermission(permissionCode string) {
	for i, p := range r.permissions {
		if p.Code() == permissionCode {
			r.permissions = append(r.permissions[:i], r.permissions[i+1:]...)
			r.updatedAt = time.Now()
			return
		}
	}
}

func (r *Role) HasPermission(permissionCode string) bool {
	for _, p := range r.permissions {
		if p.Code() == permissionCode {
			return true
		}
	}
	return false
}

func (r *Role) HasResourceAccess(resourceURL string) bool {
	for _, p := range r.permissions {
		if p.ResourceURL() == resourceURL || p.MatchResource(resourceURL) {
			return true
		}
	}
	return false
}