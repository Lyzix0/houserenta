package repo

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUserNotFound      = errors.New("user not found")

	ErrPropertyAlreadyExists = errors.New("property already exists")
	ErrPropertyNotFound      = errors.New("property not found")
	ErrLandlordNotFound      = errors.New("landlord not found")
)
