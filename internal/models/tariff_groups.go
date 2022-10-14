package models

//TariffGroup struct
type TariffGroup struct {
	Id              int    `json:"id"`
	TariffGroupName string `json:"tariffgroupname"`
}

//AddTariffGroup struct
type AddTariffGroup struct {
	TariffGroupName string `json:"tariffgroupname"`
}

//TariffGroup_count  struct
type TariffGroup_count struct {
	Values []TariffGroup `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
