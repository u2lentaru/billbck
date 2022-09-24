package models

//ContractMot struct
type ContractMot struct {
	Id                int    `json:"id"`
	ContractMotNameRu string `json:"ContractMotName_RU"`
	ContractMotNameKz string `json:"ContractMotName_KZ"`
}

//AddContractMot struct
type AddContractMot struct {
	ContractMotNameRu string `json:"ContractMotName_RU"`
	ContractMotNameKz string `json:"ContractMotName_KZ"`
}

//ContractMot_count  struct
type ContractMot_count struct {
	Values []ContractMot `json:"values"`
	Count  int           `json:"count"`
	Auth   Auth          `json:"auth"`
}
