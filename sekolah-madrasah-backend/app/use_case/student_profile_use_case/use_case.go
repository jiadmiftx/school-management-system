package student_profile_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/student_profile_repository"
	"sekolah-madrasah/database/schemas"
	"time"

	"github.com/google/uuid"
)

type StudentProfileUseCase interface {
	Create(req *CreateStudentProfileRequest) (*schemas.StudentProfile, error)
	GetById(id uuid.UUID) (*schemas.StudentProfile, error)
	GetByUserId(userId uuid.UUID) (*schemas.StudentProfile, error)
	GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.StudentProfile, int64, error)
	Update(id uuid.UUID, req *UpdateStudentProfileRequest) (*schemas.StudentProfile, error)
	Delete(id uuid.UUID) error
}

type CreateStudentProfileRequest struct {
	UserId         uuid.UUID
	UnitId         uuid.UUID
	NIS            *string
	NISN           *string
	BirthPlace     *string
	BirthDate      *string // Format: YYYY-MM-DD
	Gender         *string
	Religion       *string
	Address        *string
	FatherName     *string
	MotherName     *string
	GuardianName   *string
	ParentPhone    *string
	EnrollmentDate *string // Format: YYYY-MM-DD
}

type UpdateStudentProfileRequest struct {
	NIS            *string
	NISN           *string
	BirthPlace     *string
	BirthDate      *string
	Gender         *string
	Religion       *string
	Address        *string
	FatherName     *string
	MotherName     *string
	GuardianName   *string
	ParentPhone    *string
	EnrollmentDate *string
}

type studentProfileUseCase struct {
	repo student_profile_repository.StudentProfileRepository
}

func NewStudentProfileUseCase(repo student_profile_repository.StudentProfileRepository) StudentProfileUseCase {
	return &studentProfileUseCase{repo: repo}
}

func (uc *studentProfileUseCase) Create(req *CreateStudentProfileRequest) (*schemas.StudentProfile, error) {
	existing, _ := uc.repo.FindByUserId(req.UserId)
	if existing != nil {
		return nil, errors.New("student profile already exists for this user")
	}

	profile := &schemas.StudentProfile{
		UserId:       req.UserId,
		UnitId:       req.UnitId,
		NIS:          req.NIS,
		NISN:         req.NISN,
		BirthPlace:   req.BirthPlace,
		Gender:       req.Gender,
		Religion:     req.Religion,
		Address:      req.Address,
		FatherName:   req.FatherName,
		MotherName:   req.MotherName,
		GuardianName: req.GuardianName,
		ParentPhone:  req.ParentPhone,
	}

	if req.BirthDate != nil {
		if t, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			profile.BirthDate = &t
		}
	}
	if req.EnrollmentDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EnrollmentDate); err == nil {
			profile.EnrollmentDate = &t
		}
	}

	if err := uc.repo.Create(profile); err != nil {
		return nil, err
	}

	return uc.repo.FindById(profile.Id)
}

func (uc *studentProfileUseCase) GetById(id uuid.UUID) (*schemas.StudentProfile, error) {
	return uc.repo.FindById(id)
}

func (uc *studentProfileUseCase) GetByUserId(userId uuid.UUID) (*schemas.StudentProfile, error) {
	return uc.repo.FindByUserId(userId)
}

func (uc *studentProfileUseCase) GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.StudentProfile, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repo.FindByUnitId(unitId, page, limit)
}

func (uc *studentProfileUseCase) Update(id uuid.UUID, req *UpdateStudentProfileRequest) (*schemas.StudentProfile, error) {
	profile, err := uc.repo.FindById(id)
	if err != nil {
		return nil, errors.New("student profile not found")
	}

	if req.NIS != nil {
		profile.NIS = req.NIS
	}
	if req.NISN != nil {
		profile.NISN = req.NISN
	}
	if req.BirthPlace != nil {
		profile.BirthPlace = req.BirthPlace
	}
	if req.BirthDate != nil {
		if t, err := time.Parse("2006-01-02", *req.BirthDate); err == nil {
			profile.BirthDate = &t
		}
	}
	if req.Gender != nil {
		profile.Gender = req.Gender
	}
	if req.Religion != nil {
		profile.Religion = req.Religion
	}
	if req.Address != nil {
		profile.Address = req.Address
	}
	if req.FatherName != nil {
		profile.FatherName = req.FatherName
	}
	if req.MotherName != nil {
		profile.MotherName = req.MotherName
	}
	if req.GuardianName != nil {
		profile.GuardianName = req.GuardianName
	}
	if req.ParentPhone != nil {
		profile.ParentPhone = req.ParentPhone
	}
	if req.EnrollmentDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EnrollmentDate); err == nil {
			profile.EnrollmentDate = &t
		}
	}

	if err := uc.repo.Update(profile); err != nil {
		return nil, err
	}

	return uc.repo.FindById(id)
}

func (uc *studentProfileUseCase) Delete(id uuid.UUID) error {
	_, err := uc.repo.FindById(id)
	if err != nil {
		return errors.New("student profile not found")
	}
	return uc.repo.Delete(id)
}
