package student_profile_use_case

import (
	"errors"
	"testing"
	"time"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of StudentProfileRepository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(profile *schemas.StudentProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockRepository) FindById(id uuid.UUID) (*schemas.StudentProfile, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.StudentProfile), args.Error(1)
}

func (m *MockRepository) FindByUserId(userId uuid.UUID) (*schemas.StudentProfile, error) {
	args := m.Called(userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.StudentProfile), args.Error(1)
}

func (m *MockRepository) FindByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.StudentProfile, int64, error) {
	args := m.Called(unitId, page, limit)
	return args.Get(0).([]schemas.StudentProfile), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Update(profile *schemas.StudentProfile) error {
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
	uc := NewStudentProfileUseCase(mockRepo)

	unitId := uuid.New()
	userId := uuid.New()
	nis := "2024001"
	nisn := "0012345678"

	req := &CreateStudentProfileRequest{
		UnitId: unitId,
		UserId: userId,
		NIS:    &nis,
		NISN:   &nisn,
	}

	mockRepo.On("FindByUserId", userId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.StudentProfile")).Return(nil)

	profile, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, unitId, profile.UnitId)
	assert.Equal(t, userId, profile.UserId)
	assert.Equal(t, nis, *profile.NIS)
	mockRepo.AssertExpectations(t)
}

func TestCreate_DuplicateProfile(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	userId := uuid.New()
	existingProfile := &schemas.StudentProfile{
		Id:     uuid.New(),
		UserId: userId,
	}

	req := &CreateStudentProfileRequest{
		UnitId: uuid.New(),
		UserId: userId,
	}

	mockRepo.On("FindByUserId", userId).Return(existingProfile, nil)

	profile, err := uc.Create(req)

	assert.Error(t, err)
	assert.Nil(t, profile)
	assert.Contains(t, err.Error(), "already exists")
	mockRepo.AssertExpectations(t)
}

func TestCreate_WithBirthDate(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	birthDate := "2010-05-15"
	birthPlace := "Jakarta"

	req := &CreateStudentProfileRequest{
		UnitId:     uuid.New(),
		UserId:     uuid.New(),
		BirthPlace: &birthPlace,
		BirthDate:  &birthDate,
	}

	mockRepo.On("FindByUserId", req.UserId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.StudentProfile")).Return(nil)

	profile, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, birthPlace, *profile.BirthPlace)
	expectedDate, _ := time.Parse("2006-01-02", birthDate)
	assert.Equal(t, expectedDate, *profile.BirthDate)
}

func TestCreate_WithParentInfo(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	fatherName := "Ahmad Budi"
	motherName := "Siti Aminah"
	parentPhone := "081234567890"

	req := &CreateStudentProfileRequest{
		UnitId:      uuid.New(),
		UserId:      uuid.New(),
		FatherName:  &fatherName,
		MotherName:  &motherName,
		ParentPhone: &parentPhone,
	}

	mockRepo.On("FindByUserId", req.UserId).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*schemas.StudentProfile")).Return(nil)

	profile, err := uc.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, fatherName, *profile.FatherName)
	assert.Equal(t, motherName, *profile.MotherName)
}

func TestGetById_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	id := uuid.New()
	nis := "2024001"
	expected := &schemas.StudentProfile{
		Id:     id,
		UserId: uuid.New(),
		UnitId: uuid.New(),
		NIS:    &nis,
	}

	mockRepo.On("FindById", id).Return(expected, nil)

	profile, err := uc.GetById(id)

	assert.NoError(t, err)
	assert.Equal(t, expected.Id, profile.Id)
	mockRepo.AssertExpectations(t)
}

func TestGetById_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("FindById", id).Return(nil, errors.New("not found"))

	profile, err := uc.GetById(id)

	assert.Error(t, err)
	assert.Nil(t, profile)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	unitId := uuid.New()
	profiles := []schemas.StudentProfile{
		{Id: uuid.New(), UnitId: unitId},
		{Id: uuid.New(), UnitId: unitId},
		{Id: uuid.New(), UnitId: unitId},
	}

	mockRepo.On("FindByUnitId", unitId, 1, 10).Return(profiles, int64(3), nil)

	result, total, err := uc.GetByUnitId(unitId, 1, 10)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(result))
	assert.Equal(t, int64(3), total)
	mockRepo.AssertExpectations(t)
}

func TestGetByUnitId_Pagination(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	unitId := uuid.New()
	profiles := []schemas.StudentProfile{
		{Id: uuid.New(), UnitId: unitId},
	}

	// Page 2 with limit 5
	mockRepo.On("FindByUnitId", unitId, 2, 5).Return(profiles, int64(6), nil)

	result, total, err := uc.GetByUnitId(unitId, 2, 5)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, int64(6), total) // total 6 items
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	id := uuid.New()
	oldNIS := "2024001"
	existing := &schemas.StudentProfile{
		Id:  id,
		NIS: &oldNIS,
	}

	newNIS := "2024002"
	newGender := "L"
	req := &UpdateStudentProfileRequest{
		NIS:    &newNIS,
		Gender: &newGender,
	}

	mockRepo.On("FindById", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*schemas.StudentProfile")).Return(nil)

	profile, err := uc.Update(id, req)

	assert.NoError(t, err)
	assert.Equal(t, newNIS, *profile.NIS)
	assert.Equal(t, newGender, *profile.Gender)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := uc.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	uc := NewStudentProfileUseCase(mockRepo)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(errors.New("delete failed"))

	err := uc.Delete(id)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}
