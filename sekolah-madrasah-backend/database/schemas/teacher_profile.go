package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TeacherProfile represents extended profile data for teachers.
// Linked 1:1 with User table via UserId.
type TeacherProfile struct {
	Id               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UserId           uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`               // 1:1 with users
	UnitId           uuid.UUID      `gorm:"type:uuid;not null;index" json:"unit_id"`                     // School
	NIP              *string        `gorm:"type:varchar(30)" json:"nip"`                                 // Nomor Induk Pegawai
	NUPTK            *string        `gorm:"type:varchar(30)" json:"nuptk"`                               // Nomor Unik Pendidik
	EducationLevel   *string        `gorm:"type:varchar(20)" json:"education_level"`                     // S1/S2/S3/D3
	EducationMajor   *string        `gorm:"type:varchar(100)" json:"education_major"`                    // Jurusan
	EmploymentStatus string         `gorm:"type:varchar(20);default:'honorer'" json:"employment_status"` // PNS/Honorer/GTY/Kontrak
	JoinDate         *time.Time     `gorm:"type:date" json:"join_date"`                                  // Tanggal mulai mengajar
	Subjects         string         `gorm:"type:jsonb;default:'[]'" json:"subjects"`                     // Mata pelajaran (array)
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"foreignKey:UserId" json:"user,omitempty"`
	Unit *Unit `gorm:"foreignKey:UnitId" json:"unit,omitempty"`

	// Reverse relation: Class where this teacher is homeroom teacher
	HomeroomClass *Class `gorm:"foreignKey:HomeroomTeacherId" json:"homeroom_class,omitempty"`

	// Subjects taught by this teacher (via TeacherSubject pivot)
	TeacherSubjects []TeacherSubject `gorm:"foreignKey:TeacherProfileId" json:"teacher_subjects,omitempty"`
}

func (TeacherProfile) TableName() string { return "teacher_profiles" }

func (t *TeacherProfile) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Id == uuid.Nil {
		t.Id = uuid.New()
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return
}

func (t *TeacherProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return
}
