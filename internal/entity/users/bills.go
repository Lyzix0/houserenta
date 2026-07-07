package entity

// @Description Строка начисления в составе счета.
type BillItem struct {
	ID          string  `json:"id" example:"item-1"`
	BillID      string  `json:"bill_id" example:"bill-999"`
	Description string  `json:"description" example:"Аренда за период"`
	Amount      float64 `json:"amount" example:"35000"`
}

// @Description Счет на оплату аренды и коммунальных услуг.
type Bill struct {
	ID         string     `json:"id" example:"bill-999"`
	PropertyID string     `json:"property_id" example:"prop-z8y7x6w"`
	Date       string     `json:"date" example:"2026-06-25T14:00:00.000Z"`
	DueDate    string     `json:"due_date" example:"2026-07-05T00:00:00.000Z"`
	Status     string     `json:"status" example:"unpaid"`
	Total      float64    `json:"total" example:"35000"`
	Items      []BillItem `json:"items"`
}
