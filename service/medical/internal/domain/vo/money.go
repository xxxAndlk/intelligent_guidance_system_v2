package vo

import (
	"errors"
	"fmt"
)

// Money value object - stores amount in fen (分) for precision
// 1 Yuan (元) = 100 Fen (分)
type Money struct {
	amount   int64  // amount in fen
	currency string // currency code, default "CNY"
}

var (
	ErrNegativeAmount = errors.New("money amount cannot be negative")
	ErrInvalidCurrency = errors.New("invalid currency code")
)

// NewMoney creates a new Money value object
func NewMoney(amountInFen int64, currency string) (*Money, error) {
	if amountInFen < 0 {
		return nil, ErrNegativeAmount
	}
	if currency == "" {
		currency = "CNY"
	}
	return &Money{
		amount:   amountInFen,
		currency: currency,
	}, nil
}

// NewMoneyFromYuan creates Money from yuan amount (float)
func NewMoneyFromYuan(yuan float64, currency string) (*Money, error) {
	if yuan < 0 {
		return nil, ErrNegativeAmount
	}
	fen := int64(yuan * 100)
	return NewMoney(fen, currency)
}

// Amount returns the amount in fen
func (m *Money) Amount() int64 {
	return m.amount
}

// Yuan returns the amount in yuan (float)
func (m *Money) Yuan() float64 {
	return float64(m.amount) / 100.0
}

// Currency returns the currency code
func (m *Money) Currency() string {
	return m.currency
}

// Add adds another Money to this one
func (m *Money) Add(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrInvalidCurrency
	}
	return NewMoney(m.amount + other.amount, m.currency)
}

// Subtract subtracts another Money from this one
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, ErrInvalidCurrency
	}
	result := m.amount - other.amount
	if result < 0 {
		return nil, ErrNegativeAmount
	}
	return NewMoney(result, m.currency)
}

// IsZero checks if the amount is zero
func (m *Money) IsZero() bool {
	return m.amount == 0
}

// Equals checks if two Money objects are equal
func (m *Money) Equals(other *Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// String returns the formatted money string
func (m *Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Yuan(), m.currency)
}

// Zero returns a zero Money value object
func Zero(currency string) *Money {
	if currency == "" {
		currency = "CNY"
	}
	return &Money{amount: 0, currency: currency}
}