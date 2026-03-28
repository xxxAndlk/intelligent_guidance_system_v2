package vo

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidPhoneFormat = errors.New("invalid phone format")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)

type ContactInfo struct {
	phone string
	email string
}

func NewContactInfo(phone, email string) (*ContactInfo, error) {
	vo := &ContactInfo{}
	
	if phone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(phone) {
			return nil, ErrInvalidPhoneFormat
		}
		vo.phone = phone
	}

	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return nil, ErrInvalidEmailFormat
		}
		vo.email = email
	}

	return vo, nil
}

func (c *ContactInfo) Phone() string { return c.phone }
func (c *ContactInfo) Email() string { return c.email }

func (c *ContactInfo) UpdatePhone(phone string) error {
	if phone != "" {
		phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
		if !phoneRegex.MatchString(phone) {
			return ErrInvalidPhoneFormat
		}
	}
	c.phone = phone
	return nil
}

func (c *ContactInfo) UpdateEmail(email string) error {
	if email != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return ErrInvalidEmailFormat
		}
	}
	c.email = email
	return nil
}

type TokenInfo struct {
	token      string
	expiresAt  time.Time
	issuedAt   time.Time
}

func NewTokenInfo(token string, expiresAt time.Time) *TokenInfo {
	return &TokenInfo{
		token:     token,
		expiresAt: expiresAt,
		issuedAt:  time.Now(),
	}
}

func (t *TokenInfo) Token() string { return t.token }
func (t *TokenInfo) ExpiresAt() time.Time { return t.expiresAt }
func (t *TokenInfo) IssuedAt() time.Time { return t.issuedAt }

func (t *TokenInfo) IsExpired() bool {
	return time.Now().After(t.expiresAt)
}

func (t *TokenInfo) IsValid() bool {
	return t.token != "" && !t.IsExpired()
}