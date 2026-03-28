package dto

import "time"

type CreateDrugRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required"`
	Category    string `json:"category"`
	UseMethod   string `json:"use_method"`
}

type UpdateDrugRequest struct {
	Description string `json:"description"`
	Price       int64  `json:"price"`
	UseMethod   string `json:"use_method"`
}

type UpdateStockRequest struct {
	Quantity   int   `json:"quantity" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type AddStockRequest struct {
	Quantity   int   `json:"quantity" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type ReduceStockRequest struct {
	Quantity   int   `json:"quantity" binding:"required"`
	OperatorID int64 `json:"operator_id"`
}

type CheckAvailabilityRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

type DrugResponse struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	QuantityInStock int       `json:"quantity_in_stock"`
	Category        string    `json:"category"`
	CategoryName    string    `json:"category_name"`
	UseMethod       string    `json:"use_method"`
	Status          string    `json:"status"`
	StatusName      string    `json:"status_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type DrugListResponse struct {
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Drugs    []DrugResponse `json:"drugs"`
}

type StockResponse struct {
	ID             int64     `json:"id"`
	DrugID         int64     `json:"drug_id"`
	BatchNumber    string    `json:"batch_number"`
	Quantity       int       `json:"quantity"`
	ExpirationDate time.Time `json:"expiration_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SellingPrice   float64   `json:"selling_price"`
	Supplier       string    `json:"supplier"`
}

type AvailabilityResponse struct {
	Available bool `json:"available"`
}