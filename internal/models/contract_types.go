package models

//ContractType struct
type ContractType struct {
	Id               int    `json:"id"`
	ContractTypeName string `json:"contracttypename"`
}

//AddContractType struct
type AddContractType struct {
	ContractTypeName string `json:"contracttypename"`
}

//ContractType_count  struct
type ContractType_count struct {
	Values []ContractType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
