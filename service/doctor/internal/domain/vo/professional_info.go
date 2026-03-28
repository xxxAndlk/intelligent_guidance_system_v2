package vo

import (
	"errors"
	"time"
)

var (
	ErrInvalidLicense   = errors.New("invalid medical license number")
	ErrInvalidEducation = errors.New("invalid education level")
	ErrInvalidYears     = errors.New("invalid years of experience")
)

// EducationLevel represents the education level
type EducationLevel int

const (
	EducationUnknown EducationLevel = iota
	EducationJunior         // 大专
	EducationBachelor       // 本科
	EducationMaster         // 硕士
	EducationDoctorate      // 博士
)

func (e EducationLevel) String() string {
	switch e {
	case EducationJunior:
		return "大专"
	case EducationBachelor:
		return "本科"
	case EducationMaster:
		return "硕士"
	case EducationDoctorate:
		return "博士"
	default:
		return "未知"
	}
}

func EducationLevelFromCode(code string) EducationLevel {
	switch code {
	case "JUNIOR":
		return EducationJunior
	case "BACHELOR":
		return EducationBachelor
	case "MASTER":
		return EducationMaster
	case "DOCTORATE":
		return EducationDoctorate
	default:
		return EducationUnknown
	}
}

func (e EducationLevel) Code() string {
	switch e {
	case EducationJunior:
		return "JUNIOR"
	case EducationBachelor:
		return "BACHELOR"
	case EducationMaster:
		return "MASTER"
	case EducationDoctorate:
		return "DOCTORATE"
	default:
		return "UNKNOWN"
	}
}

// ProfessionalInfo represents the professional information of a doctor
type ProfessionalInfo struct {
	licenseNumber    string
	position         Position
	title            Title
	specialties      []string
	education        EducationLevel
	graduationSchool string
	graduationYear   int
	yearsOfExperience int
	certifications   []string
	introduction     string
}

// NewProfessionalInfo creates a new ProfessionalInfo with validation
func NewProfessionalInfo(
	licenseNumber string,
	position Position,
	title Title,
	specialties []string,
	education EducationLevel,
	graduationSchool string,
	graduationYear int,
	yearsOfExperience int,
	certifications []string,
	introduction string,
) (*ProfessionalInfo, error) {
	if licenseNumber == "" {
		return nil, ErrInvalidLicense
	}
	if yearsOfExperience < 0 {
		return nil, ErrInvalidYears
	}

	// Remove duplicates from specialties and certifications
	specialties = uniqueStrings(specialties)
	certifications = uniqueStrings(certifications)

	return &ProfessionalInfo{
		licenseNumber:     licenseNumber,
		position:          position,
		title:             title,
		specialties:       specialties,
		education:         education,
		graduationSchool:  graduationSchool,
		graduationYear:    graduationYear,
		yearsOfExperience: yearsOfExperience,
		certifications:    certifications,
		introduction:      introduction,
	}, nil
}

// Getters
func (p *ProfessionalInfo) LicenseNumber() string        { return p.licenseNumber }
func (p *ProfessionalInfo) Position() Position           { return p.position }
func (p *ProfessionalInfo) Title() Title                 { return p.title }
func (p *ProfessionalInfo) Specialties() []string        { return p.specialties }
func (p *ProfessionalInfo) Education() EducationLevel    { return p.education }
func (p *ProfessionalInfo) GraduationSchool() string     { return p.graduationSchool }
func (p *ProfessionalInfo) GraduationYear() int          { return p.graduationYear }
func (p *ProfessionalInfo) YearsOfExperience() int      { return p.yearsOfExperience }
func (p *ProfessionalInfo) Certifications() []string     { return p.certifications }
func (p *ProfessionalInfo) Introduction() string         { return p.introduction }

// UpdatePosition updates the position
func (p *ProfessionalInfo) UpdatePosition(position Position) {
	p.position = position
}

// UpdateTitle updates the title
func (p *ProfessionalInfo) UpdateTitle(title Title) {
	p.title = title
}

// AddSpecialty adds a specialty
func (p *ProfessionalInfo) AddSpecialty(specialty string) {
	for _, s := range p.specialties {
		if s == specialty {
			return
		}
	}
	p.specialties = append(p.specialties, specialty)
}

// RemoveSpecialty removes a specialty
func (p *ProfessionalInfo) RemoveSpecialty(specialty string) {
	for i, s := range p.specialties {
		if s == specialty {
			p.specialties = append(p.specialties[:i], p.specialties[i+1:]...)
			return
		}
	}
}

// UpdateEducation updates education information
func (p *ProfessionalInfo) UpdateEducation(education EducationLevel, school string, year int) {
	p.education = education
	p.graduationSchool = school
	p.graduationYear = year
}

// UpdateYearsOfExperience updates years of experience
func (p *ProfessionalInfo) UpdateYearsOfExperience(years int) error {
	if years < 0 {
		return ErrInvalidYears
	}
	p.yearsOfExperience = years
	return nil
}

// AddCertification adds a certification
func (p *ProfessionalInfo) AddCertification(cert string) {
	for _, c := range p.certifications {
		if c == cert {
			return
		}
	}
	p.certifications = append(p.certifications, cert)
}

// RemoveCertification removes a certification
func (p *ProfessionalInfo) RemoveCertification(cert string) {
	for i, c := range p.certifications {
		if c == cert {
			p.certifications = append(p.certifications[:i], p.certifications[i+1:]...)
			return
		}
	}
}

// UpdateIntroduction updates the introduction
func (p *ProfessionalInfo) UpdateIntroduction(intro string) {
	p.introduction = intro
}

// IncrementExperience increments years of experience (usually called annually)
func (p *ProfessionalInfo) IncrementExperience() {
	p.yearsOfExperience++
}

// IsExpert returns true if the doctor is at expert level
func (p *ProfessionalInfo) IsExpert() bool {
	return p.position.IsExpert()
}

// CanTreat returns true if the doctor can independently treat patients
func (p *ProfessionalInfo) CanTreat() bool {
	return p.position.CanTreat()
}

// HasSpecialty checks if the doctor has a specific specialty
func (p *ProfessionalInfo) HasSpecialty(specialty string) bool {
	for _, s := range p.specialties {
		if s == specialty {
			return true
		}
	}
	return false
}

// HasCertification checks if the doctor has a specific certification
func (p *ProfessionalInfo) HasCertification(cert string) bool {
	for _, c := range p.certifications {
		if c == cert {
			return true
		}
	}
	return false
}

// CalculatePracticeYears calculates years since graduation
func (p *ProfessionalInfo) CalculatePracticeYears() int {
	currentYear := time.Now().Year()
	return currentYear - p.graduationYear
}

// Helper function to remove duplicate strings
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(strs))
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}