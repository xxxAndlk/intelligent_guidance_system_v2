package auth

import (
	"context"
	"strings"

	"intelligent_guidance_system_v2/common/pkg/errors"
	"intelligent_guidance_system_v2/common/pkg/jwt"
	"intelligent_guidance_system_v2/common/pkg/response"

	kratos "github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	UserTypeKey  contextKey = "user_type"
	UsernameKey  contextKey = "username"
	RolesKey     contextKey = "roles"
	ClaimsKey    contextKey = "claims"
)

func AuthMiddleware(jwtManager *jwt.JWTManager) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.NewUnauthorized("missing transport context")
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return nil, errors.NewUnauthorized("missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return nil, errors.NewUnauthorized("invalid authorization header format")
			}

			token := parts[1]
			claims, err := jwtManager.VerifyToken(token)
			if err != nil {
				return nil, errors.NewUnauthorized("invalid token: " + err.Error())
			}

			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserTypeKey, claims.UserType)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			ctx = context.WithValue(ctx, RolesKey, claims.Roles)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			return handler(ctx, req)
		}
	}
}

func OptionalAuthMiddleware(jwtManager *jwt.JWTManager) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			authHeader := tr.RequestHeader().Get("Authorization")
			if authHeader == "" {
				return handler(ctx, req)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return handler(ctx, req)
			}

			token := parts[1]
			claims, err := jwtManager.VerifyToken(token)
			if err != nil {
				return handler(ctx, req)
			}

			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserTypeKey, claims.UserType)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)
			ctx = context.WithValue(ctx, RolesKey, claims.Roles)
			ctx = context.WithValue(ctx, ClaimsKey, claims)

			return handler(ctx, req)
		}
	}
}

func PermissionMiddleware(requiredPerms ...string) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			roles, ok := ctx.Value(RolesKey).([]string)
			if !ok || len(roles) == 0 {
				return nil, errors.NewForbidden("no roles assigned")
			}

			roleSet := make(map[string]struct{}, len(roles))
			for _, r := range roles {
				roleSet[r] = struct{}{}
			}

			for _, perm := range requiredPerms {
				if _, exists := roleSet[perm]; !exists {
					return nil, errors.NewForbidden("permission denied: " + perm)
				}
			}

			return handler(ctx, req)
		}
	}
}

func RoleMiddleware(requiredRoles ...string) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			roles, ok := ctx.Value(RolesKey).([]string)
			if !ok || len(roles) == 0 {
				return nil, errors.NewForbidden("no roles assigned")
			}

			for _, requiredRole := range requiredRoles {
				for _, role := range roles {
					if role == requiredRole {
						return handler(ctx, req)
					}
				}
			}

			return nil, errors.NewForbidden("role required")
		}
	}
}

func UserTypeMiddleware(allowedTypes ...int32) kratos.Middleware {
	return func(handler kratos.Handler) kratos.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			userType, ok := ctx.Value(UserTypeKey).(int32)
			if !ok {
				return nil, errors.NewForbidden("user type not found")
			}

			for _, t := range allowedTypes {
				if userType == t {
					return handler(ctx, req)
				}
			}

			return nil, errors.NewForbidden("user type not allowed")
		}
	}
}

func GetUserID(ctx context.Context) int64 {
	userID, _ := ctx.Value(UserIDKey).(int64)
	return userID
}

func GetUserType(ctx context.Context) int32 {
	userType, _ := ctx.Value(UserTypeKey).(int32)
	return userType
}

func GetUsername(ctx context.Context) string {
	username, _ := ctx.Value(UsernameKey).(string)
	return username
}

func GetRoles(ctx context.Context) []string {
	roles, _ := ctx.Value(RolesKey).([]string)
	return roles
}

func GetClaims(ctx context.Context) *jwt.UserClaims {
	claims, _ := ctx.Value(ClaimsKey).(*jwt.UserClaims)
	return claims
}

func HasRole(ctx context.Context, role string) bool {
	roles := GetRoles(ctx)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func HasAnyRole(ctx context.Context, roles ...string) bool {
	for _, role := range roles {
		if HasRole(ctx, role) {
			return true
		}
	}
	return false
}

func MustGetUserID(ctx context.Context) (int64, error) {
	userID := GetUserID(ctx)
	if userID == 0 {
		return 0, errors.NewUnauthorized("user not authenticated")
	}
	return userID, nil
}

func ResponseUnauthorized() *response.Response {
	return response.Error(errors.CodeUnauthenticated, "unauthorized")
}

func ResponseForbidden() *response.Response {
	return response.Error(errors.CodePermissionDenied, "forbidden")
}