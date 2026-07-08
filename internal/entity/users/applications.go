package entity

// @Description Отклик жильца на свободную квартиру.
type Application struct {
	ID           string `json:"id" example:"app-1a2b3c4d"`
	PropertyID   string `json:"property_id" example:"prop-z8y7x6w"`
	TenantUserID string `json:"tenant_user_id" example:"user-a9b8c7d"`
	Date         string `json:"date" example:"2026-06-25T12:00:00.000Z"`
}
