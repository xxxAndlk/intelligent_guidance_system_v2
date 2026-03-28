package entity

import "time"

type DataPermission struct {
	userID       int64
	resourceType string
	dataScope    DataScope
	deptIDs      []int64
	createdAt    time.Time
	updatedAt    time.Time
}

func NewDataPermission(userID int64, resourceType string, dataScope DataScope, deptIDs []int64) (*DataPermission, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if dataScope == DataScopeUnknown {
		return nil, ErrInvalidDataScope
	}

	now := time.Now()
	return &DataPermission{
		userID:       userID,
		resourceType: resourceType,
		dataScope:    dataScope,
		deptIDs:      deptIDs,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func ReconstructDataPermission(userID int64, resourceType string, dataScope DataScope, deptIDs []int64, createdAt, updatedAt time.Time) *DataPermission {
	return &DataPermission{
		userID:       userID,
		resourceType: resourceType,
		dataScope:    dataScope,
		deptIDs:      deptIDs,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (d *DataPermission) UserID() int64 { return d.userID }
func (d *DataPermission) ResourceType() string { return d.resourceType }
func (d *DataPermission) DataScope() DataScope { return d.dataScope }
func (d *DataPermission) DeptIDs() []int64 { return d.deptIDs }
func (d *DataPermission) CreatedAt() time.Time { return d.createdAt }
func (d *DataPermission) UpdatedAt() time.Time { return d.updatedAt }

func (d *DataPermission) UpdateDataScope(dataScope DataScope) {
	d.dataScope = dataScope
	d.updatedAt = time.Now()
}

func (d *DataPermission) UpdateDeptIDs(deptIDs []int64) {
	d.deptIDs = deptIDs
	d.updatedAt = time.Now()
}

func (d *DataPermission) CanAccess(deptID int64, targetUserID int64) bool {
	switch d.dataScope {
	case DataScopeAll:
		return true
	case DataScopeDept:
		for _, id := range d.deptIDs {
			if id == deptID {
				return true
			}
		}
		return false
	case DataScopeDeptAndSub:
		return false
	case DataScopeOwn:
		return targetUserID == d.userID
	default:
		return false
	}
}