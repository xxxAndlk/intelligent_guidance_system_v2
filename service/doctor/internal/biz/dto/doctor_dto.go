package dto

import (
	"time"
)

type DoctorCreateRequest struct {
	EmployeeID       string    `json:"employee_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birth_date"`
	Phone            string    `json:"phone" validate:"omitempty,phone"`
	Email            string    `json:"email" validate:"omitempty,email"`
	IDCard           string    `json:"id_card"`
	Address          string    `json:"address"`
	AvatarURL        string    `json:"avatar_url"`
	LicenseNumber    string    `json:"license_number" validate:"required"`
	Position         string    `json:"position" validate:"required"`
	Title            string    `json:"title"`
	Specialties      []string  `json:"specialties"`
	Education        string    `json:"education"`
	GraduationSchool string    `json:"graduation_school"`
	GraduationYear   int       `json:"graduation_year"`
	YearsOfExp       int       `json:"years_of_exp" validate:"min=0"`
	Certifications   []string  `json:"certifications"`
	Introduction     string    `json:"introduction"`
	DepartmentID     string    `json:"department_id" validate:"required"`
}

type DoctorUpdateRequest struct {
	ID               string    `json:"id" validate:"required"`
	Name             string    `json:"name" validate:"omitempty"`
	Gender           string    `json:"gender"`
	BirthDate        time.Time `json:"birth_date"`
	Phone            string    `json:"phone" validate:"omitempty,phone"`
	Email            string    `json:"email" validate:"omitempty,email"`
	Address          string    `json:"address"`
	AvatarURL        string    `json:"avatar_url"`
	Position         string    `json:"position"`
	Title            string    `json:"title"`
	Specialties      []string  `json:"specialties"`
	Education        string    `json:"education"`
	GraduationSchool string    `json:"graduation_school"`
	GraduationYear   int       `json:"graduation_year"`
	YearsOfExp       int       `json:"years_of_exp" validate:"min=0"`
	Certifications   []string  `json:"certifications"`
	Introduction     string    `json:"introduction"`
}

type DoctorDepartmentChangeRequest struct {
	DoctorID     string `json:"doctor_id" validate:"required"`
	NewDeptID    string `json:"new_department_id" validate:"required"`
	Reason       string `json:"reason"`
}

type DoctorStatusChangeRequest struct {
	DoctorID string `json:"doctor_id" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=ACTIVE INACTIVE ON_LEAVE SUSPENDED RETIRED"`
	Reason   string `json:"reason"`
}

type DoctorRoleAssignRequest struct {
	DoctorID     string     `json:"doctor_id" validate:"required"`
	RoleID       string     `json:"role_id" validate:"required"`
	DepartmentID string     `json:"department_id"`
	AssignedBy   string     `json:"assigned_by" validate:"required"`
	ExpiresAt    *time.Time `json:"expires_at"`
}

type DoctorRoleRemoveRequest struct {
	DoctorID string `json:"doctor_id" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
}

type ScheduleSetRequest struct {
	DoctorID  string           `json:"doctor_id" validate:"required"`
	Date      time.Time        `json:"date" validate:"required"`
	TimeSlots []TimeSlotInput  `json:"time_slots" validate:"required,dive"`
}

type TimeSlotInput struct {
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	MaxPatients int       `json:"max_patients" validate:"min=1"`
}

type ScheduleAddSlotRequest struct {
	DoctorID    string    `json:"doctor_id" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	MaxPatients int       `json:"max_patients" validate:"min=1"`
}

