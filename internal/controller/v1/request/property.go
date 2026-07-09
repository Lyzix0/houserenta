package request

type Property struct {
	Name        string   `json:"name" validate:"required,max=200" example:"Sunny Apartment"`
	Coordinates string   `json:"coordinates" validate:"required,max=100" example:"55.751244,37.618423"`
	Country     string   `json:"country" validate:"max=100" example:"Russia"`
	Region      string   `json:"region" validate:"required,max=100" example:"Moscow"`
	City        string   `json:"city" validate:"required,max=100" example:"Moscow"`
	Street      string   `json:"street" validate:"required,max=200" example:"Tverskaya"`
	House       string   `json:"house" validate:"required,max=20" example:"12"`
	Apartment   string   `json:"apartment" validate:"required,max=20" example:"45"`
	GvsTariff   float64  `json:"gvsTariff" validate:"required,gt=0,max=100000" example:"189.5"`
	HvsTariff   float64  `json:"hvsTariff" validate:"required,gt=0,max=100000" example:"45.2"`
	El1Tariff   float64  `json:"el1Tariff" validate:"required,gt=0,max=100000" example:"5.7"`
	El2Tariff   *float64 `json:"el2Tariff,omitempty" validate:"omitempty,gt=0,max=100000" example:"4.9"`
}
