package class_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/class_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ClassController struct {
	useCase class_use_case.ClassUseCase
}

func NewClassController(useCase class_use_case.ClassUseCase) *ClassController {
	return &ClassController{useCase: useCase}
}

type CreateClassDTO struct {
	Name              string  `json:"name" binding:"required"`
	Level             int     `json:"level" binding:"required"`
	AcademicYear      string  `json:"academic_year" binding:"required"`
	HomeroomTeacherId *string `json:"homeroom_teacher_id"`
	Capacity          int     `json:"capacity"`
}

type UpdateClassDTO struct {
	Name              *string `json:"name"`
	Level             *int    `json:"level"`
	AcademicYear      *string `json:"academic_year"`
	HomeroomTeacherId *string `json:"homeroom_teacher_id"`
	Capacity          *int    `json:"capacity"`
	IsActive          *bool   `json:"is_active"`
}

// GetAll godoc
// @Summary Get all classes in a unit
// @Tags Classes
// @Param id path string true "Unit ID"
// @Param academic_year query string false "Academic Year"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes [get]
func (c *ClassController) GetAll(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	academicYear := ctx.Query("academic_year")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	classes, total, err := c.useCase.GetByUnitId(unitId, academicYear, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "Classes retrieved successfully",
		Data: gin.H{
			"data":  classes,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetById godoc
// @Summary Get class by ID
// @Tags Classes
// @Param id path string true "Unit ID"
// @Param classId path string true "Class ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes/{classId} [get]
func (c *ClassController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("classId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid class ID"})
		return
	}

	class, err := c.useCase.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin_utils.MessageResponse{Message: "Class not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Class retrieved successfully", Data: class})
}

// Create godoc
// @Summary Create new class
// @Tags Classes
// @Param id path string true "Unit ID"
// @Param body body CreateClassDTO true "Class data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes [post]
func (c *ClassController) Create(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateClassDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var homeroomTeacherId *uuid.UUID
	if dto.HomeroomTeacherId != nil {
		id, err := uuid.Parse(*dto.HomeroomTeacherId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid homeroom teacher ID"})
			return
		}
		homeroomTeacherId = &id
	}

	req := &class_use_case.CreateClassRequest{
		UnitId:            unitId,
		Name:              dto.Name,
		Level:             dto.Level,
		AcademicYear:      dto.AcademicYear,
		HomeroomTeacherId: homeroomTeacherId,
		Capacity:          dto.Capacity,
	}

	class, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Class created successfully", Data: class})
}

// Update godoc
// @Summary Update class
// @Tags Classes
// @Param id path string true "Unit ID"
// @Param classId path string true "Class ID"
// @Param body body UpdateClassDTO true "Class data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/classes/{classId} [put]
func (c *ClassController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("classId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid class ID"})
		return
	}

	var dto UpdateClassDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	var homeroomTeacherId *uuid.UUID
	if dto.HomeroomTeacherId != nil {
		tid, err := uuid.Parse(*dto.HomeroomTeacherId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid homeroom teacher ID"})
			return
		}
		homeroomTeacherId = &tid
	}

	req := &class_use_case.UpdateClassRequest{
		Name:              dto.Name,
		Level:             dto.Level,
		AcademicYear:      dto.AcademicYear,
		HomeroomTeacherId: homeroomTeacherId,
		Capacity:          dto.Capacity,
		IsActive:          dto.IsActive,
	}

	class, err := c.useCase.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Class updated successfully", Data: class})
}

// Delete godoc
// @Summary Delete class
// @Tags Classes
// @Param id path string true "Unit ID"
// @Param classId path string true "Class ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/units/{id}/classes/{classId} [delete]
func (c *ClassController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("classId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid class ID"})
		return
	}

	if err := c.useCase.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Class deleted successfully"})
}
