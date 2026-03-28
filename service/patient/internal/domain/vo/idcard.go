// service/patient/internal/domain/vo/idcard.go
package vo

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// IDCard 身份证值对象
type IDCard struct {
	value string
}

// 身份证正则表达式（18位）
var idCardRegex = regexp.MustCompile(`^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[0-9Xx]$`)

// NewIDCard 创建身份证值对象
func NewIDCard(idCard string) (IDCard, error) {
	if !isValidIDCard(idCard) {
		return IDCard{}, errors.New("invalid ID card format")
	}
	return IDCard{value: idCard}, nil
}

// String 返回身份证号
func (i IDCard) String() string {
	return i.value
}

// Validate 验证身份证
func (i IDCard) Validate() error {
	if !isValidIDCard(i.value) {
		return errors.New("invalid ID card format")
	}
	return nil
}

// Masked 返回脱敏后的身份证号（显示前6后4位）
func (i IDCard) Masked() string {
	if len(i.value) != 18 {
		return i.value
	}
	return i.value[:6] + "********" + i.value[14:]
}

// GetBirthDate 从身份证提取出生日期
func (i IDCard) GetBirthDate() (time.Time, error) {
	if len(i.value) != 18 {
		return time.Time{}, errors.New("invalid ID card length")
	}

	yearStr := i.value[6:10]
	monthStr := i.value[10:12]
	dayStr := i.value[12:14]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// GetGender 从身份证提取性别
func (i IDCard) GetGender() (Gender, error) {
	if len(i.value) != 18 {
		return GenderUnknown, errors.New("invalid ID card length")
	}

	// 第17位为性别标识，奇数为男，偶数为女
	genderDigit := i.value[16:17]
	genderNum, err := strconv.Atoi(genderDigit)
	if err != nil {
		return GenderUnknown, err
	}

	if genderNum%2 == 1 {
		return GenderMale, nil
	}
	return GenderFemale, nil
}

// GetAge 从身份证计算年龄
func (i IDCard) GetAge() (int32, error) {
	birthDate, err := i.GetBirthDate()
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()

	// 如果还没过今年的生日，年龄减1
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return int32(age), nil
}

// isValidIDCard 验证身份证格式（包括校验位）
func isValidIDCard(idCard string) bool {
	if !idCardRegex.MatchString(idCard) {
		return false
	}

	// 校验位验证
	// 权重因子
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验码对应值
	checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	// 计算校验位
	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * weights[i]
	}

	checkCode := checkCodes[sum%11]
	lastDigit := string(idCard[17])

	// 处理大小写
	if lastDigit == "x" {
		lastDigit = "X"
	}

	return checkCode == lastDigit
}

// Equals 判断两个身份证是否相等
func (i IDCard) Equals(other IDCard) bool {
	return i.value == other.value
}