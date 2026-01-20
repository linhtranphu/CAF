package http

import (
	"log"
	"net/http"
	"expense-tracker/application/services"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *services.ExpenseService
}

func NewAdminHandler(service *services.ExpenseService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) AdminPage(c *gin.Context) {
	expenses, err := h.service.GetAllExpenses()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Convert DTOs to map for template compatibility
	var expensesMaps []map[string]interface{}
	for i, exp := range expenses {
		expenseMap := map[string]interface{}{
			"no":       exp.ID,
			"items":    exp.Items,
			"amount":   exp.Amount,
			"paidDate": exp.PaidDate,
			"paidBy":   exp.PaidBy,
		}
		log.Printf("[ADMIN] Expense %d: %+v", i, expenseMap)
		expensesMaps = append(expensesMaps, expenseMap)
	}

	summary, err := h.service.GetExpenseSummary()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	// Debug log
	log.Printf("[ADMIN] Summary data: %+v", summary)

	// Calculate grand total
	var grandTotal int64
	for _, amount := range summary {
		grandTotal += amount
	}

	log.Printf("[ADMIN] Grand total: %d", grandTotal)

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"expenses":   expensesMaps,
		"total":      len(expenses),
		"summary":    summary,
		"grandTotal": grandTotal,
	})
}

func (h *AdminHandler) DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	log.Printf("[ADMIN] Delete request for ObjectID: %s", id)
	
	if err := h.service.DeleteExpense(id); err != nil {
		log.Printf("[ADMIN] Delete error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ADMIN] Successfully deleted expense ObjectID: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}