package services

import (
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