package request

type Payment struct {
	Amount float64 `json:"amount" validate:"required,gt=0" example:"35000"`
	BillID *string `json:"billId,omitempty" example:"bill-999"`
}
