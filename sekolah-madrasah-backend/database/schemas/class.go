package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Class represents a classroom/rombel within a school unit.
// Example: "X IPA 1", "VII A"
type Class struct {
	Id                uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UnitId            uuid.UUID      `gorm:"type:uuid;not null;index" json:"unit_id"`        // School
	Name              string         `gorm:"type:varchar(50);not null" json:"name"`          // "X IPA 1", "VII A"
	Level             int            `gorm:"not null" json:"level"`                          // Tingkat (1-12)
	AcademicYear      string         `gorm:"type:varchar(20);not null" json:"academic_year"` // "2025/2026"
	HomeroomTeacherId *uuid.UUID     `gorm:"type:uuid;index" json:"homeroom_teacher_id"`     // Wali kelas (nullable)
	Capacity          int            `gorm:"default:30" json:"capacity"`                     // Kapasitas maksimal
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	Unit            *Unit           `gorm:"foreignKey:UnitId" json:"unit,omitempty"`
	HomeroomTeacher *TeacherProfile `gorm:"foreignKey:HomeroomTeacherId" json:"homeroom_teacher,omitempty"`
}

func (Class) TableName() string { return "classes" }

func (c *Class) BeforeCreate(tx *gorm.DB) (err error) {
	if c.Id == uuid.Nil {
		c.Id = uuid.New()
	}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Class) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}
