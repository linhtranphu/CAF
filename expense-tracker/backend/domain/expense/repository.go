package expense

type Repository interface {
	Save(expense *Expense) error
}

type MessageParser interface {
	Parse(message string) (items string, amount int64, error error)
}