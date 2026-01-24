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
	for _, exp := range expenses {
		expenseMap := map[string]interface{}{
			"id":       exp.ID,
			"items":    exp.Items,
			"amount":   exp.Amount,
			"quantity": exp.Quantity,
			"unit":     exp.Unit,
			"paidDate": exp.PaidDate,
			"paidBy":   exp.PaidBy,
		}
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

func (h *AdminHandler) DeletedPage(c *gin.Context) {
	expenses, err := h.service.GetDeletedExpenses()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "deleted.html", gin.H{
		"expenses": expenses,
		"total":    len(expenses),
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

func (h *AdminHandler) ExportCSV(c *gin.Context) {
	log.Printf("[ADMIN] CSV export request")
	
	data, err := h.service.ExportToCSV()
	if err != nil {
		log.Printf("[ADMIN] CSV export error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=chi-phi.csv")
	c.Data(http.StatusOK, "text/csv", data)
	log.Printf("[ADMIN] CSV export successful")
}