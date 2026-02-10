package subject_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(subject *schemas.Subject) error
	FindById(id uuid.UUID) (*schemas.Subject, error)
	FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.Subject, int64, error)
	Update(subject *schemas.Subject) error
	Delete(id uuid.UUID) error
	// Teacher-Subject relations
	AssignTeacher(ts *schemas.TeacherSubject) error
	RemoveTeacher(teacherProfileId, subjectId uuid.UUID) error
	FindByTeacher(teacherProfileId uuid.UUID) ([]schemas.Subject, error)
	FindTeachersBySubject(subjectId uuid.UUID) ([]schemas.TeacherProfile, error)
}

type subjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(subject *schemas.Subject) error {
	return r.db.Create(subject).Error
}

func (r *subjectRepository) FindById(id uuid.UUID) (*schemas.Subject, error) {
	var subject schemas.Subject
	err := r.db.Preload("Unit").First(&subject, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *subjectRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.Subject, int64, error) {
	var subjects []schemas.Subject
	var total int64

	query := r.db.Model(&schemas.Subject{}).Where("unit_id = ?", unitId)
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&subjects).Error
	return subjects, total, err
}

func (r *subjectRepository) Update(subject *schemas.Subject) error {
	return r.db.Save(subject).Error
}

func (r *subjectRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.Subject{}, "id = ?", id).Error
}

func (r *subjectRepository) AssignTeacher(ts *schemas.TeacherSubject) error {
	return r.db.Create(ts).Error
}

func (r *subjectRepository) RemoveTeacher(teacherProfileId, subjectId uuid.UUID) error {
	return r.db.Where("teacher_profile_id = ? AND subject_id = ?", teacherProfileId, subjectId).
		Delete(&schemas.TeacherSubject{}).Error
}

func (r *subjectRepository) FindByTeacher(teacherProfileId uuid.UUID) ([]schemas.Subject, error) {
	var subjects []schemas.Subject
	err := r.db.
		Joins("JOIN teacher_subjects ON teacher_subjects.subject_id = subjects.id").
		Where("teacher_subjects.teacher_profile_id = ?", teacherProfileId).
		Where("teacher_subjects.deleted_at IS NULL").
		Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepository) FindTeachersBySubject(subjectId uuid.UUID) ([]schemas.TeacherProfile, error) {
	var teachers []schemas.TeacherProfile
	err := r.db.
		Preload("User").
		Joins("JOIN teacher_subjects ON teacher_subjects.teacher_profile_id = teacher_profiles.id").
		Where("teacher_subjects.subject_id = ?", subjectId).
		Where("teacher_subjects.deleted_at IS NULL").
		Find(&teachers).Error
	return teachers, err
}
