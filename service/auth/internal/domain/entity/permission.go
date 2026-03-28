package entity

import (
	"strings"
	"time"
)

type ResourceType int

const (
	ResourceTypeUnknown ResourceType = iota
	ResourceTypeMenu
	ResourceTypeButton
	ResourceTypeAPI
)

func (t ResourceType) String() string {
	switch t {
	case ResourceTypeMenu:
		return "菜单"
	case ResourceTypeButton:
		return "按钮"
	case ResourceTypeAPI:
		return "接口"
	default:
		return "未知"
	}
}

func ResourceTypeFromCode(code string) ResourceType {
	switch code {
	case "MENU":
		return ResourceTypeMenu
	case "BUTTON":
		return ResourceTypeButton
	case "API":
		return ResourceTypeAPI
	default:
		return ResourceTypeUnknown
	}
}

func (t ResourceType) Code() string {
	switch t {
	case ResourceTypeMenu:
		return "MENU"
	case ResourceTypeButton:
		return "BUTTON"
	case ResourceTypeAPI:
		return "API"
	default:
		return "UNKNOWN"
	}
}

type Permission struct {
	id           int64
	code         string
	name         string
	resourceType ResourceType
	resourceURL  string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewPermission(code, name string, resourceType ResourceType, resourceURL string) (*Permission, error) {
	if code == "" {
		return nil, ErrInvalidPermission
	}
	if name == "" {
		return nil, ErrInvalidPermission
	}

	now := time.Now()
	return &Permission{
		id:           0,
		code:         code,
		name:         name,
		resourceType: resourceType,
		resourceURL:  resourceURL,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func ReconstructPermission(id int64, code, name string, resourceType ResourceType, resourceURL string, createdAt, updatedAt time.Time) *Permission {
	return &Permission{
		id:           id,
		code:         code,
		name:         name,
		resourceType: resourceType,
		resourceURL:  resourceURL,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (p *Permission) ID() int64 { return p.id }
func (p *Permission) Code() string { return p.code }
func (p *Permission) Name() string { return p.name }
func (p *Permission) ResourceType() ResourceType { return p.resourceType }
func (p *Permission) ResourceURL() string { return p.resourceURL }
func (p *Permission) CreatedAt() time.Time { return p.createdAt }
func (p *Permission) UpdatedAt() time.Time { return p.updatedAt }

func (p *Permission) SetID(id int64) { p.id = id }

func (p *Permission) MatchResource(url string) bool {
	if p.resourceURL == "*" {
		return true
	}
	if strings.HasSuffix(p.resourceURL, "*") {
		prefix := strings.TrimSuffix(p.resourceURL, "*")
		return strings.HasPrefix(url, prefix)
	}
	return p.resourceURL == url
}