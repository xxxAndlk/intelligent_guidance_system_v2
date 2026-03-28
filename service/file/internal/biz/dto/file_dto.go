package dto

import "time"

type UploadFileRequest struct {
	Name     string `form:"name" binding:"required"`
	URL      string `form:"url" binding:"required"`
	Size     int64  `form:"size" binding:"required"`
	MimeType string `form:"mime_type"`
}

type FileInfoResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Size       int64     `json:"size"`
	SizeKB     float64   `json:"size_kb"`
	SizeMB     float64   `json:"size_mb"`
	MimeType   string    `json:"mime_type"`
	IsImage    bool      `json:"is_image"`
	IsPDF      bool      `json:"is_pdf"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type FileListResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Files    []FileInfoResponse `json:"files"`
}