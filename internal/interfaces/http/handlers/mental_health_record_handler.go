package handlers

import (
	"net/http"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/interfaces/http/middleware"
	"github.com/atdevten/peace/internal/pkg/timeutil"

	"github.com/gin-gonic/gin"
)

type MentalHealthRecordHandler struct {
	recordUseCase usecases.MentalHealthRecordUseCase
}

type CreateMentalHealthRecordRequest struct {
	HappyLevel  int     `json:"happy_level"`
	EnergyLevel int     `json:"energy_level"`
	Notes       *string `json:"notes,omitempty"`
	Status      string  `json:"status"`
}

type UpdateMentalHealthRecordRequest struct {
	HappyLevel  *int    `json:"happy_level"`
	EnergyLevel *int    `json:"energy_level"`
	Notes       *string `json:"notes,omitempty"`
	Status      string  `json:"status"`
}

type MentalHealthRecordResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	HappyLevel  int     `json:"happy_level"`
	EnergyLevel int     `json:"energy_level"`
	Notes       *string `json:"notes,omitempty"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type HeatmapDataPoint struct {
	HappyLevel  int `json:"happy_level"`
	EnergyLevel int `json:"energy_level"`
	Count       int `json:"count"`
}

type DateRange struct {
	StartedAt *string `json:"started_at,omitempty"`
	EndedAt   *string `json:"ended_at,omitempty"`
}

type MentalHealthHeatmapResponse struct {
	Data         map[string]HeatmapDataPoint `json:"data"`
	TotalRecords int                         `json:"total_records"`
	DateRange    DateRange                   `json:"date_range"`
}

func NewMentalHealthRecordHandler(recordUseCase usecases.MentalHealthRecordUseCase) *MentalHealthRecordHandler {
	return &MentalHealthRecordHandler{
		recordUseCase: recordUseCase,
	}
}

func (h *MentalHealthRecordHandler) Create(c *gin.Context) {
	var req CreateMentalHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Create command
	command, err := commands.NewCreateMentalHealthRecordCommand(
		userID.String(),
		req.HappyLevel,
		req.EnergyLevel,
		req.Notes,
		req.Status,
	)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	record, err := h.recordUseCase.Create(ctx, command)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	response := MentalHealthRecordResponse{
		ID:          record.ID().String(),
		UserID:      record.UserID().String(),
		HappyLevel:  record.HappyLevel().Value(),
		EnergyLevel: record.EnergyLevel().Value(),
		Notes:       record.Notes(),
		Status:      record.Status().String(),
		CreatedAt:   timeutil.FormatTime(record.CreatedAt()),
		UpdatedAt:   timeutil.FormatTime(record.UpdatedAt()),
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    CodeSuccess,
		Message: "Mental health record created successfully",
		Data:    response,
	})
}

func (h *MentalHealthRecordHandler) Update(c *gin.Context) {
	var req UpdateMentalHealthRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Get record ID from URL
	recordID := c.Param("id")
	if recordID == "" {
		Error(c, CodeBadRequest, "Record ID is required")
		return
	}

	// Validate required fields for PUT
	if req.HappyLevel == nil || req.EnergyLevel == nil {
		Error(c, CodeBadRequest, "happy_level and energy_level are required")
		return
	}

	// Create command
	command, err := commands.NewUpdateMentalHealthRecordCommand(
		recordID,
		userID.String(),
		*req.HappyLevel,
		*req.EnergyLevel,
		req.Notes,
		req.Status,
	)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	record, err := h.recordUseCase.Update(ctx, command)
	if err != nil {
		if err.Error() == "record not found" {
			Error(c, CodeNotFound, "Record not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this record" {
			Error(c, CodeForbidden, "You don't have permission to update this record")
			return
		}
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	response := MentalHealthRecordResponse{
		ID:          record.ID().String(),
		UserID:      record.UserID().String(),
		HappyLevel:  record.HappyLevel().Value(),
		EnergyLevel: record.EnergyLevel().Value(),
		Notes:       record.Notes(),
		Status:      record.Status().String(),
		CreatedAt:   timeutil.FormatTime(record.CreatedAt()),
		UpdatedAt:   timeutil.FormatTime(record.UpdatedAt()),
	}

	Success(c, "Mental health record updated successfully", response)
}

func (h *MentalHealthRecordHandler) Delete(c *gin.Context) {
	// Get record ID from URL
	recordID := c.Param("id")
	if recordID == "" {
		Error(c, CodeBadRequest, "Record ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Create command
	command, err := commands.NewDeleteMentalHealthRecordCommand(recordID, userID.String())
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	err = h.recordUseCase.Delete(ctx, command)
	if err != nil {
		if err.Error() == "record not found" {
			Error(c, CodeNotFound, "Record not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this record" {
			Error(c, CodeForbidden, "You don't have permission to delete this record")
			return
		}
		Error(c, CodeBadRequest, err.Error())
		return
	}

	Success(c, "Mental health record deleted successfully", nil)
}

func (h *MentalHealthRecordHandler) GetByID(c *gin.Context) {
	// Get record ID from URL
	recordID := c.Param("id")
	if recordID == "" {
		Error(c, CodeBadRequest, "Record ID is required")
		return
	}

	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	record, err := h.recordUseCase.GetByID(ctx, recordID, userID.String())
	if err != nil {
		if err.Error() == "record not found" {
			Error(c, CodeNotFound, "Record not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this record" {
			Error(c, CodeForbidden, "You don't have permission to view this record")
			return
		}
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	response := MentalHealthRecordResponse{
		ID:          record.ID().String(),
		UserID:      record.UserID().String(),
		HappyLevel:  record.HappyLevel().Value(),
		EnergyLevel: record.EnergyLevel().Value(),
		Notes:       record.Notes(),
		Status:      record.Status().String(),
		CreatedAt:   timeutil.FormatTime(record.CreatedAt()),
		UpdatedAt:   timeutil.FormatTime(record.UpdatedAt()),
	}

	Success(c, "Mental health record retrieved successfully", response)
}

func (h *MentalHealthRecordHandler) GetByCondition(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Get query parameters
	startedAt := c.Query("started_at")
	endedAt := c.Query("ended_at")

	var startedAtPtr *string
	var endedAtPtr *string

	if startedAt != "" {
		startedAtPtr = &startedAt
	}
	if endedAt != "" {
		endedAtPtr = &endedAt
	}

	// Execute use case with pagination and order via context keys
	// For now, pass through via query params handled in usecase next iteration; keep signature
	ctx := c.Request.Context()
	records, err := h.recordUseCase.GetByCondition(ctx, userID.String(), startedAtPtr, endedAtPtr)
	if err != nil {
		if err.Error() == "invalid start date format, expected YYYY-MM-DD" ||
			err.Error() == "invalid end date format, expected YYYY-MM-DD" {
			Error(c, CodeBadRequest, err.Error())
			return
		}
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	var responses []MentalHealthRecordResponse
	for _, record := range records {
		response := MentalHealthRecordResponse{
			ID:          record.ID().String(),
			UserID:      record.UserID().String(),
			HappyLevel:  record.HappyLevel().Value(),
			EnergyLevel: record.EnergyLevel().Value(),
			Notes:       record.Notes(),
			Status:      record.Status().String(),
			CreatedAt:   timeutil.FormatTime(record.CreatedAt()),
			UpdatedAt:   timeutil.FormatTime(record.UpdatedAt()),
		}
		responses = append(responses, response)
	}

	Success(c, "Mental health records retrieved successfully", responses)
}

func (h *MentalHealthRecordHandler) GetHeatmap(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Get query parameters
	startedAt := c.Query("started_at")
	endedAt := c.Query("ended_at")

	var startedAtPtr *string
	var endedAtPtr *string

	if startedAt != "" {
		startedAtPtr = &startedAt
	}
	if endedAt != "" {
		endedAtPtr = &endedAt
	}

	// Create command
	command := commands.NewGetMentalHealthHeatmapCommand(userID.String(), startedAtPtr, endedAtPtr)

	// Execute use case
	ctx := c.Request.Context()
	heatmapResult, err := h.recordUseCase.GetHeatmap(ctx, command)
	if err != nil {
		if err.Error() == "invalid start date format, expected YYYY-MM-DD" ||
			err.Error() == "invalid end date format, expected YYYY-MM-DD" {
			Error(c, CodeBadRequest, err.Error())
			return
		}
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Simple mapping from application result to HTTP response
	httpData := make(map[string]HeatmapDataPoint)
	for date, dataPoint := range heatmapResult.Data {
		httpData[date] = HeatmapDataPoint{
			HappyLevel:  dataPoint.HappyLevel,
			EnergyLevel: dataPoint.EnergyLevel,
			Count:       dataPoint.Count,
		}
	}

	response := MentalHealthHeatmapResponse{
		Data:         httpData,
		TotalRecords: heatmapResult.TotalRecords,
		DateRange: DateRange{
			StartedAt: heatmapResult.DateRange.StartedAt,
			EndedAt:   heatmapResult.DateRange.EndedAt,
		},
	}

	Success(c, "Mental health heatmap retrieved successfully", response)
}

type MentalHealthStreakResponse struct {
	Streak        int     `json:"streak"`
	LastEntryDate *string `json:"last_entry_date,omitempty"`
}

func (h *MentalHealthRecordHandler) GetStreak(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserIDFromGinContext(c)
	if !exists {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}

	// Create command
	command := commands.NewGetMentalHealthStreakCommand(userID.String())

	// Execute use case
	ctx := c.Request.Context()
	streakResult, err := h.recordUseCase.GetStreak(ctx, command)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	response := MentalHealthStreakResponse{
		Streak:        streakResult.Streak,
		LastEntryDate: streakResult.LastEntryDate,
	}

	Success(c, "Mental health streak retrieved successfully", response)
}
