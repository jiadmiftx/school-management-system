package student_profile_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/student_profile_use_case"
	"sekolah-madrasah/app/use_case/user_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"sekolah-madrasah/pkg/paginate_utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StudentProfileController struct {
	useCase     student_profile_use_case.StudentProfileUseCase
	userUseCase user_use_case.UserUseCase
}

func NewStudentProfileController(useCase student_profile_use_case.StudentProfileUseCase, userUseCase user_use_case.UserUseCase) *StudentProfileController {
	return &StudentProfileController{useCase: useCase, userUseCase: userUseCase}
}

type CreateStudentProfileDTO struct {
	UserId         string  `json:"user_id" binding:"required"`
	NIS            *string `json:"nis"`
	NISN           *string `json:"nisn"`
	BirthPlace     *string `json:"birth_place"`
	BirthDate      *string `json:"birth_date"`
	Gender         *string `json:"gender"`
	Religion       *string `json:"religion"`
	Address        *string `json:"address"`
	FatherName     *string `json:"father_name"`
	MotherName     *string `json:"mother_name"`
	GuardianName   *string `json:"guardian_name"`
	ParentPhone    *string `json:"parent_phone"`
	EnrollmentDate *string `json:"enrollment_date"`
}

type UpdateStudentProfileDTO struct {
	NIS            *string `json:"nis"`
	NISN           *string `json:"nisn"`
	BirthPlace     *string `json:"birth_place"`
	BirthDate      *string `json:"birth_date"`
	Gender         *string `json:"gender"`
	Religion       *string `json:"religion"`
	Address        *string `json:"address"`
	FatherName     *string `json:"father_name"`
	MotherName     *string `json:"mother_name"`
	GuardianName   *string `json:"guardian_name"`
	ParentPhone    *string `json:"parent_phone"`
	EnrollmentDate *string `json:"enrollment_date"`
}

// CreateStudentWithUserDTO combines user account creation with student profile
type CreateStudentWithUserDTO struct {
	// User info
	FullName string  `json:"full_name" binding:"required"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	// Student profile info
	NIS            string  `json:"nis" binding:"required"`
	NISN           *string `json:"nisn"`
	BirthPlace     *string `json:"birth_place"`
	BirthDate      string  `json:"birth_date" binding:"required"`
	Gender         *string `json:"gender"`
	Religion       *string `json:"religion"`
	Address        *string `json:"address"`
	FatherName     *string `json:"father_name"`
	MotherName     *string `json:"mother_name"`
	GuardianName   *string `json:"guardian_name"`
	ParentPhone    *string `json:"parent_phone"`
	EnrollmentDate *string `json:"enrollment_date"`
}

// GetAll godoc
// @Summary Get all student profiles in a unit
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/students [get]
func (c *StudentProfileController) GetAll(ctx *gin.Context) {
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

	// Calculate total pages
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, gin_utils.DataWithPaginateResponse{
		DataResponse: gin_utils.DataResponse{
			Message: "Students retrieved successfully",
			Data:    profiles,
		},
		Paginate: &paginate_utils.PaginateData{
			Page:       page,
			Limit:      limit,
			TotalData:  total,
			TotalPages: totalPages,
		},
	})
}

// GetById godoc
// @Summary Get student profile by ID
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param studentId path string true "Student Profile ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/students/{studentId} [get]
func (c *StudentProfileController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("studentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid student ID"})
		return
	}

	profile, err := c.useCase.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin_utils.MessageResponse{Message: "Student not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Student retrieved successfully", Data: profile})
}

// Create godoc
// @Summary Create new student profile
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateStudentProfileDTO true "Student data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/students [post]
func (c *StudentProfileController) Create(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateStudentProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	userId, err := uuid.Parse(dto.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid user ID"})
		return
	}

	req := &student_profile_use_case.CreateStudentProfileRequest{
		UserId:         userId,
		UnitId:         unitId,
		NIS:            dto.NIS,
		NISN:           dto.NISN,
		BirthPlace:     dto.BirthPlace,
		BirthDate:      dto.BirthDate,
		Gender:         dto.Gender,
		Religion:       dto.Religion,
		Address:        dto.Address,
		FatherName:     dto.FatherName,
		MotherName:     dto.MotherName,
		GuardianName:   dto.GuardianName,
		ParentPhone:    dto.ParentPhone,
		EnrollmentDate: dto.EnrollmentDate,
	}

	profile, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Student created successfully", Data: profile})
}

// Update godoc
// @Summary Update student profile
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param studentId path string true "Student Profile ID"
// @Param body body UpdateStudentProfileDTO true "Student data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/students/{studentId} [put]
func (c *StudentProfileController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("studentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid student ID"})
		return
	}

	var dto UpdateStudentProfileDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &student_profile_use_case.UpdateStudentProfileRequest{
		NIS:            dto.NIS,
		NISN:           dto.NISN,
		BirthPlace:     dto.BirthPlace,
		BirthDate:      dto.BirthDate,
		Gender:         dto.Gender,
		Religion:       dto.Religion,
		Address:        dto.Address,
		FatherName:     dto.FatherName,
		MotherName:     dto.MotherName,
		GuardianName:   dto.GuardianName,
		ParentPhone:    dto.ParentPhone,
		EnrollmentDate: dto.EnrollmentDate,
	}

	profile, err := c.useCase.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Student updated successfully", Data: profile})
}

// Delete godoc
// @Summary Delete student profile
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param studentId path string true "Student Profile ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/units/{id}/students/{studentId} [delete]
func (c *StudentProfileController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("studentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid student ID"})
		return
	}

	if err := c.useCase.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Student deleted successfully"})
}

// CreateWithUser godoc
// @Summary Create new student with user account
// @Tags Students
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateStudentWithUserDTO true "Student and user data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/students/with-user [post]
func (c *StudentProfileController) CreateWithUser(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateStudentWithUserDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	// Generate default password: NIS + DDMMYYYY
	password := dto.NIS
	if dto.BirthDate != "" {
		// Parse date and format as DDMMYYYY
		if len(dto.BirthDate) >= 10 {
			// Assuming format YYYY-MM-DD
			password += dto.BirthDate[8:10] + dto.BirthDate[5:7] + dto.BirthDate[0:4]
		}
	}

	// Generate email if not provided
	email := ""
	if dto.Email != nil && *dto.Email != "" {
		email = *dto.Email
	} else {
		email = dto.NIS + "@siswa.sekolah.id"
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

	// Create student profile
	req := &student_profile_use_case.CreateStudentProfileRequest{
		UserId:         user.Id,
		UnitId:         unitId,
		NIS:            &dto.NIS,
		NISN:           dto.NISN,
		BirthPlace:     dto.BirthPlace,
		BirthDate:      &dto.BirthDate,
		Gender:         dto.Gender,
		Religion:       dto.Religion,
		Address:        dto.Address,
		FatherName:     dto.FatherName,
		MotherName:     dto.MotherName,
		GuardianName:   dto.GuardianName,
		ParentPhone:    dto.ParentPhone,
		EnrollmentDate: dto.EnrollmentDate,
	}

	profile, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Gagal membuat profil siswa: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{
		Message: "Siswa berhasil dibuat",
		Data: gin.H{
			"profile":  profile,
			"user_id":  user.Id,
			"email":    email,
			"password": password,
		},
	})
}
