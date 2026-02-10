package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Resource    string    `gorm:"type:varchar(50);not null;index" json:"resource"`
	Action      string    `gorm:"type:varchar(20);not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Permission) TableName() string { return "permissions" }

func (p *Permission) BeforeCreate(tx interface{}) (err error) {
	if p.Id == uuid.Nil {
		p.Id = uuid.New()
	}
	p.CreatedAt = time.Now()
	return
}
