package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OrganizationId *uuid.UUID     `gorm:"type:uuid;index" json:"organization_id"`
	Name           string         `gorm:"type:varchar(50);not null" json:"name"`
	DisplayName    string         `gorm:"type:varchar(100);not null" json:"display_name"`
	Type           string         `gorm:"type:varchar(20);not null;default:'custom'" json:"type"`
	Level          int            `gorm:"not null;default:0" json:"level"`
	Description    string         `gorm:"type:text" json:"description"`
	IsDefault      bool           `gorm:"default:false" json:"is_default"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Organization *Organization `gorm:"foreignKey:OrganizationId" json:"organization,omitempty"`
}

func (Role) TableName() string { return "roles" }

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.Id == uuid.Nil {
		r.Id = uuid.New()
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return
}

func (r *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()
	return
}
