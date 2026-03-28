package entity

import "errors"

var (
	ErrInvalidDrugID    = errors.New("invalid drug ID")
	ErrEmptyDrugName    = errors.New("drug name cannot be empty")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrDrugOutOfStock   = errors.New("drug out of stock")
)

type DrugCategory int

const (
	DrugCategoryUnknown DrugCategory = iota
	DrugCategoryWestern
	DrugCategoryChinese
	DrugCategoryBiological
	DrugCategoryOther
)

func (c DrugCategory) String() string {
	switch c {
	case DrugCategoryWestern:
		return "西药"
	case DrugCategoryChinese:
		return "中药"
	case DrugCategoryBiological:
		return "生物制品"
	case DrugCategoryOther:
		return "其他"
	default:
		return "未知"
	}
}

func DrugCategoryFromCode(code string) DrugCategory {
	switch code {
	case "WESTERN":
		return DrugCategoryWestern
	case "CHINESE":
		return DrugCategoryChinese
	case "BIOLOGICAL":
		return DrugCategoryBiological
	case "OTHER":
		return DrugCategoryOther
	default:
		return DrugCategoryUnknown
	}
}

func (c DrugCategory) Code() string {
	switch c {
	case DrugCategoryWestern:
		return "WESTERN"
	case DrugCategoryChinese:
		return "CHINESE"
	case DrugCategoryBiological:
		return "BIOLOGICAL"
	case DrugCategoryOther:
		return "OTHER"
	default:
		return "UNKNOWN"
	}
}

type DrugStatus int

const (
	DrugStatusUnknown DrugStatus = iota
	DrugStatusAvailable
	DrugStatusOutOfStock
	DrugStatusDiscontinued
)

func (s DrugStatus) String() string {
	switch s {
	case DrugStatusAvailable:
		return "可用"
	case DrugStatusOutOfStock:
		return "缺货"
	case DrugStatusDiscontinued:
		return "停用"
	default:
		return "未知"
	}
}

func DrugStatusFromCode(code string) DrugStatus {
	switch code {
	case "AVAILABLE":
		return DrugStatusAvailable
	case "OUT_OF_STOCK":
		return DrugStatusOutOfStock
	case "DISCONTINUED":
		return DrugStatusDiscontinued
	default:
		return DrugStatusUnknown
	}
}

func (s DrugStatus) Code() string {
	switch s {
	case DrugStatusAvailable:
		return "AVAILABLE"
	case DrugStatusOutOfStock:
		return "OUT_OF_STOCK"
	case DrugStatusDiscontinued:
		return "DISCONTINUED"
	default:
		return "UNKNOWN"
	}
}

func (s DrugStatus) IsAvailable() bool {
	return s == DrugStatusAvailable
}