package entity

type Property struct {
	ID          string   `json:"id" example:"7c9e6b4a-4321-4a1b-8def-abcdef123456"`
	LandlordID  string   `json:"landlord_id" example:"3f7e2b1a-1234-4c56-9abc-1234567890ab"`
	Name        string   `json:"name" example:"Sunny Apartment"`
	Coordinates string   `json:"coordinates" example:"55.751244,37.618423"`
	Country     string   `json:"country" example:"Russia"`
	Region      string   `json:"region" example:"Moscow"`
	City        string   `json:"city" example:"Moscow"`
	Street      string   `json:"street" example:"Tverskaya"`
	House       string   `json:"house" example:"12"`
	Apartment   string   `json:"apartment" example:"45"`
	GvsTariff   float64  `json:"gvs_tariff" example:"189.5"`
	HvsTariff   float64  `json:"hvs_tariff" example:"45.2"`
	El1Tariff   float64  `json:"el1_tariff" example:"5.7"`
	El2Tariff   *float64 `json:"el2_tariff" example:"4.9"`
	Balance     float64  `json:"balance" example:"1250.75"`
}
