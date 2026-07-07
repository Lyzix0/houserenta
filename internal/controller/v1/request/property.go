package request

type Property struct {
	Name        string   `json:"name" validate:"required" example:"Sunny Apartment"`
	Coordinates string   `json:"coordinates" validate:"required" example:"55.751244,37.618423"`
	Country     string   `json:"country" example:"Russia"`
	Region      string   `json:"region" validate:"required" example:"Moscow"`
	City        string   `json:"city" validate:"required" example:"Moscow"`
	Street      string   `json:"street" validate:"required" example:"Tverskaya"`
	House       string   `json:"house" validate:"required" example:"12"`
	Apartment   string   `json:"apartment" validate:"required" example:"45"`
	GvsTariff   float64  `json:"gvsTariff" validate:"required" example:"189.5"`
	HvsTariff   float64  `json:"hvsTariff" validate:"required" example:"45.2"`
	El1Tariff   float64  `json:"el1Tariff" validate:"required" example:"5.7"`
	El2Tariff   *float64 `json:"el2Tariff,omitempty" example:"4.9"`
}
