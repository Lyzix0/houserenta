package request

type Property struct {
	LandlordID  string   `json:"landlord_id" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	Coordinates string   `json:"coordinates" validate:"required"`
	Country     string   `json:"country"`
	Region      string   `json:"region" validate:"required"`
	City        string   `json:"city" validate:"required"`
	Street      string   `json:"street" validate:"required"`
	House       string   `json:"house" validate:"required"`
	Apartment   string   `json:"apartment" validate:"required"`
	GvsTariff   float64  `json:"gvs_tariff" validate:"required"`
	HvsTariff   float64  `json:"hvs_tariff" validate:"required"`
	El1Tariff   float64  `json:"el1_tariff" validate:"required"`
	El2Tariff   *float64 `json:"el2_tariff,omitempty"`
}
