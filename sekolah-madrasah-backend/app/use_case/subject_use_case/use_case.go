package subject_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/subject_repository"
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type SubjectUseCase interface {
	Create(req *CreateSubjectRequest) (*schemas.Subject, error)
	GetById(id uuid.UUID) (*schemas.Subject, error)
	GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.Subject, int64, error)
	Update(id uuid.UUID, req *UpdateSubjectRequest) (*schemas.Subject, error)
	Delete(id uuid.UUID) error
	// Teacher assignment
	AssignTeacher(req *AssignTeacherRequest) error
	RemoveTeacher(teacherProfileId, subjectId uuid.UUID) error
	GetByTeacher(teacherProfileId uuid.UUID) ([]schemas.Subject, error)
}

type CreateSubjectRequest struct {
	UnitId      uuid.UUID
	Name        string
	Code        string
	Category    string
	Description *string
}

type UpdateSubjectRequest struct {
	Name        *string
	Code        *string
	Category    *string
	Description *string
	IsActive    *bool
}

type AssignTeacherRequest struct {
	TeacherProfileId uuid.UUID
	SubjectId        uuid.UUID
	IsPrimary        bool
}

type subjectUseCase struct {
	repo subject_repository.SubjectRepository
}

func NewSubjectUseCase(repo subject_repository.SubjectRepository) SubjectUseCase {
	return &subjectUseCase{repo: repo}
}

func (uc *subjectUseCase) Create(req *CreateSubjectRequest) (*schemas.Subject, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Code == "" {
		return nil, errors.New("code is required")
	}

	subject := &schemas.Subject{
		UnitId:      req.UnitId,
		Name:        req.Name,
		Code:        req.Code,
		Category:    req.Category,
		Description: req.Description,
		IsActive:    true,
	}

	if err := uc.repo.Create(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (uc *subjectUseCase) GetById(id uuid.UUID) (*schemas.Subject, error) {
	return uc.repo.FindById(id)
}

func (uc *subjectUseCase) GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.Subject, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.FindByUnitId(unitId, page, limit)
}

func (uc *subjectUseCase) Update(id uuid.UUID, req *UpdateSubjectRequest) (*schemas.Subject, error) {
	subject, err := uc.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		subject.Name = *req.Name
	}
	if req.Code != nil {
		subject.Code = *req.Code
	}
	if req.Category != nil {
		subject.Category = *req.Category
	}
	if req.Description != nil {
		subject.Description = req.Description
	}
	if req.IsActive != nil {
		subject.IsActive = *req.IsActive
	}

	if err := uc.repo.Update(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (uc *subjectUseCase) Delete(id uuid.UUID) error {
	return uc.repo.Delete(id)
}

func (uc *subjectUseCase) AssignTeacher(req *AssignTeacherRequest) error {
	ts := &schemas.TeacherSubject{
		TeacherProfileId: req.TeacherProfileId,
		SubjectId:        req.SubjectId,
		IsPrimary:        req.IsPrimary,
	}
	return uc.repo.AssignTeacher(ts)
}

func (uc *subjectUseCase) RemoveTeacher(teacherProfileId, subjectId uuid.UUID) error {
	return uc.repo.RemoveTeacher(teacherProfileId, subjectId)
}

func (uc *subjectUseCase) GetByTeacher(teacherProfileId uuid.UUID) ([]schemas.Subject, error) {
	return uc.repo.FindByTeacher(teacherProfileId)
}
