package models

//FormType struct
type FormType struct {
	Id            int    `json:"id"`
	FormTypeName  string `json:"formtypename"`
	FormTypeDescr string `json:"formtypedescr"`
}

//AddFormType struct
type AddFormType struct {
	FormTypeName  string `json:"formtypename"`
	FormTypeDescr string `json:"formtypedescr"`
}

//FormType_count  struct
type FormType_count struct {
	Values []FormType `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
