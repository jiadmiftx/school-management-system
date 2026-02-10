package activity_use_case

import (
	"errors"
	"testing"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of ActivityRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(activity *schemas.Activity) error {
	args := m.Called(activity)
	return args.Error(0)
}

func (m *MockRepository) FindById(id uuid.UUID) (*schemas.Activity, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.Activity), args.Error(1)
}

func (m *MockRepository) FindByUnitId(unitId uuid.UUID, activityType string, page, limit int) ([]schemas.Activity, int64, error) {
	args := m.Called(unitId, activityType, page, limit)
	return args.Get(0).([]schemas.Activity), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Update(activity *schemas.Activity) error {
	args := m.Called(activity)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) AssignTeacher(at *schemas.ActivityTeacher) error {
	args := m.Called(at)
	return args.Error(0)
}

func (m *MockRepository) RemoveTeacher(activityId, teacherProfileId uuid.UUID) error {
	args := m.Called(activityId, teacherProfileId)
	return args.Error(0)
}

func (m *MockRepository) FindTeachersByActivity(activityId uuid.UUID) ([]schemas.ActivityTeacher, error) {
	args := m.Called(activityId)
	return args.Get(0).([]schemas.ActivityTeacher), args.Error(1)
}

func (m *MockRepository) EnrollStudent(as *schemas.ActivityStudent) error {
	args := m.Called(as)
	return args.Error(0)
}

func (m *MockRepository) RemoveStudent(activityId, studentProfileId uuid.UUID) error {
	args := m.Called(activityId, studentProfileId)
	return args.Error(0)
}

func (m *MockRepository) FindStudentsByActivity(activityId uuid.UUID) ([]schemas.ActivityStudent, error) {
	args := m.Called(activityId)
	return args.Get(0).([]schemas.ActivityStudent), args.Error(1)
}

// Tests

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	req := &CreateActivityRequest{
		UnitId:         uuid.New(),
		Name:           "Pramuka",
		Type:           "ekstrakurikuler",
		RecurrenceType: "weekly",
	}

	mockRepo.On("Create", mock.AnythingOfType("*schemas.Activity")).Return(nil)

	activity, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, activity)
	assert.Equal(t, "Pramuka", activity.Name)
	assert.Equal(t, "ekstrakurikuler", activity.Type)
	assert.True(t, activity.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestCreate_ValidationError_EmptyName(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	req := &CreateActivityRequest{
		UnitId: uuid.New(),
		Name:   "",
		Type:   "ekstrakurikuler",
	}

	activity, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, activity)
	assert.Equal(t, "name is required", err.Error())
}

func TestCreate_ValidationError_EmptyType(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	req := &CreateActivityRequest{
		UnitId: uuid.New(),
		Name:   "Pramuka",
		Type:   "",
	}

	activity, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, activity)
	assert.Equal(t, "type is required", err.Error())
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	id := uuid.New()
	expected := &schemas.Activity{
		Id:       id,
		Name:     "Kajian Fiqih",
		Type:     "kajian",
		IsActive: true,
	}

	mockRepo.On("FindById", id).Return(expected, nil)

	activity, err := uc.GetById(id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, activity.Id)
	assert.Equal(t, "Kajian Fiqih", activity.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("FindById", id).Return(nil, errors.New("not found"))

	activity, err := uc.GetById(id)

	assert.Error(t, err)
	assert.Nil(t, activity)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	unitId := uuid.New()
	activities := []schemas.Activity{
		{Id: uuid.New(), Name: "Pramuka", Type: "ekstrakurikuler"},
		{Id: uuid.New(), Name: "Futsal", Type: "ekstrakurikuler"},
	}

	mockRepo.On("FindByUnitId", unitId, "ekstrakurikuler", 1, 10).Return(activities, int64(2), nil)

	result, total, err := uc.GetByUnitId(unitId, "ekstrakurikuler", 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	id := uuid.New()
	existing := &schemas.Activity{
		Id:       id,
		Name:     "Pramuka",
		Type:     "ekstrakurikuler",
		IsActive: true,
	}

	newName := "Pramuka Sabtu"
	req := &UpdateActivityRequest{
		Name: &newName,
	}

	mockRepo.On("FindById", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.Activity")).Return(nil)

	activity, err := uc.Update(id, req)

	assert.NoError(t, err)
	assert.Equal(t, newName, activity.Name)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := uc.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAssignTeacher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	req := &AssignTeacherRequest{
		ActivityId:       uuid.New(),
		TeacherProfileId: uuid.New(),
		Role:             "pembina",
	}

	mockRepo.On("AssignTeacher", mock.AnythingOfType("*schemas.ActivityTeacher")).Return(nil)

	err := uc.AssignTeacher(req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveTeacher_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	activityId := uuid.New()
	teacherId := uuid.New()

	mockRepo.On("RemoveTeacher", activityId, teacherId).Return(nil)

	err := uc.RemoveTeacher(activityId, teacherId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEnrollStudent_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	req := &EnrollStudentRequest{
		ActivityId:       uuid.New(),
		StudentProfileId: uuid.New(),
		IsMandatory:      true,
	}

	mockRepo.On("EnrollStudent", mock.AnythingOfType("*schemas.ActivityStudent")).Return(nil)

	err := uc.EnrollStudent(req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRemoveStudent_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	activityId := uuid.New()
	studentId := uuid.New()

	mockRepo.On("RemoveStudent", activityId, studentId).Return(nil)

	err := uc.RemoveStudent(activityId, studentId)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTeachers_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	activityId := uuid.New()
	teachers := []schemas.ActivityTeacher{
		{Id: uuid.New(), ActivityId: activityId, Role: "pembina"},
		{Id: uuid.New(), ActivityId: activityId, Role: "pengisi"},
	}

	mockRepo.On("FindTeachersByActivity", activityId).Return(teachers, nil)

	result, err := uc.GetTeachers(activityId)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestGetStudents_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewActivityUseCase(mockRepo)

	activityId := uuid.New()
	students := []schemas.ActivityStudent{
		{Id: uuid.New(), ActivityId: activityId, IsMandatory: true},
		{Id: uuid.New(), ActivityId: activityId, IsMandatory: false},
	}

	mockRepo.On("FindStudentsByActivity", activityId).Return(students, nil)

	result, err := uc.GetStudents(activityId)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}
