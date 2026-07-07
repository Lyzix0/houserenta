package request

type CustomItem struct {
	Description string  `json:"description" validate:"required" example:"Замена смесителя на кухне"`
	Amount      float64 `json:"amount" validate:"required" example:"2500"`
}
