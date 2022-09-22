package models

//CalculationType struct
type CalculationType struct {
	Id                  int    `json:"id"`
	CalculationTypeName string `json:"calculationtypename"`
}

//AddCalculationType struct
type AddCalculationType struct {
	CalculationTypeName string `json:"calculationtypename"`
}

//CalculationType_count  struct
type CalculationType_count struct {
	Values []CalculationType `json:"values"`
	Count  int               `json:"count"`
	Auth   Auth              `json:"auth"`
}
