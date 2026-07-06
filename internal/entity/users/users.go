package entity

import (
	"fmt"
	"regexp"
)

type Role string

const (
	RoleLandlord Role = "landlord"
	RoleTenant   Role = "tenant"
	RoleAdmin    Role = "admin"
)

var (
	emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+[0-9]{9,14}$`)
	cardRegex  = regexp.MustCompile(`^[0-9]{16}$`)
)

// @Description Основная модель пользователя: арендодатель, арендатор или администратор.
type User struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	PasswordHash string  `json:"-"`
	Role         Role    `json:"role"`
	Document     string  `json:"document"`
	Phone        string  `json:"phone"`
	PaymentCard  *string `json:"payment_card,omitempty"`
}

// валидация бизнес-логики
func (u *User) Validate() error {
	nameLen := len([]rune(u.Name))
	if nameLen < 2 || nameLen > 100 {
		return fmt.Errorf(
			"invalid `Name` len: %d: %w",
			nameLen,
			ErrInvalidName,
		)
	}

	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf(
			"invalid `Email` format: %w",
			ErrInvalidEmail,
		)
	}

	if len(u.PasswordHash) < 20 {
		return fmt.Errorf(
			"invalid `PasswordHash` len: %d: %w",
			len(u.PasswordHash),
			ErrInvalidPassword,
		)
	}

	switch u.Role {
	case RoleLandlord, RoleTenant, RoleAdmin:
	default:
		return fmt.Errorf(
			"invalid `Role` value: %q: %w",
			u.Role,
			ErrInvalidRole,
		)
	}

	if len([]rune(u.Document)) == 0 {
		return fmt.Errorf(
			"invalid `Document`: %w",
			ErrInvalidDocument,
		)
	}

	phoneLen := len([]rune(u.Phone))
	if phoneLen < 10 || phoneLen > 15 {
		return fmt.Errorf(
			"invalid `Phone` len: %d: %w",
			phoneLen,
			ErrInvalidPhone,
		)
	}
	if !phoneRegex.MatchString(u.Phone) {
		return fmt.Errorf(
			"invalid `Phone` format: %w",
			ErrInvalidPhone,
		)
	}

	if u.PaymentCard != nil {
		if !cardRegex.MatchString(*u.PaymentCard) {
			return fmt.Errorf(
				"invalid `PaymentCard` format: %w",
				ErrInvalidCard,
			)
		}
	}

	return nil
}
