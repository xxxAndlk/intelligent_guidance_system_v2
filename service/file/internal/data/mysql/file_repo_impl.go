package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"intelligent-guidance-system/service/file/internal/domain/entity"
	"intelligent-guidance-system/service/file/internal/domain/repository"
)

type FileRepoImpl struct {
	db *gorm.DB
}

func NewFileRepoImpl(db *gorm.DB) repository.FileRepository {
	return &FileRepoImpl{db: db}
}

func (r *FileRepoImpl) Save(ctx context.Context, file *entity.FileInfo) error {
	po := r.toPO(file)

	if po.ID == 0 {
		if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
			return err
		}
		file.SetID(po.ID)
	} else {
		if err := r.db.WithContext(ctx).Save(&po).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *FileRepoImpl) FindByID(ctx context.Context, id int64) (*entity.FileInfo, error) {
	var po FileInfoPO
	if err := r.db.WithContext(ctx).First(&po, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return r.toEntity(&po), nil
}

func (r *FileRepoImpl) FindAll(ctx context.Context, page, pageSize int) ([]*entity.FileInfo, int64, error) {
	var pos []FileInfoPO
	var total int64

	if err := r.db.WithContext(ctx).Model(&FileInfoPO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Order("uploaded_at DESC").Find(&pos).Error; err != nil {
		return nil, 0, err
	}

	return r.toEntities(pos), total, nil
}

func (r *FileRepoImpl) FindByMimeType(ctx context.Context, mimeType string) ([]*entity.FileInfo, error) {
	var pos []FileInfoPO
	if err := r.db.WithContext(ctx).Where("mime_type = ?", mimeType).Find(&pos).Error; err != nil {
		return nil, err
	}
	return r.toEntities(pos), nil
}

func (r *FileRepoImpl) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&FileInfoPO{}, id).Error
}

func (r *FileRepoImpl) toPO(file *entity.FileInfo) *FileInfoPO {
	return &FileInfoPO{
		ID:         file.ID(),
		Name:       file.Name(),
		URL:        file.URL(),
		Size:       file.Size(),
		MimeType:   file.MimeType(),
		UploadedAt: file.UploadedAt(),
	}
}

func (r *FileRepoImpl) toEntity(po *FileInfoPO) *entity.FileInfo {
	return entity.ReconstructFileInfo(
		po.ID,
		po.Name,
		po.URL,
		po.Size,
		po.MimeType,
		po.UploadedAt,
	)
}

func (r *FileRepoImpl) toEntities(pos []FileInfoPO) []*entity.FileInfo {
	entities := make([]*entity.FileInfo, 0, len(pos))
	for _, po := range pos {
		entities = append(entities, r.toEntity(&po))
	}
	return entities
}