package models

//PuType struct
type PuType struct {
	Id         int    `json:"id"`
	PuTypeName string `json:"putypename"`
}

//AddPuType struct
type AddPuType struct {
	PuTypeName string `json:"putypename"`
}

//PuType_count  struct
type PuType_count struct {
	Values []PuType `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
