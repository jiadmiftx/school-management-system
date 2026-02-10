package teacher_profile_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherProfileRepository interface {
	Create(profile *schemas.TeacherProfile) error
	FindById(id uuid.UUID) (*schemas.TeacherProfile, error)
	FindByUserId(userId uuid.UUID) (*schemas.TeacherProfile, error)
	FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error)
	Update(profile *schemas.TeacherProfile) error
	Delete(id uuid.UUID) error
}

type teacherProfileRepository struct {
	db *gorm.DB
}

func NewTeacherProfileRepository(db *gorm.DB) TeacherProfileRepository {
	return &teacherProfileRepository{db: db}
}

func (r *teacherProfileRepository) Create(profile *schemas.TeacherProfile) error {
	return r.db.Create(profile).Error
}

func (r *teacherProfileRepository) FindById(id uuid.UUID) (*schemas.TeacherProfile, error) {
	var profile schemas.TeacherProfile
	err := r.db.Preload("User").Preload("Unit").Preload("HomeroomClass").Preload("TeacherSubjects.Subject").First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *teacherProfileRepository) FindByUserId(userId uuid.UUID) (*schemas.TeacherProfile, error) {
	var profile schemas.TeacherProfile
	err := r.db.Preload("User").Preload("Unit").Preload("HomeroomClass").Preload("TeacherSubjects.Subject").First(&profile, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *teacherProfileRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error) {
	var profiles []schemas.TeacherProfile
	var total int64

	query := r.db.Model(&schemas.TeacherProfile{}).Where("unit_id = ?", unitId)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("User").Preload("HomeroomClass").Preload("TeacherSubjects.Subject").Offset(offset).Limit(limit).Find(&profiles).Error
	return profiles, total, err
}

func (r *teacherProfileRepository) Update(profile *schemas.TeacherProfile) error {
	return r.db.Save(profile).Error
}

func (r *teacherProfileRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.TeacherProfile{}, "id = ?", id).Error
}
