package handlers

import (
	"net/http"

	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/domain/entities"
	"github.com/atdevten/peace/internal/pkg/timeutil"
	"github.com/gin-gonic/gin"
)

type QuoteHandler struct {
	quoteUseCase usecases.QuoteUseCase
}

// Request/Response structs
type CreateQuoteRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type UpdateQuoteRequest struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type QuoteResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewQuoteHandler(quoteUseCase usecases.QuoteUseCase) *QuoteHandler {
	return &QuoteHandler{
		quoteUseCase: quoteUseCase,
	}
}

// Helper function to convert entity to response
func (h *QuoteHandler) buildQuoteResponse(quote *entities.Quote) QuoteResponse {
	return QuoteResponse{
		ID:        quote.ID().String(),
		Content:   quote.Content().Value(),
		Author:    quote.Author().Value(),
		CreatedAt: timeutil.FormatTime(quote.CreatedAt()),
		UpdatedAt: timeutil.FormatTime(quote.UpdatedAt()),
	}
}

func (h *QuoteHandler) CreateQuote(c *gin.Context) {
	var req CreateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	err := h.quoteUseCase.CreateQuote(c.Request.Context(), req.Content, req.Author)
	if err != nil {
		Error(c, CodeServerError, "Failed to create quote: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Code:    CodeSuccess,
		Message: "Quote created successfully",
	})
}

func (h *QuoteHandler) GetAllQuotes(c *gin.Context) {
	quotes, err := h.quoteUseCase.GetAllQuotes(c.Request.Context())
	if err != nil {
		Error(c, CodeServerError, "Failed to get quotes: "+err.Error())
		return
	}

	var quotesData []QuoteResponse
	for _, quote := range quotes {
		quotesData = append(quotesData, h.buildQuoteResponse(quote))
	}

	Success(c, "Quotes retrieved successfully", quotesData)
}

func (h *QuoteHandler) GetRandomQuote(c *gin.Context) {
	quote, err := h.quoteUseCase.GetRandomQuote(c.Request.Context())
	if err != nil {
		if err.Error() == "no quotes found" {
			Error(c, CodeNotFound, "No quotes available")
			return
		}
		Error(c, CodeServerError, "Failed to get random quote: "+err.Error())
		return
	}

	response := h.buildQuoteResponse(quote)
	Success(c, "Random quote retrieved successfully", response)
}

func (h *QuoteHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, CodeBadRequest, "Quote ID is required")
		return
	}

	quote, err := h.quoteUseCase.GetQuoteByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "quote not found" {
			Error(c, CodeNotFound, "Quote not found")
			return
		}
		Error(c, CodeServerError, "Failed to get quote: "+err.Error())
		return
	}

	response := h.buildQuoteResponse(quote)
	Success(c, "Quote retrieved successfully", response)
}

func (h *QuoteHandler) UpdateQuote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, CodeBadRequest, "Quote ID is required")
		return
	}

	var req UpdateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	err := h.quoteUseCase.UpdateQuote(c.Request.Context(), id, req.Content, req.Author)
	if err != nil {
		if err.Error() == "quote not found" {
			Error(c, CodeNotFound, "Quote not found")
			return
		}
		Error(c, CodeServerError, "Failed to update quote: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: "Quote updated successfully",
	})
}

func (h *QuoteHandler) DeleteQuote(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		Error(c, CodeBadRequest, "Quote ID is required")
		return
	}

	err := h.quoteUseCase.DeleteQuote(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "quote not found" {
			Error(c, CodeNotFound, "Quote not found")
			return
		}
		Error(c, CodeServerError, "Failed to delete quote: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Code:    CodeSuccess,
		Message: "Quote deleted successfully",
	})
}
