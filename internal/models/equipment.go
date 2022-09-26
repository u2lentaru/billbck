package models

//Equipment struct
type Equipment struct {
	Id            int           `json:"id"`
	EquipmentType EquipmentType `json:"equipmenttype"`
	Object        Object        `json:"object"`
	Qty           int           `json:"qty"`
	WorkingHours  float32       `json:"workinghours"`
}

//AddEquipment struct
type AddEquipment struct {
	EquipmentType EquipmentType `json:"equipmenttype"`
	Object        Object        `json:"object"`
	Qty           int           `json:"qty"`
	WorkingHours  float32       `json:"workinghours"`
}

//Equipment_count  struct
type Equipment_count struct {
	Values []Equipment `json:"values"`
	Count  int         `json:"count"`
	Auth   Auth        `json:"auth"`
}
