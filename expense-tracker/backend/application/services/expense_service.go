package services

import (
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

	items, amount, err := s.parser.Parse(message)
	if err != nil {
		return err
	}

	exp, err := expense.NewExpense(items, amount, user.Name())
	if err != nil {
		return err
	}

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

func (s *ExpenseService) DeleteExpense(id string) error {
	return s.expenseRepo.Delete(id)
}