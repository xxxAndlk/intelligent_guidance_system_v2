package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("token has expired")
	ErrInvalidClaims       = errors.New("invalid token claims")
	ErrTokenNotValidYet    = errors.New("token is not valid yet")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
)

// JWTManager handles JWT token generation and verification
type JWTManager struct {
	secretKey      []byte
	accessExpiry   time.Duration
	refreshExpiry  time.Duration
	issuer         string
	signingMethod  jwt.SigningMethod
}

// Config holds JWT manager configuration
type Config struct {
	SecretKey      string        `json:"secret_key" yaml:"secret_key"`
	AccessExpiry   time.Duration `json:"access_expiry" yaml:"access_expiry"`
	RefreshExpiry  time.Duration `json:"refresh_expiry" yaml:"refresh_expiry"`
	Issuer         string        `json:"issuer" yaml:"issuer"`
}

// NewJWTManager creates a new JWT manager instance
func NewJWTManager(cfg Config) (*JWTManager, error) {
	if cfg.SecretKey == "" {
		return nil, errors.New("secret key cannot be empty")
	}
	if cfg.AccessExpiry == 0 {
		cfg.AccessExpiry = 15 * time.Minute
	}
	if cfg.RefreshExpiry == 0 {
		cfg.RefreshExpiry = 7 * 24 * time.Hour
	}
	if cfg.Issuer == "" {
		cfg.Issuer = "intelligent-guidance-system"
	}

	return &JWTManager{
		secretKey:      []byte(cfg.SecretKey),
		accessExpiry:   cfg.AccessExpiry,
		refreshExpiry:  cfg.RefreshExpiry,
		issuer:         cfg.Issuer,
		signingMethod:  jwt.SigningMethodHS256,
	}, nil
}

// GenerateToken generates both access and refresh tokens for a user
func (j *JWTManager) GenerateToken(claims UserClaims) (*TokenPair, error) {
	now := time.Now()
	
	// Generate access token
	accessExpiry := now.Add(j.accessExpiry)
	accessClaims := UserClaims{
		UserID:   claims.UserID,
		UserType: claims.UserType,
		Username: claims.Username,
		Roles:    claims.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   claims.Username,
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-access", claims.UserID),
		},
	}

	accessToken, err := j.createToken(&accessClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token
	refreshExpiry := now.Add(j.refreshExpiry)
	refreshClaims := UserClaims{
		UserID:   claims.UserID,
		UserType: claims.UserType,
		Username: claims.Username,
		Roles:    claims.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   claims.Username,
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-refresh", claims.UserID),
		},
	}

	refreshToken, err := j.createToken(&refreshClaims)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken generates only an access token
func (j *JWTManager) GenerateAccessToken(claims UserClaims) (string, error) {
	now := time.Now()
	accessExpiry := now.Add(j.accessExpiry)
	
	accessClaims := UserClaims{
		UserID:   claims.UserID,
		UserType: claims.UserType,
		Username: claims.Username,
		Roles:    claims.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   claims.Username,
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-access", claims.UserID),
		},
	}

	return j.createToken(&accessClaims)
}

// VerifyToken verifies a token and returns the claims
func (j *JWTManager) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrInvalidSigningMethod, token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// RefreshToken generates new tokens from a valid refresh token
func (j *JWTManager) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := j.VerifyToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return j.GenerateToken(*claims)
}

// ParseUnverified parses token without verification (use with caution)
func (j *JWTManager) ParseUnverified(tokenString string) (*UserClaims, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &UserClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

// createToken creates a signed token from claims
func (j *JWTManager) createToken(claims *UserClaims) (string, error) {
	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString(j.secretKey)
}

// GetAccessExpiry returns the access token expiry duration
func (j *JWTManager) GetAccessExpiry() time.Duration {
	return j.accessExpiry
}

// GetRefreshExpiry returns the refresh token expiry duration
func (j *JWTManager) GetRefreshExpiry() time.Duration {
	return j.refreshExpiry
}