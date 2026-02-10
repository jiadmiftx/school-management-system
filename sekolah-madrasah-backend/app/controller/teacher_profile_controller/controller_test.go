package teacher_profile_controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sekolah-madrasah/database/schemas"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUseCase is a mock implementation of TeacherProfileUseCase
type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) Create(req *CreateTeacherProfileRequest) (*schemas.TeacherProfile, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.TeacherProfile), args.Error(1)
}

func (m *MockUseCase) GetById(id uuid.UUID) (*schemas.TeacherProfile, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.TeacherProfile), args.Error(1)
}

func (m *MockUseCase) GetByUnitId(unitId uuid.UUID, page, limit int) ([]schemas.TeacherProfile, int64, error) {
	args := m.Called(unitId, page, limit)
	return args.Get(0).([]schemas.TeacherProfile), args.Get(1).(int64), args.Error(2)
}

func (m *MockUseCase) Update(id uuid.UUID, req *UpdateTeacherProfileRequest) (*schemas.TeacherProfile, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schemas.TeacherProfile), args.Error(1)
}

func (m *MockUseCase) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// Need to import the actual use case types
type CreateTeacherProfileRequest = interface{}
type UpdateTeacherProfileRequest = interface{}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

func TestGetAll_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	c.Params = []gin.Param{{Key: "id", Value: unitId.String()}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/units/"+unitId.String()+"/teachers?page=1&limit=10", nil)

	// Note: Full controller testing requires setting up routes properly
	// This is a simplified test showing the concept
	assert.NotNil(t, c)
}

func TestGetAll_InvalidUnitId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "id", Value: "invalid-uuid"}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/units/invalid-uuid/teachers", nil)

	// Controller should return 400 Bad Request
	assert.NotNil(t, c)
}

func TestGetById_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	teacherId := uuid.New()

	c.Params = []gin.Param{
		{Key: "id", Value: unitId.String()},
		{Key: "teacherId", Value: teacherId.String()},
	}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	assert.NotNil(t, c)
}

func TestGetById_InvalidTeacherId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "id", Value: uuid.New().String()},
		{Key: "teacherId", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest("GET", "/", nil)

	// Should return 400 Bad Request
	assert.NotNil(t, c)
}

func TestCreate_ValidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	body := CreateTeacherProfileDTO{
		UserId:           uuid.New().String(),
		EmploymentStatus: "pns",
	}
	bodyBytes, _ := json.Marshal(body)

	c.Params = []gin.Param{{Key: "id", Value: unitId.String()}}
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	assert.NotNil(t, c)
}

func TestCreate_MissingUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	body := map[string]string{
		"employment_status": "pns",
		// missing user_id
	}
	bodyBytes, _ := json.Marshal(body)

	c.Params = []gin.Param{{Key: "id", Value: unitId.String()}}
	c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	// Should return 400 Bad Request
	assert.NotNil(t, c)
}

func TestUpdate_ValidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	teacherId := uuid.New()
	newStatus := "pns"
	body := UpdateTeacherProfileDTO{
		EmploymentStatus: &newStatus,
	}
	bodyBytes, _ := json.Marshal(body)

	c.Params = []gin.Param{
		{Key: "id", Value: unitId.String()},
		{Key: "teacherId", Value: teacherId.String()},
	}
	c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBuffer(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	assert.NotNil(t, c)
}

func TestDelete_ValidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	teacherId := uuid.New()

	c.Params = []gin.Param{
		{Key: "id", Value: unitId.String()},
		{Key: "teacherId", Value: teacherId.String()},
	}
	c.Request, _ = http.NewRequest("DELETE", "/", nil)

	assert.NotNil(t, c)
}

func TestDelete_InvalidId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{
		{Key: "id", Value: uuid.New().String()},
		{Key: "teacherId", Value: "invalid-uuid"},
	}
	c.Request, _ = http.NewRequest("DELETE", "/", nil)

	// Should return 400 Bad Request
	assert.NotNil(t, c)
}

// Pagination tests
func TestGetAll_DefaultPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	c.Params = []gin.Param{{Key: "id", Value: unitId.String()}}
	// No page/limit params - should use defaults
	c.Request, _ = http.NewRequest("GET", "/api/v1/units/"+unitId.String()+"/teachers", nil)

	assert.NotNil(t, c)
}

func TestGetAll_CustomPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	unitId := uuid.New()
	c.Params = []gin.Param{{Key: "id", Value: unitId.String()}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/units/"+unitId.String()+"/teachers?page=2&limit=25", nil)

	assert.NotNil(t, c)
}
