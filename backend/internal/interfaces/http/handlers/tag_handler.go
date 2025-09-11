package handlers

import (
	"strconv"

	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/domain/repositories"
	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagUseCase usecases.TagUseCase
}

func NewTagHandler(tagUseCase usecases.TagUseCase) *TagHandler {
	return &TagHandler{
		tagUseCase: tagUseCase,
	}
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag with name and optional description
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body commands.CreateTagCommand true "Tag information"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var cmd commands.CreateTagCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		Error(c, CodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	tag, err := h.tagUseCase.CreateTag(c.Request.Context(), &cmd)
	if err != nil {
		Error(c, CodeServerError, "Failed to create tag: "+err.Error())
		return
	}

	Success(c, "Tag created successfully", tag)
}

// GetAllTags godoc
// @Summary Get all tags
// @Description Get all available tags
// @Tags tags
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/tags [get]
func (h *TagHandler) GetAllTags(c *gin.Context) {
	tags, err := h.tagUseCase.GetAllTags(c.Request.Context())
	if err != nil {
		Error(c, CodeServerError, "Failed to get tags: "+err.Error())
		return
	}

	Success(c, "Tags retrieved successfully", tags)
}

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body commands.UpdateTagCommand true "Updated tag information"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, CodeBadRequest, "Invalid tag ID")
		return
	}

	var cmd commands.UpdateTagCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		Error(c, CodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	tag, err := h.tagUseCase.UpdateTag(c.Request.Context(), id, &cmd)
	if err != nil {
		if err == repositories.ErrTagNotFound {
			Error(c, CodeNotFound, "Tag not found")
			return
		}
		Error(c, CodeServerError, "Failed to update tag: "+err.Error())
		return
	}

	Success(c, "Tag updated successfully", tag)
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a tag by its ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, CodeBadRequest, "Invalid tag ID")
		return
	}

	err = h.tagUseCase.DeleteTag(c.Request.Context(), id)
	if err != nil {
		if err == repositories.ErrTagNotFound {
			Error(c, CodeNotFound, "Tag not found")
			return
		}
		Error(c, CodeServerError, "Failed to delete tag: "+err.Error())
		return
	}

	Success(c, "Tag deleted successfully", nil)
}

// GetTagsByQuoteID godoc
// @Summary Get tags by quote ID
// @Description Get all tags associated with a specific quote
// @Tags tags
// @Produce json
// @Param id path int true "Quote ID"
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/quotes/{id}/tags [get]
func (h *TagHandler) GetTagsByQuoteID(c *gin.Context) {
	quoteIDStr := c.Param("id")
	quoteID, err := strconv.Atoi(quoteIDStr)
	if err != nil {
		Error(c, CodeBadRequest, "Invalid quote ID")
		return
	}

	tags, err := h.tagUseCase.GetTagsByQuoteID(c.Request.Context(), quoteID)
	if err != nil {
		Error(c, CodeServerError, "Failed to get tags for quote: "+err.Error())
		return
	}

	Success(c, "Tags retrieved successfully", tags)
}

// AddTagToQuote godoc
// @Summary Add tag to quote
// @Description Add a tag to a specific quote
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Quote ID"
// @Param tag body commands.AddTagToQuoteCommand true "Tag to add"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/quotes/{id}/tags [post]
func (h *TagHandler) AddTagToQuote(c *gin.Context) {
	quoteIDStr := c.Param("id")
	quoteID, err := strconv.Atoi(quoteIDStr)
	if err != nil {
		Error(c, CodeBadRequest, "Invalid quote ID")
		return
	}

	var cmd commands.AddTagToQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		Error(c, CodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	err = h.tagUseCase.AddTagToQuote(c.Request.Context(), quoteID, &cmd)
	if err != nil {
		if err == repositories.ErrTagNotFound {
			Error(c, CodeNotFound, "Tag not found")
			return
		}
		Error(c, CodeServerError, "Failed to add tag to quote: "+err.Error())
		return
	}

	Success(c, "Tag added to quote successfully", nil)
}

// RemoveTagFromQuote godoc
// @Summary Remove tag from quote
// @Description Remove a tag from a specific quote
// @Tags tags
// @Accept json
// @Produce json
// @Param id path int true "Quote ID"
// @Param tag body commands.RemoveTagFromQuoteCommand true "Tag to remove"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/quotes/{id}/tags [delete]
func (h *TagHandler) RemoveTagFromQuote(c *gin.Context) {
	quoteIDStr := c.Param("id")
	quoteID, err := strconv.Atoi(quoteIDStr)
	if err != nil {
		Error(c, CodeBadRequest, "Invalid quote ID")
		return
	}

	var cmd commands.RemoveTagFromQuoteCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		Error(c, CodeBadRequest, "Invalid request body: "+err.Error())
		return
	}

	err = h.tagUseCase.RemoveTagFromQuote(c.Request.Context(), quoteID, &cmd)
	if err != nil {
		if err == repositories.ErrTagNotFound {
			Error(c, CodeNotFound, "Tag not found")
			return
		}
		Error(c, CodeServerError, "Failed to remove tag from quote: "+err.Error())
		return
	}

	Success(c, "Tag removed from quote successfully", nil)
}
