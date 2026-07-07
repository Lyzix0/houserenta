package request

type Register struct {
	Name        string  `json:"name" validate:"required,min=2,max=100" example:"Ivan Petrov"`
	Email       string  `json:"email" validate:"required,email" example:"ivan.petrov@example.com"`
	Password    string  `json:"password" validate:"required,min=6" example:"strongPass123"`
	Role        string  `json:"role" validate:"required,oneof=landlord tenant admin" example:"landlord"`
	Document    string  `json:"document" validate:"required" example:"4510 123456"`
	Phone       string  `json:"phone" validate:"required" example:"+79161234567"`
	PaymentCard *string `json:"payment_card" validate:"omitempty,len=16,numeric" example:"1234567812345678"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"ivan.petrov@example.com"`
	Password string `json:"password" validate:"required" example:"strongPass123"`
}

type UpdateProfile struct {
	Name        *string `json:"name"`
	Document    *string `json:"document"`
	Phone       *string `json:"phone"`
	PaymentCard *string `json:"payment_card"`
	Email       *string `json:"email"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type SwitchRole struct {
	TargetRole string `json:"target_role" validate:"required,oneof=landlord tenant admin"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required"`
}

type ResetPassword struct {
	Email       string `json:"email" validate:"required"`
	Code        string `json:"code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type VerifyEmail struct {
	Email string `json:"email" validate:"required"`
	Code  string `json:"code" validate:"required"`
}
