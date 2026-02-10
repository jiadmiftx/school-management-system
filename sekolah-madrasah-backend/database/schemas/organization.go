package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents a Yayasan or Educational Institution.
// Conceptually: Organization = Yayasan (Foundation)
// Example: "Yayasan Pendidikan Islam" (code: YPI-001)
// Contains multiple Schools/Units (stored in units table)
type Organization struct {
	Id          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OwnerId     uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Type        string         `gorm:"type:varchar(50);not null" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	Address     string         `gorm:"type:text" json:"address"`
	Logo        string         `gorm:"type:varchar(500)" json:"logo"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Settings    string         `gorm:"type:jsonb;default:'{}'" json:"settings"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Owner *User `gorm:"foreignKey:OwnerId" json:"owner,omitempty"`
}

func (Organization) TableName() string { return "organizations" }

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == uuid.Nil {
		o.Id = uuid.New()
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return
}

func (o *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now()
	return
}
