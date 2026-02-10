package activity_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/activity_repository"
	"sekolah-madrasah/database/schemas"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ActivityUseCase interface {
	Create(req *CreateActivityRequest) (*schemas.Activity, error)
	GetById(id uuid.UUID) (*schemas.Activity, error)
	GetByUnitId(unitId uuid.UUID, activityType string, page, limit int) ([]schemas.Activity, int64, error)
	Update(id uuid.UUID, req *UpdateActivityRequest) (*schemas.Activity, error)
	Delete(id uuid.UUID) error
	// Teacher assignments
	AssignTeacher(req *AssignTeacherRequest) error
	RemoveTeacher(activityId, teacherProfileId uuid.UUID) error
	GetTeachers(activityId uuid.UUID) ([]schemas.ActivityTeacher, error)
	// Student enrollments
	EnrollStudent(req *EnrollStudentRequest) error
	RemoveStudent(activityId, studentProfileId uuid.UUID) error
	GetStudents(activityId uuid.UUID) ([]schemas.ActivityStudent, error)
}

type CreateActivityRequest struct {
	UnitId          uuid.UUID
	Name            string
	Type            string
	Category        *string
	Description     *string
	StartDate       *time.Time
	EndDate         *time.Time
	RecurrenceType  string
	RecurrenceDays  []int64
	StartTime       *string
	EndTime         *string
	Location        *string
	MaxParticipants *int
	Fee             *float64
}

type UpdateActivityRequest struct {
	Name            *string
	Type            *string
	Category        *string
	Description     *string
	StartDate       *time.Time
	EndDate         *time.Time
	RecurrenceType  *string
	RecurrenceDays  []int64
	StartTime       *string
	EndTime         *string
	Location        *string
	MaxParticipants *int
	Fee             *float64
	IsActive        *bool
}

type AssignTeacherRequest struct {
	ActivityId       uuid.UUID
	TeacherProfileId uuid.UUID
	Role             string
}

type EnrollStudentRequest struct {
	ActivityId       uuid.UUID
	StudentProfileId uuid.UUID
	IsMandatory      bool
}

type activityUseCase struct {
	repo activity_repository.ActivityRepository
}

func NewActivityUseCase(repo activity_repository.ActivityRepository) ActivityUseCase {
	return &activityUseCase{repo: repo}
}

func (uc *activityUseCase) Create(req *CreateActivityRequest) (*schemas.Activity, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}

	activity := &schemas.Activity{
		UnitId:          req.UnitId,
		Name:            req.Name,
		Type:            req.Type,
		Category:        req.Category,
		Description:     req.Description,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		RecurrenceType:  req.RecurrenceType,
		RecurrenceDays:  pq.Int64Array(req.RecurrenceDays),
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Location:        req.Location,
		MaxParticipants: req.MaxParticipants,
		Fee:             req.Fee,
		IsActive:        true,
	}

	if activity.RecurrenceType == "" {
		activity.RecurrenceType = "none"
	}

	if err := uc.repo.Create(activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (uc *activityUseCase) GetById(id uuid.UUID) (*schemas.Activity, error) {
	return uc.repo.FindById(id)
}

func (uc *activityUseCase) GetByUnitId(unitId uuid.UUID, activityType string, page, limit int) ([]schemas.Activity, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.FindByUnitId(unitId, activityType, page, limit)
}

func (uc *activityUseCase) Update(id uuid.UUID, req *UpdateActivityRequest) (*schemas.Activity, error) {
	activity, err := uc.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		activity.Name = *req.Name
	}
	if req.Type != nil {
		activity.Type = *req.Type
	}
	if req.Category != nil {
		activity.Category = req.Category
	}
	if req.Description != nil {
		activity.Description = req.Description
	}
	if req.StartDate != nil {
		activity.StartDate = req.StartDate
	}
	if req.EndDate != nil {
		activity.EndDate = req.EndDate
	}
	if req.RecurrenceType != nil {
		activity.RecurrenceType = *req.RecurrenceType
	}
	if req.RecurrenceDays != nil {
		activity.RecurrenceDays = pq.Int64Array(req.RecurrenceDays)
	}
	if req.StartTime != nil {
		activity.StartTime = req.StartTime
	}
	if req.EndTime != nil {
		activity.EndTime = req.EndTime
	}
	if req.Location != nil {
		activity.Location = req.Location
	}
	if req.MaxParticipants != nil {
		activity.MaxParticipants = req.MaxParticipants
	}
	if req.Fee != nil {
		activity.Fee = req.Fee
	}
	if req.IsActive != nil {
		activity.IsActive = *req.IsActive
	}

	if err := uc.repo.Update(activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (uc *activityUseCase) Delete(id uuid.UUID) error {
	return uc.repo.Delete(id)
}

func (uc *activityUseCase) AssignTeacher(req *AssignTeacherRequest) error {
	role := req.Role
	if role == "" {
		role = "pembina"
	}
	at := &schemas.ActivityTeacher{
		ActivityId:       req.ActivityId,
		TeacherProfileId: req.TeacherProfileId,
		Role:             role,
	}
	return uc.repo.AssignTeacher(at)
}

func (uc *activityUseCase) RemoveTeacher(activityId, teacherProfileId uuid.UUID) error {
	return uc.repo.RemoveTeacher(activityId, teacherProfileId)
}

func (uc *activityUseCase) GetTeachers(activityId uuid.UUID) ([]schemas.ActivityTeacher, error) {
	return uc.repo.FindTeachersByActivity(activityId)
}

func (uc *activityUseCase) EnrollStudent(req *EnrollStudentRequest) error {
	now := time.Now()
	as := &schemas.ActivityStudent{
		ActivityId:       req.ActivityId,
		StudentProfileId: req.StudentProfileId,
		IsMandatory:      req.IsMandatory,
		JoinedAt:         &now,
	}
	return uc.repo.EnrollStudent(as)
}

func (uc *activityUseCase) RemoveStudent(activityId, studentProfileId uuid.UUID) error {
	return uc.repo.RemoveStudent(activityId, studentProfileId)
}

func (uc *activityUseCase) GetStudents(activityId uuid.UUID) ([]schemas.ActivityStudent, error) {
	return uc.repo.FindStudentsByActivity(activityId)
}
