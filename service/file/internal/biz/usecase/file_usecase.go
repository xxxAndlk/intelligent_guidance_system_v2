package usecase

import (
	"context"
	"errors"

	"intelligent-guidance-system/service/file/internal/domain/entity"
	"intelligent-guidance-system/service/file/internal/domain/repository"
)

var (
	ErrFileNotFound = errors.New("file not found")
)

type FileUsecase struct {
	repo repository.FileRepository
}

func NewFileUsecase(repo repository.FileRepository) *FileUsecase {
	return &FileUsecase{repo: repo}
}

func (u *FileUsecase) UploadFile(ctx context.Context, name, url string, size int64, mimeType string) (*entity.FileInfo, error) {
	file, err := entity.NewFileInfo(name, url, size, mimeType)
	if err != nil {
		return nil, err
	}

	if err := u.repo.Save(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}

func (u *FileUsecase) GetFileInfo(ctx context.Context, id int64) (*entity.FileInfo, error) {
	file, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

func (u *FileUsecase) ListFiles(ctx context.Context, page, pageSize int) ([]*entity.FileInfo, int64, error) {
	return u.repo.FindAll(ctx, page, pageSize)
}

func (u *FileUsecase) DeleteFile(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}