package teacher_profile_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/teacher_profile_use_case"
	"sekolah-madrasah/app/use_case/user_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeacherProfileController struct {
	useCase     teacher_profile_use_case.TeacherProfileUseCase
	userUseCase user_use_case.UserUseCase
}

func NewTeacherProfileController(useCase teacher_profile_use_case.TeacherProfileUseCase, userUseCase user_use_case.UserUseCase) *TeacherProfileController {
	return &TeacherProfileController{useCase: useCase, userUseCase: userUseCase}
}

type CreateTeacherProfileDTO struct {
	UserId           string  `json:"user_id" binding:"required"`
	NIP              *string `json:"nip"`
	NUPTK            *string `json:"nuptk"`
	EducationLevel   *string `json:"education_level"`
	EducationMajor   *string `json:"education_major"`
	EmploymentStatus string  `json:"employment_status"`
	JoinDate         *string `json:"join_date"`
}

type UpdateTeacherProfileDTO struct {
	NIP              *string `json:"nip"`
	NUPTK            *string `json:"nuptk"`
	EducationLevel   *string `json:"education_level"`
	EducationMajor   *string `json:"education_major"`
	EmploymentStatus *string `json:"employment_status"`
	JoinDate         *string `json:"join_date"`
}

// CreateTeacherWithUserDTO combines user account creation with teacher profile
type CreateTeacherWithUserDTO struct {
	// User info
	FullName string  `json:"full_name" binding:"required"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	// Teacher profile info
	NIP              *string `json:"nip"`
	NUPTK            *string `json:"nuptk"`
	EducationLevel   *string `json:"education_level"`
	EducationMajor   *string `json:"education_major"`
	EmploymentStatus string  `json:"employment_status" binding:"required"`
	JoinDate         *string `json:"join_date"`
}

// GetAll godoc
// @Summary Get all teacher profiles in a unit
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/teachers [get]
func (c *TeacherProfileController) GetAll(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	profiles, total, err := c.useCase.GetByUnitId(unitId, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "Teachers retrieved successfully",
		Data: gin.H{
			"data":  profiles,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetById godoc
// @Summary Get teacher profile by ID
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param teacherId path string true "Teacher Profile ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/teachers/{teacherId} [get]
func (c *TeacherProfileController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	profile, err := c.useCase.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin_utils.MessageResponse{Message: "Teacher not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Teacher retrieved successfully", Data: profile})
}

// Create godoc
// @Summary Create new teacher profile
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateTeacherProfileDTO true "Teacher data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/teachers [post]
func (c *TeacherProfileController) Create(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateTeacherProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid user ID"})
		return
	}

	req := &teacher_profile_use_case.CreateTeacherProfileRequest{
		UserId:           userId,
		UnitId:           unitId,
		NIP:              dto.NIP,
		NUPTK:            dto.NUPTK,
		EducationLevel:   dto.EducationLevel,
		EducationMajor:   dto.EducationMajor,
		EmploymentStatus: dto.EmploymentStatus,
		JoinDate:         dto.JoinDate,
	}

	profile, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Teacher created successfully", Data: profile})
}

// CreateWithUser godoc
// @Summary Create new teacher with user account
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateTeacherWithUserDTO true "Teacher and user data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/teachers/with-user [post]
func (c *TeacherProfileController) CreateWithUser(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateTeacherWithUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	// Generate default password: guru + random 6 chars or NIP if available
	password := "guru123456"
	if dto.NIP != nil && *dto.NIP != "" {
		password = *dto.NIP
	}

	// Generate email if not provided
	email := ""
	if dto.Email != nil && *dto.Email != "" {
		email = *dto.Email
	} else {
		// Use NIP or generate random for email
		if dto.NIP != nil && *dto.NIP != "" {
			email = *dto.NIP + "@guru.sekolah.id"
		} else {
			email = uuid.New().String()[:8] + "@guru.sekolah.id"
		}
	}

	// Get phone
	phone := ""
	if dto.Phone != nil {
		phone = *dto.Phone
	}

	// Create user first
	user, code, err := c.userUseCase.CreateUser(ctx.Request.Context(), user_use_case.CreateUserRequest{
		Email:    email,
		Password: password,
		FullName: dto.FullName,
		Phone:    phone,
	})
	if err != nil {
		ctx.JSON(code, gin_utils.MessageResponse{Message: "Gagal membuat akun: " + err.Error()})
		return
	}

	// Create teacher profile
	req := &teacher_profile_use_case.CreateTeacherProfileRequest{
		UserId:           user.Id,
		UnitId:           unitId,
		NIP:              dto.NIP,
		NUPTK:            dto.NUPTK,
		EducationLevel:   dto.EducationLevel,
		EducationMajor:   dto.EducationMajor,
		EmploymentStatus: dto.EmploymentStatus,
		JoinDate:         dto.JoinDate,
	}

	profile, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Gagal membuat profil guru: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{
		Message: "Guru berhasil dibuat",
		Data: gin.H{
			"profile":  profile,
			"user_id":  user.Id,
			"email":    email,
			"password": password,
		},
	})
}

// Update godoc
// @Summary Update teacher profile
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param teacherId path string true "Teacher Profile ID"
// @Param body body UpdateTeacherProfileDTO true "Teacher data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/teachers/{teacherId} [put]
func (c *TeacherProfileController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	var dto UpdateTeacherProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &teacher_profile_use_case.UpdateTeacherProfileRequest{
		NIP:              dto.NIP,
		NUPTK:            dto.NUPTK,
		EducationLevel:   dto.EducationLevel,
		EducationMajor:   dto.EducationMajor,
		EmploymentStatus: dto.EmploymentStatus,
		JoinDate:         dto.JoinDate,
	}

	profile, err := c.useCase.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Teacher updated successfully", Data: profile})
}

// Delete godoc
// @Summary Delete teacher profile
// @Tags Teachers
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param teacherId path string true "Teacher Profile ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/units/{id}/teachers/{teacherId} [delete]
func (c *TeacherProfileController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	if err := c.useCase.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Teacher deleted successfully"})
}
