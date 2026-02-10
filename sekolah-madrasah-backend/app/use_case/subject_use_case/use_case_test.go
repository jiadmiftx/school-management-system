package subject_use_case

import (
	"errors"
	"testing"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of SubjectRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(subject *schemas.Subject) error {
	args := m.Called(subject)
	return args.Error(0)
}

func (m *MockRepository) FindById(id uuid.UUID) (*schemas.Subject, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.Subject), args.Error(1)
}

func (m *MockRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.Subject, int64, error) {
	args := m.Called(unitId, page, limit)
	return args.Get(0).([]schemas.Subject), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Update(subject *schemas.Subject) error {
	args := m.Called(subject)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) AssignTeacher(ts *schemas.TeacherSubject) error {
	args := m.Called(ts)
	return args.Error(0)
}

func (m *MockRepository) RemoveTeacher(teacherProfileId, subjectId uuid.UUID) error {
	args := m.Called(teacherProfileId, subjectId)
	return args.Error(0)
}

func (m *MockRepository) FindByTeacher(teacherProfileId uuid.UUID) ([]schemas.Subject, error) {
	args := m.Called(teacherProfileId)
	return args.Get(0).([]schemas.Subject), args.Error(1)
}

func (m *MockRepository) FindTeachersBySubject(subjectId uuid.UUID) ([]schemas.TeacherProfile, error) {
	args := m.Called(subjectId)
	return args.Get(0).([]schemas.TeacherProfile), args.Error(1)
}

// Tests

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	req := &CreateSubjectRequest{
		UnitId:   uuid.New(),
		Name:     "Matematika",
		Code:     "MTK",
		Category: "Umum",
	}

	mockRepo.On("Create", mock.AnythingOfType("*schemas.Subject")).Return(nil)

	subject, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, subject)
	assert.Equal(t, "Matematika", subject.Name)
	assert.Equal(t, "MTK", subject.Code)
	assert.Equal(t, "Umum", subject.Category)
	assert.True(t, subject.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestCreate_ValidationError_EmptyName(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	req := &CreateSubjectRequest{
		UnitId: uuid.New(),
		Name:   "",
		Code:   "MTK",
	}

	subject, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, subject)
	assert.Equal(t, "name is required", err.Error())
}

func TestCreate_ValidationError_EmptyCode(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	req := &CreateSubjectRequest{
		UnitId: uuid.New(),
		Name:   "Matematika",
		Code:   "",
	}

	subject, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, subject)
	assert.Equal(t, "code is required", err.Error())
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	id := uuid.New()
	expected := &schemas.Subject{
		Id:       id,
		Name:     "Bahasa Indonesia",
		Code:     "BIN",
		Category: "Umum",
		IsActive: true,
	}

	mockRepo.On("FindById", id).Return(expected, nil)

	subject, err := uc.GetById(id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, subject.Id)
	assert.Equal(t, "Bahasa Indonesia", subject.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("FindById", id).Return(nil, errors.New("not found"))

	subject, err := uc.GetById(id)

	assert.Error(t, err)
	assert.Nil(t, subject)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	unitId := uuid.New()
	subjects := []schemas.Subject{
		{Id: uuid.New(), Name: "Matematika", Code: "MTK"},
		{Id: uuid.New(), Name: "Fisika", Code: "FIS"},
	}

	mockRepo.On("FindByUnitId", unitId, 1, 10).Return(subjects, int64(2), nil)

	result, total, err := uc.GetByUnitId(unitId, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	id := uuid.New()
	existing := &schemas.Subject{
		Id:       id,
		Name:     "Matematika",
		Code:     "MTK",
		IsActive: true,
	}

	newName := "Matematika Wajib"
	req := &UpdateSubjectRequest{
		Name: &newName,
	}

	mockRepo.On("FindById", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.Subject")).Return(nil)

	subject, err := uc.Update(id, req)

	assert.NoError(t, err)
	assert.Equal(t, newName, subject.Name)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := uc.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssignTeacher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	req := &AssignTeacherRequest{
		TeacherProfileId: uuid.New(),
		SubjectId:        uuid.New(),
		IsPrimary:        true,
	}

	mockRepo.On("AssignTeacher", mock.AnythingOfType("*schemas.TeacherSubject")).Return(nil)

	err := uc.AssignTeacher(req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTeacher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	teacherId := uuid.New()
	subjectId := uuid.New()

	mockRepo.On("RemoveTeacher", teacherId, subjectId).Return(nil)

	err := uc.RemoveTeacher(teacherId, subjectId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetByTeacher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewSubjectUseCase(mockRepo)

	teacherId := uuid.New()
	subjects := []schemas.Subject{
		{Id: uuid.New(), Name: "Matematika"},
		{Id: uuid.New(), Name: "Fisika"},
	}

	mockRepo.On("FindByTeacher", teacherId).Return(subjects, nil)

	result, err := uc.GetByTeacher(teacherId)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}
