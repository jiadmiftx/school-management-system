package activity_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	// Activity CRUD
	Create(activity *schemas.Activity) error
	FindById(id uuid.UUID) (*schemas.Activity, error)
	FindByUnitId(unitId uuid.UUID, activityType string, page, limit int) ([]schemas.Activity, int64, error)
	Update(activity *schemas.Activity) error
	Delete(id uuid.UUID) error
	// Teacher assignments
	AssignTeacher(at *schemas.ActivityTeacher) error
	RemoveTeacher(activityId, teacherProfileId uuid.UUID) error
	FindTeachersByActivity(activityId uuid.UUID) ([]schemas.ActivityTeacher, error)
	// Student enrollments
	EnrollStudent(as *schemas.ActivityStudent) error
	RemoveStudent(activityId, studentProfileId uuid.UUID) error
	FindStudentsByActivity(activityId uuid.UUID) ([]schemas.ActivityStudent, error)
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(activity *schemas.Activity) error {
	return r.db.Create(activity).Error
}

func (r *activityRepository) FindById(id uuid.UUID) (*schemas.Activity, error) {
	var activity schemas.Activity
	err := r.db.Preload("Unit").Preload("Teachers.TeacherProfile.User").Preload("Students.StudentProfile.User").
		First(&activity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}

func (r *activityRepository) FindByUnitId(unitId uuid.UUID, activityType string, page, limit int) ([]schemas.Activity, int64, error) {
	var activities []schemas.Activity
	var total int64

	query := r.db.Model(&schemas.Activity{}).Where("unit_id = ?", unitId)
	if activityType != "" {
		query = query.Where("type = ?", activityType)
	}
	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("name ASC").Offset(offset).Limit(limit).Find(&activities).Error
	return activities, total, err
}

func (r *activityRepository) Update(activity *schemas.Activity) error {
	return r.db.Save(activity).Error
}

func (r *activityRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&schemas.Activity{}, "id = ?", id).Error
}

func (r *activityRepository) AssignTeacher(at *schemas.ActivityTeacher) error {
	return r.db.Create(at).Error
}

func (r *activityRepository) RemoveTeacher(activityId, teacherProfileId uuid.UUID) error {
	return r.db.Where("activity_id = ? AND teacher_profile_id = ?", activityId, teacherProfileId).
		Delete(&schemas.ActivityTeacher{}).Error
}

func (r *activityRepository) FindTeachersByActivity(activityId uuid.UUID) ([]schemas.ActivityTeacher, error) {
	var teachers []schemas.ActivityTeacher
	err := r.db.Preload("TeacherProfile.User").
		Where("activity_id = ?", activityId).
		Where("deleted_at IS NULL").
		Find(&teachers).Error
	return teachers, err
}

func (r *activityRepository) EnrollStudent(as *schemas.ActivityStudent) error {
	return r.db.Create(as).Error
}

func (r *activityRepository) RemoveStudent(activityId, studentProfileId uuid.UUID) error {
	return r.db.Where("activity_id = ? AND student_profile_id = ?", activityId, studentProfileId).
		Delete(&schemas.ActivityStudent{}).Error
}

func (r *activityRepository) FindStudentsByActivity(activityId uuid.UUID) ([]schemas.ActivityStudent, error) {
	var students []schemas.ActivityStudent
	err := r.db.Preload("StudentProfile.User").
		Where("activity_id = ?", activityId).
		Where("deleted_at IS NULL").
		Find(&students).Error
	return students, err
}
