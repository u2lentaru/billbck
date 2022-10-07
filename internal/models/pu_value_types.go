package models

//PuValueType struct
type PuValueType struct {
	Id              int    `json:"id"`
	PuValueTypeName string `json:"puvaluetypename"`
}

//AddPuValueType struct
type AddPuValueType struct {
	PuValueTypeName string `json:"puvaluetypename"`
}

//PuValueType_count  struct
type PuValueType_count struct {
	Values []PuValueType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
