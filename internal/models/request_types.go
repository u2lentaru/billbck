package models

//RequestType struct
type RequestType struct {
	Id              int         `json:"id"`
	RequestTypeName string      `json:"requesttypename"`
	RequestKind     RequestKind `json:"requestkind"`
}

//AddRequestType struct
type AddRequestType struct {
	RequestTypeName string      `json:"requesttypename"`
	RequestKind     RequestKind `json:"requestkind"`
}

//RequestType_count struct
type RequestType_count struct {
	Values []RequestType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
