package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentProfile represents extended profile data for students.
// Linked 1:1 with User table via UserId.
type StudentProfile struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserId         uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"` // 1:1 with users
	UnitId         uuid.UUID      `gorm:"type:uuid;not null;index" json:"unit_id"`       // School
	NIS            *string        `gorm:"type:varchar(30)" json:"nis"`                   // Nomor Induk Siswa (internal)
	NISN           *string        `gorm:"type:varchar(20)" json:"nisn"`                  // Nomor Induk Siswa Nasional
	BirthPlace     *string        `gorm:"type:varchar(100)" json:"birth_place"`          // Tempat lahir
	BirthDate      *time.Time     `gorm:"type:date" json:"birth_date"`                   // Tanggal lahir
	Gender         *string        `gorm:"type:varchar(10)" json:"gender"`                // L/P
	Religion       *string        `gorm:"type:varchar(20)" json:"religion"`              // Agama
	Address        *string        `gorm:"type:text" json:"address"`                      // Alamat lengkap
	FatherName     *string        `gorm:"type:varchar(100)" json:"father_name"`          // Nama ayah
	MotherName     *string        `gorm:"type:varchar(100)" json:"mother_name"`          // Nama ibu
	GuardianName   *string        `gorm:"type:varchar(100)" json:"guardian_name"`        // Nama wali (jika ada)
	ParentPhone    *string        `gorm:"type:varchar(20)" json:"parent_phone"`          // Telepon orang tua
	EnrollmentDate *time.Time     `gorm:"type:date" json:"enrollment_date"`              // Tanggal masuk sekolah
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Unit *Unit `gorm:"foreignKey:UnitId" json:"unit,omitempty"`
}

func (StudentProfile) TableName() string { return "student_profiles" }

func (s *StudentProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Id == uuid.Nil {
		s.Id = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return
}

func (s *StudentProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}
