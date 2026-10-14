package entity

import (
	"errors"
)

var (
	ErrUserEmailConflict            = errors.New("email already exists")
	ErrUserEmailNotFound            = errors.New("email not found")
	ErrUserEmailInternalServerError = errors.New("something went wrong")
)
