package response

// @Description Краткая информация о жильце, не привязанном ни к одной квартире.
type TenantSummary struct {
	ID       string `json:"id" example:"user-a9b8c7d"`
	Name     string `json:"name" example:"Иванов Иван Иванович"`
	Email    string `json:"email" example:"ivanov@example.com"`
	Document string `json:"document" example:"Паспорт РФ 4512 № 345678"`
	Phone    string `json:"phone" example:"+79991112233"`
}
