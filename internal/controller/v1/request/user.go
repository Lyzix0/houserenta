package request

type Register struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" example:"Ivan Petrov"`
	Email       string  `json:"email" validate:"required,email,max=254" example:"ivan.petrov@example.com"`
	Password    string  `json:"password" validate:"required,min=6,max=72" example:"strongPass123"`
	Document    string  `json:"document,omitempty" validate:"max=200" example:"4510 123456"`
	Phone       string  `json:"phone" validate:"required,max=20" example:"+79161234567"`
	InitialRole string  `json:"initialRole" validate:"required,oneof=landlord tenant" example:"landlord"`
	PaymentCard *string `json:"paymentCard,omitempty" validate:"omitempty,len=16,numeric" example:"1234567812345678"`
}

// Login identifies the user by email or phone; both are accepted in the same field.
type Login struct {
	Email    string `json:"email" validate:"required,max=254" example:"ivan.petrov@example.com"`
	Password string `json:"password" validate:"required,max=72" example:"strongPass123"`
}

// Profile is a partial update: only non-nil fields are applied.
type Profile struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100" example:"Ivan Kolesnikov"`
	Document    *string `json:"document,omitempty" validate:"omitempty,max=200" example:"4510 123456"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,max=20" example:"+79997776655"`
	PaymentCard *string `json:"paymentCard,omitempty" validate:"omitempty,len=16,numeric" example:"1111222233334444"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email,max=254" example:"newemail@example.com"`
	Password    *string `json:"password,omitempty" validate:"omitempty,min=6,max=72" example:"newSuperSecurePassword"`
}
