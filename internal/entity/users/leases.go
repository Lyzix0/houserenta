package entity

// @Description Договор аренды: юридическая связь между квартирой и конкретным жильцом.
type Lease struct {
	ID           string  `json:"id" example:"lease-1a2b3c4d"`
	PropertyID   string  `json:"property_id" example:"7c9e6b4a-4321-4a1b-8def-abcdef123456"`
	TenantUserID string  `json:"tenant_user_id" example:"user-a9b8c7d"`
	Name         string  `json:"name" example:"Иванов Иван Иванович"`
	Document     string  `json:"document" example:"Паспорт РФ 4512 № 345678"`
	Phone        string  `json:"phone" example:"+79991112233"`
	MonthsOfRent int     `json:"months_of_rent" example:"12"`
	Price        float64 `json:"price" example:"45000"`
	PaymentDay   int     `json:"payment_day" example:"5"`
	ReadingDay   int     `json:"reading_day" example:"25"`
	StartDate    string  `json:"start_date" example:"2026-01-01"`
	EndDate      string  `json:"end_date" example:"2026-12-31"`
}
