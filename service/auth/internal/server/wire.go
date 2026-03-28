package server

import (
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"intelligent-guidance-system/service/auth/internal/biz/usecase"
	"intelligent-guidance-system/service/auth/internal/data/mysql"
	"intelligent-guidance-system/service/auth/internal/domain/repository"
	"intelligent-guidance-system/service/auth/internal/service"
)

var ProviderSet = wire.NewSet(
	NewConfig,
	NewLogger,
	NewGRPCServer,
	NewHTTPServer,
	NewApp,
	NewDB,
	NewRedis,
	NewUserRepository,
	NewRoleRepository,
	NewTokenRepository,
	NewAuthUsecase,
	NewAuthService,
)

func NewDB(cfg *Config) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(cfg.Database.Url), &gorm.Config{})
}

func NewRedis(cfg *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return mysql.NewUserRepoImpl(db)
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return mysql.NewRoleRepoImpl(db)
}

func NewTokenRepository(rdb *redis.Client) repository.TokenRepository {
	return NewRedisTokenRepo(rdb)
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	tokenRepo repository.TokenRepository,
	cfg *Config,
) *usecase.AuthUsecase {
	return usecase.NewAuthUsecase(userRepo, roleRepo, nil, tokenRepo, cfg.JWT.Secret)
}

func NewAuthService(uc *usecase.AuthUsecase) *service.AuthService {
	return service.NewAuthService(uc)
}

type RedisTokenRepo struct {
	rdb *redis.Client
}

func NewRedisTokenRepo(rdb *redis.Client) repository.TokenRepository {
	return &RedisTokenRepo{rdb: rdb}
}

func (r *RedisTokenRepo) SaveToken(ctx context.Context, userID int64, token string, expiresAt int64) error {
	return r.rdb.Set(ctx, "token:"+token, userID, 0).Err()
}

func (r *RedisTokenRepo) FindToken(ctx context.Context, userID int64) (string, int64, error) {
	return "", 0, nil
}

func (r *RedisTokenRepo) DeleteToken(ctx context.Context, userID int64) error {
	return r.rdb.Del(ctx, "user_token:"+string(userID)).Err()
}

func (r *RedisTokenRepo) ExistsToken(ctx context.Context, token string) (bool, error) {
	val, err := r.rdb.Exists(ctx, "token:"+token).Result()
	return val > 0, err
}