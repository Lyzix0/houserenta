package response

// @Description Краткая информация о пользователе, возвращаемая при регистрации и входе.
type AuthUser struct {
	ID    string `json:"id" example:"user-a9b8c7d"`
	Name  string `json:"name" example:"Иванов Иван Иванович"`
	Email string `json:"email" example:"ivanov@example.com"`
	Role  string `json:"role" example:"tenant"`
}
