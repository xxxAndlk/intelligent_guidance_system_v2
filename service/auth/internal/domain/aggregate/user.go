package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/auth/internal/domain/entity"
	"intelligent-guidance-system/service/auth/internal/domain/event"
	"intelligent-guidance-system/service/auth/internal/domain/vo"
)

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserInactive        = errors.New("user is inactive")
	ErrUserLocked          = errors.New("user is locked")
	ErrUserDisabled        = errors.New("user is disabled")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrPasswordTooShort    = errors.New("password too short")
	ErrUsernameTooShort    = errors.New("username too short")
	ErrSamePassword        = errors.New("new password cannot be same as old")
	ErrRoleAlreadyAssigned = errors.New("role already assigned")
	ErrRoleNotAssigned     = errors.New("role not assigned")
)

type User struct {
	id          int64
	userType    entity.UserType
	username    string
	password    string
	realName    string
	phone       string
	email       string
	avatar      string
	status      entity.UserStatus
	roles       []*entity.Role
	dataPermissions []*entity.DataPermission
	events      []*event.AuthEvent
	createdAt   time.Time
	updatedAt   time.Time
}

func NewUser(
	userType entity.UserType,
	username string,
	password string,
	realName string,
	phone string,
	email string,
) (*User, error) {
	if userType == entity.UserTypeUnknown {
		return nil, entity.ErrInvalidUserType
	}
	if username == "" || len(username) < 3 {
		return nil, ErrUsernameTooShort
	}
	if password == "" || len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	now := time.Now()
	return &User{
		id:          0,
		userType:    userType,
		username:    username,
		password:    password,
		realName:    realName,
		phone:       phone,
		email:       email,
		avatar:      "",
		status:      entity.UserStatusInactive,
		roles:       make([]*entity.Role, 0),
		dataPermissions: make([]*entity.DataPermission, 0),
		events:      make([]*event.AuthEvent, 0),
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func ReconstructUser(
	id int64,
	userType entity.UserType,
	username string,
	password string,
	realName string,
	phone string,
	email string,
	avatar string,
	status entity.UserStatus,
	roles []*entity.Role,
	dataPermissions []*entity.DataPermission,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:          id,
		userType:    userType,
		username:    username,
		password:    password,
		realName:    realName,
		phone:       phone,
		email:       email,
		avatar:      avatar,
		status:      status,
		roles:       roles,
		dataPermissions: dataPermissions,
		events:      make([]*event.AuthEvent, 0),
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (u *User) ID() int64 { return u.id }
func (u *User) UserType() entity.UserType { return u.userType }
func (u *User) Username() string { return u.username }
func (u *User) Password() string { return u.password }
func (u *User) RealName() string { return u.realName }
func (u *User) Phone() string { return u.phone }
func (u *User) Email() string { return u.email }
func (u *User) Avatar() string { return u.avatar }
func (u *User) Status() entity.UserStatus { return u.status }
func (u *User) Roles() []*entity.Role { return u.roles }
func (u *User) DataPermissions() []*entity.DataPermission { return u.dataPermissions }
func (u *User) Events() []*event.AuthEvent { return u.events }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func (u *User) SetID(id int64) { u.id = id }

func (u *User) CanLogin() bool {
	return u.status.CanLogin()
}

func (u *User) ValidatePassword(password string) bool {
	return u.password == password
}

func (u *User) UpdatePassword(oldPassword, newPassword string) error {
	if !u.ValidatePassword(oldPassword) {
		return ErrInvalidCredentials
	}
	if newPassword == oldPassword {
		return ErrSamePassword
	}
	if len(newPassword) < 6 {
		return ErrPasswordTooShort
	}

	u.password = newPassword
	u.updatedAt = time.Now()
	u.events = append(u.events, event.PasswordChangedEvent(u.id))
	return nil
}

func (u *User) Activate(operatorID int64) error {
	u.status = entity.UserStatusActive
	u.updatedAt = time.Now()
	return nil
}

func (u *User) Lock(operatorID int64) error {
	u.status = entity.UserStatusLocked
	u.updatedAt = time.Now()
	return nil
}

func (u *User) Unlock(operatorID int64) error {
	u.status = entity.UserStatusActive
	u.updatedAt = time.Now()
	return nil
}

func (u *User) Disable(operatorID int64) error {
	u.status = entity.UserStatusDisabled
	u.updatedAt = time.Now()
	return nil
}

func (u *User) AssignRole(role *entity.Role, operatorID int64) error {
	for _, r := range u.roles {
		if r.ID() == role.ID() {
			return ErrRoleAlreadyAssigned
		}
	}
	u.roles = append(u.roles, role)
	u.updatedAt = time.Now()
	u.events = append(u.events, event.RoleAssignedEvent(u.id, operatorID, role.ID()))
	return nil
}

func (u *User) RemoveRole(roleID int64, operatorID int64) error {
	for i, r := range u.roles {
		if r.ID() == roleID {
			u.roles = append(u.roles[:i], u.roles[i+1:]...)
			u.updatedAt = time.Now()
			u.events = append(u.events, event.RoleRemovedEvent(u.id, operatorID, roleID))
			return nil
		}
	}
	return ErrRoleNotAssigned
}

func (u *User) HasRole(roleCode string) bool {
	for _, r := range u.roles {
		if r.Code() == roleCode {
			return true
		}
	}
	return false
}

func (u *User) HasPermission(permissionCode string) bool {
	for _, r := range u.roles {
		if r.HasPermission(permissionCode) {
			return true
		}
	}
	return false
}

func (u *User) HasResourceAccess(resourceURL string) bool {
	for _, r := range u.roles {
		if r.HasResourceAccess(resourceURL) {
			return true
		}
	}
	return false
}

func (u *User) UpdateRealName(realName string) {
	u.realName = realName
	u.updatedAt = time.Now()
}

func (u *User) UpdateContactInfo(contact *vo.ContactInfo) {
	u.phone = contact.Phone()
	u.email = contact.Email()
	u.updatedAt = time.Now()
}

func (u *User) UpdateAvatar(avatar string) {
	u.avatar = avatar
	u.updatedAt = time.Now()
}

func (u *User) ClearEvents() {
	u.events = make([]*event.AuthEvent, 0)
}

func (u *User) GetDataScope() entity.DataScope {
	for _, r := range u.roles {
		if r.DataScope() != entity.DataScopeUnknown {
			return r.DataScope()
		}
	}
	return entity.DataScopeOwn
}

func (u *User) IsAdmin() bool {
	return u.userType == entity.UserTypeAdmin
}

func (u *User) IsDoctor() bool {
	return u.userType == entity.UserTypeDoctor
}

func (u *User) IsPatient() bool {
	return u.userType == entity.UserTypePatient
}