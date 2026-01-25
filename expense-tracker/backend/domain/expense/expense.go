package expense

import (
	"errors"
	"time"
)

type Expense struct {
	id              int
	items           string
	amount          Money
	quantity        string
	unit            string
	baseQuantity    string
	baseUnit        string
	originalMessage string
	paidDate        time.Time
	paidBy          string
	status          Status
}

type Status string

const (
	StatusActive  Status = "active"
	StatusDeleted Status = "deleted"
)

func NewExpense(items string, amount int64, paidBy string) (*Expense, error) {
	if items == "" {
		return nil, errors.New("items cannot be empty")
	}
	if paidBy == "" {
		return nil, errors.New("paidBy cannot be empty")
	}

	money, err := NewMoney(amount)
	if err != nil {
		return nil, err
	}

	return &Expense{
		items:    items,
		amount:   money,
		quantity: "",  // Sẽ được AI phân tích sau
		unit:     "",  // Sẽ được AI phân tích sau
		paidDate: time.Now(),
		paidBy:   paidBy,
		status:   StatusActive,
	}, nil
}

func NewExpenseWithQuantityUnit(items string, amount int64, quantity, unit, paidBy string) (*Expense, error) {
	if items == "" {
		return nil, errors.New("items cannot be empty")
	}
	if paidBy == "" {
		return nil, errors.New("paidBy cannot be empty")
	}

	money, err := NewMoney(amount)
	if err != nil {
		return nil, err
	}

	return &Expense{
		items:    items,
		amount:   money,
		quantity: quantity,
		unit:     unit,
		paidDate: time.Now(),
		paidBy:   paidBy,
		status:   StatusActive,
	}, nil
}

func NewExpenseWithDate(items string, amount int64, paidBy string, paidDate time.Time) *Expense {
	money, _ := NewMoney(amount)
	return &Expense{
		items:    items,
		amount:   money,
		quantity: "",
		unit:     "",
		paidDate: paidDate,
		paidBy:   paidBy,
		status:   StatusActive,
	}
}

func (e *Expense) Items() string           { return e.items }
func (e *Expense) Amount() int64            { return e.amount.Value() }
func (e *Expense) Quantity() string         { return e.quantity }
func (e *Expense) Unit() string             { return e.unit }
func (e *Expense) BaseQuantity() string     { return e.baseQuantity }
func (e *Expense) BaseUnit() string         { return e.baseUnit }
func (e *Expense) OriginalMessage() string  { return e.originalMessage }
func (e *Expense) PaidDate() time.Time      { return e.paidDate }
func (e *Expense) PaidBy() string           { return e.paidBy }
func (e *Expense) Status() Status           { return e.status }
func (e *Expense) ID() int                  { return e.id }

// Business logic methods
func (e *Expense) Delete() {
	e.status = StatusDeleted
}

func (e *Expense) IsActive() bool {
	return e.status == StatusActive
}

func (e *Expense) IsDeleted() bool {
	return e.status == StatusDeleted
}

func (e *Expense) SetQuantityUnit(quantity, unit string) {
	e.quantity = quantity
	e.unit = unit
}

func (e *Expense) SetBaseQuantityUnit(baseQuantity, baseUnit string) {
	e.baseQuantity = baseQuantity
	e.baseUnit = baseUnit
}

func (e *Expense) SetOriginalMessage(message string) {
	e.originalMessage = message
}