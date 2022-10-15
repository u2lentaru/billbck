package models

//Tgu struct
type Tgu struct {
	Id             int     `json:"id"`
	PId            int     `json:"pid"`
	TguName        string  `json:"tguname"`
	TguType        TguType `json:"tgutype"`
	InvNumber      *string `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
}

//AddTgu struct
type AddTgu struct {
	PId            int     `json:"pid"`
	TguName        string  `json:"tguname"`
	TguType        TguType `json:"tgutype"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
}

//Tgu_count  struct
type Tgu_count struct {
	Values []Tgu `json:"values"`
	Count  int   `json:"count"`
	Auth   Auth  `json:"auth"`
}

//TguType struct
type TguType struct {
	Id          int    `json:"id"`
	TguTypeName string `json:"tgutypename"`
}
