package request

type Lease struct {
	TenantUserID string  `json:"tenantUserId" validate:"required" example:"user-a9b8c7d"`
	Price        float64 `json:"price" validate:"required" example:"32000"`
	MonthsOfRent int     `json:"monthsOfRent" validate:"required" example:"11"`
	PaymentDay   int     `json:"paymentDay" validate:"required,min=1,max=28" example:"10"`
	ReadingDay   int     `json:"readingDay" validate:"required,min=1,max=28" example:"28"`
}
