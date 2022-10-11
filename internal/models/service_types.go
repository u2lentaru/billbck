package models

//ServiceType struct
type ServiceType struct {
	Id              int    `json:"id"`
	ServiceTypeName string `json:"servicetypename"`
}

//AddServiceType struct
type AddServiceType struct {
	ServiceTypeName string `json:"servicetypename"`
}

//ServiceType_count struct
type ServiceType_count struct {
	Values []ServiceType `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
