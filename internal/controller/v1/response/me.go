package response

// @Description Полный профиль авторизованного пользователя текущей сессии.
type Me struct {
	ID               string  `json:"id" example:"user-a9b8c7d"`
	Name             string  `json:"name" example:"Иванов Иван Иванович"`
	Email            string  `json:"email" example:"ivanov@example.com"`
	Role             string  `json:"role" example:"tenant"`
	Document         string  `json:"document" example:"Паспорт РФ 4512 № 345678"`
	Phone            string  `json:"phone" example:"+79991112233"`
	PaymentCard      *string `json:"paymentCard" example:"4276111122223333"`
	TenantPropertyID *string `json:"tenantPropertyId" example:"prop-z8y7x6w"`
}
