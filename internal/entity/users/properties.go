package entity

type Property struct {
	ID          string   `json:"id"`
	LandlordID  string   `json:"landlord_id"`
	Name        string   `json:"name"`
	Coordinates string   `json:"coordinates"`
	Country     string   `json:"country"`
	Region      string   `json:"region"`
	City        string   `json:"city"`
	Street      string   `json:"street"`
	House       string   `json:"house"`
	Apartment   string   `json:"apartment"`
	GvsTariff   float64  `json:"gvs_tariff"`
	HvsTariff   float64  `json:"hvs_tariff"`
	El1Tariff   float64  `json:"el1_tariff"`
	El2Tariff   *float64 `json:"el2_tariff,omitempty"`
	Balance     float64  `json:"balance"`
}
