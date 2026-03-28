package mysql

import "time"

type FileInfoPO struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"column:name;size:200"`
	URL        string    `gorm:"column:url;size:500"`
	Size       int64     `gorm:"column:size"`
	MimeType   string    `gorm:"column:mime_type;size:100"`
	UploadedAt time.Time `gorm:"column:uploaded_at"`
}

func (FileInfoPO) TableName() string {
	return "file_infos"
}