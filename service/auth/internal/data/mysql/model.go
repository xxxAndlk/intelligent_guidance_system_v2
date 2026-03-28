package mysql

import "time"

type UserPO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	UserType  string    `gorm:"column:user_type;size:20;not null"`
	Username  string    `gorm:"column:username;size:50;unique;not null"`
	Password  string    `gorm:"column:password;size:100;not null"`
	RealName  string    `gorm:"column:real_name;size:50"`
	Phone     string    `gorm:"column:phone;size:20;unique"`
	Email     string    `gorm:"column:email;size:50;unique"`
	Avatar    string    `gorm:"column:avatar;size:200"`
	Status    string    `gorm:"column:status;size:20;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserPO) TableName() string {
	return "users"
}

type RolePO struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Code      string    `gorm:"column:code;size:50;unique;not null"`
	Name      string    `gorm:"column:name;size:50;not null"`
	DataScope string    `gorm:"column:data_scope;size:20;not null"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (RolePO) TableName() string {
	return "roles"
}

type PermissionPO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Code         string    `gorm:"column:code;size:50;unique;not null"`
	Name         string    `gorm:"column:name;size:50;not null"`
	ResourceType string    `gorm:"column:resource_type;size:20;not null"`
	ResourceURL  string    `gorm:"column:resource_url;size:200;not null"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (PermissionPO) TableName() string {
	return "permissions"
}

type UserRolePO struct {
	ID     int64 `gorm:"primaryKey;autoIncrement"`
	UserID int64 `gorm:"column:user_id;index"`
	RoleID int64 `gorm:"column:role_id;index"`
}

func (UserRolePO) TableName() string {
	return "user_roles"
}

type RolePermissionPO struct {
	ID           int64 `gorm:"primaryKey;autoIncrement"`
	RoleID       int64 `gorm:"column:role_id;index"`
	PermissionID int64 `gorm:"column:permission_id;index"`
}

func (RolePermissionPO) TableName() string {
	return "role_permissions"
}

type DataPermissionPO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	UserID       int64     `gorm:"column:user_id;index"`
	ResourceType string    `gorm:"column:resource_type;size:50"`
	DataScope    string    `gorm:"column:data_scope;size:20"`
	DeptIDs      string    `gorm:"column:dept_ids;size:500"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (DataPermissionPO) TableName() string {
	return "data_permissions"
}