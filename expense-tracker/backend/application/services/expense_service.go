package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"expense-tracker/domain/expense"
	"expense-tracker/domain/user"
)

type ExpenseService struct {
	expenseRepo expense.Repository
	parser      expense.MessageParser
}

func NewExpenseService(repo expense.Repository, parser expense.MessageParser) *ExpenseService {
	return &ExpenseService{
		expenseRepo: repo,
		parser:      parser,
	}
}

func (s *ExpenseService) CreateExpenseFromMessage(message, userName string) error {
	_, err := s.CreateExpenseFromMessageWithDetails(message, userName)
	return err
}

func (s *ExpenseService) CreateExpenseFromMessageWithDetails(message, userName string) (map[string]interface{}, error) {
	user, err := user.NewUser(userName)
	if err != nil {
		return nil, err
	}

	items, amount, quantity, unit, paidDate, err := s.parser.Parse(message)
	if err != nil {
		return nil, err
	}

	exp := expense.NewExpenseWithDate(items, amount, user.Name(), paidDate)
	exp.SetQuantityUnit(quantity, unit)
	
	if err := s.expenseRepo.Save(exp); err != nil {
		return nil, err
	}

	// Return parsed data
	parsedData := map[string]interface{}{
		"items":    items,
		"amount":   amount,
		"quantity": quantity,
		"unit":     unit,
		"paidDate": paidDate.Format("2006-01-02"),
		"paidBy":   user.Name(),
	}

	return parsedData, nil
}

func (s *ExpenseService) GetAllExpenses() ([]expense.ExpenseDTO, error) {
	expenses, err := s.expenseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var dtos []expense.ExpenseDTO
	for _, exp := range expenses {
		dto := expense.ExpenseDTO{
			ID:       exp["no"].(string), // Use ObjectID string
			Items:    exp["items"].(string),
			Amount:   exp["amount"].(int64),
			Quantity: getStringField(exp, "quantity"),
			Unit:     getStringField(exp, "unit"),
			PaidDate: exp["paidDate"].(string),
			PaidBy:   exp["paidBy"].(string),
		}
		log.Printf("[SERVICE] DTO: ID=%s, Items=%s, Quantity=%s, Unit=%s", dto.ID, dto.Items, dto.Quantity, dto.Unit)
		dtos = append(dtos, dto)
	}

	return dtos, nil
}

func (s *ExpenseService) GetExpenseSummary() (map[string]int64, error) {
	return s.expenseRepo.GetSummaryByPaidBy()
}

func (s *ExpenseService) GetDeletedExpenses() ([]map[string]interface{}, error) {
	return s.expenseRepo.GetDeleted()
}

func (s *ExpenseService) DeleteExpense(id string) error {
	return s.expenseRepo.Delete(id)
}

func (s *ExpenseService) ExportToCSV() ([]byte, error) {
	expenses, err := s.expenseRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write headers
	headers := []string{"Mô tả", "Số lượng", "Đơn vị", "Số tiền (VND)", "Ngày", "Người trả"}
	writer.Write(headers)

	// Write data
	for _, expense := range expenses {
		record := []string{
			expense["items"].(string),
			getStringField(expense, "quantity"),
			getStringField(expense, "unit"),
			fmt.Sprintf("%d", expense["amount"].(int64)),
			expense["paidDate"].(string),
			expense["paidBy"].(string),
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// Helper function to safely get string field from map
func getStringField(data map[string]interface{}, field string) string {
	if val, exists := data[field]; exists && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}