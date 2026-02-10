package schemas

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Activity represents a school activity (ekstrakurikuler, kajian, event).
type Activity struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UnitId      uuid.UUID `gorm:"type:uuid;not null;index" json:"unit_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`     // "Pramuka", "Kajian Fiqih"
	Type        string    `gorm:"type:varchar(50);not null" json:"type"`      // ekstrakurikuler/kajian/event
	Category    *string   `gorm:"type:varchar(50)" json:"category,omitempty"` // halaqah/tahsin/daurah/olahraga/seni/akademik
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	// Date Range
	StartDate *time.Time `gorm:"type:date" json:"start_date,omitempty"` // Activity start date
	EndDate   *time.Time `gorm:"type:date" json:"end_date,omitempty"`   // Activity end date (null = ongoing)

	// Recurrence
	RecurrenceType string        `gorm:"type:varchar(20);default:'none'" json:"recurrence_type"` // none/daily/weekly/monthly
	RecurrenceDays pq.Int64Array `gorm:"type:integer[]" json:"recurrence_days,omitempty"`        // Array of days: weekly [0-6], monthly [1-31]

	// Time
	StartTime *string `gorm:"type:varchar(10)" json:"start_time,omitempty"` // "14:00"
	EndTime   *string `gorm:"type:varchar(10)" json:"end_time,omitempty"`   // "16:00"

	// Additional Info
	Location        *string  `gorm:"type:varchar(200)" json:"location,omitempty"` // Venue/place
	MaxParticipants *int     `gorm:"type:int" json:"max_participants,omitempty"`  // Max capacity
	Fee             *float64 `gorm:"type:decimal(12,2)" json:"fee,omitempty"`     // Cost (if any)

	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Unit     *Unit             `gorm:"foreignKey:UnitId" json:"unit,omitempty"`
	Teachers []ActivityTeacher `gorm:"foreignKey:ActivityId" json:"teachers,omitempty"`
	Students []ActivityStudent `gorm:"foreignKey:ActivityId" json:"students,omitempty"`
}

func (Activity) TableName() string { return "activities" }

func (a *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	if a.Id == uuid.Nil {
		a.Id = uuid.New()
	}
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return
}

func (a *Activity) BeforeUpdate(tx *gorm.DB) (err error) {
	a.UpdatedAt = time.Now()
	return
}
