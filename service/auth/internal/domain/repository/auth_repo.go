package repository

import (
	"context"

	"intelligent-guidance-system/service/auth/internal/domain/aggregate"
	"intelligent-guidance-system/service/auth/internal/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *aggregate.User) error
	FindByID(ctx context.Context, id int64) (*aggregate.User, error)
	FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
	FindByPhone(ctx context.Context, phone string) (*aggregate.User, error)
	FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.User, int64, error)
	FindByType(ctx context.Context, userType entity.UserType) ([]*aggregate.User, error)
	FindByStatus(ctx context.Context, status entity.UserStatus) ([]*aggregate.User, error)
	Delete(ctx context.Context, id int64) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type RoleRepository interface {
	Save(ctx context.Context, role *entity.Role) error
	FindByID(ctx context.Context, id int64) (*entity.Role, error)
	FindByCode(ctx context.Context, code string) (*entity.Role, error)
	FindAll(ctx context.Context) ([]*entity.Role, error)
	Delete(ctx context.Context, id int64) error
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

type PermissionRepository interface {
	Save(ctx context.Context, permission *entity.Permission) error
	FindByID(ctx context.Context, id int64) (*entity.Permission, error)
	FindByCode(ctx context.Context, code string) (*entity.Permission, error)
	FindAll(ctx context.Context) ([]*entity.Permission, error)
	FindByResourceType(ctx context.Context, resourceType entity.ResourceType) ([]*entity.Permission, error)
	Delete(ctx context.Context, id int64) error
}

type DataPermissionRepository interface {
	Save(ctx context.Context, dp *entity.DataPermission) error
	FindByUserID(ctx context.Context, userID int64) ([]*entity.DataPermission, error)
	FindByUserIDAndResource(ctx context.Context, userID int64, resourceType string) (*entity.DataPermission, error)
	Delete(ctx context.Context, userID int64, resourceType string) error
}

type TokenRepository interface {
	SaveToken(ctx context.Context, userID int64, token string, expiresAt int64) error
	FindToken(ctx context.Context, userID int64) (string, int64, error)
	DeleteToken(ctx context.Context, userID int64) error
	ExistsToken(ctx context.Context, token string) (bool, error)
}