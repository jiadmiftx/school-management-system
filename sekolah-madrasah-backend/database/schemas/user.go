package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"`
	FullName        string         `gorm:"type:varchar(100)" json:"full_name"`
	Phone           string         `gorm:"type:varchar(20)" json:"phone"`
	Avatar          string         `gorm:"type:varchar(500)" json:"avatar"`
	IsSuperAdmin    bool           `gorm:"default:false" json:"is_super_admin"`
	IsActive        bool           `gorm:"default:true" json:"is_active"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == uuid.Nil {
		u.Id = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}
