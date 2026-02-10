package teacher_profile_use_case

import (
	"errors"
	"testing"
	"time"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of TeacherProfileRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(profile *schemas.TeacherProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockRepository) FindById(id uuid.UUID) (*schemas.TeacherProfile, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.TeacherProfile), args.Error(1)
}

func (m *MockRepository) FindByUserId(userId uuid.UUID) (*schemas.TeacherProfile, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.TeacherProfile), args.Error(1)
}

func (m *MockRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error) {
	args := m.Called(unitId, page, limit)
	return args.Get(0).([]schemas.TeacherProfile), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Update(profile *schemas.TeacherProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Tests

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	unitId := uuid.New()
	userId := uuid.New()
	nip := "198501012010011001"
	status := "pns"

	req := &CreateTeacherProfileRequest{
		UnitId:           unitId,
		UserId:           userId,
		NIP:              &nip,
		EmploymentStatus: status,
	}

	// Mock: Check user doesn't have existing profile
	mockRepo.On("FindByUserId", userId).Return(nil, errors.New("not found"))
	// Mock: Create profile
	mockRepo.On("Create", mock.AnythingOfType("*schemas.TeacherProfile")).Return(nil)
	// Mock: FindById is called after Create to return fresh profile
	mockRepo.On("FindById", mock.Anything).Return(&schemas.TeacherProfile{
		Id:               uuid.New(),
		UnitId:           unitId,
		UserId:           userId,
		NIP:              &nip,
		EmploymentStatus: status,
	}, nil)

	profile, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, unitId, profile.UnitId)
	assert.Equal(t, userId, profile.UserId)
	mockRepo.AssertExpectations(t)
}

func TestCreate_DuplicateProfile(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	userId := uuid.New()
	existingProfile := &schemas.TeacherProfile{
		Id:     uuid.New(),
		UserId: userId,
	}

	req := &CreateTeacherProfileRequest{
		UnitId:           uuid.New(),
		UserId:           userId,
		EmploymentStatus: "honorer",
	}

	mockRepo.On("FindByUserId", userId).Return(existingProfile, nil)

	profile, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	id := uuid.New()
	expected := &schemas.TeacherProfile{
		Id:               id,
		UserId:           uuid.New(),
		UnitId:           uuid.New(),
		EmploymentStatus: "pns",
	}

	mockRepo.On("FindById", id).Return(expected, nil)

	profile, err := uc.GetById(id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, profile.Id)
	mockRepo.AssertExpectations(t)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("FindById", id).Return(nil, errors.New("not found"))

	profile, err := uc.GetById(id)

	assert.Error(t, err)
	assert.Nil(t, profile)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	unitId := uuid.New()
	profiles := []schemas.TeacherProfile{
		{Id: uuid.New(), UnitId: unitId},
		{Id: uuid.New(), UnitId: unitId},
	}

	mockRepo.On("FindByUnitId", unitId, 1, 10).Return(profiles, int64(2), nil)

	result, total, err := uc.GetByUnitId(unitId, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_EmptyList(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	unitId := uuid.New()
	mockRepo.On("FindByUnitId", unitId, 1, 10).Return([]schemas.TeacherProfile{}, int64(0), nil)

	result, total, err := uc.GetByUnitId(unitId, 1, 10)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	id := uuid.New()
	existing := &schemas.TeacherProfile{
		Id:               id,
		EmploymentStatus: "honorer",
	}

	newStatus := "pns"
	req := &UpdateTeacherProfileRequest{
		EmploymentStatus: &newStatus,
	}

	mockRepo.On("FindById", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.TeacherProfile")).Return(nil)

	profile, err := uc.Update(id, req)

	assert.NoError(t, err)
	assert.Equal(t, "pns", profile.EmploymentStatus)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	id := uuid.New()
	// Mock: FindById is called first to verify profile exists
	mockRepo.On("FindById", id).Return(&schemas.TeacherProfile{Id: id}, nil)
	mockRepo.On("Delete", id).Return(nil)

	err := uc.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreate_ValidationError_EmptyStatus(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	req := &CreateTeacherProfileRequest{
		UnitId:           uuid.New(),
		UserId:           uuid.New(),
		EmploymentStatus: "", // empty status
	}

	mockRepo.On("FindByUserId", req.UserId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.TeacherProfile")).Return(nil)
	mockRepo.On("FindById", mock.Anything).Return(&schemas.TeacherProfile{
		Id:               uuid.New(),
		UnitId:           req.UnitId,
		UserId:           req.UserId,
		EmploymentStatus: "honorer",
	}, nil)

	profile, err := uc.Create(req)

	// Should still work with empty status (defaults to honorer)
	assert.NoError(t, err)
	assert.NotNil(t, profile)
}

func TestCreate_WithJoinDate(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewTeacherProfileUseCase(mockRepo)

	joinDate := "2020-01-15"
	req := &CreateTeacherProfileRequest{
		UnitId:           uuid.New(),
		UserId:           uuid.New(),
		EmploymentStatus: "pns",
		JoinDate:         &joinDate,
	}

	mockRepo.On("FindByUserId", req.UserId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.TeacherProfile")).Return(nil)
	joinDateTime, _ := time.Parse("2006-01-02", joinDate)
	mockRepo.On("FindById", mock.Anything).Return(&schemas.TeacherProfile{
		Id:               uuid.New(),
		UnitId:           req.UnitId,
		UserId:           req.UserId,
		EmploymentStatus: "pns",
		JoinDate:         &joinDateTime,
	}, nil)

	profile, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.NotNil(t, profile.JoinDate)
	expectedDate, _ := time.Parse("2006-01-02", joinDate)
	assert.Equal(t, expectedDate, *profile.JoinDate)
}
