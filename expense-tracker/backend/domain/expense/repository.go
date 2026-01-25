package expense

import "time"

type Repository interface {
	Save(expense *Expense) error
	FindByID(id int) (*Expense, error)
	FindAll() ([]*Expense, error)
	FindActiveExpenses() ([]*Expense, error)
	GetSummaryByPaidBy() (map[string]int64, error)
	Delete(id string) error
	ClearAll() error
	GetAll() ([]map[string]interface{}, error)
	GetDeleted() ([]map[string]interface{}, error)
}

type MessageParser interface {
	Parse(message string) (items string, amount int64, quantity string, unit string, baseQuantity string, baseUnit string, originalMessage string, paidDate time.Time, error error)
}

// DTOs for presentation layer
type ExpenseDTO struct {
	ID              string `json:"id"`
	Items           string `json:"items"`
	Amount          int64  `json:"amount"`
	Quantity        string `json:"quantity,omitempty"`
	Unit            string `json:"unit,omitempty"`
	BaseQuantity    string `json:"baseQuantity,omitempty"`
	BaseUnit        string `json:"baseUnit,omitempty"`
	OriginalMessage string `json:"originalMessage,omitempty"`
	PaidDate        string `json:"paidDate"`
	PaidBy          string `json:"paidBy"`
}