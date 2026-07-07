package entity

// @Description Разовое начисление, которое будет автоматически добавлено в следующий ежемесячный счет.
type CustomNextItem struct {
	ID          string  `json:"id" example:"custom-1"`
	PropertyID  string  `json:"property_id" example:"prop-z8y7x6w"`
	Description string  `json:"description" example:"Ремонт сантехники"`
	Amount      float64 `json:"amount" example:"2500"`
}
