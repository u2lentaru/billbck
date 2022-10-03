package models

//ObjType struct
type ObjType struct {
	Id          int    `json:"id"`
	ObjTypeName string `json:"objtypename"`
}

//AddObjType struct
type AddObjType struct {
	ObjTypeName string `json:"objtypename"`
}

//ObjType_count  struct
type ObjType_count struct {
	Values []ObjType `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
