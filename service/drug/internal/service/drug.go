package service

import (
	"context"

	"intelligent-guidance-system/service/drug/internal/biz/dto"
	"intelligent-guidance-system/service/drug/internal/biz/usecase"
	"intelligent-guidance-system/service/drug/internal/domain/aggregate"
	"intelligent-guidance-system/service/drug/internal/domain/entity"
)

type DrugService struct {
	usecase *usecase.DrugUsecase
}

func NewDrugService(usecase *usecase.DrugUsecase) *DrugService {
	return &DrugService{usecase: usecase}
}

func (s *DrugService) CreateDrug(ctx context.Context, req *dto.CreateDrugRequest) (*dto.DrugResponse, error) {
	category := entity.DrugCategoryFromCode(req.Category)

	drug, err := s.usecase.CreateDrug(ctx, req.Name, req.Description, req.Price, category, req.UseMethod)
	if err != nil {
		return nil, err
	}

	return s.toResponse(drug), nil
}

func (s *DrugService) GetDrug(ctx context.Context, id int64) (*dto.DrugResponse, error) {
	drug, err := s.usecase.GetDrug(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(drug), nil
}

func (s *DrugService) UpdateDrug(ctx context.Context, id int64, req *dto.UpdateDrugRequest) error {
	return s.usecase.UpdateDrug(ctx, id, req.Description, req.Price, req.UseMethod)
}

func (s *DrugService) UpdateStock(ctx context.Context, id int64, req *dto.UpdateStockRequest) error {
	return s.usecase.UpdateStock(ctx, id, req.Quantity, req.OperatorID)
}

func (s *DrugService) AddStock(ctx context.Context, id int64, req *dto.AddStockRequest) error {
	return s.usecase.AddStock(ctx, id, req.Quantity, req.OperatorID)
}

func (s *DrugService) ReduceStock(ctx context.Context, id int64, req *dto.ReduceStockRequest) error {
	return s.usecase.ReduceStock(ctx, id, req.Quantity, req.OperatorID)
}

func (s *DrugService) CheckAvailability(ctx context.Context, id int64, req *dto.CheckAvailabilityRequest) (*dto.AvailabilityResponse, error) {
	available, err := s.usecase.CheckAvailability(ctx, id, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &dto.AvailabilityResponse{Available: available}, nil
}

func (s *DrugService) ListDrugs(ctx context.Context, page, pageSize int) (*dto.DrugListResponse, error) {
	drugs, total, err := s.usecase.ListDrugs(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.DrugListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Drugs:    s.toResponseList(drugs),
	}, nil
}

func (s *DrugService) SearchDrugs(ctx context.Context, keyword string) ([]*dto.DrugResponse, error) {
	drugs, err := s.usecase.SearchDrugs(ctx, keyword)
	if err != nil {
		return nil, err
	}
	return s.toResponses(drugs), nil
}

func (s *DrugService) DiscontinueDrug(ctx context.Context, id int64) error {
	return s.usecase.DiscontinueDrug(ctx, id)
}

func (s *DrugService) ActivateDrug(ctx context.Context, id int64) error {
	return s.usecase.ActivateDrug(ctx, id)
}

func (s *DrugService) DeleteDrug(ctx context.Context, id int64) error {
	return s.usecase.DeleteDrug(ctx, id)
}

func (s *DrugService) toResponse(drug *aggregate.Drug) *dto.DrugResponse {
	return &dto.DrugResponse{
		ID:              drug.ID(),
		Name:            drug.Name(),
		Description:     drug.Description(),
		Price:           drug.PriceYuan(),
		QuantityInStock: drug.QuantityInStock(),
		Category:        drug.Category().Code(),
		CategoryName:    drug.Category().String(),
		UseMethod:       drug.UseMethod(),
		Status:          drug.Status().Code(),
		StatusName:      drug.Status().String(),
		CreatedAt:       drug.CreatedAt(),
		UpdatedAt:       drug.UpdatedAt(),
	}
}

func (s *DrugService) toResponses(drugs []*aggregate.Drug) []*dto.DrugResponse {
	responses := make([]*dto.DrugResponse, 0, len(drugs))
	for _, drug := range drugs {
		responses = append(responses, s.toResponse(drug))
	}
	return responses
}

func (s *DrugService) toResponseList(drugs []*aggregate.Drug) []dto.DrugResponse {
	responses := make([]dto.DrugResponse, 0, len(drugs))
	for _, drug := range drugs {
		responses = append(responses, *s.toResponse(drug))
	}
	return responses
}