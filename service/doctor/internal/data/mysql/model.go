package mysql

import (
	"time"
)

type DoctorPO struct {
	ID             string    `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	EmployeeID     string    `gorm:"uniqueIndex;type:varchar(50);comment:工号"`
	Name           string    `gorm:"type:varchar(100);comment:姓名"`
	Gender         string    `gorm:"type:varchar(10);comment:性别"`
	BirthDate      time.Time `gorm:"type:date;comment:出生日期"`
	Phone          string    `gorm:"type:varchar(20);comment:手机号"`
	Email          string    `gorm:"type:varchar(100);comment:邮箱"`
	IDCard         string    `gorm:"type:varchar(20);comment:身份证号"`
	Address        string    `gorm:"type:varchar(200);comment:地址"`
	AvatarURL      string    `gorm:"type:varchar(500);comment:头像URL"`
	LicenseNumber  string    `gorm:"uniqueIndex;type:varchar(50);comment:执业证号"`
	Position       string    `gorm:"type:varchar(20);comment:职位代码"`
	Title          string    `gorm:"type:varchar(20);comment:职称代码"`
	Specialties    string    `gorm:"type:text;comment:专业特长(JSON数组)"`
	Education      string    `gorm:"type:varchar(20);comment:学历代码"`
	GraduationSchool string `gorm:"type:varchar(100);comment:毕业院校"`
	GraduationYear int       `gorm:"type:int;comment:毕业年份"`
	YearsOfExp     int       `gorm:"type:int;comment:从业年限"`
	Certifications string    `gorm:"type:text;comment:资质证书(JSON数组)"`
	Introduction   string    `gorm:"type:text;comment:简介"`
	DepartmentID   string    `gorm:"type:varchar(36);index;comment:科室ID"`
	Status         string    `gorm:"type:varchar(20);index;comment:状态"`
	CreatedAt      time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime;comment:更新时间"`
	Version        int       `gorm:"type:int;default:1;comment:版本号"`
}

func (DoctorPO) TableName() string {
	return "doctors"
}

type SchedulePO struct {
	ID        string    `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	DoctorID  string    `gorm:"type:varchar(36);index;comment:医生ID"`
	Date      time.Time `gorm:"type:date;uniqueIndex:idx_doctor_date;comment:排班日期"`
	Status    string    `gorm:"type:varchar(20);index;comment:状态"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

func (SchedulePO) TableName() string {
	return "doctor_schedules"
}

type TimeSlotPO struct {
	ID          string    `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	ScheduleID  string    `gorm:"type:varchar(36);index;comment:排班ID"`
	StartTime   time.Time `gorm:"type:datetime;comment:开始时间"`
	EndTime     time.Time `gorm:"type:datetime;comment:结束时间"`
	IsAvailable bool      `gorm:"type:tinyint(1);comment:是否可用"`
	MaxPatients int       `gorm:"type:int;comment:最大患者数"`
	BookedCount int       `gorm:"type:int;default:0;comment:已预约数"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

func (TimeSlotPO) TableName() string {
	return "doctor_time_slots"
}

type RolePO struct {
	ID          string    `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	Name        string    `gorm:"type:varchar(50);comment:角色名称"`
	Code        string    `gorm:"uniqueIndex;type:varchar(50);comment:角色代码"`
	Description string    `gorm:"type:varchar(200);comment:角色描述"`
	RoleType    string    `gorm:"type:varchar(20);index;comment:角色类型"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

func (RolePO) TableName() string {
	return "doctor_roles"
}

type RolePermissionPO struct {
	ID       string `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	RoleID   string `gorm:"type:varchar(36);index;comment:角色ID"`
	PermCode string `gorm:"type:varchar(50);comment:权限代码"`
	PermName string `gorm:"type:varchar(100);comment:权限名称"`
	Resource string `gorm:"type:varchar(100);comment:资源"`
	Action   string `gorm:"type:varchar(50);comment:操作"`
}

func (RolePermissionPO) TableName() string {
	return "doctor_role_permissions"
}

type DoctorRoleAssignmentPO struct {
	ID           string     `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	DoctorID     string     `gorm:"type:varchar(36);index;comment:医生ID"`
	RoleID       string     `gorm:"type:varchar(36);index;comment:角色ID"`
	DepartmentID string     `gorm:"type:varchar(36);comment:科室ID"`
	AssignedBy   string     `gorm:"type:varchar(36);comment:分配人ID"`
	AssignedAt   time.Time  `gorm:"type:datetime;comment:分配时间"`
	ExpiresAt    *time.Time `gorm:"type:datetime;comment:过期时间"`
	IsActive     bool       `gorm:"type:tinyint(1);default:1;comment:是否激活"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;comment:更新时间"`
}

func (DoctorRoleAssignmentPO) TableName() string {
	return "doctor_role_assignments"
}

type DoctorEventPO struct {
	ID           string    `gorm:"primaryKey;type:varchar(36);comment:主键ID"`
	EventType    string    `gorm:"type:varchar(50);index;comment:事件类型"`
	DoctorID     string    `gorm:"type:varchar(36);index;comment:医生ID"`
	EmployeeID   string    `gorm:"type:varchar(50);comment:工号"`
	DepartmentID string    `gorm:"type:varchar(36);comment:科室ID"`
	Payload      string    `gorm:"type:text;comment:事件数据(JSON)"`
	OccurredAt   time.Time `gorm:"type:datetime;index;comment:发生时间"`
	CreatedAt    time.Time `gorm:"autoCreateTime;comment:创建时间"`
}

func (DoctorEventPO) TableName() string {
	return "doctor_events"
}