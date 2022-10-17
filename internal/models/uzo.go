package models

//Uzo struct
type Uzo struct {
	Id       int    `json:"id"`
	UzoName  string `json:"uzoname"`
	UzoValue int    `json:"uzovalue"`
}

//AddUzo struct
type AddUzo struct {
	UzoName  string `json:"uzoname"`
	UzoValue int    `json:"uzovalue"`
}

//Uzo_count  struct
type Uzo_count struct {
	Values []Uzo `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}
