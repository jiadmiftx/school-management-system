package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationMember struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserId         uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	OrganizationId uuid.UUID      `gorm:"type:uuid;not null;index" json:"organization_id"`
	RoleId         uuid.UUID      `gorm:"type:uuid;not null;index" json:"role_id"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	JoinedAt       time.Time      `json:"joined_at"`
	InvitedBy      *uuid.UUID     `gorm:"type:uuid" json:"invited_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	User         *User         `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Organization *Organization `gorm:"foreignKey:OrganizationId" json:"organization,omitempty"`
	Role         *Role         `gorm:"foreignKey:RoleId" json:"role,omitempty"`
	Inviter      *User         `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}

func (OrganizationMember) TableName() string { return "organization_members" }

func (om *OrganizationMember) BeforeCreate(tx *gorm.DB) (err error) {
	if om.Id == uuid.Nil {
		om.Id = uuid.New()
	}
	om.JoinedAt = time.Now()
	om.CreatedAt = time.Now()
	om.UpdatedAt = time.Now()
	return
}

func (om *OrganizationMember) BeforeUpdate(tx *gorm.DB) (err error) {
	om.UpdatedAt = time.Now()
	return
}
