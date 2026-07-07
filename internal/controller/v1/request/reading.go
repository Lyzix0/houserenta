package request

type Reading struct {
	Gvs float64  `json:"gvs" validate:"required" example:"15.2"`
	Hvs float64  `json:"hvs" validate:"required" example:"29.8"`
	El1 float64  `json:"el1" validate:"required" example:"395.5"`
	El2 *float64 `json:"el2,omitempty" example:"110.2"`
}
