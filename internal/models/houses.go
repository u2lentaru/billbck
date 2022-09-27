package models

//House struct
type House struct {
	Id             int          `json:"id"`
	BuildingType   BuildingType `json:"buildingtype"`
	Street         Street       `json:"street"`
	HouseNumber    string       `json:"housenumber"`
	BuildingNumber *string      `json:"buildingnumber"`
	RP             Rp           `json:"tgu"`
	Area           Area         `json:"area"`
	Ksk            Ksk          `json:"ksk"`
	Sector         Sector       `json:"sector"`
	Connector      Connector    `json:"connector"`
	InputType      InputType    `json:"inputtype"`
	Reliability    Reliability  `json:"reliability"`
	Voltage        Voltage      `json:"voltage"`
	Notes          *string      `json:"notes"`
}

type AddHouse struct {
	BuildingType   BuildingType `json:"buildingtype"`
	Street         Street       `json:"street"`
	HouseNumber    string       `json:"housenumber"`
	BuildingNumber *string      `json:"buildingnumber"`
	RP             Rp           `json:"tgu"`
	Area           Area         `json:"area"`
	Ksk            Ksk          `json:"ksk"`
	Sector         Sector       `json:"sector"`
	Connector      Connector    `json:"connector"`
	InputType      InputType    `json:"inputtype"`
	Reliability    Reliability  `json:"reliability"`
	Voltage        Voltage      `json:"voltage"`
	Notes          *string      `json:"notes"`
}

//House_count  struct
type House_count struct {
	Values []House `json:"values"`
	Count  int     `json:"count"`
	Auth   Auth    `json:"auth"`
}
