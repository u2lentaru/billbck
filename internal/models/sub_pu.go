package models

//SubPu struct
type SubPu struct {
	Id    int `json:"id"`
	ParId int `json:"parid"`
	SubId int `json:"subid"`
}

//AddSubPu struct
type AddSubPu struct {
	ParId int `json:"parid"`
	SubId int `json:"subid"`
}

//SubPu_count  struct
type SubPu_count struct {
	Values []SubPu `json:"values"`
	Count  int     `json:"count"`
	Auth   Auth    `json:"auth"`
}
