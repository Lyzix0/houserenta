package entity

import "errors"

var (
	ErrInvalidRole     = errors.New("role must be one of: landlord, tenant, admin")
	ErrInvalidName     = errors.New("name must be between 2 and 100 characters")
	ErrInvalidEmail    = errors.New("email must be a valid email address")
	ErrInvalidPhone    = errors.New("phone must be a valid phone number")
	ErrInvalidPassword = errors.New("password hash must be at least 20 characters")
	ErrInvalidCard     = errors.New("payment card must be exactly 16 digits")
)
