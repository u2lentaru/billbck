package models

//TransPwrType struct
type TransPwrType struct {
	Id                int     `json:"id"`
	TransPwrTypeName  string  `json:"transpwrtypename"`
	ShortCircuitPower float32 `json:"shortcircuitpower"`
	IdlingLossPower   float32 `json:"idlinglosspower"`
	NominalPower      int     `json:"nominalpower"`
}

//AddTransPwrType struct
type AddTransPwrType struct {
	TransPwrTypeName  string  `json:"transpwrtypename"`
	ShortCircuitPower float32 `json:"shortcircuitpower"`
	IdlingLossPower   float32 `json:"idlinglosspower"`
	NominalPower      int     `json:"nominalpower"`
}

//TransPwrType_count  struct
type TransPwrType_count struct {
	Values []TransPwrType `json:"values"`
	Count  int            `json:"count"`
	Auth   Auth           `json:"auth"`
}
