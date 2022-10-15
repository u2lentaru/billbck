package models

//TransPwr struct
type TransPwr struct {
	Id           int          `json:"id"`
	TransPwrName string       `json:"transpwrname"`
	TransPwrType TransPwrType `json:"tranpwrstype"`
}

//AddTransPwr struct
type AddTransPwr struct {
	TransPwrName string       `json:"transpwrname"`
	TransPwrType TransPwrType `json:"tranpwrstype"`
}

//TransPwr_count  struct
type TransPwr_count struct {
	Values []TransPwr `json:"values"`
	Count  int        `json:"count"`
	Auth   Auth       `json:"auth"`
}
