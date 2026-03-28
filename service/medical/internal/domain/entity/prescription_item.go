package entity

import (
	"errors"

	"github.com/google/uuid"
	"intelligent_guidance_system_v2/service/medical/internal/domain/vo"
)

var (
	ErrEmptyDrugName  = errors.New("drug name cannot be empty")
	ErrInvalidQuantity = errors.New("quantity must be positive")
	ErrInvalidPrice    = errors.New("price must be non-negative")
	ErrEmptyUsage      = errors.New("usage instructions cannot be empty")
)

// PrescriptionItem represents a prescription drug item entity
type PrescriptionItem struct {
	id          string    // unique identifier
	drugID      int64     // reference to drug catalog
	drugName    string    // drug name
	quantity    int32     // quantity prescribed
	price       *vo.Money // unit price
	usage       string    // usage and dosage instructions
}

// NewPrescriptionItem creates a new PrescriptionItem entity
func NewPrescriptionItem(drugID int64, drugName string, quantity int32, priceInFen int64, usage string) (*PrescriptionItem, error) {
	if drugName == "" {
		return nil, ErrEmptyDrugName
	}
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if priceInFen < 0 {
		return nil, ErrInvalidPrice
	}
	if usage == "" {
		return nil, ErrEmptyUsage
	}

	price, err := vo.NewMoney(priceInFen, "CNY")
	if err != nil {
		return nil, err
	}

	return &PrescriptionItem{
		id:       uuid.New().String(),
		drugID:   drugID,
		drugName: drugName,
		quantity: quantity,
		price:    price,
		usage:    usage,
	}, nil
}

// ReconstructPrescriptionItem reconstructs a PrescriptionItem from persistence
func ReconstructPrescriptionItem(id string, drugID int64, drugName string, quantity int32, price *vo.Money, usage string) *PrescriptionItem {
	return &PrescriptionItem{
		id:       id,
		drugID:   drugID,
		drugName: drugName,
		quantity: quantity,
		price:    price,
		usage:    usage,
	}
}

// ID returns the prescription item ID
func (p *PrescriptionItem) ID() string {
	return p.id
}

// DrugID returns the drug catalog ID
func (p *PrescriptionItem) DrugID() int64 {
	return p.drugID
}

// DrugName returns the drug name
func (p *PrescriptionItem) DrugName() string {
	return p.drugName
}

// Quantity returns the quantity
func (p *PrescriptionItem) Quantity() int32 {
	return p.quantity
}

// Price returns the unit price
func (p *PrescriptionItem) Price() *vo.Money {
	return p.price
}

// Usage returns the usage instructions
func (p *PrescriptionItem) Usage() string {
	return p.usage
}

// TotalPrice calculates the total price for this prescription item
func (p *PrescriptionItem) TotalPrice() (*vo.Money, error) {
	total := int64(p.quantity) * p.price.Amount()
	return vo.NewMoney(total, p.price.Currency())
}

// UpdateQuantity updates the quantity
func (p *PrescriptionItem) UpdateQuantity(quantity int32) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	p.quantity = quantity
	return nil
}

// UpdateUsage updates the usage instructions
func (p *PrescriptionItem) UpdateUsage(usage string) error {
	if usage == "" {
		return ErrEmptyUsage
	}
	p.usage = usage
	return nil
}