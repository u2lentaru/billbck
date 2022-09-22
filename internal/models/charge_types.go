package models

//ChargeType struct
type ChargeType struct {
	Id             int    `json:"id"`
	ChargeTypeName string `json:"chargetypename"`
}

//AddChargeType struct
type AddChargeType struct {
	ChargeTypeName string `json:"chargetypename"`
}

//ChargeType_count  struct
type ChargeType_count struct {
	Values []ChargeType `json:"values"`
	Count  int          `json:"count"`
	Auth   Auth         `json:"auth"`
}
