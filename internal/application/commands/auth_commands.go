package commands

import (
	"errors"
)

type RegisterCommand struct {
	Email     string
	Username  string
	Password  string
	FirstName *string
	LastName  *string
}

func NewRegisterCommand(email string, username string, password string, firstName *string, lastName *string) (RegisterCommand, error) {
	if email == "" {
		return RegisterCommand{}, errors.New("email is required")
	}

	if username == "" {
		return RegisterCommand{}, errors.New("username is required")
	}

	if password == "" {
		return RegisterCommand{}, errors.New("password is required")
	}

	return RegisterCommand{
		Email:     email,
		Username:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

type LoginCommand struct {
	Email    string
	Password string
}

func NewLoginCommand(email, password string) (LoginCommand, error) {
	if email == "" {
		return LoginCommand{}, errors.New("email is required")
	}

	if password == "" {
		return LoginCommand{}, errors.New("password is required")
	}

	return LoginCommand{
		Email:    email,
		Password: password,
	}, nil
}
