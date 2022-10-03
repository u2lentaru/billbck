package models

//OrderType struct
type OrderType struct {
	Id            int    `json:"id"`
	OrderTypeName string `json:"ordertypename"`
}

//AddOrderType struct
type AddOrderType struct {
	OrderTypeName string `json:"ordertypename"`
}

//OrderType_count  struct
type OrderType_count struct {
	Values []OrderType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
