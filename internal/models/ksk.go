package models

//Ksk struct
type Ksk struct {
	Id         int    `json:"id"`
	KskName    string `json:"kskname"`
	KskAddress string `json:"kskaddress"`
	KskHead    string `json:"kskhead"`
	KskPhone   string `json:"kskphone"`
}

//AddKsk struct
type AddKsk struct {
	KskName    string `json:"kskname"`
	KskAddress string `json:"kskaddress"`
	KskHead    string `json:"kskhead"`
	KskPhone   string `json:"kskphone"`
}

//Ksk_count  struct
type Ksk_count struct {
	Values []Ksk `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}
