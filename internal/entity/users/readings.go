package entity

// @Description Показания приборов учета квартиры.
type Reading struct {
	ID          string   `json:"id" example:"read-123"`
	PropertyID  string   `json:"property_id" example:"prop-z8y7x6w"`
	Date        string   `json:"date" example:"2026-06-25T12:00:00.000Z"`
	Gvs         float64  `json:"gvs" example:"12.5"`
	Hvs         float64  `json:"hvs" example:"24.1"`
	El1         float64  `json:"el1" example:"340"`
	El2         *float64 `json:"el2"`
	IsAccounted int      `json:"is_accounted" example:"1"`
}
