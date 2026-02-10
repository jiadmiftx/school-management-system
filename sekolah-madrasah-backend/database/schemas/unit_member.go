package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnitMemberRole defines the role of a user within a unit
type UnitMemberRole string

const (
	UnitMemberRoleOwner    UnitMemberRole = "owner"
	UnitMemberRoleAdmin    UnitMemberRole = "admin"
	UnitMemberRolePengurus UnitMemberRole = "pengurus"
	UnitMemberRoleStaff    UnitMemberRole = "staff"
	UnitMemberRoleParent   UnitMemberRole = "parent"
	UnitMemberRoleAnggota  UnitMemberRole = "anggota" // siswa/member
)

type UnitMember struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	UnitId    uuid.UUID      `gorm:"type:uuid;not null;index" json:"unit_id"`
	Role      UnitMemberRole `gorm:"type:varchar(20);not null;default:'staff'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	JoinedAt  time.Time      `json:"joined_at"`
	InvitedBy *uuid.UUID     `gorm:"type:uuid" json:"invited_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User    *User `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Unit    *Unit `gorm:"foreignKey:UnitId" json:"unit,omitempty"`
	Inviter *User `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
}

func (UnitMember) TableName() string { return "unit_members" }

func (sm *UnitMember) BeforeCreate(tx *gorm.DB) (err error) {
	if sm.Id == uuid.Nil {
		sm.Id = uuid.New()
	}
	sm.JoinedAt = time.Now()
	sm.CreatedAt = time.Now()
	sm.UpdatedAt = time.Now()
	return
}

func (sm *UnitMember) BeforeUpdate(tx *gorm.DB) (err error) {
	sm.UpdatedAt = time.Now()
	return
}
