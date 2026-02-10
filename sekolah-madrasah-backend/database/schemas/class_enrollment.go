package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ClassEnrollmentStatus represents the status of a student in a class
type ClassEnrollmentStatus string

const (
	EnrollmentStatusActive      ClassEnrollmentStatus = "active"
	EnrollmentStatusGraduated   ClassEnrollmentStatus = "graduated"
	EnrollmentStatusTransferred ClassEnrollmentStatus = "transferred"
	EnrollmentStatusDropped     ClassEnrollmentStatus = "dropped"
)

// ClassEnrollment represents the relationship between students and classes.
// Tracks which class a student is enrolled in for a specific academic year.
type ClassEnrollment struct {
	Id               uuid.UUID             `gorm:"type:uuid;primaryKey" json:"id"`
	StudentProfileId uuid.UUID             `gorm:"type:uuid;not null;index" json:"student_profile_id"` // FK to student_profiles
	ClassId          uuid.UUID             `gorm:"type:uuid;not null;index" json:"class_id"`           // FK to classes
	AcademicYear     string                `gorm:"type:varchar(20);not null" json:"academic_year"`     // "2025/2026"
	Status           ClassEnrollmentStatus `gorm:"type:varchar(20);default:'active'" json:"status"`    // active/graduated/transferred
	EnrolledAt       time.Time             `gorm:"type:date;not null" json:"enrolled_at"`              // Tanggal masuk kelas
	LeftAt           *time.Time            `gorm:"type:date" json:"left_at"`                           // Tanggal keluar (nullable)
	Notes            *string               `gorm:"type:text" json:"notes"`                             // Catatan
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`

	StudentProfile *StudentProfile `gorm:"foreignKey:StudentProfileId" json:"student_profile,omitempty"`
	Class          *Class          `gorm:"foreignKey:ClassId" json:"class,omitempty"`
}

func (ClassEnrollment) TableName() string { return "class_enrollments" }

func (ce *ClassEnrollment) BeforeCreate(tx *gorm.DB) (err error) {
	if ce.Id == uuid.Nil {
		ce.Id = uuid.New()
	}
	ce.CreatedAt = time.Now()
	ce.UpdatedAt = time.Now()
	return
}

func (ce *ClassEnrollment) BeforeUpdate(tx *gorm.DB) (err error) {
	ce.UpdatedAt = time.Now()
	return
}
