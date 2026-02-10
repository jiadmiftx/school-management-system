package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subject represents a subject/mata pelajaran in a school unit.
type Subject struct {
	Id          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UnitId      uuid.UUID      `gorm:"type:uuid;not null;index" json:"unit_id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"` // "Matematika"
	Code        string         `gorm:"type:varchar(20);not null" json:"code"`  // "MTK"
	Category    string         `gorm:"type:varchar(50)" json:"category"`       // Umum/Jurusan/Mulok
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Unit     *Unit            `gorm:"foreignKey:UnitId" json:"unit,omitempty"`
	Teachers []TeacherSubject `gorm:"foreignKey:SubjectId" json:"teachers,omitempty"`
}

func (Subject) TableName() string { return "subjects" }

func (s *Subject) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Id == uuid.Nil {
		s.Id = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return
}

func (s *Subject) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}
