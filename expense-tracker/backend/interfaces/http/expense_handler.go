package http

import (
	"net/http"

	"expense-tracker/application/services"
	"github.com/gin-gonic/gin"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

type ExpenseRequest struct {
	Message string `json:"message" binding:"required"`
	UserID  string `json:"userId" binding:"required"`
}

func NewExpenseHandler(service *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	var req ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateExpenseFromMessage(req.Message, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}