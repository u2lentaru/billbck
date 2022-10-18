package models

//Voltage struct
type Voltage struct {
	Id           int    `json:"id"`
	VoltageName  string `json:"voltagename"`
	VoltageValue int    `json:"voltagevalue"`
}

//AddVoltage struct
type AddVoltage struct {
	VoltageName  string `json:"voltagename"`
	VoltageValue int    `json:"voltagevalue"`
}

//Voltage_count  struct
type Voltage_count struct {
	Values []Voltage `json:"values"`
	Count  int       `json:"count"`
	Auth   Auth      `json:"auth"`
}
