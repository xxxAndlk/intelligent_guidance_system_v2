package repository

import (
	"context"

	"intelligent-guidance-system/service/file/internal/domain/entity"
)

type FileRepository interface {
	Save(ctx context.Context, file *entity.FileInfo) error
	FindByID(ctx context.Context, id int64) (*entity.FileInfo, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.FileInfo, int64, error)
	FindByMimeType(ctx context.Context, mimeType string) ([]*entity.FileInfo, error)
	Delete(ctx context.Context, id int64) error
}