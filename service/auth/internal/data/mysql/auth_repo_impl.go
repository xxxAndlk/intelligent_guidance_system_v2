package mysql

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/auth/internal/domain/aggregate"
	"intelligent-guidance-system/service/auth/internal/domain/entity"
	"intelligent-guidance-system/service/auth/internal/domain/repository"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) repository.UserRepository {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) Save(ctx context.Context, user *aggregate.User) error {
	po := r.userToPO(user)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if po.ID == 0 {
			if err := tx.Create(&po).Error; err != nil {
				return err
			}
			user.SetID(po.ID)
		} else {
			if err := tx.Save(&po).Error; err != nil {
				return err
			}
		}

		tx.Where("user_id = ?", user.ID()).Delete(&UserRolePO{})

		for _, role := range user.Roles() {
			userRole := UserRolePO{UserID: user.ID(), RoleID: role.ID()}
			if err := tx.Create(&userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *UserRepoImpl) FindByID(ctx context.Context, id int64) (*aggregate.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	roles, err := r.findUserRoles(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	return r.poToUser(&po, roles), nil
}

func (r *UserRepoImpl) FindByUsername(ctx context.Context, username string) (*aggregate.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	roles, err := r.findUserRoles(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	return r.poToUser(&po, roles), nil
}

func (r *UserRepoImpl) FindByPhone(ctx context.Context, phone string) (*aggregate.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	roles, err := r.findUserRoles(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	return r.poToUser(&po, roles), nil
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*aggregate.User, error) {
	var po UserPO
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	roles, err := r.findUserRoles(ctx, po.ID)
	if err != nil {
		return nil, err
	}

	return r.poToUser(&po, roles), nil
}

func (r *UserRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*aggregate.User, int64, error) {
	var pos []UserPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&UserPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	users := make([]*aggregate.User, 0, len(pos))
	for _, po := range pos {
		roles, err := r.findUserRoles(ctx, po.ID)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, r.poToUser(&po, roles))
	}

	return users, total, nil
}

func (r *UserRepoImpl) FindByType(ctx context.Context, userType entity.UserType) ([]*aggregate.User, error) {
	var pos []UserPO
	if err := r.db.WithContext(ctx).Where("user_type = ?", userType.Code()).Find(&pos).Error; err != nil {
		return nil, err
	}

	users := make([]*aggregate.User, 0, len(pos))
	for _, po := range pos {
		roles, err := r.findUserRoles(ctx, po.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, r.poToUser(&po, roles))
	}

	return users, nil
}

func (r *UserRepoImpl) FindByStatus(ctx context.Context, status entity.UserStatus) ([]*aggregate.User, error) {
	var pos []UserPO
	if err := r.db.WithContext(ctx).Where("status = ?", status.Code()).Find(&pos).Error; err != nil {
		return nil, err
	}

	users := make([]*aggregate.User, 0, len(pos))
	for _, po := range pos {
		roles, err := r.findUserRoles(ctx, po.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, r.poToUser(&po, roles))
	}

	return users, nil
}

func (r *UserRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("user_id = ?", id).Delete(&UserRolePO{})
		return tx.Delete(&UserPO{}, id).Error
	})
}

func (r *UserRepoImpl) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserPO{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepoImpl) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserPO{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepoImpl) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserPO{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepoImpl) findUserRoles(ctx context.Context, userID int64) ([]*entity.Role, error) {
	var userRoles []UserRolePO
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return []*entity.Role{}, nil
	}

	roleIDs := make([]int64, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	var rolePOs []RolePO
	if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&rolePOs).Error; err != nil {
		return nil, err
	}

	roles := make([]*entity.Role, 0, len(rolePOs))
	for _, rp := range rolePOs {
		permissions, err := r.findRolePermissions(ctx, rp.ID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r.poToRole(&rp, permissions))
	}

	return roles, nil
}

func (r *UserRepoImpl) findRolePermissions(ctx context.Context, roleID int64) ([]*entity.Permission, error) {
	var rolePermissions []RolePermissionPO
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	if len(rolePermissions) == 0 {
		return []*entity.Permission{}, nil
	}

	permIDs := make([]int64, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		permIDs = append(permIDs, rp.PermissionID)
	}

	var permPOs []PermissionPO
	if err := r.db.WithContext(ctx).Where("id IN ?", permIDs).Find(&permPOs).Error; err != nil {
		return nil, err
	}

	permissions := make([]*entity.Permission, 0, len(permPOs))
	for _, pp := range permPOs {
		permissions = append(permissions, r.poToPermission(&pp))
	}

	return permissions, nil
}

func (r *UserRepoImpl) userToPO(user *aggregate.User) *UserPO {
	return &UserPO{
		ID:       user.ID(),
		UserType: user.UserType().Code(),
		Username: user.Username(),
		Password: user.Password(),
		RealName: user.RealName(),
		Phone:    user.Phone(),
		Email:    user.Email(),
		Avatar:   user.Avatar(),
		Status:   user.Status().Code(),
		CreatedAt: user.CreatedAt(),
		UpdatedAt: user.UpdatedAt(),
	}
}

func (r *UserRepoImpl) poToUser(po *UserPO, roles []*entity.Role) *aggregate.User {
	return aggregate.ReconstructUser(
		po.ID,
		entity.UserTypeFromCode(po.UserType),
		po.Username,
		po.Password,
		po.RealName,
		po.Phone,
		po.Email,
		po.Avatar,
		entity.UserStatusFromCode(po.Status),
		roles,
		nil,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *UserRepoImpl) poToRole(po *RolePO, permissions []*entity.Permission) *entity.Role {
	return entity.ReconstructRole(
		po.ID,
		po.Code,
		po.Name,
		entity.DataScopeFromCode(po.DataScope),
		permissions,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

func (r *UserRepoImpl) poToPermission(po *PermissionPO) *entity.Permission {
	return entity.ReconstructPermission(
		po.ID,
		po.Code,
		po.Name,
		entity.ResourceTypeFromCode(po.ResourceType),
		po.ResourceURL,
		po.CreatedAt,
		po.UpdatedAt,
	)
}

type RoleRepoImpl struct {
	db *gorm.DB
}

func NewRoleRepoImpl(db *gorm.DB) repository.RoleRepository {
	return &RoleRepoImpl{db: db}
}

func (r *RoleRepoImpl) Save(ctx context.Context, role *entity.Role) error {
	po := RolePO{
		ID:        role.ID(),
		Code:      role.Code(),
		Name:      role.Name(),
		DataScope: role.DataScope().Code(),
		CreatedAt: role.CreatedAt(),
		UpdatedAt: role.UpdatedAt(),
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if po.ID == 0 {
			if err := tx.Create(&po).Error; err != nil {
				return err
			}
			role.SetID(po.ID)
		} else {
			if err := tx.Save(&po).Error; err != nil {
				return err
			}
		}

		tx.Where("role_id = ?", role.ID()).Delete(&RolePermissionPO{})

		for _, perm := range role.Permissions() {
			rp := RolePermissionPO{RoleID: role.ID(), PermissionID: perm.ID()}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *RoleRepoImpl) FindByID(ctx context.Context, id int64) (*entity.Role, error) {
	var po RolePO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var rolePermissions []RolePermissionPO
	r.db.WithContext(ctx).Where("role_id = ?", po.ID).Find(&rolePermissions)

	permIDs := make([]int64, 0, len(rolePermissions))
	for _, rp := range rolePermissions {
		permIDs = append(permIDs, rp.PermissionID)
	}

	var permPOs []PermissionPO
	if len(permIDs) > 0 {
		r.db.WithContext(ctx).Where("id IN ?", permIDs).Find(&permPOs)
	}

	permissions := make([]*entity.Permission, 0, len(permPOs))
	for _, pp := range permPOs {
		permissions = append(permissions, entity.ReconstructPermission(
			pp.ID, pp.Code, pp.Name,
			entity.ResourceTypeFromCode(pp.ResourceType),
			pp.ResourceURL,
			pp.CreatedAt, pp.UpdatedAt,
		))
	}

	return entity.ReconstructRole(
		po.ID, po.Code, po.Name,
		entity.DataScopeFromCode(po.DataScope),
		permissions,
		po.CreatedAt, po.UpdatedAt,
	), nil
}

func (r *RoleRepoImpl) FindByCode(ctx context.Context, code string) (*entity.Role, error) {
	var po RolePO
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.FindByID(ctx, po.ID)
}

func (r *RoleRepoImpl) FindAll(ctx context.Context) ([]*entity.Role, error) {
	var pos []RolePO
	if err := r.db.WithContext(ctx).Find(&pos).Error; err != nil {
		return nil, err
	}

	roles := make([]*entity.Role, 0, len(pos))
	for _, po := range pos {
		role, err := r.FindByID(ctx, po.ID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RoleRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tx.Where("role_id = ?", id).Delete(&RolePermissionPO{})
		tx.Where("role_id = ?", id).Delete(&UserRolePO{})
		return tx.Delete(&RolePO{}, id).Error
	})
}

func (r *RoleRepoImpl) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RolePO{}).Where("code = ?", code).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func parseDeptIDs(deptIDsStr string) []int64 {
	if deptIDsStr == "" {
		return []int64{}
	}
	parts := strings.Split(deptIDsStr, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		id, _ := strconv.ParseInt(p, 10, 64)
		ids = append(ids, id)
	}
	return ids
}

func deptIDsToString(deptIDs []int64) string {
	parts := make([]string, 0, len(deptIDs))
	for _, id := range deptIDs {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}