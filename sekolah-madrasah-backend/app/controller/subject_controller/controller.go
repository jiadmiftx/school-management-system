package subject_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/subject_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubjectController struct {
	useCase subject_use_case.SubjectUseCase
}

func NewSubjectController(useCase subject_use_case.SubjectUseCase) *SubjectController {
	return &SubjectController{useCase: useCase}
}

type CreateSubjectDTO struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Category    string  `json:"category"`
	Description *string `json:"description"`
}

type UpdateSubjectDTO struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Category    *string `json:"category"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type AssignTeacherDTO struct {
	TeacherProfileId string `json:"teacher_profile_id" binding:"required"`
	IsPrimary        bool   `json:"is_primary"`
}

// GetAll godoc
// @Summary Get all subjects in a unit
// @Tags Subjects
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/subjects [get]
func (c *SubjectController) GetAll(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	subjects, total, err := c.useCase.GetByUnitId(unitId, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "Subjects retrieved successfully",
		Data: gin.H{
			"data":  subjects,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetById godoc
// @Summary Get subject by ID
// @Tags Subjects
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param subjectId path string true "Subject ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/subjects/{subjectId} [get]
func (c *SubjectController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("subjectId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid subject ID"})
		return
	}

	subject, err := c.useCase.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin_utils.MessageResponse{Message: "Subject not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Subject retrieved successfully", Data: subject})
}

// Create godoc
// @Summary Create new subject
// @Tags Subjects
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateSubjectDTO true "Subject data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/subjects [post]
func (c *SubjectController) Create(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateSubjectDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &subject_use_case.CreateSubjectRequest{
		UnitId:      unitId,
		Name:        dto.Name,
		Code:        dto.Code,
		Category:    dto.Category,
		Description: dto.Description,
	}

	subject, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Subject created successfully", Data: subject})
}

// Update godoc
// @Summary Update subject
// @Tags Subjects
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param subjectId path string true "Subject ID"
// @Param body body UpdateSubjectDTO true "Subject data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/subjects/{subjectId} [put]
func (c *SubjectController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("subjectId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid subject ID"})
		return
	}

	var dto UpdateSubjectDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &subject_use_case.UpdateSubjectRequest{
		Name:        dto.Name,
		Code:        dto.Code,
		Category:    dto.Category,
		Description: dto.Description,
		IsActive:    dto.IsActive,
	}

	subject, err := c.useCase.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Subject updated successfully", Data: subject})
}

// Delete godoc
// @Summary Delete subject
// @Tags Subjects
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param subjectId path string true "Subject ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/units/{id}/subjects/{subjectId} [delete]
func (c *SubjectController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("subjectId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid subject ID"})
		return
	}

	if err := c.useCase.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Subject deleted successfully"})
}

// AssignTeacher godoc
// @Summary Assign teacher to subject
// @Tags Subjects
// @Security BearerAuth
// @Param subjectId path string true "Subject ID"
// @Param body body AssignTeacherDTO true "Teacher assignment data"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/subjects/{subjectId}/teachers [post]
func (c *SubjectController) AssignTeacher(ctx *gin.Context) {
	subjectId, err := uuid.Parse(ctx.Param("subjectId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid subject ID"})
		return
	}

	var dto AssignTeacherDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	teacherProfileId, err := uuid.Parse(dto.TeacherProfileId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher profile ID"})
		return
	}

	req := &subject_use_case.AssignTeacherRequest{
		TeacherProfileId: teacherProfileId,
		SubjectId:        subjectId,
		IsPrimary:        dto.IsPrimary,
	}

	if err := c.useCase.AssignTeacher(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Teacher assigned to subject successfully"})
}

// RemoveTeacher godoc
// @Summary Remove teacher from subject
// @Tags Subjects
// @Security BearerAuth
// @Param subjectId path string true "Subject ID"
// @Param teacherId path string true "Teacher Profile ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/subjects/{subjectId}/teachers/{teacherId} [delete]
func (c *SubjectController) RemoveTeacher(ctx *gin.Context) {
	subjectId, err := uuid.Parse(ctx.Param("subjectId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid subject ID"})
		return
	}

	teacherId, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	if err := c.useCase.RemoveTeacher(teacherId, subjectId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Teacher removed from subject successfully"})
}

// GetByTeacher godoc
// @Summary Get subjects assigned to a teacher
// @Tags Subjects
// @Security BearerAuth
// @Param teacherId path string true "Teacher Profile ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/teachers/{teacherId}/subjects [get]
func (c *SubjectController) GetByTeacher(ctx *gin.Context) {
	teacherId, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	subjects, err := c.useCase.GetByTeacher(teacherId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Teacher subjects retrieved successfully", Data: subjects})
}
