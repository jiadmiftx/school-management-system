package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActivityTeacher represents teachers assigned to an activity (pembina/pengisi).
type ActivityTeacher struct {
	Id               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ActivityId       uuid.UUID      `gorm:"type:uuid;not null;index" json:"activity_id"`
	TeacherProfileId uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_profile_id"`
	Role             string         `gorm:"type:varchar(50);default:'pembina'" json:"role"` // pembina/pengisi/koordinator
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Activity       *Activity       `gorm:"foreignKey:ActivityId" json:"activity,omitempty"`
	TeacherProfile *TeacherProfile `gorm:"foreignKey:TeacherProfileId" json:"teacher_profile,omitempty"`
}

func (ActivityTeacher) TableName() string { return "activity_teachers" }

func (at *ActivityTeacher) BeforeCreate(tx *gorm.DB) (err error) {
	if at.Id == uuid.Nil {
		at.Id = uuid.New()
	}
	at.CreatedAt = time.Now()
	return
}
