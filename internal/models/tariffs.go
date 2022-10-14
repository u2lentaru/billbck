package models

//Tariff struct
type Tariff struct {
	Id          int         `json:"id"`
	TariffName  string      `json:"tariffname"`
	TariffGroup TariffGroup `json:"tariffgroup"`
	Norma       float32     `json:"norma"`
	Tariff      float32     `json:"tariff"`
	Startdate   string      `json:"startdate"`
	//*string correct scan null data!
	Enddate *string `json:"enddate"`
}

//AddTariff struct
type AddTariff struct {
	TariffName  string      `json:"tariffname"`
	TariffGroup TariffGroup `json:"tariffgroup"`
	Norma       float32     `json:"norma"`
	Tariff      float32     `json:"tariff"`
	Startdate   string      `json:"startdate"`
}

//Tariff_count  struct
type Tariff_count struct {
	Values []Tariff `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
