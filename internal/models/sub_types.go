package models

//SubType struct
type SubType struct {
	Id           int    `json:"id"`
	SubTypeName  string `json:"subtypename"`
	SubTypeDescr string `json:"subtypedescr"`
}

//SubType_count  struct
type SubType_count struct {
	Values []SubType `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}

//AddSubType struct
type AddSubType struct {
	SubTypeName  string `json:"subtypename"`
	SubTypeDescr string `json:"subtypedescr"`
}
