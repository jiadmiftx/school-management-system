package class_enrollment_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/class_enrollment_use_case"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/gin_utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClassEnrollmentController struct {
	useCase class_enrollment_use_case.ClassEnrollmentUseCase
}

func NewClassEnrollmentController(useCase class_enrollment_use_case.ClassEnrollmentUseCase) *ClassEnrollmentController {
	return &ClassEnrollmentController{useCase: useCase}
}

type EnrollStudentDTO struct {
	StudentProfileId string  `json:"student_profile_id" binding:"required"`
	AcademicYear     string  `json:"academic_year" binding:"required"`
	EnrolledAt       *string `json:"enrolled_at"`
}

type UpdateStatusDTO struct {
	Status string  `json:"status" binding:"required"`
	Notes  *string `json:"notes"`
}

type TransferDTO struct {
	NewClassId string `json:"new_class_id" binding:"required"`
}

// GetByClass godoc
// @Summary Get all students enrolled in a class
// @Tags Class Enrollments
// @Param id path string true "Unit ID"
// @Param classId path string true "Class ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes/{classId}/students [get]
func (c *ClassEnrollmentController) GetByClass(ctx *gin.Context) {
	classId, err := uuid.Parse(ctx.Param("classId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid class ID"})
		return
	}

	enrollments, err := c.useCase.GetByClassId(classId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Enrollments retrieved successfully", Data: enrollments})
}

// Enroll godoc
// @Summary Enroll a student in a class
// @Tags Class Enrollments
// @Param id path string true "Unit ID"
// @Param classId path string true "Class ID"
// @Param body body EnrollStudentDTO true "Enrollment data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes/{classId}/enroll [post]
func (c *ClassEnrollmentController) Enroll(ctx *gin.Context) {
	classId, err := uuid.Parse(ctx.Param("classId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid class ID"})
		return
	}

	var dto EnrollStudentDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	studentProfileId, err := uuid.Parse(dto.StudentProfileId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid student profile ID"})
		return
	}

	req := &class_enrollment_use_case.EnrollStudentRequest{
		StudentProfileId: studentProfileId,
		ClassId:          classId,
		AcademicYear:     dto.AcademicYear,
		EnrolledAt:       dto.EnrolledAt,
	}

	enrollment, err := c.useCase.Enroll(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Student enrolled successfully", Data: enrollment})
}

// UpdateStatus godoc
// @Summary Update enrollment status
// @Tags Class Enrollments
// @Param enrollmentId path string true "Enrollment ID"
// @Param body body UpdateStatusDTO true "Status data"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/class-enrollments/{enrollmentId} [put]
func (c *ClassEnrollmentController) UpdateStatus(ctx *gin.Context) {
	enrollmentId, err := uuid.Parse(ctx.Param("enrollmentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid enrollment ID"})
		return
	}

	var dto UpdateStatusDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	status := schemas.ClassEnrollmentStatus(dto.Status)
	if status != schemas.EnrollmentStatusActive &&
		status != schemas.EnrollmentStatusGraduated &&
		status != schemas.EnrollmentStatusTransferred &&
		status != schemas.EnrollmentStatusDropped {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid status"})
		return
	}

	if err := c.useCase.UpdateStatus(enrollmentId, status, dto.Notes); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Enrollment status updated successfully"})
}

// Transfer godoc
// @Summary Transfer student to another class
// @Tags Class Enrollments
// @Param enrollmentId path string true "Enrollment ID"
// @Param body body TransferDTO true "Transfer data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/class-enrollments/{enrollmentId}/transfer [post]
func (c *ClassEnrollmentController) Transfer(ctx *gin.Context) {
	enrollmentId, err := uuid.Parse(ctx.Param("enrollmentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid enrollment ID"})
		return
	}

	var dto TransferDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	newClassId, err := uuid.Parse(dto.NewClassId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid new class ID"})
		return
	}

	enrollment, err := c.useCase.Transfer(enrollmentId, newClassId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Student transferred successfully", Data: enrollment})
}

// Remove godoc
// @Summary Remove student from class
// @Tags Class Enrollments
// @Param enrollmentId path string true "Enrollment ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/class-enrollments/{enrollmentId} [delete]
func (c *ClassEnrollmentController) Remove(ctx *gin.Context) {
	enrollmentId, err := uuid.Parse(ctx.Param("enrollmentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid enrollment ID"})
		return
	}

	if err := c.useCase.Remove(enrollmentId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Enrollment removed successfully"})
}
