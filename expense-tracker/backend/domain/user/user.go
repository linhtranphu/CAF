package user

import "errors"

type User struct {
	id   string
	name string
}

func NewUser(name string) (*User, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	return &User{name: name}, nil
}

func (u *User) Name() string { return u.name }