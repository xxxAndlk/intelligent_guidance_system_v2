package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidStockID    = errors.New("invalid stock ID")
	ErrStockNotFound     = errors.New("stock not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type DrugStock struct {
	id              int64
	drugID          int64
	batchNumber     string
	quantity        int
	expirationDate  time.Time
	purchasePrice   int64
	sellingPrice    int64
	supplier        string
	createdAt       time.Time
	updatedAt       time.Time
}

func NewDrugStock(
	drugID int64,
	batchNumber string,
	quantity int,
	expirationDate time.Time,
	purchasePrice int64,
	sellingPrice int64,
	supplier string,
) (*DrugStock, error) {
	if drugID <= 0 {
		return nil, ErrInvalidDrugID
	}
	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	now := time.Now()
	return &DrugStock{
		id:             0,
		drugID:         drugID,
		batchNumber:    batchNumber,
		quantity:       quantity,
		expirationDate: expirationDate,
		purchasePrice:  purchasePrice,
		sellingPrice:   sellingPrice,
		supplier:       supplier,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func ReconstructDrugStock(
	id int64,
	drugID int64,
	batchNumber string,
	quantity int,
	expirationDate time.Time,
	purchasePrice int64,
	sellingPrice int64,
	supplier string,
	createdAt time.Time,
	updatedAt time.Time,
) *DrugStock {
	return &DrugStock{
		id:             id,
		drugID:         drugID,
		batchNumber:    batchNumber,
		quantity:       quantity,
		expirationDate: expirationDate,
		purchasePrice:  purchasePrice,
		sellingPrice:   sellingPrice,
		supplier:       supplier,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (s *DrugStock) ID() int64 { return s.id }
func (s *DrugStock) DrugID() int64 { return s.drugID }
func (s *DrugStock) BatchNumber() string { return s.batchNumber }
func (s *DrugStock) Quantity() int { return s.quantity }
func (s *DrugStock) ExpirationDate() time.Time { return s.expirationDate }
func (s *DrugStock) PurchasePrice() int64 { return s.purchasePrice }
func (s *DrugStock) SellingPrice() int64 { return s.sellingPrice }
func (s *DrugStock) Supplier() string { return s.supplier }
func (s *DrugStock) CreatedAt() time.Time { return s.createdAt }
func (s *DrugStock) UpdatedAt() time.Time { return s.updatedAt }

func (s *DrugStock) SetID(id int64) { s.id = id }

func (s *DrugStock) AddQuantity(quantity int) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	s.quantity += quantity
	s.updatedAt = time.Now()
	return nil
}

func (s *DrugStock) ReduceQuantity(quantity int) error {
	if quantity < 0 {
		return ErrInvalidQuantity
	}
	if s.quantity < quantity {
		return ErrInsufficientStock
	}
	s.quantity -= quantity
	s.updatedAt = time.Now()
	return nil
}

func (s *DrugStock) IsExpired() bool {
	return time.Now().After(s.expirationDate)
}

func (s *DrugStock) IsNearExpiry(days int) bool {
	return time.Now().AddDate(0, 0, days).After(s.expirationDate)
}