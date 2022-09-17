package models

//ActType struct
type ActType struct {
	Id          int    `json:"id"`
	ActTypeName string `json:"acttypename"`
}

//AddActType struct
type AddActType struct {
	ActTypeName string `json:"acttypename"`
}

//ActType_count  struct
type ActType_count struct {
	Values []ActType `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
