package models

//Position struct
type Position struct {
	Id           int    `json:"id"`
	PositionName string `json:"positionname"`
}

//AddPosition struct
type AddPosition struct {
	PositionName string `json:"positionname"`
}

//Position_count  struct
type Position_count struct {
	Values []Position `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
