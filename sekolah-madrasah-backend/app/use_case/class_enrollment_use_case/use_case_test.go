package class_enrollment_use_case

import (
	"errors"
	"testing"
	"time"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of ClassEnrollmentRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(enrollment *schemas.ClassEnrollment) error {
	args := m.Called(enrollment)
	return args.Error(0)
}

func (m *MockRepository) FindById(id uuid.UUID) (*schemas.ClassEnrollment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.ClassEnrollment), args.Error(1)
}

func (m *MockRepository) FindByClassId(classId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	args := m.Called(classId)
	return args.Get(0).([]schemas.ClassEnrollment), args.Error(1)
}

func (m *MockRepository) FindByStudentId(studentId uuid.UUID) ([]schemas.ClassEnrollment, error) {
	args := m.Called(studentId)
	return args.Get(0).([]schemas.ClassEnrollment), args.Error(1)
}

func (m *MockRepository) FindActiveByStudentId(studentId uuid.UUID) (*schemas.ClassEnrollment, error) {
	args := m.Called(studentId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.ClassEnrollment), args.Error(1)
}

func (m *MockRepository) Update(enrollment *schemas.ClassEnrollment) error {
	args := m.Called(enrollment)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Tests

func TestEnroll_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	studentId := uuid.New()
	classId := uuid.New()

	req := &EnrollStudentRequest{
		StudentProfileId: studentId,
		ClassId:          classId,
		AcademicYear:     "2025/2026",
	}

	mockRepo.On("FindActiveByStudentId", studentId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)

	enrollment, err := uc.Enroll(req)

	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	assert.Equal(t, studentId, enrollment.StudentProfileId)
	assert.Equal(t, classId, enrollment.ClassId)
	assert.Equal(t, schemas.EnrollmentStatusActive, enrollment.Status)
	mockRepo.AssertExpectations(t)
}

func TestEnroll_AlreadyEnrolled(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	studentId := uuid.New()
	existingEnrollment := &schemas.ClassEnrollment{
		Id:               uuid.New(),
		StudentProfileId: studentId,
		ClassId:          uuid.New(),
		Status:           schemas.EnrollmentStatusActive,
	}

	req := &EnrollStudentRequest{
		StudentProfileId: studentId,
		ClassId:          uuid.New(),
		AcademicYear:     "2025/2026",
	}

	mockRepo.On("FindActiveByStudentId", studentId).Return(existingEnrollment, nil)

	enrollment, err := uc.Enroll(req)

	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "already enrolled")
	mockRepo.AssertExpectations(t)
}

func TestEnroll_WithEnrollmentDate(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	studentId := uuid.New()
	enrolledAt := "2025-07-15"

	req := &EnrollStudentRequest{
		StudentProfileId: studentId,
		ClassId:          uuid.New(),
		AcademicYear:     "2025/2026",
		EnrolledAt:       &enrolledAt,
	}

	mockRepo.On("FindActiveByStudentId", studentId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)

	enrollment, err := uc.Enroll(req)

	assert.NoError(t, err)
	assert.NotNil(t, enrollment)
	expectedDate, _ := time.Parse("2006-01-02", enrolledAt)
	assert.Equal(t, expectedDate, enrollment.EnrolledAt)
}

func TestGetByClassId_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	classId := uuid.New()
	enrollments := []schemas.ClassEnrollment{
		{Id: uuid.New(), ClassId: classId, Status: schemas.EnrollmentStatusActive},
		{Id: uuid.New(), ClassId: classId, Status: schemas.EnrollmentStatusActive},
	}

	mockRepo.On("FindByClassId", classId).Return(enrollments, nil)

	result, err := uc.GetByClassId(classId)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestGetByClassId_EmptyList(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	classId := uuid.New()
	mockRepo.On("FindByClassId", classId).Return([]schemas.ClassEnrollment{}, nil)

	result, err := uc.GetByClassId(classId)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStatus_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	id := uuid.New()
	existing := &schemas.ClassEnrollment{
		Id:     id,
		Status: schemas.EnrollmentStatusActive,
	}

	notes := "Graduated with honors"

	mockRepo.On("FindById", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)

	err := uc.UpdateStatus(id, schemas.EnrollmentStatusGraduated, &notes)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateStatus_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("FindById", id).Return(nil, errors.New("not found"))

	err := uc.UpdateStatus(id, schemas.EnrollmentStatusGraduated, nil)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	enrollmentId := uuid.New()
	newClassId := uuid.New()
	studentId := uuid.New()

	existing := &schemas.ClassEnrollment{
		Id:               enrollmentId,
		StudentProfileId: studentId,
		ClassId:          uuid.New(),
		AcademicYear:     "2025/2026",
		Status:           schemas.EnrollmentStatusActive,
	}

	mockRepo.On("FindById", enrollmentId).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)
	mockRepo.On("Create", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)

	newEnrollment, err := uc.Transfer(enrollmentId, newClassId)

	assert.NoError(t, err)
	assert.NotNil(t, newEnrollment)
	assert.Equal(t, newClassId, newEnrollment.ClassId)
	assert.Equal(t, studentId, newEnrollment.StudentProfileId)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	enrollmentId := uuid.New()
	newClassId := uuid.New()

	mockRepo.On("FindById", enrollmentId).Return(nil, errors.New("not found"))

	newEnrollment, err := uc.Transfer(enrollmentId, newClassId)

	assert.Error(t, err)
	assert.Nil(t, newEnrollment)
	mockRepo.AssertExpectations(t)
}

func TestRemove_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := uc.Remove(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemove_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewClassEnrollmentUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(errors.New("delete failed"))

	err := uc.Remove(id)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// Table-driven test for status transitions
func TestUpdateStatus_Transitions(t *testing.T) {
	tests := []struct {
		name      string
		newStatus schemas.ClassEnrollmentStatus
		wantErr   bool
	}{
		{"to graduated", schemas.EnrollmentStatusGraduated, false},
		{"to transferred", schemas.EnrollmentStatusTransferred, false},
		{"to dropped", schemas.EnrollmentStatusDropped, false},
		{"to active", schemas.EnrollmentStatusActive, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			uc := NewClassEnrollmentUseCase(mockRepo)

			id := uuid.New()
			existing := &schemas.ClassEnrollment{
				Id:     id,
				Status: schemas.EnrollmentStatusActive,
			}

			mockRepo.On("FindById", id).Return(existing, nil)
			mockRepo.On("Update", mock.AnythingOfType("*schemas.ClassEnrollment")).Return(nil)

			err := uc.UpdateStatus(id, tt.newStatus, nil)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
