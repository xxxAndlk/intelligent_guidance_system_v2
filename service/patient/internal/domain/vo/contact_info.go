// service/patient/internal/domain/vo/contact_info.go
package vo

import (
	"errors"
	"regexp"
)

// ContactInfo 联系信息值对象
type ContactInfo struct {
	Phone           Phone
	Email           string
	Address         string
	EmergencyContact EmergencyContact
}

// EmergencyContact 紧急联系人
type EmergencyContact struct {
	Name   string
	Phone  Phone
	Relation string
}

// 邮箱正则表达式
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// NewContactInfo 创建联系信息
func NewContactInfo(phone Phone, email string, address string) ContactInfo {
	return ContactInfo{
		Phone:   phone,
		Email:   email,
		Address: address,
	}
}

// Validate 验证联系信息
func (c ContactInfo) Validate() error {
	if err := c.Phone.Validate(); err != nil {
		return err
	}
	if c.Email != "" && !isValidEmail(c.Email) {
		return errors.New("invalid email format")
	}
	if len(c.Address) > 200 {
		return errors.New("address length exceeds limit")
	}
	if c.EmergencyContact.Name != "" {
		if err := c.EmergencyContact.Phone.Validate(); err != nil {
			return errors.New("invalid emergency contact phone")
		}
		if c.EmergencyContact.Relation == "" {
			return errors.New("emergency contact relation is required")
		}
	}
	return nil
}

// UpdatePhone 更新手机号
func (c *ContactInfo) UpdatePhone(phone Phone) {
	c.Phone = phone
}

// UpdateEmail 更新邮箱
func (c *ContactInfo) UpdateEmail(email string) error {
	if email != "" && !isValidEmail(email) {
		return errors.New("invalid email format")
	}
	c.Email = email
	return nil
}

// UpdateAddress 更新地址
func (c *ContactInfo) UpdateAddress(address string) error {
	if len(address) > 200 {
		return errors.New("address length exceeds limit")
	}
	c.Address = address
	return nil
}

// SetEmergencyContact 设置紧急联系人
func (c *ContactInfo) SetEmergencyContact(name string, phone Phone, relation string) error {
	if name == "" {
		return errors.New("emergency contact name is required")
	}
	if err := phone.Validate(); err != nil {
		return errors.New("invalid emergency contact phone")
	}
	if relation == "" {
		return errors.New("emergency contact relation is required")
	}

	c.EmergencyContact = EmergencyContact{
		Name:     name,
		Phone:    phone,
		Relation: relation,
	}
	return nil
}

// GetMaskedPhone 获取脱敏手机号
func (c ContactInfo) GetMaskedPhone() string {
	return c.Phone.Masked()
}

// GetMaskedEmergencyPhone 获取脱敏紧急联系人手机号
func (c ContactInfo) GetMaskedEmergencyPhone() string {
	if c.EmergencyContact.Phone.String() == "" {
		return ""
	}
	return c.EmergencyContact.Phone.Masked()
}

// isValidEmail 验证邮箱格式
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}