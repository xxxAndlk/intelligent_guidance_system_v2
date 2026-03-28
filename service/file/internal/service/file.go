package service

import (
	"context"

	"intelligent-guidance-system/service/file/internal/biz/dto"
	"intelligent-guidance-system/service/file/internal/biz/usecase"
	"intelligent-guidance-system/service/file/internal/domain/entity"
)

type FileService struct {
	usecase *usecase.FileUsecase
}

func NewFileService(usecase *usecase.FileUsecase) *FileService {
	return &FileService{usecase: usecase}
}

func (s *FileService) UploadFile(ctx context.Context, req *dto.UploadFileRequest) (*dto.FileInfoResponse, error) {
	file, err := s.usecase.UploadFile(ctx, req.Name, req.URL, req.Size, req.MimeType)
	if err != nil {
		return nil, err
	}

	return s.toResponse(file), nil
}

func (s *FileService) GetFileInfo(ctx context.Context, id int64) (*dto.FileInfoResponse, error) {
	file, err := s.usecase.GetFileInfo(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(file), nil
}

func (s *FileService) ListFiles(ctx context.Context, page, pageSize int) (*dto.FileListResponse, error) {
	files, total, err := s.usecase.ListFiles(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.FileListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Files:    s.toResponseList(files),
	}, nil
}

func (s *FileService) DeleteFile(ctx context.Context, id int64) error {
	return s.usecase.DeleteFile(ctx, id)
}

func (s *FileService) toResponse(file *entity.FileInfo) *dto.FileInfoResponse {
	return &dto.FileInfoResponse{
		ID:         file.ID(),
		Name:       file.Name(),
		URL:        file.URL(),
		Size:       file.Size(),
		SizeKB:     file.SizeKB(),
		SizeMB:     file.SizeMB(),
		MimeType:   file.MimeType(),
		IsImage:    file.IsImage(),
		IsPDF:      file.IsPDF(),
		UploadedAt: file.UploadedAt(),
	}
}

func (s *FileService) toResponseList(files []*entity.FileInfo) []dto.FileInfoResponse {
	responses := make([]dto.FileInfoResponse, 0, len(files))
	for _, file := range files {
		responses = append(responses, *s.toResponse(file))
	}
	return responses
}