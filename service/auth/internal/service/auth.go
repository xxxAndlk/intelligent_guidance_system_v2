package service

import (
	"context"

	"intelligent-guidance-system/service/auth/internal/biz/dto"
	"intelligent-guidance-system/service/auth/internal/biz/usecase"
	"intelligent-guidance-system/service/auth/internal/domain/aggregate"
	"intelligent-guidance-system/service/auth/internal/domain/entity"
)

type AuthService struct {
	usecase *usecase.AuthUsecase
}

func NewAuthService(usecase *usecase.AuthUsecase) *AuthService {
	return &AuthService{usecase: usecase}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, token, expiresAt, err := s.usecase.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      s.toUserResponse(user),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	return s.usecase.Logout(ctx, userID)
}

func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	token, expiresAt, err := s.usecase.RefreshToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) CheckPermission(ctx context.Context, userID int64, req *dto.CheckPermissionRequest) (bool, error) {
	return s.usecase.CheckPermission(ctx, userID, req.ResourceURL)
}

func (s *AuthService) GetUser(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.usecase.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	userType := entity.UserTypeFromCode(req.UserType)

	user, err := s.usecase.CreateUser(ctx, userType, req.Username, req.Password, req.RealName, req.Phone, req.Email)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *AuthService) UpdatePassword(ctx context.Context, userID int64, req *dto.UpdatePasswordRequest) error {
	return s.usecase.UpdatePassword(ctx, userID, req.OldPassword, req.NewPassword)
}

func (s *AuthService) AssignRole(ctx context.Context, userID int64, req *dto.AssignRoleRequest) error {
	return s.usecase.AssignRole(ctx, userID, req.RoleID, req.OperatorID)
}

func (s *AuthService) RemoveRole(ctx context.Context, userID int64, req *dto.RemoveRoleRequest) error {
	return s.usecase.RemoveRole(ctx, userID, req.RoleID, req.OperatorID)
}

func (s *AuthService) ListUsers(ctx context.Context, page, pageSize int) (*dto.UserListResponse, error) {
	users, total, err := s.usecase.ListUsers(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, *s.toUserResponse(user))
	}

	return &dto.UserListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Users:    userResponses,
	}, nil
}

func (s *AuthService) ActivateUser(ctx context.Context, userID int64) error {
	return s.usecase.ActivateUser(ctx, userID)
}

func (s *AuthService) LockUser(ctx context.Context, userID int64) error {
	return s.usecase.LockUser(ctx, userID)
}

func (s *AuthService) UnlockUser(ctx context.Context, userID int64) error {
	return s.usecase.UnlockUser(ctx, userID)
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (int64, error) {
	return s.usecase.ValidateToken(ctx, token)
}

func (s *AuthService) toUserResponse(user *aggregate.User) *dto.UserResponse {
	roles := make([]dto.RoleResponse, 0, len(user.Roles()))
	for _, r := range user.Roles() {
		permissions := make([]dto.PermissionResponse, 0, len(r.Permissions()))
		for _, p := range r.Permissions() {
			permissions = append(permissions, dto.PermissionResponse{
				ID:           p.ID(),
				Code:         p.Code(),
				Name:         p.Name(),
				ResourceType: p.ResourceType().Code(),
				ResourceURL:  p.ResourceURL(),
			})
		}
		roles = append(roles, dto.RoleResponse{
			ID:          r.ID(),
			Code:        r.Code(),
			Name:        r.Name(),
			DataScope:   r.DataScope().Code(),
			Permissions: permissions,
		})
	}

	return &dto.UserResponse{
		ID:        user.ID(),
		UserType:  user.UserType().Code(),
		Username:  user.Username(),
		RealName:  user.RealName(),
		Phone:     user.Phone(),
		Email:     user.Email(),
		Avatar:    user.Avatar(),
		Status:    user.Status().Code(),
		Roles:     roles,
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}