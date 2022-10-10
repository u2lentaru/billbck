package models

//SealType struct
type SealType struct {
	Id           int    `json:"id"`
	SealTypeName string `json:"sealtypename"`
}

//AddSealType struct
type AddSealType struct {
	SealTypeName string `json:"sealtypename"`
}

//SealType_count  struct
type SealType_count struct {
	Values []SealType `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
