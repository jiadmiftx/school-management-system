package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ActivityStudent represents students enrolled in an activity.
type ActivityStudent struct {
	Id               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	ActivityId       uuid.UUID      `gorm:"type:uuid;not null;index" json:"activity_id"`
	StudentProfileId uuid.UUID      `gorm:"type:uuid;not null;index" json:"student_profile_id"`
	IsMandatory      bool           `gorm:"default:false" json:"is_mandatory"` // Wajib/Pilihan
	JoinedAt         *time.Time     `gorm:"type:date" json:"joined_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	Activity       *Activity       `gorm:"foreignKey:ActivityId" json:"activity,omitempty"`
	StudentProfile *StudentProfile `gorm:"foreignKey:StudentProfileId" json:"student_profile,omitempty"`
}

func (ActivityStudent) TableName() string { return "activity_students" }

func (as *ActivityStudent) BeforeCreate(tx *gorm.DB) (err error) {
	if as.Id == uuid.Nil {
		as.Id = uuid.New()
	}
	as.CreatedAt = time.Now()
	return
}
