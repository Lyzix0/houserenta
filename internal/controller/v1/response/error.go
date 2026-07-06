package response

// @Description Стандартный формат ошибки API
type Error struct {
	Error string `json:"error"`
}
