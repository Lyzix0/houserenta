package repo

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUserNotFound      = errors.New("user not found")
)
