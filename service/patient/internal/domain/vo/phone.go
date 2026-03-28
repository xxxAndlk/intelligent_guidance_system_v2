// service/patient/internal/domain/vo/phone.go
package vo

import (
	"errors"
	"regexp"
)

// Phone 手机号值对象
type Phone struct {
	value string
}

// 手机号正则表达式（中国大陆手机号）
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// NewPhone 创建手机号值对象
func NewPhone(phone string) (Phone, error) {
	if !isValidPhone(phone) {
		return Phone{}, errors.New("invalid phone number format")
	}
	return Phone{value: phone}, nil
}

// String 返回手机号
func (p Phone) String() string {
	return p.value
}

// Validate 验证手机号
func (p Phone) Validate() error {
	if !isValidPhone(p.value) {
		return errors.New("invalid phone number format")
	}
	return nil
}

// Masked 返回脱敏后的手机号（显示前3后4位）
func (p Phone) Masked() string {
	if len(p.value) != 11 {
		return p.value
	}
	return p.value[:3] + "****" + p.value[7:]
}

// isValidPhone 验证手机号格式
func isValidPhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// Equals 判断两个手机号是否相等
func (p Phone) Equals(other Phone) bool {
	return p.value == other.value
}