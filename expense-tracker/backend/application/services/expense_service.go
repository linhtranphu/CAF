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
	user, err := user.NewUser(userName)
	if err != nil {
		return err
	}

	items, amount, paidDate, err := s.parser.Parse(message)
	if err != nil {
		return err
	}

	exp := expense.NewExpenseWithDate(items, amount, user.Name(), paidDate)
	return s.expenseRepo.Save(exp)
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
			PaidDate: exp["paidDate"].(string),
			PaidBy:   exp["paidBy"].(string),
		}
		log.Printf("[SERVICE] DTO: ID=%s, Items=%s", dto.ID, dto.Items)
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
	headers := []string{"Mô tả", "Số tiền (VND)", "Ngày", "Người trả"}
	writer.Write(headers)

	// Write data
	for _, expense := range expenses {
		record := []string{
			expense["items"].(string),
			fmt.Sprintf("%d", expense["amount"].(int64)),
			expense["paidDate"].(string),
			expense["paidBy"].(string),
		}
		writer.Write(record)
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}