package class_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(class *schemas.Class) error
	FindById(id uuid.UUID) (*schemas.Class, error)
	FindByUnitId(unitId uuid.UUID, academicYear string, page, limit int) ([]schemas.Class, int64, error)
	Update(class *schemas.Class) error
	Delete(id uuid.UUID) error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *schemas.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) FindById(id uuid.UUID) (*schemas.Class, error) {
	var class schemas.Class
	err := r.db.Preload("Unit").Preload("HomeroomTeacher").Preload("HomeroomTeacher.User").First(&class, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) FindByUnitId(unitId uuid.UUID, academicYear string, page, limit int) ([]schemas.Class, int64, error) {
	var classes []schemas.Class
	var total int64

	query := r.db.Model(&schemas.Class{}).Where("unit_id = ?", unitId)

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Preload("HomeroomTeacher").Preload("HomeroomTeacher.User").
		Order("level ASC, name ASC").
		Offset(offset).
		Limit(limit).
		Find(&classes).Error

	if err != nil {
		return nil, 0, err
	}

	return classes, total, nil
}

func (r *classRepository) Update(class *schemas.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.Class{}, "id = ?", id).Error
}
