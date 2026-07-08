package repo

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrUserNotFound      = errors.New("user not found")

	ErrPropertyAlreadyExists = errors.New("property already exists")
	ErrPropertyNotFound      = errors.New("property not found")
	ErrLandlordNotFound      = errors.New("landlord not found")

	ErrLeaseNotFound  = errors.New("lease not found")
	ErrTenantNotFound = errors.New("tenant not found")

	ErrBillNotFound    = errors.New("bill not found")
	ErrReadingNotFound = errors.New("reading not found")

	ErrApplicationAlreadyExists = errors.New("application already exists")
)
