package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidFileID   = errors.New("invalid file ID")
	ErrFileNotFound    = errors.New("file not found")
	ErrEmptyFileName   = errors.New("file name cannot be empty")
	ErrEmptyFileURL    = errors.New("file URL cannot be empty")
)

type FileInfo struct {
	id         int64
	name       string
	url        string
	size       int64
	mimeType   string
	uploadedAt time.Time
}

func NewFileInfo(name, url string, size int64, mimeType string) (*FileInfo, error) {
	if name == "" {
		return nil, ErrEmptyFileName
	}
	if url == "" {
		return nil, ErrEmptyFileURL
	}

	return &FileInfo{
		id:         0,
		name:       name,
		url:        url,
		size:       size,
		mimeType:   mimeType,
		uploadedAt: time.Now(),
	}, nil
}

func ReconstructFileInfo(id int64, name, url string, size int64, mimeType string, uploadedAt time.Time) *FileInfo {
	return &FileInfo{
		id:         id,
		name:       name,
		url:        url,
		size:       size,
		mimeType:   mimeType,
		uploadedAt: uploadedAt,
	}
}

func (f *FileInfo) ID() int64        { return f.id }
func (f *FileInfo) Name() string     { return f.name }
func (f *FileInfo) URL() string      { return f.url }
func (f *FileInfo) Size() int64      { return f.size }
func (f *FileInfo) MimeType() string { return f.mimeType }
func (f *FileInfo) UploadedAt() time.Time { return f.uploadedAt }

func (f *FileInfo) SetID(id int64) { f.id = id }

func (f *FileInfo) SizeKB() float64 {
	return float64(f.size) / 1024.0
}

func (f *FileInfo) SizeMB() float64 {
	return float64(f.size) / (1024.0 * 1024.0)
}

func (f *FileInfo) IsImage() bool {
	return f.mimeType == "image/jpeg" || f.mimeType == "image/png" || f.mimeType == "image/gif"
}

func (f *FileInfo) IsPDF() bool {
	return f.mimeType == "application/pdf"
}