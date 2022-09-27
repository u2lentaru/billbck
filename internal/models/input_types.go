package models

//InputType struct
type InputType struct {
	Id            int    `json:"id"`
	InputTypeName string `json:"inputtypename"`
}

//AddInputType struct
type AddInputType struct {
	InputTypeName string `json:"inputtypename"`
}

//InputType_count  struct
type InputType_count struct {
	Values []InputType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
