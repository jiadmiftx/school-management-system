package student_profile_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentProfileRepository interface {
	Create(profile *schemas.StudentProfile) error
	FindById(id uuid.UUID) (*schemas.StudentProfile, error)
	FindByUserId(userId uuid.UUID) (*schemas.StudentProfile, error)
	FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.StudentProfile, int64, error)
	Update(profile *schemas.StudentProfile) error
	Delete(id uuid.UUID) error
}

type studentProfileRepository struct {
	db *gorm.DB
}

func NewStudentProfileRepository(db *gorm.DB) StudentProfileRepository {
	return &studentProfileRepository{db: db}
}

func (r *studentProfileRepository) Create(profile *schemas.StudentProfile) error {
	return r.db.Create(profile).Error
}

func (r *studentProfileRepository) FindById(id uuid.UUID) (*schemas.StudentProfile, error) {
	var profile schemas.StudentProfile
	err := r.db.Preload("User").Preload("Unit").First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *studentProfileRepository) FindByUserId(userId uuid.UUID) (*schemas.StudentProfile, error) {
	var profile schemas.StudentProfile
	err := r.db.Preload("User").Preload("Unit").First(&profile, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *studentProfileRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.StudentProfile, int64, error) {
	var profiles []schemas.StudentProfile
	var total int64

	query := r.db.Model(&schemas.StudentProfile{}).Where("unit_id = ?", unitId)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Preload("User").Offset(offset).Limit(limit).Find(&profiles).Error
	return profiles, total, err
}

func (r *studentProfileRepository) Update(profile *schemas.StudentProfile) error {
	return r.db.Save(profile).Error
}

func (r *studentProfileRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.StudentProfile{}, "id = ?", id).Error
}
