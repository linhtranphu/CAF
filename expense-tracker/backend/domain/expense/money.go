package expense

import "errors"

// Money is a value object representing monetary amount in VND
type Money struct {
	value int64
}

func NewMoney(value int64) (Money, error) {
	if value < 0 {
		return Money{}, errors.New("money value cannot be negative")
	}
	return Money{value: value}, nil
}

func (m Money) Value() int64 {
	return m.value
}

func (m Money) Add(other Money) Money {
	return Money{value: m.value + other.value}
}

func (m Money) IsZero() bool {
	return m.value == 0
}

func (m Money) IsPositive() bool {
	return m.value > 0
}

func (m Money) Equals(other Money) bool {
	return m.value == other.value
}