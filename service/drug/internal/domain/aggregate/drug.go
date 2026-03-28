package aggregate

import (
	"errors"
	"time"

	"intelligent-guidance-system/service/drug/internal/domain/entity"
	"intelligent-guidance-system/service/drug/internal/domain/event"
)

var (
	ErrDrugNotFound     = errors.New("drug not found")
	ErrDrugOutOfStock   = errors.New("drug out of stock")
	ErrDrugDiscontinued = errors.New("drug discontinued")
)

type Drug struct {
	id               int64
	name             string
	description      string
	price            int64
	quantityInStock  int
	category         entity.DrugCategory
	useMethod        string
	expirationDate   time.Time
	status           entity.DrugStatus
	stocks           []*entity.DrugStock
	events           []*event.DrugEvent
	createdAt        time.Time
	updatedAt        time.Time
}

func NewDrug(
	name string,
	description string,
	price int64,
	category entity.DrugCategory,
	useMethod string,
) (*Drug, error) {
	if name == "" {
		return nil, entity.ErrEmptyDrugName
	}
	if price < 0 {
		return nil, entity.ErrInvalidPrice
	}
	if category == entity.DrugCategoryUnknown {
		category = entity.DrugCategoryOther
	}

	now := time.Now()
	return &Drug{
		id:              0,
		name:            name,
		description:     description,
		price:           price,
		quantityInStock: 0,
		category:        category,
		useMethod:       useMethod,
		status:          entity.DrugStatusAvailable,
		stocks:          make([]*entity.DrugStock, 0),
		events:          make([]*event.DrugEvent, 0),
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructDrug(
	id int64,
	name string,
	description string,
	price int64,
	quantityInStock int,
	category entity.DrugCategory,
	useMethod string,
	expirationDate time.Time,
	status entity.DrugStatus,
	stocks []*entity.DrugStock,
	createdAt time.Time,
	updatedAt time.Time,
) *Drug {
	return &Drug{
		id:              id,
		name:            name,
		description:     description,
		price:           price,
		quantityInStock: quantityInStock,
		category:        category,
		useMethod:       useMethod,
		expirationDate:  expirationDate,
		status:          status,
		stocks:          stocks,
		events:          make([]*event.DrugEvent, 0),
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}
}

func (d *Drug) ID() int64 { return d.id }
func (d *Drug) Name() string { return d.name }
func (d *Drug) Description() string { return d.description }
func (d *Drug) Price() int64 { return d.price }
func (d *Drug) QuantityInStock() int { return d.quantityInStock }
func (d *Drug) Category() entity.DrugCategory { return d.category }
func (d *Drug) UseMethod() string { return d.useMethod }
func (d *Drug) ExpirationDate() time.Time { return d.expirationDate }
func (d *Drug) Status() entity.DrugStatus { return d.status }
func (d *Drug) Stocks() []*entity.DrugStock { return d.stocks }
func (d *Drug) Events() []*event.DrugEvent { return d.events }
func (d *Drug) CreatedAt() time.Time { return d.createdAt }
func (d *Drug) UpdatedAt() time.Time { return d.updatedAt }

func (d *Drug) SetID(id int64) { d.id = id }

func (d *Drug) UpdateStock(quantity int, operatorID int64) error {
	d.quantityInStock = quantity
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StockUpdatedEvent(d.id, operatorID, quantity))
	
	if d.quantityInStock <= 0 {
		d.status = entity.DrugStatusOutOfStock
	} else if d.status == entity.DrugStatusOutOfStock {
		d.status = entity.DrugStatusAvailable
	}
	
	return nil
}

func (d *Drug) AddStock(quantity int, operatorID int64) error {
	if quantity < 0 {
		return entity.ErrInvalidQuantity
	}
	d.quantityInStock += quantity
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StockAddedEvent(d.id, operatorID, quantity))
	
	if d.status == entity.DrugStatusOutOfStock && d.quantityInStock > 0 {
		d.status = entity.DrugStatusAvailable
	}
	
	return nil
}

func (d *Drug) ReduceStock(quantity int, operatorID int64) error {
	if quantity < 0 {
		return entity.ErrInvalidQuantity
	}
	if d.quantityInStock < quantity {
		return entity.ErrDrugOutOfStock
	}
	d.quantityInStock -= quantity
	d.updatedAt = time.Now()
	d.events = append(d.events, event.StockReducedEvent(d.id, operatorID, quantity))
	
	if d.quantityInStock <= 0 {
		d.status = entity.DrugStatusOutOfStock
	}
	
	return nil
}

func (d *Drug) CheckAvailability(quantity int) bool {
	return d.status.IsAvailable() && d.quantityInStock >= quantity
}

func (d *Drug) Discontinue() {
	d.status = entity.DrugStatusDiscontinued
	d.updatedAt = time.Now()
}

func (d *Drug) Activate() {
	if d.quantityInStock > 0 {
		d.status = entity.DrugStatusAvailable
	} else {
		d.status = entity.DrugStatusOutOfStock
	}
	d.updatedAt = time.Now()
}

func (d *Drug) UpdatePrice(newPrice int64) error {
	if newPrice < 0 {
		return entity.ErrInvalidPrice
	}
	d.price = newPrice
	d.updatedAt = time.Now()
	return nil
}

func (d *Drug) UpdateDescription(description string) {
	d.description = description
	d.updatedAt = time.Now()
}

func (d *Drug) ClearEvents() {
	d.events = make([]*event.DrugEvent, 0)
}

func (d *Drug) PriceYuan() float64 {
	return float64(d.price) / 100.0
}