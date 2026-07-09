package request

type CustomItem struct {
	Description string  `json:"description" validate:"required,max=500" example:"Замена смесителя на кухне"`
	Amount      float64 `json:"amount" validate:"required,gt=0,max=100000000" example:"2500"`
}
