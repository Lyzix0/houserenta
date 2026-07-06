package request

type Register struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" example:"Ivan Petrov"`
	Email       string  `json:"email" validate:"required,email" example:"ivan.petrov@example.com"`
	Password    string  `json:"password" validate:"required,min=6" example:"strongPass123"`
	Role        string  `json:"role" validate:"required,oneof=landlord tenant admin" example:"landlord"`
	Document    string  `json:"document" validate:"required" example:"4510 123456"`
	Phone       string  `json:"phone" validate:"required" example:"+79161234567"`
	PaymentCard *string `json:"payment_card,omitempty" validate:"omitempty,len=16,numeric" example:"1234567812345678"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"ivan.petrov@example.com"`
	Password string `json:"password" validate:"required" example:"strongPass123"`
}
