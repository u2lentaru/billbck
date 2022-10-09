package models

//Rp struct
type Rp struct {
	Id             int     `json:"id"`
	RpName         string  `json:"tguname"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
	Tp             Tp      `json:"tp"`
}

//AddRp struct
type AddRp struct {
	RpName         string  `json:"tguname"`
	InvNumber      string  `json:"invnumber"`
	InputVoltage   Voltage `json:"inputvoltage"`
	OutputVoltage1 Voltage `json:"outputvoltage1"`
	OutputVoltage2 Voltage `json:"outputvoltage2"`
	Tp             Tp      `json:"tp"`
}

//Rp_count  struct
type Rp_count struct {
	Values []Rp `json:"values"`
	Count  int  `json:"count"`
	Auth   Auth `json:"auth"`
}
