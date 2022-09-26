package models

//Eso struct
type Eso struct {
	Id      int    `json:"id"`
	EsoName string `json:"esoname"`
}

//AddEso struct
type AddEso struct {
	EsoName string `json:"esoname"`
}

//Eso_count  struct
type Eso_count struct {
	Values []Eso `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}
