package class_enrollment_use_case

import (
	"errors"
	"sekolah-madrasah/app/repository/class_enrollment_repository"
	"sekolah-madrasah/database/schemas"
	"time"

	"github.com/google/uuid"
)

type ClassEnrollmentUseCase interface {
	Enroll(req *EnrollStudentRequest) (*schemas.ClassEnrollment, error)
	GetById(id uuid.UUID) (*schemas.ClassEnrollment, error)
	GetByClassId(classId uuid.UUID) ([]schemas.ClassEnrollment, error)
	GetByStudentProfileId(studentProfileId uuid.UUID) ([]schemas.ClassEnrollment, error)
	UpdateStatus(id uuid.UUID, status schemas.ClassEnrollmentStatus, notes *string) error
	Transfer(enrollmentId uuid.UUID, newClassId uuid.UUID) (*schemas.ClassEnrollment, error)
	Remove(id uuid.UUID) error
}

type EnrollStudentRequest struct {
	StudentProfileId uuid.UUID
	ClassId          uuid.UUID
	AcademicYear     string
	EnrolledAt       *string // Format: YYYY-MM-DD
}

type classEnrollmentUseCase struct {
	repo class_enrollment_repository.ClassEnrollmentRepository
}

func NewClassEnrollmentUseCase(repo class_enrollment_repository.ClassEnrollmentRepository) ClassEnrollmentUseCase {
	return &classEnrollmentUseCase{repo: repo}
}

func (uc *classEnrollmentUseCase) Enroll(req *EnrollStudentRequest) (*schemas.ClassEnrollment, error) {
	// Check if student already enrolled in a class for this academic year
	existing, _ := uc.repo.FindActiveByStudentAndYear(req.StudentProfileId, req.AcademicYear)
	if existing != nil {
		return nil, errors.New("student is already enrolled in a class for this academic year")
	}

	enrolledAt := time.Now()
	if req.EnrolledAt != nil {
		if t, err := time.Parse("2006-01-02", *req.EnrolledAt); err == nil {
			enrolledAt = t
		}
	}

	enrollment := &schemas.ClassEnrollment{
		StudentProfileId: req.StudentProfileId,
		ClassId:          req.ClassId,
		AcademicYear:     req.AcademicYear,
		Status:           schemas.EnrollmentStatusActive,
		EnrolledAt:       enrolledAt,
	}

	if err := uc.repo.Create(enrollment); err != nil {
		return nil, err
	}

	return uc.repo.FindById(enrollment.Id)
}

func (uc *classEnrollmentUseCase) GetById(id uuid.UUID) (*schemas.ClassEnrollment, error) {
	return uc.repo.FindById(id)
}

func (uc *classEnrollmentUseCase) GetByClassId(classId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	return uc.repo.FindByClassId(classId)
}

func (uc *classEnrollmentUseCase) GetByStudentProfileId(studentProfileId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	return uc.repo.FindByStudentProfileId(studentProfileId)
}

func (uc *classEnrollmentUseCase) UpdateStatus(id uuid.UUID, status schemas.ClassEnrollmentStatus, notes *string) error {
	enrollment, err := uc.repo.FindById(id)
	if err != nil {
		return errors.New("enrollment not found")
	}

	enrollment.Status = status
	if notes != nil {
		enrollment.Notes = notes
	}
	if status != schemas.EnrollmentStatusActive {
		now := time.Now()
		enrollment.LeftAt = &now
	}

	return uc.repo.Update(enrollment)
}

func (uc *classEnrollmentUseCase) Transfer(enrollmentId uuid.UUID, newClassId uuid.UUID) (*schemas.ClassEnrollment, error) {
	oldEnrollment, err := uc.repo.FindById(enrollmentId)
	if err != nil {
		return nil, errors.New("enrollment not found")
	}

	// Mark old enrollment as transferred
	now := time.Now()
	oldEnrollment.Status = schemas.EnrollmentStatusTransferred
	oldEnrollment.LeftAt = &now
	notes := "Transferred to new class"
	oldEnrollment.Notes = &notes
	if err := uc.repo.Update(oldEnrollment); err != nil {
		return nil, err
	}

	// Create new enrollment
	newEnrollment := &schemas.ClassEnrollment{
		StudentProfileId: oldEnrollment.StudentProfileId,
		ClassId:          newClassId,
		AcademicYear:     oldEnrollment.AcademicYear,
		Status:           schemas.EnrollmentStatusActive,
		EnrolledAt:       now,
	}

	if err := uc.repo.Create(newEnrollment); err != nil {
		return nil, err
	}

	return uc.repo.FindById(newEnrollment.Id)
}

func (uc *classEnrollmentUseCase) Remove(id uuid.UUID) error {
	_, err := uc.repo.FindById(id)
	if err != nil {
		return errors.New("enrollment not found")
	}
	return uc.repo.Delete(id)
}
