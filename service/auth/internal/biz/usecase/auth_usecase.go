package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"intelligent-guidance-system/service/auth/internal/domain/aggregate"
	"intelligent-guidance-system/service/auth/internal/domain/entity"
	"intelligent-guidance-system/service/auth/internal/domain/repository"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrTokenExpired       = errors.New("token expired")
	ErrPermissionDenied   = errors.New("permission denied")
)

type AuthUsecase struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	tokenRepo      repository.TokenRepository
	jwtSecret      string
	tokenExpiry    time.Duration
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	tokenRepo repository.TokenRepository,
	jwtSecret string,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		tokenRepo:      tokenRepo,
		jwtSecret:      jwtSecret,
		tokenExpiry:    24 * time.Hour,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, username, password string) (*aggregate.User, string, time.Time, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	if user == nil {
		return nil, "", time.Time{}, ErrUserNotFound
	}

	if !user.CanLogin() {
		return nil, "", time.Time{}, ErrUserInactive
	}

	if !user.ValidatePassword(password) {
		return nil, "", time.Time{}, ErrInvalidCredentials
	}

	token, expiresAt, err := u.generateToken(user)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if err := u.tokenRepo.SaveToken(ctx, user.ID(), token, expiresAt.Unix()); err != nil {
		return nil, "", time.Time{}, err
	}

	return user, token, expiresAt, nil
}

func (u *AuthUsecase) Logout(ctx context.Context, userID int64) error {
	return u.tokenRepo.DeleteToken(ctx, userID)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, oldToken string) (string, time.Time, error) {
	claims, err := u.parseToken(oldToken)
	if err != nil {
		return "", time.Time{}, err
	}

	userID := claims["user_id"].(int64)
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", time.Time{}, err
	}
	if user == nil {
		return "", time.Time{}, ErrUserNotFound
	}

	token, expiresAt, err := u.generateToken(user)
	if err != nil {
		return "", time.Time{}, err
	}

	if err := u.tokenRepo.SaveToken(ctx, user.ID(), token, expiresAt.Unix()); err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

func (u *AuthUsecase) CheckPermission(ctx context.Context, userID int64, resourceURL string) (bool, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, ErrUserNotFound
	}

	return user.HasResourceAccess(resourceURL), nil
}

func (u *AuthUsecase) GetUser(ctx context.Context, userID int64) (*aggregate.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (u *AuthUsecase) CreateUser(ctx context.Context,
	userType entity.UserType,
	username string,
	password string,
	realName string,
	phone string,
	email string,
) (*aggregate.User, error) {
	exists, err := u.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	user, err := aggregate.NewUser(userType, username, password, realName, phone, email)
	if err != nil {
		return nil, err
	}

	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) UpdatePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := user.UpdatePassword(oldPassword, newPassword); err != nil {
		return err
	}

	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) AssignRole(ctx context.Context, userID, roleID, operatorID int64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	role, err := u.roleRepo.FindByID(ctx, roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	if err := user.AssignRole(role, operatorID); err != nil {
		return err
	}

	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) RemoveRole(ctx context.Context, userID, roleID, operatorID int64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := user.RemoveRole(roleID, operatorID); err != nil {
		return err
	}

	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) ListUsers(ctx context.Context, page, pageSize int) ([]*aggregate.User, int64, error) {
	return u.userRepo.FindAll(ctx, page, pageSize)
}

func (u *AuthUsecase) ActivateUser(ctx context.Context, userID int64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.Activate(userID)
	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) LockUser(ctx context.Context, userID int64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.Lock(userID)
	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) UnlockUser(ctx context.Context, userID int64) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	user.Unlock(userID)
	return u.userRepo.Save(ctx, user)
}

func (u *AuthUsecase) generateToken(user *aggregate.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(u.tokenExpiry)

	claims := jwt.MapClaims{
		"user_id":   user.ID(),
		"user_type": user.UserType().Code(),
		"username":  user.Username(),
		"exp":       expiresAt.Unix(),
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (u *AuthUsecase) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (u *AuthUsecase) ValidateToken(ctx context.Context, tokenString string) (int64, error) {
	claims, err := u.parseToken(tokenString)
	if err != nil {
		return 0, err
	}

	userID := int64(claims["user_id"].(float64))
	
	exists, err := u.tokenRepo.ExistsToken(ctx, tokenString)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("token not found in storage")
	}

	return userID, nil
}