type SchedulePublishRequest struct {
	DoctorID string    `json:"doctor_id" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type ScheduleCancelRequest struct {
	DoctorID string    `json:"doctor_id" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type DoctorFilterRequest struct {
	DepartmentID string `json:"department_id" form:"department_id"`
	Status       string `json:"status" form:"status"`
	Position     string `json:"position" form:"position"`
	Specialty    string `json:"specialty" form:"specialty"`
	Name         string `json:"name" form:"name"`
	EmployeeID   string `json:"employee_id" form:"employee_id"`
	IsExpert     bool   `json:"is_expert" form:"is_expert"`
	Page         int    `json:"page" form:"page" validate:"min=1"`
	PageSize     int    `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

type DoctorBasicInfoDTO struct {
	Name      string    `json:"name"`
	Gender    string    `json:"gender"`
	BirthDate time.Time `json:"birth_date"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	IDCard    string    `json:"id_card"`
	Address   string    `json:"address"`
	AvatarURL string    `json:"avatar_url"`
}

type ProfessionalInfoDTO struct {
	LicenseNumber     string   `json:"license_number"`
	Position          string   `json:"position"`
	PositionDisplay   string   `json:"position_display"`
	Title             string   `json:"title"`
	TitleDisplay      string   `json:"title_display"`
	Specialties       []string `json:"specialties"`
	Education         string   `json:"education"`
	EducationDisplay  string   `json:"education_display"`
	GraduationSchool  string   `json:"graduation_school"`
	GraduationYear    int      `json:"graduation_year"`
	YearsOfExperience int      `json:"years_of_experience"`
	Certifications    []string `json:"certifications"`
	Introduction      string   `json:"introduction"`
	IsExpert          bool     `json:"is_expert"`
	CanTreat          bool     `json:"can_treat"`
}

type DoctorDTO struct {
	ID             string              `json:"id"`
	EmployeeID     string              `json:"employee_id"`
	BasicInfo      *DoctorBasicInfoDTO `json:"basic_info"`
	Professional   *ProfessionalInfoDTO `json:"professional"`
	DepartmentID   string              `json:"department_id"`
	Status         string              `json:"status"`
	StatusDisplay  string              `json:"status_display"`
	IsExpert       bool                `json:"is_expert"`
	CanTreat       bool                `json:"can_treat"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

type DoctorDetailDTO struct {
	*DoctorDTO
	Schedules []*ScheduleDTO      `json:"schedules"`
	Roles     []*RoleAssignmentDTO `json:"roles"`
}

type TimeSlotDTO struct {
	ID           string    `json:"id"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	IsAvailable  bool      `json:"is_available"`
	MaxPatients  int       `json:"max_patients"`
	BookedCount  int       `json:"booked_count"`
	AvailableSlots int     `json:"available_slots"`
	Duration     int       `json:"duration_minutes"`
}

type ScheduleDTO struct {
	ID           string         `json:"id"`
	DoctorID     string         `json:"doctor_id"`
	Date         time.Time      `json:"date"`
	TimeSlots    []*TimeSlotDTO `json:"time_slots"`
	Status       string         `json:"status"`
	StatusDisplay string        `json:"status_display"`
	TotalCapacity int           `json:"total_capacity"`
	TotalBooked  int            `json:"total_booked"`
	AvailableCapacity int       `json:"available_capacity"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type RoleDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	TypeDisplay string   `json:"type_display"`
	Permissions []string `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

type RoleAssignmentDTO struct {
	ID           string     `json:"id"`
	RoleID       string     `json:"role_id"`
	RoleName     string     `json:"role_name"`
	RoleCode     string     `json:"role_code"`
	DepartmentID string     `json:"department_id"`
	AssignedBy   string     `json:"assigned_by"`
	AssignedAt   time.Time  `json:"assigned_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	IsActive     bool       `json:"is_active"`
	IsExpired    bool       `json:"is_expired"`
	IsValid      bool       `json:"is_valid"`
}

type DoctorListResponse struct {
	Items    []*DoctorDTO `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
 PageSize int          `json:"page_size"`
}

type DoctorDetailResponse struct {
	Data *DoctorDetailDTO `json:"data"`
}

type ScheduleResponse struct {
	Data *ScheduleDTO `json:"data"`
}

type ScheduleListResponse struct {
	Items []*ScheduleDTO `json:"items"`
	Total int64         `json:"total"`
}

type AvailableSlotsResponse struct {
	DoctorID string         `json:"doctor_id"`
	Date     time.Time      `json:"date"`
	Slots    []*TimeSlotDTO `json:"slots"`
}