package class_enrollment_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClassEnrollmentRepository interface {
	Create(enrollment *schemas.ClassEnrollment) error
	FindById(id uuid.UUID) (*schemas.ClassEnrollment, error)
	FindByClassId(classId uuid.UUID) ([]schemas.ClassEnrollment, error)
	FindByStudentProfileId(studentProfileId uuid.UUID) ([]schemas.ClassEnrollment, error)
	FindActiveByStudentAndYear(studentProfileId uuid.UUID, academicYear string) (*schemas.ClassEnrollment, error)
	Update(enrollment *schemas.ClassEnrollment) error
	Delete(id uuid.UUID) error
}

type classEnrollmentRepository struct {
	db *gorm.DB
}

func NewClassEnrollmentRepository(db *gorm.DB) ClassEnrollmentRepository {
	return &classEnrollmentRepository{db: db}
}

func (r *classEnrollmentRepository) Create(enrollment *schemas.ClassEnrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *classEnrollmentRepository) FindById(id uuid.UUID) (*schemas.ClassEnrollment, error) {
	var enrollment schemas.ClassEnrollment
	err := r.db.Preload("StudentProfile").Preload("StudentProfile.User").Preload("Class").First(&enrollment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *classEnrollmentRepository) FindByClassId(classId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	var enrollments []schemas.ClassEnrollment
	err := r.db.Preload("StudentProfile").Preload("StudentProfile.User").
		Where("class_id = ? AND status = ?", classId, schemas.EnrollmentStatusActive).
		Find(&enrollments).Error
	return enrollments, err
}

func (r *classEnrollmentRepository) FindByStudentProfileId(studentProfileId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	var enrollments []schemas.ClassEnrollment
	err := r.db.Preload("Class").
		Where("student_profile_id = ?", studentProfileId).
		Order("academic_year DESC").
		Find(&enrollments).Error
	return enrollments, err
}

func (r *classEnrollmentRepository) FindActiveByStudentAndYear(studentProfileId uuid.UUID, academicYear string) (*schemas.ClassEnrollment, error) {
	var enrollment schemas.ClassEnrollment
	err := r.db.Preload("Class").
		Where("student_profile_id = ? AND academic_year = ? AND status = ?", studentProfileId, academicYear, schemas.EnrollmentStatusActive).
		First(&enrollment).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

func (r *classEnrollmentRepository) Update(enrollment *schemas.ClassEnrollment) error {
	return r.db.Save(enrollment).Error
}

func (r *classEnrollmentRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.ClassEnrollment{}, "id = ?", id).Error
}
