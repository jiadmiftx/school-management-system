package class_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/class_repository"
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type ClassUseCase interface {
	Create(req *CreateClassRequest) (*schemas.Class, error)
	GetById(id uuid.UUID) (*schemas.Class, error)
	GetByUnitId(unitId uuid.UUID, academicYear string, page, limit int) ([]schemas.Class, int64, error)
	Update(id uuid.UUID, req *UpdateClassRequest) (*schemas.Class, error)
	Delete(id uuid.UUID) error
}

type CreateClassRequest struct {
	UnitId            uuid.UUID
	Name              string
	Level             int
	AcademicYear      string
	HomeroomTeacherId *uuid.UUID
	Capacity          int
}

type UpdateClassRequest struct {
	Name              *string
	Level             *int
	AcademicYear      *string
	HomeroomTeacherId *uuid.UUID
	Capacity          *int
	IsActive          *bool
}

type classUseCase struct {
	repo class_repository.ClassRepository
}

func NewClassUseCase(repo class_repository.ClassRepository) ClassUseCase {
	return &classUseCase{repo: repo}
}

func (uc *classUseCase) Create(req *CreateClassRequest) (*schemas.Class, error) {
	if req.Name == "" {
		return nil, errors.New("class name is required")
	}
	if req.AcademicYear == "" {
		return nil, errors.New("academic year is required")
	}
	if req.Level < 1 || req.Level > 12 {
		return nil, errors.New("level must be between 1 and 12")
	}

	capacity := req.Capacity
	if capacity <= 0 {
		capacity = 30
	}

	class := &schemas.Class{
		UnitId:            req.UnitId,
		Name:              req.Name,
		Level:             req.Level,
		AcademicYear:      req.AcademicYear,
		HomeroomTeacherId: req.HomeroomTeacherId,
		Capacity:          capacity,
		IsActive:          true,
	}

	if err := uc.repo.Create(class); err != nil {
		return nil, err
	}

	return uc.repo.FindById(class.Id)
}

func (uc *classUseCase) GetById(id uuid.UUID) (*schemas.Class, error) {
	return uc.repo.FindById(id)
}

func (uc *classUseCase) GetByUnitId(unitId uuid.UUID, academicYear string, page, limit int) ([]schemas.Class, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.FindByUnitId(unitId, academicYear, page, limit)
}

func (uc *classUseCase) Update(id uuid.UUID, req *UpdateClassRequest) (*schemas.Class, error) {
	class, err := uc.repo.FindById(id)
	if err != nil {
		return nil, errors.New("class not found")
	}

	if req.Name != nil {
		class.Name = *req.Name
	}
	if req.Level != nil {
		if *req.Level < 1 || *req.Level > 12 {
			return nil, errors.New("level must be between 1 and 12")
		}
		class.Level = *req.Level
	}
	if req.AcademicYear != nil {
		class.AcademicYear = *req.AcademicYear
	}
	if req.HomeroomTeacherId != nil {
		class.HomeroomTeacherId = req.HomeroomTeacherId
	}
	if req.Capacity != nil {
		class.Capacity = *req.Capacity
	}
	if req.IsActive != nil {
		class.IsActive = *req.IsActive
	}

	if err := uc.repo.Update(class); err != nil {
		return nil, err
	}

	return uc.repo.FindById(id)
}

func (uc *classUseCase) Delete(id uuid.UUID) error {
	_, err := uc.repo.FindById(id)
	if err != nil {
		return errors.New("class not found")
	}
	return uc.repo.Delete(id)
}
