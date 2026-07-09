package request

type Reading struct {
	Gvs float64  `json:"gvs" validate:"required,gt=0,max=10000000" example:"15.2"`
	Hvs float64  `json:"hvs" validate:"required,gt=0,max=10000000" example:"29.8"`
	El1 float64  `json:"el1" validate:"required,gt=0,max=10000000" example:"395.5"`
	El2 *float64 `json:"el2,omitempty" validate:"omitempty,gt=0,max=10000000" example:"110.2"`
}
