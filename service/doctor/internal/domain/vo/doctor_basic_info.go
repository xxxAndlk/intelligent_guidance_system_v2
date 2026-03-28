package vo

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidPhone    = errors.New("invalid phone number")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidIDCard   = errors.New("invalid ID card number")
)

// DoctorBasicInfo represents the basic personal information of a doctor
type DoctorBasicInfo struct {
	name       string
	gender     string
	birthDate  time.Time
	phone      string
	email      string
	idCard     string
	address    string
	avatarURL  string
}

// NewDoctorBasicInfo creates a new DoctorBasicInfo with validation
func NewDoctorBasicInfo(
	name string,
	gender string,
	birthDate time.Time,
	phone string,
	email string,
	idCard string,
	address string,
	avatarURL string,
) (*DoctorBasicInfo, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if phone != "" && !isValidPhone(phone) {
		return nil, ErrInvalidPhone
	}
	if email != "" && !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	return &DoctorBasicInfo{
		name:      name,
		gender:    gender,
		birthDate: birthDate,
		phone:     phone,
		email:     email,
		idCard:    idCard,
		address:   address,
		avatarURL: avatarURL,
	}, nil
}

// Getters
func (b *DoctorBasicInfo) Name() string       { return b.name }
func (b *DoctorBasicInfo) Gender() string     { return b.gender }
func (b *DoctorBasicInfo) BirthDate() time.Time { return b.birthDate }
func (b *DoctorBasicInfo) Phone() string      { return b.phone }
func (b *DoctorBasicInfo) Email() string      { return b.email }
func (b *DoctorBasicInfo) IDCard() string     { return b.idCard }
func (b *DoctorBasicInfo) Address() string    { return b.address }
func (b *DoctorBasicInfo) AvatarURL() string  { return b.avatarURL }

// UpdatePhone updates the phone number with validation
func (b *DoctorBasicInfo) UpdatePhone(phone string) error {
	if phone != "" && !isValidPhone(phone) {
		return ErrInvalidPhone
	}
	b.phone = phone
	return nil
}

// UpdateEmail updates the email with validation
func (b *DoctorBasicInfo) UpdateEmail(email string) error {
	if email != "" && !isValidEmail(email) {
		return ErrInvalidEmail
	}
	b.email = email
	return nil
}

// UpdateAddress updates the address
func (b *DoctorBasicInfo) UpdateAddress(address string) {
	b.address = address
}

// UpdateAvatarURL updates the avatar URL
func (b *DoctorBasicInfo) UpdateAvatarURL(url string) {
	b.avatarURL = url
}

// Age calculates the doctor's age
func (b *DoctorBasicInfo) Age() int {
	now := time.Now()
	age := now.Year() - b.birthDate.Year()
	if now.Month() < b.birthDate.Month() ||
		(now.Month() == b.birthDate.Month() && now.Day() < b.birthDate.Day()) {
		age--
	}
	return age
}

// Validation helpers
func isValidPhone(phone string) bool {
	// Chinese mobile phone number pattern
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	return matched
}

func isValidEmail(email string) bool {
	// Simple email pattern
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return matched
}