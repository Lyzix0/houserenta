package entity

// @Description Объект недвижимости с автоматически собранными связанными сущностями:
// @Description показаниями счетчиков, счетами, будущими начислениями, данными аренды и арендодателя.
type PropertyDetail struct {
	Property

	Readings        []Reading        `json:"readings"`
	Bills           []Bill           `json:"bills"`
	CustomNextItems []CustomNextItem `json:"customNextItems"`
	Tenant          *Lease           `json:"tenant"`
	LandlordName    string           `json:"landlordName" example:"Алексей Петров"`
	LandlordPhone   string           `json:"landlordPhone" example:"+79001234567"`
	Applications    []Application    `json:"applications"`
}
