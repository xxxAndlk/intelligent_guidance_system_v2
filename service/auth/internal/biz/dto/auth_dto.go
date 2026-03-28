package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      *UserResponse `json:"user"`
}

type RegisterRequest struct {
	UserType string `json:"user_type" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateProfileRequest struct {
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type AssignRoleRequest struct {
	RoleID     int64 `json:"role_id" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type RemoveRoleRequest struct {
	RoleID     int64 `json:"role_id" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type CheckPermissionRequest struct {
	ResourceURL string `json:"resource_url" binding:"required"`
}

type UserResponse struct {
	ID        int64    `json:"id"`
	UserType  string   `json:"user_type"`
	Username  string   `json:"username"`
	RealName  string   `json:"real_name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Avatar    string   `json:"avatar"`
	Status    string   `json:"status"`
	Roles     []RoleResponse `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleResponse struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	DataScope   string           `json:"data_scope"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	ID           int64  `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	ResourceType string `json:"resource_type"`
	ResourceURL  string `json:"resource_url"`
}

type UserListResponse struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Users    []UserResponse `json:"users"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type RefreshTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}