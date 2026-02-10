package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeacherSubject represents the many-to-many relationship between teachers and subjects.
type TeacherSubject struct {
	Id               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	TeacherProfileId uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_profile_id"`
	SubjectId        uuid.UUID      `gorm:"type:uuid;not null;index" json:"subject_id"`
	IsPrimary        bool           `gorm:"default:false" json:"is_primary"` // Mapel utama
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	TeacherProfile *TeacherProfile `gorm:"foreignKey:TeacherProfileId" json:"teacher_profile,omitempty"`
	Subject        *Subject        `gorm:"foreignKey:SubjectId" json:"subject,omitempty"`
}

func (TeacherSubject) TableName() string { return "teacher_subjects" }

func (ts *TeacherSubject) BeforeCreate(tx *gorm.DB) (err error) {
	if ts.Id == uuid.Nil {
		ts.Id = uuid.New()
	}
	ts.CreatedAt = time.Now()
	return
}
