package models

//Reason struct
type Reason struct {
	Id         int    `json:"id"`
	ReasonName string `json:"reasonname"`
}

//AddReason struct
type AddReason struct {
	ReasonName string `json:"reasonname"`
}

//Reason_count  struct
type Reason_count struct {
	Values []Reason `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
