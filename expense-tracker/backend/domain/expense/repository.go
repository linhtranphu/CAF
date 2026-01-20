package expense

type Repository interface {
	Save(expense *Expense) error
	FindByID(id int) (*Expense, error)
	FindAll() ([]*Expense, error)
	FindActiveExpenses() ([]*Expense, error)
	GetSummaryByPaidBy() (map[string]int64, error)
	Delete(id string) error
	GetAll() ([]map[string]interface{}, error)
}

type MessageParser interface {
	Parse(message string) (items string, amount int64, error error)
}

// DTOs for presentation layer
type ExpenseDTO struct {
	ID       string `json:"id"`
	Items    string `json:"items"`
	Amount   int64  `json:"amount"`
	PaidDate string `json:"paidDate"`
	PaidBy   string `json:"paidBy"`
}