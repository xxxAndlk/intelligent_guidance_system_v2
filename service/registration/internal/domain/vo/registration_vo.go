package vo

import (
	"errors"
	"time"
)

var (
	ErrInvalidFeeAmount = errors.New("invalid fee amount")
	ErrInvalidAppointmentTime = errors.New("invalid appointment time")
)

type RegistrationFee struct {
	amount   int64
	currency string
}

func NewRegistrationFee(amount int64, currency string) (*RegistrationFee, error) {
	if amount < 0 {
		return nil, ErrInvalidFeeAmount
	}
	if currency == "" {
		currency = "CNY"
	}
	return &RegistrationFee{amount: amount, currency: currency}, nil
}

func (f *RegistrationFee) Amount() int64 { return f.amount }
func (f *RegistrationFee) AmountYuan() float64 { return float64(f.amount) / 100.0 }
func (f *RegistrationFee) Currency() string { return f.currency }

type AppointmentTime struct {
	scheduledTime time.Time
}

func NewAppointmentTime(scheduledTime time.Time) (*AppointmentTime, error) {
	if scheduledTime.Before(time.Now()) {
		return nil, ErrInvalidAppointmentTime
	}
	return &AppointmentTime{scheduledTime: scheduledTime}, nil
}

func (a *AppointmentTime) ScheduledTime() time.Time { return a.scheduledTime }
func (a *AppointmentTime) IsToday() bool {
	now := time.Now()
	return a.scheduledTime.Year() == now.Year() && 
		a.scheduledTime.Month() == now.Month() && 
		a.scheduledTime.Day() == now.Day()
}