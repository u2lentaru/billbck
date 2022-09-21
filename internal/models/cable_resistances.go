package models

//CableResistance struct
type CableResistance struct {
	Id                  int     `json:"id"`
	CableResistanceName string  `json:"cableresistancename"`
	Resistance          float32 `json:"resistance"`
	MaterialType        bool    `json:"materialtype"`
}

//AddCableResistance struct
type AddCableResistance struct {
	CableResistanceName string  `json:"cableresistancename"`
	Resistance          float32 `json:"resistance"`
	MaterialType        bool    `json:"materialtype"`
}

//CableResistance_count  struct
type CableResistance_count struct {
	Values []CableResistance `json:"values"`
	Count  int               `json:"count"`
	Auth   Auth              `json:"auth"`
}
