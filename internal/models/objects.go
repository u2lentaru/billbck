package models

//Object struct
type Object struct {
	Id              int             `json:"id"`
	ObjectName      string          `json:"objectname"`
	House           House           `json:"house"`
	FlatNumber      *string         `json:"flatnumber"`
	ObjType         ObjType         `json:"objtype"`
	RegQty          int             `json:"regqty"`
	Uzo             Uzo             `json:"uzo"`
	TariffGroup     TariffGroup     `json:"tariffgroup"`
	CalculationType CalculationType `json:"calculationtype"`
	ObjStatus       ObjStatus       `json:"objstatus"`
	Notes           *string         `json:"notes"`
	MffId           *int            `json:"mffid"`
}

//AddObject struct
type AddObject struct {
	ObjectName      string          `json:"objectname"`
	House           House           `json:"house"`
	FlatNumber      *string         `json:"flatnumber"`
	ObjType         ObjType         `json:"objtype"`
	RegQty          int             `json:"regqty"`
	Uzo             Uzo             `json:"uzo"`
	TariffGroup     TariffGroup     `json:"tariffgroup"`
	CalculationType CalculationType `json:"calculationtype"`
	ObjStatus       ObjStatus       `json:"objstatus"`
	Notes           *string         `json:"notes"`
	MffId           *int            `json:"mffid"`
}

//Object_count  struct
type Object_count struct {
	Values []Object `json:"values"`
	Count  int      `json:"count"`
	Auth   Auth     `json:"auth"`
}
