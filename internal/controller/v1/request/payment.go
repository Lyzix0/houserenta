package request

type Payment struct {
	Amount float64 `json:"amount" validate:"required,gt=0,max=100000000" example:"35000"`
	BillID *string `json:"billId,omitempty" validate:"omitempty,max=100" example:"bill-999"`
}
