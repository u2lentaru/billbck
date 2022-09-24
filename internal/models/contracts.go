package models

//Contract struct
type Contract struct {
	Id                int           `json:"id"`
	BarCode           string        `json:"barcode"`
	PersonalAccount   int           `json:"personalaccount"`
	ContractNumber    string        `json:"contractnumber"`
	Startdate         string        `json:"startdate"`
	Enddate           *string       `json:"enddate"`
	Customer          Subject       `json:"customer"`
	Consignee         Subject       `json:"consignee"`
	EsoContractNumber string        `json:"esocontractnumber"`
	Eso               Eso           `json:"eso"`
	Area              Area          `json:"area"`
	CustomerGroup     CustomerGroup `json:"customergroup"`
	ContractMot       ContractMot   `json:"contractmot"`
	Notes             *string       `json:"notes"`
	MotNotes          *string       `json:"motnotes"`
}

//AddContract struct
type AddContract struct {
	BarCode           string        `json:"barcode"`
	PersonalAccount   int           `json:"personalaccount"`
	ContractNumber    string        `json:"contractnumber"`
	Startdate         string        `json:"startdate"`
	Enddate           *string       `json:"enddate"`
	Customer          Subject       `json:"customer"`
	Consignee         Subject       `json:"consignee"`
	EsoContractNumber string        `json:"esocontractnumber"`
	Eso               Eso           `json:"eso"`
	Area              Area          `json:"area"`
	CustomerGroup     CustomerGroup `json:"customergroup"`
	ContractMot       ContractMot   `json:"contractmot"`
	Notes             *string       `json:"notes"`
	MotNotes          *string       `json:"motnotes"`
}

//Contract_count  struct
type Contract_count struct {
	Values []Contract `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
