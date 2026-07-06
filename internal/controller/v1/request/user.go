package request

type Register struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=6"`
	Role        string  `json:"role" validate:"required,oneof=landlord tenant admin"`
	Document    string  `json:"document" validate:"required"`
	Phone       string  `json:"phone" validate:"required"`
	PaymentCard *string `json:"payment_card,omitempty" validate:"omitempty,len=16,numeric"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
