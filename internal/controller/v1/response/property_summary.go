package response

// @Description Краткая информация о только что созданном объекте недвижимости.
type PropertySummary struct {
	ID   string `json:"id" example:"prop-k8r3f9d"`
	Name string `json:"name" example:"Студия у метро"`
}
