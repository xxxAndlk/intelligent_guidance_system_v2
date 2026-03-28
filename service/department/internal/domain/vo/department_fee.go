package vo

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidFeeAmount = errors.New("invalid fee amount")
	ErrNegativeFee      = errors.New("fee cannot be negative")
)

// DepartmentFee represents department registration fees
type DepartmentFee struct {
	registrationFee       int64 // in fen (cents)
	expertRegistrationFee int64 // in fen (cents)
	currency              string
}

func NewDepartmentFee(registrationFee, expertRegistrationFee int64, currency string) (*DepartmentFee, error) {
	if registrationFee < 0 {
		return nil, ErrNegativeFee
	}
	if expertRegistrationFee < 0 {
		return nil, ErrNegativeFee
	}
	if currency == "" {
		currency = "CNY"
	}
	return &DepartmentFee{
		registrationFee:       registrationFee,
		expertRegistrationFee: expertRegistrationFee,
		currency:              currency,
	}, nil
}

func ZeroDepartmentFee(currency string) *DepartmentFee {
	return &DepartmentFee{
		registrationFee:       0,
		expertRegistrationFee: 0,
		currency:              currency,
	}
}

func (f *DepartmentFee) RegistrationFee() int64 {
	return f.registrationFee
}

func (f *DepartmentFee) ExpertRegistrationFee() int64 {
	return f.expertRegistrationFee
}

func (f *DepartmentFee) Currency() string {
	return f.currency
}

func (f *DepartmentFee) RegistrationFeeYuan() float64 {
	return float64(f.registrationFee) / 100.0
}

func (f *DepartmentFee) ExpertRegistrationFeeYuan() float64 {
	return float64(f.expertRegistrationFee) / 100.0
}

func (f *DepartmentFee) SetRegistrationFee(amount int64) error {
	if amount < 0 {
		return ErrNegativeFee
	}
	f.registrationFee = amount
	return nil
}

func (f *DepartmentFee) SetExpertRegistrationFee(amount int64) error {
	if amount < 0 {
		return ErrNegativeFee
	}
	f.expertRegistrationFee = amount
	return nil
}

func (f *DepartmentFee) String() string {
	return fmt.Sprintf("普通挂号费: %.2f %s, 专家挂号费: %.2f %s", 
		f.RegistrationFeeYuan(), f.currency, 
		f.ExpertRegistrationFeeYuan(), f.currency)
}