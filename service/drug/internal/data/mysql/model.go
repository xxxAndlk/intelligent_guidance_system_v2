package mysql

import "time"

type DrugPO struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	Name            string    `gorm:"column:name;size:100;unique;not null"`
	Description     string    `gorm:"column:description;size:500"`
	Price           int64     `gorm:"column:price"`
	QuantityInStock int       `gorm:"column:quantity_in_stock"`
	Category        string    `gorm:"column:category;size:20"`
	UseMethod       string    `gorm:"column:use_method;size:200"`
	ExpirationDate  time.Time `gorm:"column:expiration_date"`
	Status          string    `gorm:"column:status;size:20"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (DrugPO) TableName() string {
	return "drugs"
}

type DrugStockPO struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	DrugID         int64     `gorm:"column:drug_id;index"`
	BatchNumber    string    `gorm:"column:batch_number;size:50"`
	Quantity       int       `gorm:"column:quantity"`
	ExpirationDate time.Time `gorm:"column:expiration_date"`
	PurchasePrice  int64     `gorm:"column:purchase_price"`
	SellingPrice   int64     `gorm:"column:selling_price"`
	Supplier       string    `gorm:"column:supplier;size:100"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (DrugStockPO) TableName() string {
	return "drug_stocks"
}