package http

import (
	"net/http"
	"log"
	"time"

	"expense-tracker/application/services"
	"github.com/gin-contrib/sessions"
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
	start := time.Now()
	log.Printf("[REQUEST] POST /api/expense from %s", c.ClientIP())
	
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not logged in"})
		return
	}
	
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ERROR] Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[INFO] Processing expense: user=%s, message=%s", username, req.Message)
	
	if err := h.service.CreateExpenseFromMessage(req.Message, username.(string)); err != nil {
		log.Printf("[ERROR] Failed to create expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Expense created in %v", time.Since(start))
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	start := time.Now()
	log.Printf("[REQUEST] GET /api/expenses from %s", c.ClientIP())
	
	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		log.Printf("[ERROR] Failed to get expenses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[SUCCESS] Retrieved %d expenses in %v", len(expenses), time.Since(start))
	c.JSON(http.StatusOK, gin.H{"data": expenses})
}