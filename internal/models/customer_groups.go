package models

//CustomerGroup struct
type CustomerGroup struct {
	Id                int    `json:"id"`
	CustomerGroupName string `json:"customergroupname"`
}

//AddCustomerGroup struct
type AddCustomerGroup struct {
	CustomerGroupName string `json:"customergroupname"`
}

//CustomerGroup_count  struct
type CustomerGroup_count struct {
	Values []CustomerGroup `json:"values"`
	Count  int             `json:"count"`
	Auth   Auth            `json:"auth"`
}
