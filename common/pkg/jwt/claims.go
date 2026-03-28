package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims represents the custom JWT claims for user authentication
type UserClaims struct {
	UserID   int64    `json:"user_id"`
	UserType int32    `json:"user_type"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenPair contains both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds until access token expires
	TokenType    string `json:"token_type"` // usually "Bearer"
}

// NewUserClaims creates a new UserClaims instance
func NewUserClaims(userID int64, userType int32, username string, roles []string, issuer string, expiry time.Time) *UserClaims {
	return &UserClaims{
		UserID:   userID,
		UserType: userType,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}

// HasRole checks if the user has a specific role
func (c *UserClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func (c *UserClaims) HasAnyRole(roles ...string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}

// HasAllRoles checks if the user has all of the specified roles
func (c *UserClaims) HasAllRoles(roles ...string) bool {
	for _, role := range roles {
		if !c.HasRole(role) {
			return false
		}
	}
	return true
}

// IsExpired checks if the token is expired
func (c *UserClaims) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}