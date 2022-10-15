package models

//TransType struct
type TransType struct {
	Id            int    `json:"id"`
	TransTypeName string `json:"transtypename"`
	Ratio         int    `json:"ratio"`
	Class         int    `json:"class"`
	MaxCurr       int    `json:"maxcurr"`
	NomCurr       int    `json:"nomcurr"`
}

//AddTransType struct
type AddTransType struct {
	TransTypeName string `json:"transtypename"`
	Ratio         int    `json:"ratio"`
	Class         int    `json:"class"`
	MaxCurr       int    `json:"maxcurr"`
	NomCurr       int    `json:"nomcurr"`
}

//TransType_count  struct
type TransType_count struct {
	Values []TransType `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
