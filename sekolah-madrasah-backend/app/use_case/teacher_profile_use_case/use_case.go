package teacher_profile_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/teacher_profile_repository"
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type TeacherProfileUseCase interface {
	Create(req *CreateTeacherProfileRequest) (*schemas.TeacherProfile, error)
	GetById(id uuid.UUID) (*schemas.TeacherProfile, error)
	GetByUserId(userId uuid.UUID) (*schemas.TeacherProfile, error)
	GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error)
	Update(id uuid.UUID, req *UpdateTeacherProfileRequest) (*schemas.TeacherProfile, error)
	Delete(id uuid.UUID) error
}

type CreateTeacherProfileRequest struct {
	UserId           uuid.UUID
	UnitId           uuid.UUID
	NIP              *string
	NUPTK            *string
	EducationLevel   *string
	EducationMajor   *string
	EmploymentStatus string
	JoinDate         *string // Format: YYYY-MM-DD
	Subjects         []string
}

type UpdateTeacherProfileRequest struct {
	NIP              *string
	NUPTK            *string
	EducationLevel   *string
	EducationMajor   *string
	EmploymentStatus *string
	JoinDate         *string
	Subjects         []string
}

type teacherProfileUseCase struct {
	repo teacher_profile_repository.TeacherProfileRepository
}

func NewTeacherProfileUseCase(repo teacher_profile_repository.TeacherProfileRepository) TeacherProfileUseCase {
	return &teacherProfileUseCase{repo: repo}
}

func (uc *teacherProfileUseCase) Create(req *CreateTeacherProfileRequest) (*schemas.TeacherProfile, error) {
	// Check if profile already exists for user
	existing, _ := uc.repo.FindByUserId(req.UserId)
	if existing != nil {
		return nil, errors.New("teacher profile already exists for this user")
	}

	profile := &schemas.TeacherProfile{
		UserId:           req.UserId,
		UnitId:           req.UnitId,
		NIP:              req.NIP,
		NUPTK:            req.NUPTK,
		EducationLevel:   req.EducationLevel,
		EducationMajor:   req.EducationMajor,
		EmploymentStatus: req.EmploymentStatus,
		Subjects:         "[]", // Default empty JSON array
	}

	if err := uc.repo.Create(profile); err != nil {
		return nil, err
	}

	return uc.repo.FindById(profile.Id)
}

func (uc *teacherProfileUseCase) GetById(id uuid.UUID) (*schemas.TeacherProfile, error) {
	return uc.repo.FindById(id)
}

func (uc *teacherProfileUseCase) GetByUserId(userId uuid.UUID) (*schemas.TeacherProfile, error) {
	return uc.repo.FindByUserId(userId)
}

func (uc *teacherProfileUseCase) GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.FindByUnitId(unitId, page, limit)
}

func (uc *teacherProfileUseCase) Update(id uuid.UUID, req *UpdateTeacherProfileRequest) (*schemas.TeacherProfile, error) {
	profile, err := uc.repo.FindById(id)
	if err != nil {
		return nil, errors.New("teacher profile not found")
	}

	if req.NIP != nil {
		profile.NIP = req.NIP
	}
	if req.NUPTK != nil {
		profile.NUPTK = req.NUPTK
	}
	if req.EducationLevel != nil {
		profile.EducationLevel = req.EducationLevel
	}
	if req.EducationMajor != nil {
		profile.EducationMajor = req.EducationMajor
	}
	if req.EmploymentStatus != nil {
		profile.EmploymentStatus = *req.EmploymentStatus
	}

	if err := uc.repo.Update(profile); err != nil {
		return nil, err
	}

	return uc.repo.FindById(id)
}

func (uc *teacherProfileUseCase) Delete(id uuid.UUID) error {
	_, err := uc.repo.FindById(id)
	if err != nil {
		return errors.New("teacher profile not found")
	}
	return uc.repo.Delete(id)
}
