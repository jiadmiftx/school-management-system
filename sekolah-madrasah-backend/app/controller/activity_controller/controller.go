package activity_controller

import (
	"net/http"
	"sekolah-madrasah/app/use_case/activity_use_case"
	"sekolah-madrasah/pkg/gin_utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

type ActivityController struct {
	useCase activity_use_case.ActivityUseCase
}

func NewActivityController(useCase activity_use_case.ActivityUseCase) *ActivityController {
	return &ActivityController{useCase: useCase}
}

type CreateActivityDTO struct {
	Name            string   `json:"name" binding:"required"`
	Type            string   `json:"type" binding:"required"`
	Category        *string  `json:"category"`
	Description     *string  `json:"description"`
	StartDate       *string  `json:"start_date"` // Format: YYYY-MM-DD
	EndDate         *string  `json:"end_date"`   // Format: YYYY-MM-DD
	RecurrenceType  string   `json:"recurrence_type"`
	RecurrenceDays  []int64  `json:"recurrence_days"` // Array of days
	StartTime       *string  `json:"start_time"`
	EndTime         *string  `json:"end_time"`
	Location        *string  `json:"location"`
	MaxParticipants *int     `json:"max_participants"`
	Fee             *float64 `json:"fee"`
}

type UpdateActivityDTO struct {
	Name            *string  `json:"name"`
	Type            *string  `json:"type"`
	Category        *string  `json:"category"`
	Description     *string  `json:"description"`
	StartDate       *string  `json:"start_date"` // Format: YYYY-MM-DD
	EndDate         *string  `json:"end_date"`   // Format: YYYY-MM-DD
	RecurrenceType  *string  `json:"recurrence_type"`
	RecurrenceDays  []int64  `json:"recurrence_days"`
	StartTime       *string  `json:"start_time"`
	EndTime         *string  `json:"end_time"`
	Location        *string  `json:"location"`
	MaxParticipants *int     `json:"max_participants"`
	Fee             *float64 `json:"fee"`
	IsActive        *bool    `json:"is_active"`
}

type AssignTeacherDTO struct {
	TeacherProfileId string `json:"teacher_profile_id" binding:"required"`
	Role             string `json:"role"`
}

type EnrollStudentDTO struct {
	StudentProfileId string `json:"student_profile_id" binding:"required"`
	IsMandatory      bool   `json:"is_mandatory"`
}

// GetAll godoc
// @Summary Get all activities in a unit
// @Tags Activities
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param type query string false "Filter by type (ekstrakurikuler/kajian/event)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/activities [get]
func (c *ActivityController) GetAll(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	activityType := ctx.Query("type")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	activities, total, err := c.useCase.GetByUnitId(unitId, activityType, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{
		Message: "Activities retrieved successfully",
		Data: gin.H{
			"data":  activities,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetById godoc
// @Summary Get activity by ID
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/activities/{activityId} [get]
func (c *ActivityController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	activity, err := c.useCase.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin_utils.MessageResponse{Message: "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Activity retrieved successfully", Data: activity})
}

// Create godoc
// @Summary Create new activity
// @Tags Activities
// @Security BearerAuth
// @Param id path string true "Unit ID"
// @Param body body CreateActivityDTO true "Activity data"
// @Success 201 {object} gin_utils.DataResponse
// @Router /api/v1/units/{id}/activities [post]
func (c *ActivityController) Create(ctx *gin.Context) {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid unit ID"})
		return
	}

	var dto CreateActivityDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &activity_use_case.CreateActivityRequest{
		UnitId:          unitId,
		Name:            dto.Name,
		Type:            dto.Type,
		Category:        dto.Category,
		Description:     dto.Description,
		RecurrenceType:  dto.RecurrenceType,
		RecurrenceDays:  dto.RecurrenceDays,
		StartTime:       dto.StartTime,
		EndTime:         dto.EndTime,
		Location:        dto.Location,
		MaxParticipants: dto.MaxParticipants,
		Fee:             dto.Fee,
	}

	// Parse dates if provided
	if dto.StartDate != nil {
		if t, err := parseDate(*dto.StartDate); err == nil {
			req.StartDate = &t
		}
	}
	if dto.EndDate != nil {
		if t, err := parseDate(*dto.EndDate); err == nil {
			req.EndDate = &t
		}
	}

	activity, err := c.useCase.Create(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin_utils.DataResponse{Message: "Activity created successfully", Data: activity})
}

// Update godoc
// @Summary Update activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Param body body UpdateActivityDTO true "Activity data"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/activities/{activityId} [put]
func (c *ActivityController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	var dto UpdateActivityDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	req := &activity_use_case.UpdateActivityRequest{
		Name:            dto.Name,
		Type:            dto.Type,
		Category:        dto.Category,
		Description:     dto.Description,
		RecurrenceType:  dto.RecurrenceType,
		RecurrenceDays:  dto.RecurrenceDays,
		StartTime:       dto.StartTime,
		EndTime:         dto.EndTime,
		Location:        dto.Location,
		MaxParticipants: dto.MaxParticipants,
		Fee:             dto.Fee,
		IsActive:        dto.IsActive,
	}

	// Parse dates if provided
	if dto.StartDate != nil {
		if t, err := parseDate(*dto.StartDate); err == nil {
			req.StartDate = &t
		}
	}
	if dto.EndDate != nil {
		if t, err := parseDate(*dto.EndDate); err == nil {
			req.EndDate = &t
		}
	}

	activity, err := c.useCase.Update(id, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Activity updated successfully", Data: activity})
}

// Delete godoc
// @Summary Delete activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/activities/{activityId} [delete]
func (c *ActivityController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	if err := c.useCase.Delete(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Activity deleted successfully"})
}

// AssignTeacher godoc
// @Summary Assign teacher to activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Param body body AssignTeacherDTO true "Teacher assignment data"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/activities/{activityId}/teachers [post]
func (c *ActivityController) AssignTeacher(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
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

	req := &activity_use_case.AssignTeacherRequest{
		ActivityId:       activityId,
		TeacherProfileId: teacherProfileId,
		Role:             dto.Role,
	}

	if err := c.useCase.AssignTeacher(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Teacher assigned to activity successfully"})
}

// RemoveTeacher godoc
// @Summary Remove teacher from activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Param teacherId path string true "Teacher Profile ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/activities/{activityId}/teachers/{teacherId} [delete]
func (c *ActivityController) RemoveTeacher(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	teacherId, err := uuid.Parse(ctx.Param("teacherId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid teacher ID"})
		return
	}

	if err := c.useCase.RemoveTeacher(activityId, teacherId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Teacher removed from activity successfully"})
}

// GetTeachers godoc
// @Summary Get teachers assigned to activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/activities/{activityId}/teachers [get]
func (c *ActivityController) GetTeachers(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	teachers, err := c.useCase.GetTeachers(activityId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Teachers retrieved successfully", Data: teachers})
}

// EnrollStudent godoc
// @Summary Enroll student in activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Param body body EnrollStudentDTO true "Student enrollment data"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/activities/{activityId}/students [post]
func (c *ActivityController) EnrollStudent(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
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

	req := &activity_use_case.EnrollStudentRequest{
		ActivityId:       activityId,
		StudentProfileId: studentProfileId,
		IsMandatory:      dto.IsMandatory,
	}

	if err := c.useCase.EnrollStudent(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Student enrolled in activity successfully"})
}

// RemoveStudent godoc
// @Summary Remove student from activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Param studentId path string true "Student Profile ID"
// @Success 200 {object} gin_utils.MessageResponse
// @Router /api/v1/activities/{activityId}/students/{studentId} [delete]
func (c *ActivityController) RemoveStudent(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	studentId, err := uuid.Parse(ctx.Param("studentId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid student ID"})
		return
	}

	if err := c.useCase.RemoveStudent(activityId, studentId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.MessageResponse{Message: "Student removed from activity successfully"})
}

// GetStudents godoc
// @Summary Get students enrolled in activity
// @Tags Activities
// @Security BearerAuth
// @Param activityId path string true "Activity ID"
// @Success 200 {object} gin_utils.DataResponse
// @Router /api/v1/activities/{activityId}/students [get]
func (c *ActivityController) GetStudents(ctx *gin.Context) {
	activityId, err := uuid.Parse(ctx.Param("activityId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin_utils.MessageResponse{Message: "Invalid activity ID"})
		return
	}

	students, err := c.useCase.GetStudents(activityId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin_utils.MessageResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin_utils.DataResponse{Message: "Students retrieved successfully", Data: students})
}
