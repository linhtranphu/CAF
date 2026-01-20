package expense

import (
	"errors"
	"time"
)

type Expense struct {
	id       string
	items    string
	amount   Money
	paidDate time.Time
	paidBy   string
}

type Money struct {
	value int64
}

func NewExpense(items string, amount int64, paidBy string) (*Expense, error) {
	if items == "" {
		return nil, errors.New("items cannot be empty")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if paidBy == "" {
		return nil, errors.New("paidBy cannot be empty")
	}

	return &Expense{
		items:    items,
		amount:   Money{value: amount},
		paidDate: time.Now(),
		paidBy:   paidBy,
	}, nil
}

func (e *Expense) Items() string     { return e.items }
func (e *Expense) Amount() int64     { return e.amount.value }
func (e *Expense) PaidDate() time.Time { return e.paidDate }
func (e *Expense) PaidBy() string    { return e.paidBy }